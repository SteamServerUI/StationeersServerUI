package web

import (
	"fmt"
	"io/fs"
	"net/http"
	"strconv"
	"text/template"

	"github.com/SteamServerUI/StationeersServerUI/v6/src/config"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/core/security"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/localization"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/logger"
)

func ServeConfigPage(w http.ResponseWriter, r *http.Request) {

	htmlFS, err := fs.Sub(config.V1UIFS, "SSUI/onboard_bundled/ui")
	if err != nil {
		http.Error(w, "Error accessing Virt FS: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFS(htmlFS, "config.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Core.Error("failed to serve config.html")
		return
	}

	// Determine selected attributes for boolean fields
	upnpTrueSelected := ""
	upnpFalseSelected := ""
	if config.GetUPNPEnabled() {
		upnpTrueSelected = "selected"
	} else {
		upnpFalseSelected = "selected"
	}

	discordTrueSelected := ""
	discordFalseSelected := ""
	if config.GetIsDiscordEnabled() {
		discordTrueSelected = "selected"
	} else {
		discordFalseSelected = "selected"
	}

	autoSaveTrueSelected := ""
	autoSaveFalseSelected := ""
	if config.GetAutoSave() {
		autoSaveTrueSelected = "selected"
	} else {
		autoSaveFalseSelected = "selected"
	}

	autoPauseTrueSelected := ""
	autoPauseFalseSelected := ""
	if config.GetAutoPauseServer() {
		autoPauseTrueSelected = "selected"
	} else {
		autoPauseFalseSelected = "selected"
	}

	startLocalTrueSelected := ""
	startLocalFalseSelected := ""
	if config.GetStartLocalHost() {
		startLocalTrueSelected = "selected"
	} else {
		startLocalFalseSelected = "selected"
	}

	serverVisibleTrueSelected := ""
	serverVisibleFalseSelected := ""
	if config.GetServerVisible() {
		serverVisibleTrueSelected = "selected"
	} else {
		serverVisibleFalseSelected = "selected"
	}

	isNewTerrainAndSaveSystemTrueSelected := ""

	if config.GetIsNewTerrainAndSaveSystem() {
		isNewTerrainAndSaveSystemTrueSelected = "selected"
	}

	autoStartServerTrueSelected := ""
	autoStartServerFalseSelected := ""
	if config.GetAutoStartServerOnStartup() {
		autoStartServerTrueSelected = "selected"
	} else {
		autoStartServerFalseSelected = "selected"
	}

	connectivityCheckEnabledTrueSelected := ""
	connectivityCheckEnabledFalseSelected := ""
	if config.GetConnectivityCheckEnabled() {
		connectivityCheckEnabledTrueSelected = "selected"
	} else {
		connectivityCheckEnabledFalseSelected = "selected"
	}

	steamP2PTrueSelected := ""
	steamP2PFalseSelected := ""
	if config.GetUseSteamP2P() {
		steamP2PTrueSelected = "selected"
	} else {
		steamP2PFalseSelected = "selected"
	}

	autoGameServerUpdatesTrueSelected := ""
	autoGameServerUpdatesFalseSelected := ""
	if config.GetAllowAutoGameServerUpdates() {
		autoGameServerUpdatesTrueSelected = "selected"
	} else {
		autoGameServerUpdatesFalseSelected = "selected"
	}

	createSSUILogFileTrueSelected := ""
	createSSUILogFileFalseSelected := ""
	if config.GetCreateSSUILogFile() {
		createSSUILogFileTrueSelected = "selected"
	} else {
		createSSUILogFileFalseSelected = "selected"
	}

	createGameServerLogFileTrueSelected := ""
	createGameServerLogFileFalseSelected := ""
	if config.GetCreateGameServerLogFile() {
		createGameServerLogFileTrueSelected = "selected"
	} else {
		createGameServerLogFileFalseSelected = "selected"
	}

	isStationeersLaunchPadEnabled := "false"
	if config.GetIsStationeersLaunchPadEnabled() {
		isStationeersLaunchPadEnabled = "true"
	}

	rotateServerPasswordTrueSelected := ""
	rotateServerPasswordFalseSelected := ""
	if config.GetRotateServerPassword() {
		rotateServerPasswordTrueSelected = "selected"
	} else {
		rotateServerPasswordFalseSelected = "selected"
	}
	discordRestartVoteEnabledTrueSelected, discordRestartVoteEnabledFalseSelected := "", ""
	if config.GetDiscordRestartVoteEnabled() {
		discordRestartVoteEnabledTrueSelected = "selected"
	} else {
		discordRestartVoteEnabledFalseSelected = "selected"
	}
	discordRestoreVoteEnabledTrueSelected, discordRestoreVoteEnabledFalseSelected := "", ""
	if config.GetDiscordRestoreVoteEnabled() {
		discordRestoreVoteEnabledTrueSelected = "selected"
	} else {
		discordRestoreVoteEnabledFalseSelected = "selected"
	}

	// Expert Settings toggle
	showExpertSettingsTrueSelected := ""
	showExpertSettingsFalseSelected := ""
	if config.GetShowExpertSettings() {
		showExpertSettingsTrueSelected = "selected"
	} else {
		showExpertSettingsFalseSelected = "selected"
	}

	// Expert Settings booleans
	debugTrueSelected := ""
	debugFalseSelected := ""
	if config.GetIsDebugMode() {
		debugTrueSelected = "selected"
	} else {
		debugFalseSelected = "selected"
	}

	logClutterToConsoleTrueSelected := ""
	logClutterToConsoleFalseSelected := ""
	if config.GetLogClutterToConsole() {
		logClutterToConsoleTrueSelected = "selected"
	} else {
		logClutterToConsoleFalseSelected = "selected"
	}

	isSSCMEnabledTrueSelected := ""
	isSSCMEnabledFalseSelected := ""
	if config.GetIsSSCMEnabled() {
		isSSCMEnabledTrueSelected = "selected"
	} else {
		isSSCMEnabledFalseSelected = "selected"
	}

	isConsoleEnabledTrueSelected := ""
	isConsoleEnabledFalseSelected := ""
	if config.GetIsConsoleEnabled() {
		isConsoleEnabledTrueSelected = "selected"
	} else {
		isConsoleEnabledFalseSelected = "selected"
	}

	isUpdateEnabledTrueSelected := ""
	isUpdateEnabledFalseSelected := ""
	if config.GetIsUpdateEnabled() {
		isUpdateEnabledTrueSelected = "selected"
	} else {
		isUpdateEnabledFalseSelected = "selected"
	}

	allowPrereleaseUpdatesTrueSelected := ""
	allowPrereleaseUpdatesFalseSelected := ""
	if config.GetAllowPrereleaseUpdates() {
		allowPrereleaseUpdatesTrueSelected = "selected"
	} else {
		allowPrereleaseUpdatesFalseSelected = "selected"
	}

	allowMajorUpdatesTrueSelected := ""
	allowMajorUpdatesFalseSelected := ""
	if config.GetAllowMajorUpdates() {
		allowMajorUpdatesTrueSelected = "selected"
	} else {
		allowMajorUpdatesFalseSelected = "selected"
	}

	isStationeersLaunchPadAutoUpdatesEnabledTrueSelected := ""
	isStationeersLaunchPadAutoUpdatesEnabledFalseSelected := ""
	if config.GetIsStationeersLaunchPadAutoUpdatesEnabled() {
		isStationeersLaunchPadAutoUpdatesEnabledTrueSelected = "selected"
	} else {
		isStationeersLaunchPadAutoUpdatesEnabledFalseSelected = "selected"
	}

	backupRetentionEnabledTrueSelected := ""
	backupRetentionEnabledFalseSelected := ""
	if config.GetBackupRetentionEnabled() {
		backupRetentionEnabledTrueSelected = "selected"
	} else {
		backupRetentionEnabledFalseSelected = "selected"
	}

	data := ConfigTemplateData{
		Permissions:         pagePermissions(r),
		CanViewSettings:     pageCan(r, security.PermissionSettingsView) || pageCan(r, security.PermissionSettingsManage),
		CanManageSettings:   pageCan(r, security.PermissionSettingsManage),
		CanManageSLP:        pageCan(r, security.PermissionSLPManage),
		CanManageDetections: pageCan(r, security.PermissionDetectionsManage),

		// Config values
		DiscordToken:                            config.GetDiscordToken(),
		DiscordAdminRoleID:                      config.GetDiscordAdminRoleID(),
		EventLogChannelID:                       config.GetEventLogChannelID(),
		StatusPanelChannelID:                    config.GetStatusPanelChannelID(),
		LogChannelID:                            config.GetLogChannelID(),
		BlackListFilePath:                       config.GetBlackListFilePath(),
		IsDiscordEnabled:                        fmt.Sprintf("%v", config.GetIsDiscordEnabled()),
		IsDiscordEnabledTrueSelected:            discordTrueSelected,
		IsDiscordEnabledFalseSelected:           discordFalseSelected,
		RotateServerPassword:                    fmt.Sprintf("%v", config.GetRotateServerPassword()),
		RotateServerPasswordTrueSelected:        rotateServerPasswordTrueSelected,
		RotateServerPasswordFalseSelected:       rotateServerPasswordFalseSelected,
		DiscordRestartVoteEnabledTrueSelected:   discordRestartVoteEnabledTrueSelected,
		DiscordRestartVoteEnabledFalseSelected:  discordRestartVoteEnabledFalseSelected,
		DiscordRestoreVoteEnabledTrueSelected:   discordRestoreVoteEnabledTrueSelected,
		DiscordRestoreVoteEnabledFalseSelected:  discordRestoreVoteEnabledFalseSelected,
		DiscordVoteDurationMinutes:              strconv.Itoa(config.GetDiscordVoteDurationMinutes()),
		DiscordRestartVoteThreshold:             strconv.Itoa(config.GetDiscordRestartVoteThreshold()),
		DiscordRestartVoteMinimum:               strconv.Itoa(config.GetDiscordRestartVoteMinimum()),
		DiscordRestartVoteCooldownMinutes:       strconv.Itoa(config.GetDiscordRestartVoteCooldownMinutes()),
		DiscordRestoreVoteThreshold:             strconv.Itoa(config.GetDiscordRestoreVoteThreshold()),
		DiscordRestoreVoteMinimum:               strconv.Itoa(config.GetDiscordRestoreVoteMinimum()),
		DiscordRestoreVoteCooldownMinutes:       strconv.Itoa(config.GetDiscordRestoreVoteCooldownMinutes()),
		GameBranch:                              config.GetGameBranch(),
		Difficulty:                              config.GetDifficulty(),
		StartCondition:                          config.GetStartCondition(),
		StartLocation:                           config.GetStartLocation(),
		ServerName:                              config.GetServerName(),
		SaveName:                                config.GetSaveName(),
		WorldID:                                 config.GetWorldID(),
		ServerMaxPlayers:                        config.GetServerMaxPlayers(),
		ServerPassword:                          config.GetServerPassword(),
		ServerAuthSecret:                        config.GetServerAuthSecret(),
		AdminPassword:                           config.GetAdminPassword(),
		GamePort:                                config.GetGamePort(),
		UpdatePort:                              config.GetUpdatePort(),
		UPNPEnabled:                             fmt.Sprintf("%v", config.GetUPNPEnabled()),
		UPNPEnabledTrueSelected:                 upnpTrueSelected,
		UPNPEnabledFalseSelected:                upnpFalseSelected,
		AutoSave:                                fmt.Sprintf("%v", config.GetAutoSave()),
		AutoSaveTrueSelected:                    autoSaveTrueSelected,
		AutoSaveFalseSelected:                   autoSaveFalseSelected,
		SaveInterval:                            config.GetSaveInterval(),
		AutoPauseServer:                         fmt.Sprintf("%v", config.GetAutoPauseServer()),
		AutoPauseServerTrueSelected:             autoPauseTrueSelected,
		AutoPauseServerFalseSelected:            autoPauseFalseSelected,
		LocalIpAddress:                          config.GetLocalIpAddress(),
		StartLocalHost:                          fmt.Sprintf("%v", config.GetStartLocalHost()),
		StartLocalHostTrueSelected:              startLocalTrueSelected,
		StartLocalHostFalseSelected:             startLocalFalseSelected,
		ServerVisible:                           fmt.Sprintf("%v", config.GetServerVisible()),
		ServerVisibleTrueSelected:               serverVisibleTrueSelected,
		ServerVisibleFalseSelected:              serverVisibleFalseSelected,
		UseSteamP2P:                             fmt.Sprintf("%v", config.GetUseSteamP2P()),
		UseSteamP2PTrueSelected:                 steamP2PTrueSelected,
		UseSteamP2PFalseSelected:                steamP2PFalseSelected,
		ExePath:                                 config.GetExePath(),
		AdditionalParams:                        config.GetAdditionalParams(),
		AutoRestartServerTimer:                  config.GetAutoRestartServerTimer(),
		AutoRestartCountdown:                    config.GetAutoRestartCountdown(),
		IsNewTerrainAndSaveSystem:               fmt.Sprintf("%v", config.GetIsNewTerrainAndSaveSystem()),
		IsNewTerrainAndSaveSystemTrueSelected:   isNewTerrainAndSaveSystemTrueSelected,
		AutoStartServerOnStartup:                fmt.Sprintf("%v", config.GetAutoStartServerOnStartup()),
		AutoStartServerOnStartupTrueSelected:    autoStartServerTrueSelected,
		AutoStartServerOnStartupFalseSelected:   autoStartServerFalseSelected,
		ConnectivityCheckEnabled:                fmt.Sprintf("%v", config.GetConnectivityCheckEnabled()),
		ConnectivityCheckEnabledTrueSelected:    connectivityCheckEnabledTrueSelected,
		ConnectivityCheckEnabledFalseSelected:   connectivityCheckEnabledFalseSelected,
		AllowAutoGameServerUpdates:              fmt.Sprintf("%v", config.GetAllowAutoGameServerUpdates()),
		AllowAutoGameServerUpdatesTrueSelected:  autoGameServerUpdatesTrueSelected,
		AllowAutoGameServerUpdatesFalseSelected: autoGameServerUpdatesFalseSelected,
		CreateSSUILogFile:                       fmt.Sprintf("%v", config.GetCreateSSUILogFile()),
		CreateSSUILogFileTrueSelected:           createSSUILogFileTrueSelected,
		CreateSSUILogFileFalseSelected:          createSSUILogFileFalseSelected,
		CreateGameServerLogFile:                 fmt.Sprintf("%v", config.GetCreateGameServerLogFile()),
		CreateGameServerLogFileTrueSelected:     createGameServerLogFileTrueSelected,
		CreateGameServerLogFileFalseSelected:    createGameServerLogFileFalseSelected,
		BackupKeepNewestCount:                   strconv.Itoa(config.GetBackupKeepNewestCount()),
		BackupRetentionEnabled:                  fmt.Sprintf("%v", config.GetBackupRetentionEnabled()),
		BackupRetentionEnabledTrueSelected:      backupRetentionEnabledTrueSelected,
		BackupRetentionEnabledFalseSelected:     backupRetentionEnabledFalseSelected,
		BackupDailyRetentionDays:                strconv.Itoa(config.GetBackupDailyRetentionDays()),
		BackupWeeklyRetentionWeeks:              strconv.Itoa(config.GetBackupWeeklyRetentionWeeks()),
		BackupMonthlyRetentionMonths:            strconv.Itoa(config.GetBackupMonthlyRetentionMonths()),
		BackupCleanupIntervalHours:              strconv.Itoa(int(config.GetBackupCleanupInterval().Hours())),

		// Localized UI text
		UIText_ConfigHeadline:            localization.GetString("UIText_ConfigHeadline"),
		UIText_ServerConfig:              localization.GetString("UIText_ServerConfig"),
		UIText_BackToDashboard:           localization.GetString("UIText_BackToDashboard"),
		UIText_DiscordIntegration:        localization.GetString("UIText_DiscordIntegration"),
		UIText_SLPModIntegration:         localization.GetString("UIText_SLPModIntegration"),
		UIText_DetectionManager:          localization.GetString("UIText_DetectionManager"),
		UIText_ConfigurationWizard:       localization.GetString("UIText_ConfigurationWizard"),
		UIText_PleaseSelectSection:       localization.GetString("UIText_PleaseSelectSection"),
		UIText_UseWizardAlternative:      localization.GetString("UIText_UseWizardAlternative"),
		UIText_BasicSettings:             localization.GetString("UIText_BasicSettings"),
		UIText_NetworkSettings:           localization.GetString("UIText_NetworkSettings"),
		UIText_AdvancedSettings:          localization.GetString("UIText_AdvancedSettings"),
		UIText_TerrainSettings:           localization.GetString("UIText_TerrainSettings"),
		UIText_ConfigSaving:              localization.GetString("UIText_ConfigSaving"),
		UIText_ConfigSaveSuccess:         localization.GetString("UIText_ConfigSaveSuccess"),
		UIText_ConfigSaveFailed:          localization.GetString("UIText_ConfigSaveFailed"),
		UIText_SaveConfiguration:         localization.GetString("UIText_SaveConfiguration"),
		UIText_BackupRetention:           localization.GetString("UIText_BackupRetention"),
		UIText_BackupRetentionTitle:      localization.GetString("UIText_BackupRetentionTitle"),
		UIText_BackupRetentionIntro:      localization.GetString("UIText_BackupRetentionIntro"),
		UIText_BackupCleanupEnabled:      localization.GetString("UIText_BackupCleanupEnabled"),
		UIText_BackupCleanupInfo:         localization.GetString("UIText_BackupCleanupInfo"),
		UIText_BackupKeepNewest:          localization.GetString("UIText_BackupKeepNewest"),
		UIText_BackupKeepNewestInfo:      localization.GetString("UIText_BackupKeepNewestInfo"),
		UIText_BackupKeepDaily:           localization.GetString("UIText_BackupKeepDaily"),
		UIText_BackupKeepDailyInfo:       localization.GetString("UIText_BackupKeepDailyInfo"),
		UIText_BackupKeepWeekly:          localization.GetString("UIText_BackupKeepWeekly"),
		UIText_BackupKeepWeeklyInfo:      localization.GetString("UIText_BackupKeepWeeklyInfo"),
		UIText_BackupKeepMonthly:         localization.GetString("UIText_BackupKeepMonthly"),
		UIText_BackupKeepMonthlyInfo:     localization.GetString("UIText_BackupKeepMonthlyInfo"),
		UIText_BackupCleanupInterval:     localization.GetString("UIText_BackupCleanupInterval"),
		UIText_BackupCleanupIntervalInfo: localization.GetString("UIText_BackupCleanupIntervalInfo"),
		UIText_BackupHours:               localization.GetString("UIText_BackupHours"),
		UIText_BackupDays:                localization.GetString("UIText_BackupDays"),
		UIText_BackupWeeks:               localization.GetString("UIText_BackupWeeks"),
		UIText_BackupMonths:              localization.GetString("UIText_BackupMonths"),
		UIText_BackupFiles:               localization.GetString("UIText_BackupFiles"),
		UIText_BackupWarningTitle:        localization.GetString("UIText_BackupWarningTitle"),
		UIText_BackupWarning:             localization.GetString("UIText_BackupWarning"),
		UIText_BackupCopyDelayInfo:       localization.GetString("UIText_BackupCopyDelayInfo"),
		UIText_BasicServerSettings:       localization.GetString("UIText_BasicServerSettings"),

		UIText_ServerName:                     localization.GetString("UIText_ServerName"),
		UIText_ServerNameInfo:                 localization.GetString("UIText_ServerNameInfo"),
		UIText_SaveName:                       localization.GetString("UIText_SaveName"),
		UIText_SaveNameInfo:                   localization.GetString("UIText_SaveNameInfo"),
		UIText_WorldID:                        localization.GetString("UIText_WorldID"),
		UIText_WorldIDInfo:                    localization.GetString("UIText_WorldIDInfo"),
		UIText_MaxPlayers:                     localization.GetString("UIText_MaxPlayers"),
		UIText_MaxPlayersInfo:                 localization.GetString("UIText_MaxPlayersInfo"),
		UIText_ServerPassword:                 localization.GetString("UIText_ServerPassword"),
		UIText_ServerPasswordInfo:             localization.GetString("UIText_ServerPasswordInfo"),
		UIText_AdminPassword:                  localization.GetString("UIText_AdminPassword"),
		UIText_AdminPasswordInfo:              localization.GetString("UIText_AdminPasswordInfo"),
		UIText_AutoSave:                       localization.GetString("UIText_AutoSave"),
		UIText_AutoSaveInfo:                   localization.GetString("UIText_AutoSaveInfo"),
		UIText_SaveInterval:                   localization.GetString("UIText_SaveInterval"),
		UIText_SaveIntervalInfo:               localization.GetString("UIText_SaveIntervalInfo"),
		UIText_AutoPauseServer:                localization.GetString("UIText_AutoPauseServer"),
		UIText_AutoPauseServerInfo:            localization.GetString("UIText_AutoPauseServerInfo"),
		UIText_NetworkConfiguration:           localization.GetString("UIText_NetworkConfiguration"),
		UIText_GamePort:                       localization.GetString("UIText_GamePort"),
		UIText_GamePortInfo:                   localization.GetString("UIText_GamePortInfo"),
		UIText_UpdatePort:                     localization.GetString("UIText_UpdatePort"),
		UIText_UpdatePortInfo:                 localization.GetString("UIText_UpdatePortInfo"),
		UIText_UPNPEnabled:                    localization.GetString("UIText_UPNPEnabled"),
		UIText_UPNPEnabledInfo:                localization.GetString("UIText_UPNPEnabledInfo"),
		UIText_LocalIpAddress:                 localization.GetString("UIText_LocalIpAddress"),
		UIText_LocalIpAddressInfo:             localization.GetString("UIText_LocalIpAddressInfo"),
		UIText_StartLocalHost:                 localization.GetString("UIText_StartLocalHost"),
		UIText_StartLocalHostInfo:             localization.GetString("UIText_StartLocalHostInfo"),
		UIText_ServerVisible:                  localization.GetString("UIText_ServerVisible"),
		UIText_ServerVisibleInfo:              localization.GetString("UIText_ServerVisibleInfo"),
		UIText_UseSteamP2P:                    localization.GetString("UIText_UseSteamP2P"),
		UIText_UseSteamP2PInfo:                localization.GetString("UIText_UseSteamP2PInfo"),
		UIText_AdvertiserTitle:                localization.GetString("UIText_AdvertiserTitle"),
		UIText_AdvertiserSummary:              localization.GetString("UIText_AdvertiserSummary"),
		UIText_AdvertiserConfigure:            localization.GetString("UIText_AdvertiserConfigure"),
		UIText_AdvertiserExplanation:          localization.GetString("UIText_AdvertiserExplanation"),
		UIText_AdvertiserMasterServer:         localization.GetString("UIText_AdvertiserMasterServer"),
		UIText_AdvertiserDisabled:             localization.GetString("UIText_AdvertiserDisabled"),
		UIText_AdvertiserAuto:                 localization.GetString("UIText_AdvertiserAuto"),
		UIText_AdvertiserIPv4:                 localization.GetString("UIText_AdvertiserIPv4"),
		UIText_AdvertiserDNS:                  localization.GetString("UIText_AdvertiserDNS"),
		UIText_AdvertiserAddress:              localization.GetString("UIText_AdvertiserAddress"),
		UIText_AdvertiserRestartWarning:       localization.GetString("UIText_AdvertiserRestartWarning"),
		UIText_AdvertiserKicker:               localization.GetString("UIText_AdvertiserKicker"),
		UIText_AdvertiserRestartRequired:      localization.GetString("UIText_AdvertiserRestartRequired"),
		UIText_DebugMode:                      localization.GetString("UIText_DebugMode"),
		UIText_DebugModeInfo:                  localization.GetString("UIText_DebugModeInfo"),
		UIText_LogLevel:                       localization.GetString("UIText_LogLevel"),
		UIText_LogLevelInfo:                   localization.GetString("UIText_LogLevelInfo"),
		UIText_LogClutterToConsole:            localization.GetString("UIText_LogClutterToConsole"),
		UIText_LogClutterToConsoleInfo:        localization.GetString("UIText_LogClutterToConsoleInfo"),
		UIText_SSUIWebPort:                    localization.GetString("UIText_SSUIWebPort"),
		UIText_SSUIWebPortInfo:                localization.GetString("UIText_SSUIWebPortInfo"),
		UIText_ConsoleEnabled:                 localization.GetString("UIText_ConsoleEnabled"),
		UIText_ConsoleEnabledInfo:             localization.GetString("UIText_ConsoleEnabledInfo"),
		UIText_SSCMEnabled:                    localization.GetString("UIText_SSCMEnabled"),
		UIText_SSCMEnabledInfo:                localization.GetString("UIText_SSCMEnabledInfo"),
		UIText_UpdateEnabled:                  localization.GetString("UIText_UpdateEnabled"),
		UIText_UpdateEnabledInfo:              localization.GetString("UIText_UpdateEnabledInfo"),
		UIText_AllowPrereleaseUpdates:         localization.GetString("UIText_AllowPrereleaseUpdates"),
		UIText_AllowPrereleaseUpdatesInfo:     localization.GetString("UIText_AllowPrereleaseUpdatesInfo"),
		UIText_AllowMajorUpdates:              localization.GetString("UIText_AllowMajorUpdates"),
		UIText_AllowMajorUpdatesInfo:          localization.GetString("UIText_AllowMajorUpdatesInfo"),
		UIText_SLPAutoUpdates:                 localization.GetString("UIText_SLPAutoUpdates"),
		UIText_SLPAutoUpdatesInfo:             localization.GetString("UIText_SLPAutoUpdatesInfo"),
		UIText_AdvertiserSaveRestart:          localization.GetString("UIText_AdvertiserSaveRestart"),
		UIText_AdvertiserRestarting:           localization.GetString("UIText_AdvertiserRestarting"),
		UIText_Cancel:                         localization.GetString("UIText_Cancel"),
		UIText_Close:                          localization.GetString("UIText_Close"),
		UIText_AdvancedConfiguration:          localization.GetString("UIText_AdvancedConfiguration"),
		UIText_TLSCertificateTitle:            localization.GetString("UIText_TLSCertificateTitle"),
		UIText_TLSCertificateSummary:          localization.GetString("UIText_TLSCertificateSummary"),
		UIText_TLSConfigure:                   localization.GetString("UIText_TLSConfigure"),
		UIText_TLSCertificateExplanation:      localization.GetString("UIText_TLSCertificateExplanation"),
		UIText_TLSCertificateFormats:          localization.GetString("UIText_TLSCertificateFormats"),
		UIText_TLSCertificateFile:             localization.GetString("UIText_TLSCertificateFile"),
		UIText_TLSPrivateKeyFile:              localization.GetString("UIText_TLSPrivateKeyFile"),
		UIText_TLSRestartRequired:             localization.GetString("UIText_TLSRestartRequired"),
		UIText_TLSRestartWarning:              localization.GetString("UIText_TLSRestartWarning"),
		UIText_TLSSaveRestart:                 localization.GetString("UIText_TLSSaveRestart"),
		UIText_TLSRestarting:                  localization.GetString("UIText_TLSRestarting"),
		UIText_ServerAuthSecret:               localization.GetString("UIText_ServerAuthSecret"),
		UIText_ServerAuthSecretInfo:           localization.GetString("UIText_ServerAuthSecretInfo"),
		UIText_ServerExePath:                  localization.GetString("UIText_ServerExePath"),
		UIText_ServerExePathInfo:              localization.GetString("UIText_ServerExePathInfo"),
		UIText_ServerExePathInfo2:             localization.GetString("UIText_ServerExePathInfo2"),
		UIText_AdditionalParams:               localization.GetString("UIText_AdditionalParams"),
		UIText_AdditionalParamsInfo:           localization.GetString("UIText_AdditionalParamsInfo"),
		UIText_ShowExpertSettings:             localization.GetString("UIText_ShowExpertSettings"),
		UIText_ShowExpertSettingsInfo:         localization.GetString("UIText_ShowExpertSettingsInfo"),
		UIText_ExpertSettings:                 localization.GetString("UIText_ExpertSettings"),
		UIText_ExpertSettingsWarning:          localization.GetString("UIText_ExpertSettingsWarning"),
		UIText_AutoRestartServerTimer:         localization.GetString("UIText_AutoRestartServerTimer"),
		UIText_AutoRestartServerTimerInfo:     localization.GetString("UIText_AutoRestartServerTimerInfo"),
		UIText_AutoRestartCountdown:           localization.GetString("UIText_AutoRestartCountdown"),
		UIText_AutoRestartCountdownInfo:       localization.GetString("UIText_AutoRestartCountdownInfo"),
		UIText_GameBranch:                     localization.GetString("UIText_GameBranch"),
		UIText_GameBranchInfo:                 localization.GetString("UIText_GameBranchInfo"),
		UIText_TerrainSettingsHeader:          localization.GetString("UIText_TerrainSettingsHeader"),
		UIText_TerrainWarning:                 localization.GetString("UIText_TerrainWarning"),
		UIText_Difficulty:                     localization.GetString("UIText_Difficulty"),
		UIText_DifficultyInfo:                 localization.GetString("UIText_DifficultyInfo"),
		UIText_StartCondition:                 localization.GetString("UIText_StartCondition"),
		UIText_StartConditionInfo:             localization.GetString("UIText_StartConditionInfo"),
		UIText_StartLocation:                  localization.GetString("UIText_StartLocation"),
		UIText_StartLocationInfo:              localization.GetString("UIText_StartLocationInfo"),
		UIText_TerrainSettingsFillHint:        localization.GetString("UIText_TerrainSettingsFillHint"),
		UIText_AutoStartServerOnStartup:       localization.GetString("UIText_AutoStartServerOnStartup"),
		UIText_AutoStartServerOnStartupInfo:   localization.GetString("UIText_AutoStartServerOnStartupInfo"),
		UIText_ConnectivityCheck:              localization.GetString("UIText_ConnectivityCheck"),
		UIText_ConnectivityCheckInfo:          localization.GetString("UIText_ConnectivityCheckInfo"),
		UIText_AllowAutoGameServerUpdates:     localization.GetString("UIText_AllowAutoGameServerUpdates"),
		UIText_AllowAutoGameServerUpdatesInfo: localization.GetString("UIText_AllowAutoGameServerUpdatesInfo"),
		UIText_CreateSSUILogFile:              localization.GetString("UIText_CreateSSUILogFile"),
		UIText_CreateSSUILogFileInfo:          localization.GetString("UIText_CreateSSUILogFileInfo"),
		UIText_CreateGameServerLogFile:        localization.GetString("UIText_CreateGameServerLogFile"),
		UIText_CreateGameServerLogFileInfo:    localization.GetString("UIText_CreateGameServerLogFileInfo"),

		UIText_DiscordIntegrationTitle:    localization.GetString("UIText_DiscordIntegrationTitle"),
		UIText_DiscordBotToken:            localization.GetString("UIText_DiscordBotToken"),
		UIText_DiscordBotTokenInfo:        localization.GetString("UIText_DiscordBotTokenInfo"),
		UIText_ChannelConfiguration:       localization.GetString("UIText_ChannelConfiguration"),
		UIText_DiscordAdminRole:           localization.GetString("UIText_DiscordAdminRole"),
		UIText_DiscordAdminRoleInfo:       localization.GetString("UIText_DiscordAdminRoleInfo"),
		UIText_RotateServerPassword:       localization.GetString("UIText_RotateServerPassword"),
		UIText_RotateServerPasswordInfo:   localization.GetString("UIText_RotateServerPasswordInfo"),
		UIText_EventLogChannel:            localization.GetString("UIText_EventLogChannel"),
		UIText_EventLogChannelInfo:        localization.GetString("UIText_EventLogChannelInfo"),
		UIText_StatusPanelChannel:         localization.GetString("UIText_StatusPanelChannel"),
		UIText_StatusPanelChannelInfo:     localization.GetString("UIText_StatusPanelChannelInfo"),
		UIText_LogChannel:                 localization.GetString("UIText_LogChannel"),
		UIText_LogChannelInfo:             localization.GetString("UIText_LogChannelInfo"),
		UIText_BannedPlayersListPath:      localization.GetString("UIText_BannedPlayersListPath"),
		UIText_BannedPlayersListPathInfo:  localization.GetString("UIText_BannedPlayersListPathInfo"),
		UIText_DiscordIntegrationBenefits: localization.GetString("UIText_DiscordIntegrationBenefits"),
		UIText_DiscordBenefit1:            localization.GetString("UIText_DiscordBenefit1"),
		UIText_DiscordBenefit2:            localization.GetString("UIText_DiscordBenefit2"),
		UIText_DiscordBenefit3:            localization.GetString("UIText_DiscordBenefit3"),
		UIText_DiscordBenefit4:            localization.GetString("UIText_DiscordBenefit4"),
		UIText_DiscordBenefit5:            localization.GetString("UIText_DiscordBenefit5"),
		UIText_DiscordSetupInstructions:   localization.GetString("UIText_DiscordSetupInstructions"),
		UIText_DiscordVotingTitle:         localization.GetString("UIText_DiscordVotingTitle"),
		UIText_DiscordVotingInfo:          localization.GetString("UIText_DiscordVotingInfo"),
		UIText_DiscordRestartVoteEnabled:  localization.GetString("UIText_DiscordRestartVoteEnabled"),
		UIText_DiscordRestoreVoteEnabled:  localization.GetString("UIText_DiscordRestoreVoteEnabled"),
		UIText_DiscordVoteDuration:        localization.GetString("UIText_DiscordVoteDuration"),
		UIText_DiscordVoteThreshold:       localization.GetString("UIText_DiscordVoteThreshold"),
		UIText_DiscordVoteMinimum:         localization.GetString("UIText_DiscordVoteMinimum"),
		UIText_DiscordVoteCooldown:        localization.GetString("UIText_DiscordVoteCooldown"),

		UIText_CopyrightConfig1: localization.GetString("UIText_Copyright1"),
		UIText_CopyrightConfig2: localization.GetString("UIText_Copyright2"),

		// SLP Section
		UIText_SLP_Title:                  localization.GetString("UIText_SLP_Title"),
		UIText_SLP_Description:            localization.GetString("UIText_SLP_Description"),
		UIText_SLP_ReadyToInstall:         localization.GetString("UIText_SLP_ReadyToInstall"),
		UIText_SLP_InstallButton:          localization.GetString("UIText_SLP_InstallButton"),
		UIText_SLP_UploadModPackage:       localization.GetString("UIText_SLP_UploadModPackage"),
		UIText_SLP_UploadDescription:      localization.GetString("UIText_SLP_UploadDescription"),
		UIText_SLP_UploadDescriptionLink:  localization.GetString("UIText_SLP_UploadDescriptionLink"),
		UIText_SLP_InstallFirst:           localization.GetString("UIText_SLP_InstallFirst"),
		UIText_SLP_InstallFirstSubtext:    localization.GetString("UIText_SLP_InstallFirstSubtext"),
		UIText_SLP_DragDropHere:           localization.GetString("UIText_SLP_DragDropHere"),
		UIText_SLP_OrClickToSelect:        localization.GetString("UIText_SLP_OrClickToSelect"),
		UIText_SLP_UploadButton:           localization.GetString("UIText_SLP_UploadButton"),
		UIText_SLP_ManageInstallation:     localization.GetString("UIText_SLP_ManageInstallation"),
		UIText_SLP_UninstallWarning:       localization.GetString("UIText_SLP_UninstallWarning"),
		UIText_SLP_UninstallButton:        localization.GetString("UIText_SLP_UninstallButton"),
		UIText_SLP_UpdateWorkshopMods:     localization.GetString("UIText_SLP_UpdateWorkshopMods"),
		UIText_SLP_UpdateWorkshopModsDesc: localization.GetString("UIText_SLP_UpdateWorkshopModsDesc"),
		UIText_SLP_UpdateButton:           localization.GetString("UIText_SLP_UpdateButton"),
		UIText_SLP_InstalledMods:          localization.GetString("UIText_SLP_InstalledMods"),

		UIText_DM_Intro:                  localization.GetString("UIText_DM_Intro"),
		UIText_DM_Routines:               localization.GetString("UIText_DM_Routines"),
		UIText_DM_NoRoutines:             localization.GetString("UIText_DM_NoRoutines"),
		UIText_DM_NewRoutine:             localization.GetString("UIText_DM_NewRoutine"),
		UIText_DM_Name:                   localization.GetString("UIText_DM_Name"),
		UIText_DM_NamePlaceholder:        localization.GetString("UIText_DM_NamePlaceholder"),
		UIText_DM_When:                   localization.GetString("UIText_DM_When"),
		UIText_DM_Then:                   localization.GetString("UIText_DM_Then"),
		UIText_DM_ModeKeyword:            localization.GetString("UIText_DM_ModeKeyword"),
		UIText_DM_ModeRegex:              localization.GetString("UIText_DM_ModeRegex"),
		UIText_DM_ModeEvent:              localization.GetString("UIText_DM_ModeEvent"),
		UIText_DM_KeywordPlaceholder:     localization.GetString("UIText_DM_KeywordPlaceholder"),
		UIText_DM_EventPlayerReady:       localization.GetString("UIText_DM_EventPlayerReady"),
		UIText_DM_EventPlayerConnecting:  localization.GetString("UIText_DM_EventPlayerConnecting"),
		UIText_DM_EventPlayerDisconnect:  localization.GetString("UIText_DM_EventPlayerDisconnect"),
		UIText_DM_EventWorldSaved:        localization.GetString("UIText_DM_EventWorldSaved"),
		UIText_DM_EventSessionRegistered: localization.GetString("UIText_DM_EventSessionRegistered"),
		UIText_DM_EventServerRunning:     localization.GetString("UIText_DM_EventServerRunning"),
		UIText_DM_EventException:         localization.GetString("UIText_DM_EventException"),
		UIText_DM_EventSettingsChanged:   localization.GetString("UIText_DM_EventSettingsChanged"),
		UIText_DM_EventNewGameStarted:    localization.GetString("UIText_DM_EventNewGameStarted"),
		UIText_DM_HintKeyword:            localization.GetString("UIText_DM_HintKeyword"),
		UIText_DM_HintRegex:              localization.GetString("UIText_DM_HintRegex"),
		UIText_DM_HintEvent:              localization.GetString("UIText_DM_HintEvent"),
		UIText_DM_AddAction:              localization.GetString("UIText_DM_AddAction"),
		UIText_DM_Cooldown:               localization.GetString("UIText_DM_Cooldown"),
		UIText_DM_CooldownHint:           localization.GetString("UIText_DM_CooldownHint"),
		UIText_DM_SaveRoutine:            localization.GetString("UIText_DM_SaveRoutine"),
		UIText_DM_NoPermission:           localization.GetString("UIText_DM_NoPermission"),
		UIText_DM_GroupMessages:          localization.GetString("UIText_DM_GroupMessages"),
		UIText_DM_GroupInGame:            localization.GetString("UIText_DM_GroupInGame"),
		UIText_DM_GroupFlow:              localization.GetString("UIText_DM_GroupFlow"),
		UIText_DM_GroupServer:            localization.GetString("UIText_DM_GroupServer"),
		UIText_DM_ActionEventlog:         localization.GetString("UIText_DM_ActionEventlog"),
		UIText_DM_ActionDiscordHub:       localization.GetString("UIText_DM_ActionDiscordHub"),
		UIText_DM_ActionAnnounce:         localization.GetString("UIText_DM_ActionAnnounce"),
		UIText_DM_ActionCommand:          localization.GetString("UIText_DM_ActionCommand"),
		UIText_DM_ActionSave:             localization.GetString("UIText_DM_ActionSave"),
		UIText_DM_ActionWait:             localization.GetString("UIText_DM_ActionWait"),
		UIText_DM_ActionWebhook:          localization.GetString("UIText_DM_ActionWebhook"),
		UIText_DM_ActionRestart:          localization.GetString("UIText_DM_ActionRestart"),
		UIText_DM_ActionStop:             localization.GetString("UIText_DM_ActionStop"),
		UIText_DM_ActionStart:            localization.GetString("UIText_DM_ActionStart"),
		UIText_DM_ActionRestore:          localization.GetString("UIText_DM_ActionRestore"),
		UIText_DM_PlaceholderEventlog:    localization.GetString("UIText_DM_PlaceholderEventlog"),
		UIText_DM_PlaceholderDiscord:     localization.GetString("UIText_DM_PlaceholderDiscord"),
		UIText_DM_PlaceholderAnnounce:    localization.GetString("UIText_DM_PlaceholderAnnounce"),
		UIText_DM_WebhookMessage:         localization.GetString("UIText_DM_WebhookMessage"),
		UIText_DM_Seconds:                localization.GetString("UIText_DM_Seconds"),
		UIText_DM_NoSettings:             localization.GetString("UIText_DM_NoSettings"),
		UIText_DM_MoveUp:                 localization.GetString("UIText_DM_MoveUp"),
		UIText_DM_MoveDown:               localization.GetString("UIText_DM_MoveDown"),
		UIText_DM_Remove:                 localization.GetString("UIText_DM_Remove"),
		UIText_DM_Delete:                 localization.GetString("UIText_DM_Delete"),
		UIText_DM_Untitled:               localization.GetString("UIText_DM_Untitled"),
		UIText_DM_CooldownMeta:           localization.GetString("UIText_DM_CooldownMeta"),
		UIText_DM_WaitSummary:            localization.GetString("UIText_DM_WaitSummary"),
		UIText_DM_BadgeKeyword:           localization.GetString("UIText_DM_BadgeKeyword"),
		UIText_DM_BadgeRegex:             localization.GetString("UIText_DM_BadgeRegex"),
		UIText_DM_BadgeEvent:             localization.GetString("UIText_DM_BadgeEvent"),

		UIText_WS_Kicker:              localization.GetString("UIText_WS_Kicker"),
		UIText_WS_Title:               localization.GetString("UIText_WS_Title"),
		UIText_WS_Description:         localization.GetString("UIText_WS_Description"),
		UIText_WS_Info:                localization.GetString("UIText_WS_Info"),
		UIText_WS_Install:             localization.GetString("UIText_WS_Install"),
		UIText_SUP_Title:              localization.GetString("UIText_SUP_Title"),
		UIText_SUP_DownloadPackage:    localization.GetString("UIText_SUP_DownloadPackage"),
		UIText_SUP_StartRecording:     localization.GetString("UIText_SUP_StartRecording"),
		UIText_SUP_StopRecording:      localization.GetString("UIText_SUP_StopRecording"),
		UIText_SUP_RecordingFull:      localization.GetString("UIText_SUP_RecordingFull"),
		UIText_SUP_DownloadRecording:  localization.GetString("UIText_SUP_DownloadRecording"),
		UIText_SUP_PackageInfo:        localization.GetString("UIText_SUP_PackageInfo"),
		UIText_SUP_Running:            localization.GetString("UIText_SUP_Running"),
		UIText_SUP_AllGood:            localization.GetString("UIText_SUP_AllGood"),
		UIText_SUP_Intro:              localization.GetString("UIText_SUP_Intro"),
		UIText_SUP_PackageTitle:       localization.GetString("UIText_SUP_PackageTitle"),
		UIText_SUP_Refresh:            localization.GetString("UIText_SUP_Refresh"),
		UIText_SUP_LinkTitle:          localization.GetString("UIText_SUP_LinkTitle"),
		UIText_SUP_LinkInfo:           localization.GetString("UIText_SUP_LinkInfo"),
		UIText_SUP_LinkDuration:       localization.GetString("UIText_SUP_LinkDuration"),
		UIText_SUP_Link30:             localization.GetString("UIText_SUP_Link30"),
		UIText_SUP_Link60:             localization.GetString("UIText_SUP_Link60"),
		UIText_SUP_Link120:            localization.GetString("UIText_SUP_Link120"),
		UIText_SUP_LinkOpen:           localization.GetString("UIText_SUP_LinkOpen"),
		UIText_SUP_LinkClose:          localization.GetString("UIText_SUP_LinkClose"),
		UIText_SUP_LinkCode:           localization.GetString("UIText_SUP_LinkCode"),
		UIText_SUP_LinkCopy:           localization.GetString("UIText_SUP_LinkCopy"),
		UIText_SUP_LinkWaiting:        localization.GetString("UIText_SUP_LinkWaiting"),
		UIText_SUP_LinkClaimed:        localization.GetString("UIText_SUP_LinkClaimed"),
		UIText_SUP_LinkOffline:        localization.GetString("UIText_SUP_LinkOffline"),
		UIText_SUP_LinkExpires:        localization.GetString("UIText_SUP_LinkExpires"),
		UIText_SUP_LinkActivity:       localization.GetString("UIText_SUP_LinkActivity"),
		UIText_SUP_LinkNoActivity:     localization.GetString("UIText_SUP_LinkNoActivity"),
		UIText_AC_Tab:                 localization.GetString("UIText_AC_Tab"),
		UIText_AC_Kicker:              localization.GetString("UIText_AC_Kicker"),
		UIText_AC_Intro:               localization.GetString("UIText_AC_Intro"),
		UIText_AC_LoadingAccount:      localization.GetString("UIText_AC_LoadingAccount"),
		UIText_AC_NavAccount:          localization.GetString("UIText_AC_NavAccount"),
		UIText_AC_NavPeople:           localization.GetString("UIText_AC_NavPeople"),
		UIText_AC_NavGroups:           localization.GetString("UIText_AC_NavGroups"),
		UIText_AC_NavAPI:              localization.GetString("UIText_AC_NavAPI"),
		UIText_AC_AccountKicker:       localization.GetString("UIText_AC_AccountKicker"),
		UIText_AC_ChangeMyPassword:    localization.GetString("UIText_AC_ChangeMyPassword"),
		UIText_AC_ChangePasswordInfo:  localization.GetString("UIText_AC_ChangePasswordInfo"),
		UIText_AC_CurrentPassword:     localization.GetString("UIText_AC_CurrentPassword"),
		UIText_AC_NewPassword:         localization.GetString("UIText_AC_NewPassword"),
		UIText_AC_RepeatPassword:      localization.GetString("UIText_AC_RepeatPassword"),
		UIText_AC_ChangePassword:      localization.GetString("UIText_AC_ChangePassword"),
		UIText_AC_DevicesKicker:       localization.GetString("UIText_AC_DevicesKicker"),
		UIText_AC_Sessions:            localization.GetString("UIText_AC_Sessions"),
		UIText_AC_LoadingSessions:     localization.GetString("UIText_AC_LoadingSessions"),
		UIText_AC_NoUsersPermission:   localization.GetString("UIText_AC_NoUsersPermission"),
		UIText_AC_NewAccountKicker:    localization.GetString("UIText_AC_NewAccountKicker"),
		UIText_AC_AddUser:             localization.GetString("UIText_AC_AddUser"),
		UIText_AC_Username:            localization.GetString("UIText_AC_Username"),
		UIText_AC_TempPassword:        localization.GetString("UIText_AC_TempPassword"),
		UIText_AC_AccessGroup:         localization.GetString("UIText_AC_AccessGroup"),
		UIText_AC_NoAccessGroup:       localization.GetString("UIText_AC_NoAccessGroup"),
		UIText_AC_AccountsKicker:      localization.GetString("UIText_AC_AccountsKicker"),
		UIText_AC_Users:               localization.GetString("UIText_AC_Users"),
		UIText_AC_LoadingUsers:        localization.GetString("UIText_AC_LoadingUsers"),
		UIText_AC_NoGroupsPermission:  localization.GetString("UIText_AC_NoGroupsPermission"),
		UIText_AC_RolesKicker:         localization.GetString("UIText_AC_RolesKicker"),
		UIText_AC_Group:               localization.GetString("UIText_AC_Group"),
		UIText_AC_CreateGroup:         localization.GetString("UIText_AC_CreateGroup"),
		UIText_AC_Preset:              localization.GetString("UIText_AC_Preset"),
		UIText_AC_PresetCustom:        localization.GetString("UIText_AC_PresetCustom"),
		UIText_AC_PresetViewer:        localization.GetString("UIText_AC_PresetViewer"),
		UIText_AC_PresetOperator:      localization.GetString("UIText_AC_PresetOperator"),
		UIText_AC_PresetBackups:       localization.GetString("UIText_AC_PresetBackups"),
		UIText_AC_PresetAdministrator: localization.GetString("UIText_AC_PresetAdministrator"),
		UIText_AC_Name:                localization.GetString("UIText_AC_Name"),
		UIText_AC_Description:         localization.GetString("UIText_AC_Description"),
		UIText_AC_FineTuneKicker:      localization.GetString("UIText_AC_FineTuneKicker"),
		UIText_AC_Permissions:         localization.GetString("UIText_AC_Permissions"),
		UIText_AC_ClearAll:            localization.GetString("UIText_AC_ClearAll"),
		UIText_AC_SaveGroup:           localization.GetString("UIText_AC_SaveGroup"),
		UIText_AC_DeleteGroup:         localization.GetString("UIText_AC_DeleteGroup"),
		UIText_AC_NoTokensPermission:  localization.GetString("UIText_AC_NoTokensPermission"),
		UIText_AC_TokensKicker:        localization.GetString("UIText_AC_TokensKicker"),
		UIText_AC_CreateTokenTitle:    localization.GetString("UIText_AC_CreateTokenTitle"),
		UIText_AC_TokenInfo:           localization.GetString("UIText_AC_TokenInfo"),
		UIText_AC_Expires:             localization.GetString("UIText_AC_Expires"),
		UIText_AC_Days30:              localization.GetString("UIText_AC_Days30"),
		UIText_AC_Days90:              localization.GetString("UIText_AC_Days90"),
		UIText_AC_Year1:               localization.GetString("UIText_AC_Year1"),
		UIText_AC_Never:               localization.GetString("UIText_AC_Never"),
		UIText_AC_CreateToken:         localization.GetString("UIText_AC_CreateToken"),
		UIText_AC_CopyTokenNow:        localization.GetString("UIText_AC_CopyTokenNow"),
		UIText_AC_Copy:                localization.GetString("UIText_AC_Copy"),
		UIText_AC_IssuedKicker:        localization.GetString("UIText_AC_IssuedKicker"),
		UIText_AC_ApiTokens:           localization.GetString("UIText_AC_ApiTokens"),
		UIText_AC_LoadingTokens:       localization.GetString("UIText_AC_LoadingTokens"),
		UIText_AC_NoSessions:          localization.GetString("UIText_AC_NoSessions"),
		UIText_AC_ThisSession:         localization.GetString("UIText_AC_ThisSession"),
		UIText_AC_LastUsed:            localization.GetString("UIText_AC_LastUsed"),
		UIText_AC_SignOut:             localization.GetString("UIText_AC_SignOut"),
		UIText_AC_Unknown:             localization.GetString("UIText_AC_Unknown"),
		UIText_AC_SystemGroup:         localization.GetString("UIText_AC_SystemGroup"),
		UIText_AC_Enabled:             localization.GetString("UIText_AC_Enabled"),
		UIText_AC_Disabled:            localization.GetString("UIText_AC_Disabled"),
		UIText_AC_Manage:              localization.GetString("UIText_AC_Manage"),
		UIText_AC_Close:               localization.GetString("UIText_AC_Close"),
		UIText_AC_Status:              localization.GetString("UIText_AC_Status"),
		UIText_AC_AccountEnabled:      localization.GetString("UIText_AC_AccountEnabled"),
		UIText_AC_AccountEnabledInfo:  localization.GetString("UIText_AC_AccountEnabledInfo"),
		UIText_AC_ResetPassword:       localization.GetString("UIText_AC_ResetPassword"),
		UIText_AC_KeepPassword:        localization.GetString("UIText_AC_KeepPassword"),
		UIText_AC_SaveUser:            localization.GetString("UIText_AC_SaveUser"),
		UIText_AC_DeleteUser:          localization.GetString("UIText_AC_DeleteUser"),
		UIText_AC_GroupLocked:         localization.GetString("UIText_AC_GroupLocked"),
		UIText_AC_NoGrantableGroups:   localization.GetString("UIText_AC_NoGrantableGroups"),
		UIText_AC_ExpiresOn:           localization.GetString("UIText_AC_ExpiresOn"),
		UIText_AC_NeverExpires:        localization.GetString("UIText_AC_NeverExpires"),
		UIText_AC_Revoke:              localization.GetString("UIText_AC_Revoke"),
		UIText_AC_NoTokens:            localization.GetString("UIText_AC_NoTokens"),
		UIText_TH_Tab:                 localization.GetString("UIText_TH_Tab"),
		UIText_TH_Presets:             localization.GetString("UIText_TH_Presets"),
		UIText_TH_GroupCoreColors:     localization.GetString("UIText_TH_GroupCoreColors"),
		UIText_TH_GroupTextColors:     localization.GetString("UIText_TH_GroupTextColors"),
		UIText_TH_GroupSurfaces:       localization.GetString("UIText_TH_GroupSurfaces"),
		UIText_TH_GroupControlPanel:   localization.GetString("UIText_TH_GroupControlPanel"),
		UIText_TH_GroupConsoleColors:  localization.GetString("UIText_TH_GroupConsoleColors"),
		UIText_TH_LivePreview:         localization.GetString("UIText_TH_LivePreview"),
		UIText_TH_Actions:             localization.GetString("UIText_TH_Actions"),
		UIText_TH_SaveTheme:           localization.GetString("UIText_TH_SaveTheme"),
		UIText_TH_ResetDefaults:       localization.GetString("UIText_TH_ResetDefaults"),
		UIText_TH_ExportTheme:         localization.GetString("UIText_TH_ExportTheme"),
		UIText_TH_ImportTheme:         localization.GetString("UIText_TH_ImportTheme"),
		UIText_TH_ImportPlaceholder:   localization.GetString("UIText_TH_ImportPlaceholder"),
		UIText_TH_ApplyImport:         localization.GetString("UIText_TH_ApplyImport"),
		UIText_TH_Cancel:              localization.GetString("UIText_TH_Cancel"),
		UIText_TH_CopyClipboard:       localization.GetString("UIText_TH_CopyClipboard"),
		UIText_TH_Close:               localization.GetString("UIText_TH_Close"),
		UIText_TH_PrimaryLabel:        localization.GetString("UIText_TH_PrimaryLabel"),
		UIText_TH_PrimaryDesc:         localization.GetString("UIText_TH_PrimaryDesc"),
		UIText_TH_BgDarkLabel:         localization.GetString("UIText_TH_BgDarkLabel"),
		UIText_TH_BgDarkDesc:          localization.GetString("UIText_TH_BgDarkDesc"),
		UIText_TH_BgPanelLabel:        localization.GetString("UIText_TH_BgPanelLabel"),
		UIText_TH_BgPanelDesc:         localization.GetString("UIText_TH_BgPanelDesc"),
		UIText_TH_AccentLabel:         localization.GetString("UIText_TH_AccentLabel"),
		UIText_TH_AccentDesc:          localization.GetString("UIText_TH_AccentDesc"),
		UIText_TH_DangerLabel:         localization.GetString("UIText_TH_DangerLabel"),
		UIText_TH_DangerDesc:          localization.GetString("UIText_TH_DangerDesc"),
		UIText_TH_SuccessLabel:        localization.GetString("UIText_TH_SuccessLabel"),
		UIText_TH_SuccessDesc:         localization.GetString("UIText_TH_SuccessDesc"),
		UIText_TH_WarningLabel:        localization.GetString("UIText_TH_WarningLabel"),
		UIText_TH_WarningDesc:         localization.GetString("UIText_TH_WarningDesc"),
		UIText_TH_TextHeaderLabel:     localization.GetString("UIText_TH_TextHeaderLabel"),
		UIText_TH_TextHeaderDesc:      localization.GetString("UIText_TH_TextHeaderDesc"),
		UIText_TH_TextBrightLabel:     localization.GetString("UIText_TH_TextBrightLabel"),
		UIText_TH_TextBrightDesc:      localization.GetString("UIText_TH_TextBrightDesc"),
		UIText_TH_TextDimLabel:        localization.GetString("UIText_TH_TextDimLabel"),
		UIText_TH_TextDimDesc:         localization.GetString("UIText_TH_TextDimDesc"),
		UIText_TH_TextMutedLabel:      localization.GetString("UIText_TH_TextMutedLabel"),
		UIText_TH_TextMutedDesc:       localization.GetString("UIText_TH_TextMutedDesc"),
		UIText_TH_SurfaceDarkLabel:    localization.GetString("UIText_TH_SurfaceDarkLabel"),
		UIText_TH_SurfaceDarkDesc:     localization.GetString("UIText_TH_SurfaceDarkDesc"),
		UIText_TH_SurfaceHoverLabel:   localization.GetString("UIText_TH_SurfaceHoverLabel"),
		UIText_TH_SurfaceHoverDesc:    localization.GetString("UIText_TH_SurfaceHoverDesc"),
		UIText_TH_SurfaceOverlayLabel: localization.GetString("UIText_TH_SurfaceOverlayLabel"),
		UIText_TH_SurfaceOverlayDesc:  localization.GetString("UIText_TH_SurfaceOverlayDesc"),
		UIText_TH_PanelFrameLabel:     localization.GetString("UIText_TH_PanelFrameLabel"),
		UIText_TH_PanelFrameDesc:      localization.GetString("UIText_TH_PanelFrameDesc"),
		UIText_TH_PanelSurfaceLabel:   localization.GetString("UIText_TH_PanelSurfaceLabel"),
		UIText_TH_PanelSurfaceDesc:    localization.GetString("UIText_TH_PanelSurfaceDesc"),
		UIText_TH_PanelInsetLabel:     localization.GetString("UIText_TH_PanelInsetLabel"),
		UIText_TH_PanelInsetDesc:      localization.GetString("UIText_TH_PanelInsetDesc"),
		UIText_TH_PanelLampLabel:      localization.GetString("UIText_TH_PanelLampLabel"),
		UIText_TH_PanelLampDesc:       localization.GetString("UIText_TH_PanelLampDesc"),
		UIText_TH_ConsoleInfoLabel:    localization.GetString("UIText_TH_ConsoleInfoLabel"),
		UIText_TH_ConsoleInfoDesc:     localization.GetString("UIText_TH_ConsoleInfoDesc"),
		UIText_TH_ConsoleWarningLabel: localization.GetString("UIText_TH_ConsoleWarningLabel"),
		UIText_TH_ConsoleWarningDesc:  localization.GetString("UIText_TH_ConsoleWarningDesc"),
		UIText_TH_ConsoleErrorLabel:   localization.GetString("UIText_TH_ConsoleErrorLabel"),
		UIText_TH_ConsoleErrorDesc:    localization.GetString("UIText_TH_ConsoleErrorDesc"),
		UIText_TH_ConsoleSuccessLabel: localization.GetString("UIText_TH_ConsoleSuccessLabel"),
		UIText_TH_ConsoleSuccessDesc:  localization.GetString("UIText_TH_ConsoleSuccessDesc"),

		IsStationeersLaunchPadEnabled: isStationeersLaunchPadEnabled,

		// Expert Settings
		ShowExpertSettings:              fmt.Sprintf("%v", config.GetShowExpertSettings()),
		ShowExpertSettingsTrueSelected:  showExpertSettingsTrueSelected,
		ShowExpertSettingsFalseSelected: showExpertSettingsFalseSelected,

		// Expert Settings values
		Debug:                                    fmt.Sprintf("%v", config.GetIsDebugMode()),
		DebugTrueSelected:                        debugTrueSelected,
		DebugFalseSelected:                       debugFalseSelected,
		LogLevel:                                 fmt.Sprintf("%d", config.GetLogLevel()),
		LogClutterToConsole:                      fmt.Sprintf("%v", config.GetLogClutterToConsole()),
		LogClutterToConsoleTrueSelected:          logClutterToConsoleTrueSelected,
		LogClutterToConsoleFalseSelected:         logClutterToConsoleFalseSelected,
		IsSSCMEnabled:                            fmt.Sprintf("%v", config.GetIsSSCMEnabled()),
		IsSSCMEnabledTrueSelected:                isSSCMEnabledTrueSelected,
		IsSSCMEnabledFalseSelected:               isSSCMEnabledFalseSelected,
		IsConsoleEnabled:                         fmt.Sprintf("%v", config.GetIsConsoleEnabled()),
		IsConsoleEnabledTrueSelected:             isConsoleEnabledTrueSelected,
		IsConsoleEnabledFalseSelected:            isConsoleEnabledFalseSelected,
		SSUIWebPort:                              config.GetSSUIWebPort(),
		IsUpdateEnabled:                          fmt.Sprintf("%v", config.GetIsUpdateEnabled()),
		IsUpdateEnabledTrueSelected:              isUpdateEnabledTrueSelected,
		IsUpdateEnabledFalseSelected:             isUpdateEnabledFalseSelected,
		AllowPrereleaseUpdates:                   fmt.Sprintf("%v", config.GetAllowPrereleaseUpdates()),
		AllowPrereleaseUpdatesTrueSelected:       allowPrereleaseUpdatesTrueSelected,
		AllowPrereleaseUpdatesFalseSelected:      allowPrereleaseUpdatesFalseSelected,
		AllowMajorUpdates:                        fmt.Sprintf("%v", config.GetAllowMajorUpdates()),
		AllowMajorUpdatesTrueSelected:            allowMajorUpdatesTrueSelected,
		AllowMajorUpdatesFalseSelected:           allowMajorUpdatesFalseSelected,
		DiscordCharBufferSize:                    fmt.Sprintf("%d", config.GetDiscordCharBufferSize()),
		AdvertiserOverride:                       config.GetAdvertiserOverride(),
		IsStationeersLaunchPadAutoUpdatesEnabled: fmt.Sprintf("%v", config.GetIsStationeersLaunchPadAutoUpdatesEnabled()),
		IsStationeersLaunchPadAutoUpdatesEnabledTrueSelected:  isStationeersLaunchPadAutoUpdatesEnabledTrueSelected,
		IsStationeersLaunchPadAutoUpdatesEnabledFalseSelected: isStationeersLaunchPadAutoUpdatesEnabledFalseSelected,
	}
	if !data.CanViewSettings {
		data.AdvertiserOverride = ""
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
