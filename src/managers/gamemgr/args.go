package gamemgr

import (
	"runtime"
	"strconv"
	"strings"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
)

type Arg struct {
	Flag          string
	Value         string
	RequiresValue bool
	Condition     func() bool
	NoQuote       bool
}

func buildCommandArgs() []string {
	var argOrder []Arg

	if config.GetIsNewTerrainAndSaveSystem() {
		argOrder = []Arg{
			{Flag: "-nographics", RequiresValue: false},
			{Flag: "-batchmode", RequiresValue: false},
			/* file start: (expects up to four optional args for:
			-worldid (Optional to LOAD, required to CREATE save, if start value is not found, tries to create map with worldid) -> config.BackupWorldName for legacy reasons
			-difficulty (Optional, defaults to "Normal" if not provided)
			-startcondition  (Optional, defaults to the default start condition for the world setting if not provided.)
			-startlocation (Optional, defaults to "DefaultStartLocation" if not provided.)
			*/
			{Flag: "-file", RequiresValue: false},
			{Flag: "start", Value: config.GetSaveName(), RequiresValue: true},
			{Flag: config.GetWorldID(), RequiresValue: false},
			{Flag: config.GetDifficulty(), RequiresValue: false, Condition: func() bool { return config.GetDifficulty() != "" }},
			{Flag: config.GetStartCondition(), RequiresValue: false, Condition: func() bool { return config.GetStartCondition() != "" }},
			{Flag: config.GetStartLocation(), RequiresValue: false, Condition: func() bool { return config.GetStartLocation() != "" }},
			// file start end
			{Flag: "-logFile", Value: "./debug.log", Condition: func() bool { return runtime.GOOS == "linux" }, RequiresValue: true},
			{Flag: "-settings", RequiresValue: false},
			{Flag: "StartLocalHost", Value: strconv.FormatBool(config.GetStartLocalHost()), RequiresValue: true},
			{Flag: "ServerVisible", Value: strconv.FormatBool(config.ServerVisible.Get()), RequiresValue: true},
			{Flag: "GamePort", Value: config.GamePort.Get(), RequiresValue: true},
			{Flag: "UPNPEnabled", Value: strconv.FormatBool(config.GetUPNPEnabled()), RequiresValue: true},
			{Flag: "ServerName", Value: config.ServerName.Get(), RequiresValue: true},
			{Flag: "ServerPassword", Value: config.ServerPassword.Get(), Condition: func() bool { return config.ServerPassword.Get() != "" }, RequiresValue: true},
			{Flag: "ServerMaxPlayers", Value: config.ServerMaxPlayers.Get(), RequiresValue: true},
			{Flag: "AutoSave", Value: strconv.FormatBool(config.GetAutoSave()), RequiresValue: true},
			{Flag: "SaveInterval", Value: config.GetSaveInterval(), RequiresValue: true},
			{Flag: "ServerAuthSecret", Value: config.ServerAuthSecret.Get(), Condition: func() bool { return config.ServerAuthSecret.Get() != "" }, RequiresValue: true},
			{Flag: "UpdatePort", Value: config.UpdatePort.Get(), RequiresValue: true},
			{Flag: "AutoPauseServer", Value: strconv.FormatBool(config.GetAutoPauseServer()), RequiresValue: true},
			{Flag: "UseSteamP2P", Value: strconv.FormatBool(config.GetUseSteamP2P()), RequiresValue: true},
			{Flag: "AdminPassword", Value: config.AdminPassword.Get(), Condition: func() bool { return config.AdminPassword.Get() != "" }, RequiresValue: true},
		}
	}
	if !config.GetIsNewTerrainAndSaveSystem() {
		argOrder = []Arg{
			{Flag: "-nographics", RequiresValue: false},
			{Flag: "-batchmode", RequiresValue: false},
			{Flag: "-LOAD", Value: config.GetLegacySaveInfo(), RequiresValue: true, NoQuote: true}, // LOAD has special handling because the gameserver expects 2 parameters
			{Flag: "-logFile", Value: "./debug.log", Condition: func() bool { return runtime.GOOS == "linux" }, RequiresValue: true},
			{Flag: "-settings", RequiresValue: false},
			{Flag: "StartLocalHost", Value: strconv.FormatBool(config.GetStartLocalHost()), RequiresValue: true},
			{Flag: "ServerVisible", Value: strconv.FormatBool(config.ServerVisible.Get()), RequiresValue: true},
			{Flag: "GamePort", Value: config.GamePort.Get(), RequiresValue: true},
			{Flag: "UPNPEnabled", Value: strconv.FormatBool(config.GetUPNPEnabled()), RequiresValue: true},
			{Flag: "ServerName", Value: config.ServerName.Get(), RequiresValue: true},
			{Flag: "ServerPassword", Value: config.ServerPassword.Get(), Condition: func() bool { return config.ServerPassword.Get() != "" }, RequiresValue: true},
			{Flag: "ServerMaxPlayers", Value: config.ServerMaxPlayers.Get(), RequiresValue: true},
			{Flag: "AutoSave", Value: strconv.FormatBool(config.GetAutoSave()), RequiresValue: true},
			{Flag: "SaveInterval", Value: config.GetSaveInterval(), RequiresValue: true},
			{Flag: "ServerAuthSecret", Value: config.ServerAuthSecret.Get(), Condition: func() bool { return config.ServerAuthSecret.Get() != "" }, RequiresValue: true},
			{Flag: "UpdatePort", Value: config.UpdatePort.Get(), RequiresValue: true},
			{Flag: "AutoPauseServer", Value: strconv.FormatBool(config.GetAutoPauseServer()), RequiresValue: true},
			{Flag: "UseSteamP2P", Value: strconv.FormatBool(config.GetUseSteamP2P()), RequiresValue: true},
			{Flag: "AdminPassword", Value: config.AdminPassword.Get(), Condition: func() bool { return config.AdminPassword.Get() != "" }, RequiresValue: true},
		}
	}

	var args []string
	for _, arg := range argOrder {
		if arg.Condition != nil && !arg.Condition() {
			continue
		}
		if arg.RequiresValue && arg.Value == "" {
			continue
		}

		args = append(args, arg.Flag)

		// handling of Legacy SaveInfo: Split on semicolon and add each part as a separate arg. This is a hack to continue to support the old saveinfo format for preterrain servers.
		if arg.Flag == "-LOAD" && arg.Value != "" {
			parts := strings.SplitN(arg.Value, ";", 2)
			for _, part := range parts {
				if part != "" {
					args = append(args, part)
				}
			}
			continue
		}

		if arg.Value != "" {
			args = append(args, arg.Value)
		}
	}

	if config.GetAdditionalParams() != "" {
		args = append(args, strings.Fields(config.GetAdditionalParams())...)
	}

	if config.LocalIpAddress.Get() != "" {
		args = append(args, "LocalIpAddress")
		args = append(args, config.LocalIpAddress.Get())
	}

	return args
}
