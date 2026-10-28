package supportlink

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/SteamServerUI/StationeersServerUI/v6/src/connectivity"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/diagnostics"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/discordbot"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/logger"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/managers/gamemgr"
)

var commands = map[string]func() response{
	"snapshot":          snapshotCommand,
	"status":            statusCommand,
	"diagnostics":       diagnosticsCommand,
	"report":            reportCommand,
	"config":            configCommand,
	"players":           playersCommand,
	"backups":           backupsCommand,
	"mods":              modsCommand,
	"routines":          routinesCommand,
	"logs":              logsCommand,
	"connectivity":      connectivityCommand,
	"package":           packageCommand,
	"record_start":      recordStartCommand,
	"record_stop":       recordStopCommand,
	"recording_package": recordingPackageCommand,
	"server_state":      serverStateCommand,
	"server_start":      serverActionCommand("start"),
	"server_stop":       serverActionCommand("stop"),
	"server_restart":    serverActionCommand("restart"),
}

func execute(command string) response {
	run, ok := commands[command]
	if !ok {
		return response{Error: "this command is not allowed by SSUI"}
	}
	return run()
}

func statusCommand() response {
	report := diagnostics.Run()
	system := report.System
	var text strings.Builder
	fmt.Fprintf(&text, "SSUI %s (%s), %s/%s", system.Version, system.Branch, system.OS, system.Architecture)
	if system.OSRelease != "" {
		fmt.Fprintf(&text, " %s", system.OSRelease)
	}
	if system.Container {
		text.WriteString(", container")
	}
	fmt.Fprintf(&text, "\nGame server: %s", system.ServerState)
	if system.UptimeSeconds > 0 {
		fmt.Fprintf(&text, ", up %s", time.Duration(system.UptimeSeconds)*time.Second)
	}
	if system.GameBuild != "" {
		fmt.Fprintf(&text, ", build %s", system.GameBuild)
	}
	fmt.Fprintf(&text, "\nDiagnostics: %s", report.Overall)
	for _, finding := range report.Findings {
		fmt.Fprintf(&text, "\n- %s (%s): %s", finding.Code, finding.Status, finding.Summary)
	}
	return response{OK: true, Text: text.String(), Data: map[string]any{"system": system, "overall": report.Overall, "findings": report.Findings}}
}

func diagnosticsCommand() response {
	report := diagnostics.RunFull()
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return response{Error: err.Error()}
	}
	summary := statusCommand()
	return response{OK: true, Text: summary.Text, Data: report, File: &file{Name: "diagnosis.json", Data: data}}
}

func snapshotCommand() response {
	return response{OK: true, Text: "Full server snapshot", Data: buildSnapshot()}
}

func reportCommand() response {
	report := diagnostics.RunFull()
	var buffer bytes.Buffer
	if err := diagnostics.WriteReportHTML(&buffer, report); err != nil {
		return response{Error: err.Error()}
	}
	return response{OK: true, Text: fmt.Sprintf("Support report (%s)", report.Overall), Data: report, File: &file{Name: "report.html", Data: buffer.Bytes()}}
}

func configCommand() response {
	report := diagnostics.Run()
	data, err := json.MarshalIndent(report.Config, "", "  ")
	if err != nil {
		return response{Error: err.Error()}
	}
	return response{OK: true, Text: string(data), Data: report.Config}
}

func playersCommand() response {
	players := connectedPlayers()
	var text strings.Builder
	fmt.Fprintf(&text, "%d connected", len(players))
	for _, player := range players {
		fmt.Fprintf(&text, "\n- %s (%s)", player.Name, player.SteamID)
	}
	return response{OK: true, Text: text.String(), Data: players}
}

func backupsCommand() response {
	backups := recentBackups(25)
	var text strings.Builder
	fmt.Fprintf(&text, "%d newest backups", len(backups))
	for _, backup := range backups {
		fmt.Fprintf(&text, "\n- %s  %s", backup.SavedAt.Format("2006-01-02 15:04"), backup.Name)
	}
	return response{OK: true, Text: text.String(), Data: backups}
}

func modsCommand() response {
	mods := installedMods()
	var text strings.Builder
	fmt.Fprintf(&text, "LaunchPad enabled: %v, %d mods", mods.LaunchPadEnabled, len(mods.Installed))
	for _, mod := range mods.Installed {
		fmt.Fprintf(&text, "\n- %s %s", mod.Name, mod.Version)
	}
	return response{OK: true, Text: text.String(), Data: mods}
}

func routinesCommand() response {
	list := routines()
	var text strings.Builder
	fmt.Fprintf(&text, "%d routines", len(list))
	for _, routine := range list {
		fmt.Fprintf(&text, "\n- %s: %s %q -> %s", routine.Name, routine.Trigger, routine.Pattern, strings.Join(routine.Actions, ", "))
	}
	return response{OK: true, Text: text.String(), Data: list}
}

func logsCommand() response {
	logs := recentLogs(logsCommandLines)
	text := "== SSUI ==\n" + strings.Join(logs.SSUI, "\n") + "\n\n== Game server ==\n" + strings.Join(logs.Game, "\n")
	return response{OK: true, Text: fmt.Sprintf("Last %d SSUI and game server lines", logsCommandLines), Data: logs, File: &file{Name: "recent-logs.txt", Data: []byte(text)}}
}

func connectivityCommand() response {
	result, err := connectivity.Check()
	if errors.Is(err, connectivity.ErrCheckInProgress) {
		return response{Error: "a connectivity check is already running"}
	}
	var text strings.Builder
	fmt.Fprintf(&text, "Spacecat: %s", result.State)
	if result.Address != "" {
		fmt.Fprintf(&text, " (port %d)", result.Port)
	}
	for _, probe := range result.Probes {
		fmt.Fprintf(&text, "\n- %d bytes: received %v", probe.Size, probe.Received)
		if probe.Received {
			fmt.Fprintf(&text, ", %d ms", probe.LatencyMillis)
		}
	}
	if result.Message != "" {
		fmt.Fprintf(&text, "\n%s", result.Message)
	}
	if result.Error != "" {
		fmt.Fprintf(&text, "\nError: %s", result.Error)
	}
	return response{OK: err == nil, Text: text.String(), Data: result}
}

func packageCommand() response {
	return zipResponse("ssui-support", diagnostics.WritePackage)
}

func recordingPackageCommand() response {
	reply := zipResponse("ssui-support-recording", diagnostics.WriteRecordingPackage)
	if reply.Error == diagnostics.ErrNoRecording.Error() {
		reply.Error = "there is no stopped recording; run record_start, let the user reproduce the issue, then record_stop"
	}
	return reply
}

func recordStartCommand() response {
	_, err := logger.StartSupportRecording()
	if errors.Is(err, logger.ErrSupportRecordingActive) {
		return response{Error: "a support recording is already running"}
	}
	if err != nil {
		return response{Error: err.Error()}
	}
	return response{OK: true, Text: "Support recording started. Ask the user to reproduce the issue, then run record_stop.", Data: logger.SupportRecordingStatus()}
}

func recordStopCommand() response {
	recording, err := diagnostics.StopRecording()
	if errors.Is(err, logger.ErrNoSupportRecording) {
		return response{Error: "no support recording is running"}
	}
	text := fmt.Sprintf("Support recording stopped (%.1f MB). Run recording_package to fetch it.", float64(recording.Bytes)/(1<<20))
	if err != nil {
		text += "\nThe recording stopped with an error: " + err.Error()
	}
	return response{OK: true, Text: text, Data: recording}
}

func serverStateCommand() response {
	state := map[string]any{"state": gamemgr.GetServerState(), "running": gamemgr.InternalIsServerRunning()}
	text := fmt.Sprintf("Game server: %s", state["state"])
	if uptime := gamemgr.GetServerUptime(); uptime > 0 {
		state["uptimeSeconds"] = int64(uptime / time.Second)
		text += fmt.Sprintf(", up %s", uptime.Round(time.Second))
	}
	return response{OK: true, Text: text, Data: state}
}

func serverActionCommand(action string) func() response {
	return func() response {
		message, err := discordbot.RunServerAction(action)
		if err != nil {
			return response{Error: err.Error()}
		}
		reply := serverStateCommand()
		reply.Text = message + "\n" + reply.Text
		return reply
	}
}

func zipResponse(name string, write func(io.Writer) error) response {
	var buffer bytes.Buffer
	if err := write(&limitedWriter{buffer: &buffer, limit: maxFileBytes + 1}); err != nil {
		if errors.Is(err, errTooLarge) {
			return response{Error: "the package is larger than 8 MB; ask the user to download it from the WebUI"}
		}
		return response{Error: err.Error()}
	}
	stamp := time.Now().Format("20060102-150405")
	return response{OK: true, Text: fmt.Sprintf("%s (%.1f MB)", name, float64(buffer.Len())/(1<<20)), File: &file{Name: name + "-" + stamp + ".zip", Data: buffer.Bytes()}}
}

var errTooLarge = errors.New("support response too large")

type limitedWriter struct {
	buffer *bytes.Buffer
	limit  int
}

func (w *limitedWriter) Write(data []byte) (int, error) {
	if w.buffer.Len()+len(data) > w.limit {
		return 0, errTooLarge
	}
	return w.buffer.Write(data)
}
