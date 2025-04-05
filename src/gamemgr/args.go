package gamemgr

import (
	"StationeersServerUI/src/config"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

type Arg struct {
	Flag          string
	Value         string
	RequiresValue bool
	Condition     func() bool
	NoQuote       bool
}

func buildCommandArgs() []string {
	var argOrder = []Arg{
		{Flag: "-nographics", RequiresValue: false},
		{Flag: "-batchmode", RequiresValue: false},
		{Flag: "-LOAD", Value: config.SaveInfo, RequiresValue: true, NoQuote: true}, // LOAD has special handling because the gameserver expects 2 parameters
		{Flag: "-logFile", Value: "./debug.log", Condition: func() bool { return runtime.GOOS == "linux" }, RequiresValue: true},
		{Flag: "-settings", RequiresValue: false},
		{Flag: "StartLocalHost", Value: strconv.FormatBool(config.StartLocalHost), RequiresValue: true},
		{Flag: "ServerVisible", Value: strconv.FormatBool(config.ServerVisible), RequiresValue: true},
		{Flag: "GamePort", Value: config.GamePort, RequiresValue: true},
		{Flag: "UPNPEnabled", Value: strconv.FormatBool(config.UPNPEnabled), RequiresValue: true},
		{Flag: "ServerName", Value: config.ServerName, RequiresValue: true},
		{Flag: "ServerPassword", Value: config.ServerPassword, Condition: func() bool { return config.ServerPassword != "" }, RequiresValue: true},
		{Flag: "ServerMaxPlayers", Value: config.ServerMaxPlayers, RequiresValue: true},
		{Flag: "AutoSave", Value: strconv.FormatBool(config.AutoSave), RequiresValue: true},
		{Flag: "SaveInterval", Value: config.SaveInterval, RequiresValue: true},
		{Flag: "ServerAuthSecret", Value: config.ServerAuthSecret, Condition: func() bool { return config.ServerAuthSecret != "" }, RequiresValue: true},
		{Flag: "UpdatePort", Value: config.UpdatePort, RequiresValue: true},
		{Flag: "AutoPauseServer", Value: strconv.FormatBool(config.AutoPauseServer), RequiresValue: true},
		{Flag: "UseSteamP2P", Value: strconv.FormatBool(config.UseSteamP2P), RequiresValue: true},
		{Flag: "AdminPassword", Value: config.AdminPassword, Condition: func() bool { return config.AdminPassword != "" }, RequiresValue: true},
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

		if arg.Flag == "-LOAD" && arg.Value != "" {
			parts := strings.SplitN(arg.Value, " ", 2)
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

	if config.AdditionalParams != "" {
		args = append(args, strings.Fields(config.AdditionalParams)...)
	}

	if config.LocalIpAddress != "" {
		args = append(args, "LocalIpAddress")
		args = append(args, config.LocalIpAddress)
	}

	if config.IsDebugMode {
		fmt.Println("=== DEBUG: Raw arguments passed to exec.Command ===")
		for i, arg := range args {
			fmt.Printf("Arg[%d]: %q\n", i, arg)
		}
	}
	return args
}
