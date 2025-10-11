package discordrpc

import (
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
	"github.com/jacksonthemaster/discordrichpresence"
)

// StartDiscordRPC starts the Discord Rich Presence client in a non-blocking manner.
// It returns the client so the caller can manage its lifecycle (e.g., call Close()).
func StartDiscordRPC() (*discordrichpresence.Client, error) {
	client := discordrichpresence.NewClient("1408848834875887669")
	activity := discordrichpresence.NewActivity().
		State("Managing a Gameserver").
		Details("Your one-stop-shop for running a Gameserver").
		StartTime(time.Now()).
		LargeImage("logo", "The easy to use Dedicated Server Manager").
		SmallImage("rocket", "Online and Active").
		Type(0).
		Build()
	if err := client.StartWithActivity(activity, 30*time.Second); err != nil {
		return nil, err
	}

	logger.Core.Debug("Discord Rich Presence started successfully!")
	return client, nil
}
