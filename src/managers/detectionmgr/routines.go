package detectionmgr

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/SteamServerUI/StationeersServerUI/v6/src/config"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/core/ssestream"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/discordbot"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/logger"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/managers/commandmgr"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/managers/gamemgr"
)

type RoutineAction struct {
	Type  string `json:"type"`
	Value string `json:"value,omitempty"`
	Text  string `json:"text,omitempty"`
}

const (
	maxRoutineActions    = 10
	maxWaitSeconds       = 3600
	maxCooldownSeconds   = 86400
	serverActionCooldown = 2 * time.Minute
)

var serverActions = map[string]bool{"start": true, "stop": true, "restart": true, "restore": true}

var textActions = map[string]bool{"eventlog": true, "discord_hub": true, "announce": true, "command": true}

var otherActions = map[string]bool{"save": true, "wait": true, "webhook": true}

var TriggerEvents = []EventType{
	EventPlayerConnecting, EventPlayerReady, EventPlayerDisconnect, EventWorldSaved,
	EventSessionRegistered, EventServerRunning, EventException, EventSettingsChanged, EventNewGameStarted,
}

var routineRuns = struct {
	sync.Mutex
	last map[string]time.Time
}{last: map[string]time.Time{}}

var webhookClient = &http.Client{
	Timeout:       10 * time.Second,
	CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
}

func ValidateDetection(detection CustomDetection) error {
	switch detection.Type {
	case "keyword":
	case "regex":
		if _, err := regexp.Compile(detection.Pattern); err != nil {
			return fmt.Errorf("invalid regex pattern: %w", err)
		}
	case "event":
		if !isTriggerEvent(EventType(detection.Pattern)) {
			return errors.New("unknown trigger event")
		}
	default:
		return errors.New("type must be keyword, regex or event")
	}
	if strings.TrimSpace(detection.Pattern) == "" {
		return errors.New("pattern cannot be empty")
	}
	if detection.CooldownSeconds < 0 || detection.CooldownSeconds > maxCooldownSeconds {
		return fmt.Errorf("cooldown must be between 0 and %d seconds", maxCooldownSeconds)
	}
	if len(detection.Actions) == 0 || len(detection.Actions) > maxRoutineActions {
		return fmt.Errorf("add between 1 and %d actions", maxRoutineActions)
	}
	for i, action := range detection.Actions {
		if err := validateAction(action); err != nil {
			return fmt.Errorf("action %d: %w", i+1, err)
		}
	}
	return nil
}

func validateAction(action RoutineAction) error {
	switch {
	case textActions[action.Type]:
		if strings.TrimSpace(action.Value) == "" {
			return errors.New("text cannot be empty")
		}
	case action.Type == "wait":
		seconds, err := strconv.Atoi(action.Value)
		if err != nil || seconds < 1 || seconds > maxWaitSeconds {
			return fmt.Errorf("wait must be between 1 and %d seconds", maxWaitSeconds)
		}
	case action.Type == "webhook":
		target, err := url.Parse(action.Value)
		if err != nil || (target.Scheme != "http" && target.Scheme != "https") || target.Host == "" {
			return errors.New("webhook needs an http or https URL")
		}
	case serverActions[action.Type], otherActions[action.Type]:
	default:
		return errors.New("unknown action")
	}
	return nil
}

func isTriggerEvent(eventType EventType) bool {
	for _, candidate := range TriggerEvents {
		if candidate == eventType {
			return true
		}
	}
	return false
}

func runRoutine(routine CustomPattern, values []string) {
	if !claimRoutine(routine) {
		return
	}
	go func() {
		for _, action := range routine.Actions {
			if err := runAction(routine, action, values); err != nil {
				logger.Detection.Warnf("Routine %q stopped at %s: %v", routine.Name, action.Type, err)
				discordbot.SendMessageToEventLogChannel(fmt.Sprintf("❌ Routine %s stopped at %s: %v", routine.Name, action.Type, err))
				return
			}
		}
	}()
}

func claimRoutine(routine CustomPattern) bool {
	cooldown := routine.Cooldown
	for _, action := range routine.Actions {
		if serverActions[action.Type] && cooldown < serverActionCooldown {
			cooldown = serverActionCooldown
		}
	}
	routineRuns.Lock()
	defer routineRuns.Unlock()
	if last, ok := routineRuns.last[routine.ID]; ok && time.Since(last) < cooldown {
		return false
	}
	routineRuns.last[routine.ID] = time.Now()
	return true
}

func runAction(routine CustomPattern, action RoutineAction, values []string) error {
	text := formatMessage(action.Value, values)
	switch action.Type {
	case "eventlog":
		message := "🎮 [Custom Detection] " + text
		logger.Detection.Info(message)
		ssestream.BroadcastDetectionEvent(message)
		discordbot.SendMessageToEventLogChannel(message)
	case "discord_hub":
		discordbot.SendMessageToHubChannel(text)
	case "announce":
		return sscmCommand("announce " + text)
	case "command":
		return sscmCommand(text)
	case "save":
		return sscmCommand("save")
	case "wait":
		seconds, _ := strconv.Atoi(action.Value)
		time.Sleep(time.Duration(seconds) * time.Second)
	case "webhook":
		return postWebhook(action.Value, routine.Name, formatMessage(action.Text, values), values)
	default:
		message, err := discordbot.RunServerAction(action.Type)
		if err == nil {
			logger.Detection.Infof("Routine %q: %s", routine.Name, message)
		}
		return err
	}
	return nil
}

func sscmCommand(command string) error {
	if !config.GetIsSSCMEnabled() || !gamemgr.InternalIsServerRunning() {
		return errors.New("the game server must be running with SSCM enabled")
	}
	return commandmgr.WriteCommand(command)
}

func postWebhook(target, routine, text string, values []string) error {
	body, err := json.Marshal(map[string]any{"content": text, "routine": routine, "values": values})
	if err != nil {
		return err
	}
	response, err := webhookClient.Post(target, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	response.Body.Close()
	if response.StatusCode >= 300 {
		return fmt.Errorf("webhook answered %s", response.Status)
	}
	return nil
}
