package gamemgr

import (
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	"github.com/google/uuid"
)

func clearGameServerUUID() {
	config.ConfigMu.Lock()
	defer config.ConfigMu.Unlock()
	config.GameServerUUID = uuid.Nil
}

func createGameServerUUID() {
	config.ConfigMu.Lock()
	defer config.ConfigMu.Unlock()
	config.GameServerUUID = uuid.New()
	logger.Core.Debug("Created Game Server with internal UUID: " + config.GameServerUUID.String())
}
