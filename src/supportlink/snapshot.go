package supportlink

import (
	"sort"
	"time"

	"github.com/SteamServerUI/StationeersServerUI/v6/src/config"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/diagnostics"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/logger"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/managers/backupmgr"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/managers/detectionmgr"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/modding"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/setup/update"
)

const (
	snapshotBackups  = 10
	snapshotLogLines = 50
	logsCommandLines = 200
)

type ServerSnapshot struct {
	GeneratedAt  time.Time               `json:"generatedAt"`
	Diagnostics  diagnostics.Report      `json:"diagnostics"`
	Players      []Player                `json:"players"`
	MaxPlayers   string                  `json:"maxPlayers"`
	Backups      []Backup                `json:"backups"`
	Mods         Mods                    `json:"mods"`
	Routines     []Routine               `json:"routines"`
	Integrations Integrations            `json:"integrations"`
	Update       UpdateInfo              `json:"update"`
	Recording    logger.SupportRecording `json:"recording"`
	Logs         Logs                    `json:"logs"`
}

type Player struct {
	Name    string `json:"name"`
	SteamID string `json:"steamId"`
}

type Backup struct {
	Name        string    `json:"name"`
	SavedAt     time.Time `json:"savedAt"`
	WorldName   string    `json:"worldName,omitempty"`
	GameVersion string    `json:"gameVersion,omitempty"`
	DaysPlayed  int64     `json:"daysPlayed,omitempty"`
	ArchiveSize int64     `json:"archiveSize,omitempty"`
}

type Mods struct {
	LaunchPadEnabled bool  `json:"launchPadEnabled"`
	Installed        []Mod `json:"installed"`
}

type Mod struct {
	Name       string `json:"name"`
	Author     string `json:"author,omitempty"`
	Version    string `json:"version,omitempty"`
	WorkshopID string `json:"workshopId,omitempty"`
}

type Routine struct {
	Name            string   `json:"name,omitempty"`
	Trigger         string   `json:"trigger"`
	Pattern         string   `json:"pattern"`
	Actions         []string `json:"actions"`
	CooldownSeconds int      `json:"cooldownSeconds,omitempty"`
}

type Integrations struct {
	Discord     bool `json:"discord"`
	SSCM        bool `json:"sscm"`
	CLIConsole  bool `json:"cliConsole"`
	AutoUpdates bool `json:"autoUpdates"`
	Auth        bool `json:"auth"`
}

type UpdateInfo struct {
	CurrentVersion   string     `json:"currentVersion"`
	UpdateAvailable  bool       `json:"updateAvailable"`
	LatestStable     string     `json:"latestStable,omitempty"`
	LatestPrerelease string     `json:"latestPrerelease,omitempty"`
	ContainerManaged bool       `json:"containerManaged"`
	LastError        string     `json:"lastError,omitempty"`
	CheckedAt        *time.Time `json:"checkedAt,omitempty"`
}

type Logs struct {
	SSUI []string `json:"ssui"`
	Game []string `json:"game"`
}

func buildSnapshot() ServerSnapshot {
	return ServerSnapshot{
		GeneratedAt:  time.Now().UTC(),
		Diagnostics:  diagnostics.Run(),
		Players:      connectedPlayers(),
		MaxPlayers:   config.GetServerMaxPlayers(),
		Backups:      recentBackups(snapshotBackups),
		Mods:         installedMods(),
		Routines:     routines(),
		Integrations: integrations(),
		Update:       updateInfo(),
		Recording:    logger.SupportRecordingStatus(),
		Logs:         recentLogs(snapshotLogLines),
	}
}

func connectedPlayers() (players []Player) {
	players = []Player{}
	defer func() {
		if recover() != nil {
			players = []Player{}
		}
	}()
	for steamID, name := range detectionmgr.GetPlayers(detectionmgr.GetDetector()) {
		players = append(players, Player{Name: name, SteamID: steamID})
	}
	sort.Slice(players, func(i, j int) bool { return players[i].Name < players[j].Name })
	return players
}

func recentBackups(limit int) []Backup {
	backups := []Backup{}
	manager := backupmgr.CurrentBackupManager()
	if manager == nil {
		return backups
	}
	saves, err := manager.ListBackups(limit)
	if err != nil {
		return backups
	}
	for _, save := range saves {
		backup := Backup{Name: save.Name, SavedAt: save.SaveTime}
		if save.SummaryReady {
			backup.WorldName = save.Summary.WorldName
			backup.GameVersion = save.Summary.GameVersion
			backup.DaysPlayed = save.Summary.DaysPlayed
			backup.ArchiveSize = save.Summary.ArchiveSize
		}
		backups = append(backups, backup)
	}
	return backups
}

func installedMods() Mods {
	mods := Mods{LaunchPadEnabled: config.GetIsStationeersLaunchPadEnabled(), Installed: []Mod{}}
	for _, mod := range modding.GetModList() {
		mods.Installed = append(mods.Installed, Mod{Name: mod.Name, Author: mod.Author, Version: mod.Version, WorkshopID: mod.WorkshopHandle})
	}
	return mods
}

func routines() []Routine {
	result := []Routine{}
	for _, detection := range detectionmgr.GetCustomDetections() {
		routine := Routine{Name: detection.Name, Trigger: detection.Type, Pattern: detection.Pattern, Actions: []string{}, CooldownSeconds: detection.CooldownSeconds}
		for _, action := range detection.Actions {
			routine.Actions = append(routine.Actions, action.Type)
		}
		result = append(result, routine)
	}
	return result
}

func integrations() Integrations {
	return Integrations{
		Discord:     config.GetIsDiscordEnabled(),
		SSCM:        config.GetIsSSCMEnabled(),
		CLIConsole:  config.GetIsConsoleEnabled(),
		AutoUpdates: config.GetIsUpdateEnabled(),
		Auth:        config.GetAuthEnabled(),
	}
}

func updateInfo() UpdateInfo {
	status := update.StatusSnapshot()
	info := UpdateInfo{
		CurrentVersion:   status.CurrentVersion,
		UpdateAvailable:  status.UpdateAvailable,
		ContainerManaged: status.ContainerManaged,
		LastError:        status.LastError,
		CheckedAt:        status.CheckedAt,
	}
	if status.LatestStable != nil {
		info.LatestStable = status.LatestStable.Version
	}
	if status.LatestPrerelease != nil {
		info.LatestPrerelease = status.LatestPrerelease.Version
	}
	return info
}

func recentLogs(lines int) Logs {
	return Logs{SSUI: redactLines(logger.RecentLines(false, lines)), Game: redactLines(logger.RecentLines(true, lines))}
}

func redactLines(lines []string) []string {
	redacted := make([]string, len(lines))
	for i, line := range lines {
		redacted[i] = diagnostics.Redact(line)
	}
	return redacted
}
