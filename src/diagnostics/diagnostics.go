package diagnostics

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/SteamServerUI/StationeersServerUI/v6/src/config"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/connectivity"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/managers/gamemgr"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/steamcmd"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/worldgen"
)

type Status string

const (
	StatusWarn    Status = "warn"
	StatusFail    Status = "fail"
	StatusUnknown Status = "unknown"
)

const (
	codeRoot        = "RootCat"
	codeConfigFile  = "HairballConfig"
	codeGame        = "EmptyBoxCat"
	codeWorldArgs   = "TangledYarnCat"
	codeWorld       = "LostWorldCat"
	codeSSUIFolder  = "Folder404Cat"
	codeLogFolder   = "Log4Cat"
	codeSteamCMD    = "SteamedCat"
	codeSSCM        = "CatnipMissing"
	codeSave        = "SaveMeow"
	codePorts       = "Portcat"
	codeUnreachable = "Spacecat"
	codeJitter      = "Jittercat"
	codeUnverified  = "SchrodingersCat"
	codeReadOnly    = "LockedBoxCat"
	codeDisk        = "FullBowlCat"
	codeOwner       = "StrayCat"
	codeAccess      = "ScratchedCat"
	codeDNS         = "LostCat"
)

type Finding struct {
	Code    string `json:"code"`
	Status  Status `json:"status"`
	Summary string `json:"summary"`
	Fix     string `json:"fix,omitempty"`
}

type Report struct {
	GeneratedAt time.Time           `json:"generatedAt"`
	Overall     string              `json:"overall"`
	Findings    []Finding           `json:"findings"`
	System      System              `json:"system"`
	Config      EffectiveConfig     `json:"config"`
	Paths       map[string]string   `json:"paths"`
	World       World               `json:"world"`
	Spacecat    connectivity.Status `json:"spacecat"`
	Environment Environment         `json:"environment"`
	Access      *AccessScan         `json:"access,omitempty"`
	DNS         []DNSResult         `json:"dns,omitempty"`
}

type System struct {
	Version       string `json:"version"`
	Branch        string `json:"branch"`
	GameBuild     string `json:"gameBuild,omitempty"`
	OS            string `json:"os"`
	OSRelease     string `json:"osRelease,omitempty"`
	Architecture  string `json:"architecture"`
	GoVersion     string `json:"goVersion"`
	Container     bool   `json:"container"`
	Root          bool   `json:"root"`
	AuthEnabled   bool   `json:"authEnabled"`
	ServerState   string `json:"serverState"`
	UptimeSeconds int64  `json:"uptimeSeconds,omitempty"`
}

type EffectiveConfig struct {
	GameBranch         string `json:"gameBranch"`
	GamePort           string `json:"gamePort"`
	UpdatePort         string `json:"updatePort"`
	SSUIWebPort        string `json:"ssuiWebPort"`
	LocalIP            string `json:"localIpAddress"`
	AdvertiserOverride string `json:"advertiserOverride,omitempty"`
	SaveName           string `json:"saveName"`
	WorldID            string `json:"worldId"`
	Difficulty         string `json:"difficulty,omitempty"`
	StartCondition     string `json:"startCondition,omitempty"`
	StartLocation      string `json:"startLocation,omitempty"`
	ServerMaxPlayers   string `json:"serverMaxPlayers"`
	SaveInterval       string `json:"saveInterval"`
	AutoSave           bool   `json:"autoSave"`
	SSCMEnabled        bool   `json:"sscmEnabled"`
	DiscordEnabled     bool   `json:"discordEnabled"`
	ConnectivityCheck  bool   `json:"connectivityCheckEnabled"`
	SSUILogFile        bool   `json:"ssuiLogFileEnabled"`
	GameServerLogFile  bool   `json:"gameServerLogFileEnabled"`
	HasServerPassword  bool   `json:"serverPasswordConfigured"`
	HasAuthSecret      bool   `json:"serverAuthSecretConfigured"`
	HasAdminPassword   bool   `json:"adminPasswordConfigured"`
	HasDiscordToken    bool   `json:"discordTokenConfigured"`
}

type World struct {
	Catalog string `json:"catalog"`
	Worlds  int    `json:"worlds"`
}

func Run() Report {
	cfg := config.Snapshot()
	report := Report{
		GeneratedAt: time.Now().UTC(),
		Findings:    []Finding{},
		System:      collectSystem(),
		Config:      summarizeConfig(cfg),
		Paths:       map[string]string{},
		Spacecat:    connectivity.StatusSnapshot(),
		Environment: collectEnvironment(),
	}
	report.checkSystem()
	report.checkGame(cfg)
	report.checkFolders(cfg)
	report.checkPorts(cfg)
	report.checkSpacecat()
	report.checkEnvironment()
	report.summarize()
	return report
}

func RunFull() Report {
	report := Run()
	scan := scanAccess(accessScanLimit)
	report.Access = &scan
	if scan.Problems > 0 {
		report.add(codeAccess, StatusFail, fmt.Sprintf("SSUI cannot fully access %d files or folders in its data or save folders.", scan.Problems), "Fix the owner and permissions of the listed paths, for example with chown -R.")
	}
	report.DNS = checkDNS()
	report.checkDNS()
	report.summarize()
	return report
}

func (report *Report) RecheckSpacecat() {
	report.Spacecat, _ = connectivity.Check()
	kept := report.Findings[:0]
	for _, finding := range report.Findings {
		if finding.Code != codeUnreachable && finding.Code != codeJitter && finding.Code != codeUnverified {
			kept = append(kept, finding)
		}
	}
	report.Findings = kept
	report.checkSpacecat()
	report.summarize()
}

func (report *Report) summarize() {
	report.Overall = "healthy"
	for _, finding := range report.Findings {
		switch {
		case finding.Status == StatusFail:
			report.Overall = "failed"
		case finding.Status == StatusWarn && report.Overall != "failed":
			report.Overall = "degraded"
		case finding.Status == StatusUnknown && report.Overall == "healthy":
			report.Overall = "unknown"
		}
	}
}

func (report *Report) add(code string, status Status, summary, fix string) {
	report.Findings = append(report.Findings, Finding{Code: code, Status: status, Summary: summary, Fix: fix})
}

func collectSystem() System {
	system := System{
		Version:      config.GetVersion(),
		Branch:       config.GetBranch(),
		GameBuild:    config.GetCurrentBranchBuildID(),
		OS:           runtime.GOOS,
		OSRelease:    osRelease(),
		Architecture: runtime.GOARCH,
		GoVersion:    runtime.Version(),
		Container:    config.GetIsDockerContainer(),
		Root:         runtime.GOOS != "windows" && os.Geteuid() == 0,
		AuthEnabled:  config.GetAuthEnabled(),
		ServerState:  string(gamemgr.GetServerState()),
	}
	if uptime := gamemgr.GetServerUptime(); uptime > 0 {
		system.UptimeSeconds = int64(uptime / time.Second)
	}
	return system
}

func osRelease() string {
	switch runtime.GOOS {
	case "linux":
		data, _ := os.ReadFile("/etc/os-release")
		for _, line := range strings.Split(string(data), "\n") {
			if value, ok := strings.CutPrefix(line, "PRETTY_NAME="); ok {
				return strings.Trim(value, `"`)
			}
		}
	case "windows":
		output, _ := exec.Command("cmd", "/c", "ver").Output()
		return strings.TrimSpace(string(output))
	}
	return ""
}

func summarizeConfig(cfg config.JsonConfig) EffectiveConfig {
	return EffectiveConfig{
		GameBranch:         cfg.GameBranch,
		GamePort:           cfg.GamePort,
		UpdatePort:         cfg.UpdatePort,
		SSUIWebPort:        cfg.SSUIWebPort,
		LocalIP:            cfg.LocalIpAddress,
		AdvertiserOverride: cfg.AdvertiserOverride,
		SaveName:           cfg.SaveName,
		WorldID:            cfg.WorldID,
		Difficulty:         cfg.Difficulty,
		StartCondition:     cfg.StartCondition,
		StartLocation:      cfg.StartLocation,
		ServerMaxPlayers:   cfg.ServerMaxPlayers,
		SaveInterval:       cfg.SaveInterval,
		AutoSave:           boolValue(cfg.AutoSave),
		SSCMEnabled:        boolValue(cfg.IsSSCMEnabled),
		DiscordEnabled:     boolValue(cfg.IsDiscordEnabled),
		ConnectivityCheck:  boolValue(cfg.ConnectivityCheckEnabled),
		SSUILogFile:        boolValue(cfg.CreateSSUILogFile),
		GameServerLogFile:  boolValue(cfg.CreateGameServerLogFile),
		HasServerPassword:  cfg.ServerPassword != "",
		HasAuthSecret:      cfg.ServerAuthSecret != "",
		HasAdminPassword:   cfg.AdminPassword != "",
		HasDiscordToken:    cfg.DiscordToken != "",
	}
}

func (report *Report) checkSystem() {
	if report.System.Root && !report.System.Container {
		report.add(codeRoot, StatusWarn, "SSUI is running as root.", "Run SSUI as a regular user.")
	}
}

func (report *Report) checkGame(cfg config.JsonConfig) {
	exe := strings.TrimSpace(cfg.ExePath)
	info, err := os.Stat(exe)
	switch {
	case exe == "":
		report.add(codeGame, StatusFail, "The Stationeers executable is not configured.", "Configure ExePath or run the installation flow.")
	case err != nil:
		report.add(codeGame, StatusFail, "The configured Stationeers executable does not exist.", "Check ExePath and the game installation.")
	case !info.Mode().IsRegular() || !executable(info.Mode(), exe):
		report.add(codeGame, StatusFail, "The configured game path is not an executable file.", "Point ExePath at the Stationeers server executable.")
	case inspectPath(worldgen.DataRootFromExecutable(exe), true) != "ok":
		report.add(codeGame, StatusWarn, "The game data folder next to the executable is missing.", "Validate the Stationeers installation.")
	}

	switch {
	case cfg.WorldID == "":
		report.add(codeWorldArgs, StatusFail, "No WorldID is configured.", "Configure a WorldID before creating a world.")
	case cfg.StartCondition != "" && cfg.Difficulty == "":
		report.add(codeWorldArgs, StatusFail, "StartCondition is set without Difficulty.", "Set Difficulty or clear StartCondition.")
	case cfg.StartLocation != "" && cfg.StartCondition == "":
		report.add(codeWorldArgs, StatusFail, "StartLocation is set without StartCondition.", "Set StartCondition and Difficulty, or clear StartLocation.")
	}

	report.World.Catalog = "unavailable"
	if err != nil || exe == "" {
		return
	}
	catalog, err := worldgen.DefaultLoader.Load(exe)
	if err != nil || catalog.Source != "game-files" {
		report.add(codeWorld, StatusUnknown, "The installed world catalog could not be read, so the world settings are unverified.", "Check ExePath and the game data.")
		return
	}
	report.World = World{Catalog: catalog.Source, Worlds: len(catalog.Worlds)}
	known := cfg.WorldID == "" || worldHas(catalog, cfg.WorldID, "", nil)
	known = known && (cfg.Difficulty == "" || hasOption(catalog.Difficulties, cfg.Difficulty))
	known = known && (cfg.StartCondition == "" || worldHas(catalog, cfg.WorldID, cfg.StartCondition, func(world worldgen.World) []worldgen.Option { return world.StartConditions }))
	known = known && (cfg.StartLocation == "" || worldHas(catalog, cfg.WorldID, cfg.StartLocation, func(world worldgen.World) []worldgen.Option { return world.StartLocations }))
	if !known {
		report.add(codeWorld, StatusWarn, "Some world settings are not in the installed game's catalog.", "Check them against the installed game version. Modded or newer values can still be valid.")
	}
}

func (report *Report) checkFolders(cfg config.JsonConfig) {
	paths := report.Paths
	paths["ssui"] = inspectPath(config.GetSSUIFolder(), true)
	paths["config"] = inspectPath(config.GetConfigPath(), false)
	paths["logs"] = inspectPath(config.GetLogFolder(), true)
	paths["saves"] = inspectPath("./saves", true)
	paths["sscm"] = inspectPath(config.GetSSCMPluginDir(), true)
	if runtime.GOOS == "windows" {
		paths["steamcmd"] = inspectPath(filepath.Join(steamcmd.SteamCMDWindowsDir, "steamcmd.exe"), false)
	} else {
		paths["steamcmd"] = inspectPath(filepath.Join(steamcmd.SteamCMDLinuxDir, "linux32", "steamcmd"), false)
	}

	if paths["ssui"] != "ok" {
		report.add(codeSSUIFolder, StatusFail, "The SSUI data folder is "+paths["ssui"]+".", "Check the SSUI folder and its permissions.")
	}
	switch paths["config"] {
	case "ok":
	case "missing":
		report.add(codeConfigFile, StatusWarn, "The config file does not exist, so settings will not survive a restart.", "Check that the SSUI config folder is persistent.")
	default:
		report.add(codeConfigFile, StatusFail, "The config file is "+paths["config"]+".", "Check the config file and its permissions.")
	}
	if boolValue(cfg.CreateSSUILogFile) && paths["logs"] != "ok" {
		report.add(codeLogFolder, StatusFail, "File logging is on, but the log folder is "+paths["logs"]+".", "Check the log folder and its permissions.")
	}
	if !config.GetSkipSteamCMD() && paths["steamcmd"] != "ok" {
		report.add(codeSteamCMD, StatusWarn, "SteamCMD is "+paths["steamcmd"]+".", "Run the SteamCMD update if SSUI should manage game updates.")
	}
	if boolValue(cfg.IsSSCMEnabled) && paths["sscm"] != "ok" {
		report.add(codeSSCM, StatusFail, "SSCM is enabled, but its plugin folder is "+paths["sscm"]+".", "Reinstall SSCM or disable it.")
	}

	if !safeSaveName(cfg.SaveName) {
		report.add(codeSave, StatusFail, "SaveName is empty or not a plain folder name.", "Use a save name without slashes or '..'.")
		return
	}
	paths["save"] = inspectPath(filepath.Join("./saves", cfg.SaveName), true)
	if paths["save"] != "ok" {
		report.add(codeSave, StatusWarn, "The save folder is "+paths["save"]+".", "Expected before the first start of a new world; otherwise check the saves folder.")
	}
}

func (report *Report) checkPorts(cfg config.JsonConfig) {
	seen := map[int]bool{}
	valid := true
	for _, value := range []string{cfg.GamePort, cfg.UpdatePort, cfg.SSUIWebPort} {
		port, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil || port < 1 || port > 65535 || seen[port] {
			valid = false
		}
		seen[port] = true
	}
	ip := strings.TrimSpace(cfg.LocalIpAddress)
	if !valid || (ip != "" && net.ParseIP(ip) == nil) {
		report.add(codePorts, StatusFail, "A port or the local IP address is invalid, or two ports collide.", "Use distinct ports between 1 and 65535 and a valid local IP.")
	}
}

func (report *Report) checkSpacecat() {
	switch report.Spacecat.State {
	case "reachable", "disabled":
	case "partial":
		report.add(codeJitter, StatusWarn, "The last connectivity check reached only some UDP packet sizes.", "Check port forwarding, firewall rules and MTU filtering.")
	case "service_unavailable":
		report.add(codeUnverified, StatusUnknown, "The connectivity check service could not be reached.", "Try the connectivity check again later.")
	case "unreachable", "bind_failed", "configuration_error":
		report.add(codeUnreachable, StatusWarn, "The last connectivity check could not reach the server from outside.", "Check port forwarding, the advertised address and the firewall.")
	default:
		report.add(codeUnverified, StatusUnknown, "There is no connectivity check result yet.", "Run the connectivity check while the server is stopped.")
	}
}

func inspectPath(path string, wantDirectory bool) string {
	if strings.TrimSpace(path) == "" {
		return "not configured"
	}
	info, err := os.Stat(path)
	if err != nil {
		return "missing"
	}
	if info.IsDir() != wantDirectory {
		return "wrong type"
	}
	file, err := os.Open(path)
	if err != nil {
		return "unreadable"
	}
	file.Close()
	return "ok"
}

func executable(mode os.FileMode, path string) bool {
	if runtime.GOOS == "windows" {
		return strings.EqualFold(filepath.Ext(path), ".exe")
	}
	return mode&0111 != 0
}

func safeSaveName(value string) bool {
	value = strings.TrimSpace(value)
	return value != "" && value != "." && value != ".." && filepath.Base(value) == value &&
		!strings.ContainsAny(value, `/\`) && !strings.ContainsFunc(value, func(r rune) bool { return r < 0x20 || r == 0x7f })
}

func worldHas(catalog worldgen.Catalog, worldID, value string, options func(worldgen.World) []worldgen.Option) bool {
	for _, world := range catalog.Worlds {
		if world.ID == worldID {
			return options == nil || hasOption(options(world), value)
		}
	}
	return false
}

func hasOption(options []worldgen.Option, value string) bool {
	for _, option := range options {
		if option.Value == value {
			return true
		}
	}
	return false
}

func boolValue(value *bool) bool {
	return value != nil && *value
}
