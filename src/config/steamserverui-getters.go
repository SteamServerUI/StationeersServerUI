package config

// GetRunFilesFolder returns the RunFilesFolder
func GetRunFilesFolder() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return RunFilesFolder
}

// GetRunfileGame returns the RunfileGame
func GetRunfileGame() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return RunfileGame
}
