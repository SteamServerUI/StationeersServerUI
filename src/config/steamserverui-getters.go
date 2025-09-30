package config

// GetIsSteamServerUI returns if the system is in SteamServerUI mode
func GetIsSteamServerUI() bool {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return IsSteamServerUI
}

// GetRunFilesFolder returns the RunFilesFolder
func GetRunFilesFolder() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return RunFilesFolder
}

// GetRunfileGame returns the RunfileGame
func GetRunfileIdentifier() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return RunfileIdentifier
}
