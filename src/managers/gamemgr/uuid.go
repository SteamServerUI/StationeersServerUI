package gamemgr

import (
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
	"github.com/google/uuid"
)

var GameServerUUID uuid.UUID

func clearGameServerUUID() {
	GameServerUUID = uuid.Nil
}

func createGameServerUUID() {
	GameServerUUID = uuid.New()
	logger.Core.Debug("Created Game Server with internal UUID: " + GameServerUUID.String())
}
