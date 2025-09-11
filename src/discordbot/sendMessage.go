package discordbot

import (
	"strings"
	"time"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"

	"github.com/bwmarrin/discordgo"
)

var ControlMessageID string

func SendMessageToControlChannel(message string) {
	if !config.GetIsDiscordEnabled() {
		return
	}
	if config.DiscordSession == nil {
		logger.Discord.Error("Discord Error: Discord is enabled but session is not initialized")
		return
	}
	//clearMessagesAboveLastN(config.ControlChannelID, 20)
	_, err := config.DiscordSession.ChannelMessageSend(config.GetControlChannelID(), message)
	if err != nil {
		logger.Discord.Error("Error sending message to control channel: " + err.Error())
	}
}

func SendMessageToStatusChannel(message string) {
	if !config.GetIsDiscordEnabled() {
		return
	}
	if config.DiscordSession == nil {
		logger.Discord.Error("Discord Error: Discord is enabled but session is not initialized")
		return
	}
	//clearMessagesAboveLastN(config.StatusChannelID, 10)
	_, err := config.DiscordSession.ChannelMessageSend(config.GetStatusChannelID(), message)
	if err != nil {
		logger.Discord.Error("Error sending message to status channel: " + err.Error())
	}
}

func SendMessageToSavesChannel(message string) {
	if !config.GetIsDiscordEnabled() {
		return
	}
	if config.DiscordSession == nil {
		logger.Discord.Error("Discord Error: Discord is enabled but session is not initialized")
		return
	}
	//clearMessagesAboveLastN(config.SaveChannelID, 300)
	_, err := config.DiscordSession.ChannelMessageSend(config.GetSaveChannelID(), message)
	if err != nil {
		logger.Discord.Error("Error sending message to saves channel: " + err.Error())
	}
}

func SendUntrackedMessageToErrorChannel(message string) {
	if !config.GetIsDiscordEnabled() {
		return
	}
	if config.DiscordSession == nil {
		logger.Discord.Error("Discord Error: Discord is enabled but session is not initialized")
		return
	}

	maxMessageLength := 2000 // Discord's message character limit

	// Function to split the message into chunks and send each one
	for len(message) > 0 {
		if len(message) > maxMessageLength {
			// Find a safe split point, for example, the last newline before the limit
			splitIndex := strings.LastIndex(message[:maxMessageLength], "\n")
			if splitIndex == -1 {
				splitIndex = maxMessageLength // No newline found, force split at max length
			}

			// Send the chunk
			_, err := config.DiscordSession.ChannelMessageSend(config.GetErrorChannelID(), message[:splitIndex])
			if err != nil {
				logger.Discord.Error("Error sending message to error channel: " + err.Error())
				return
			}

			// Remove the sent chunk from the message
			message = message[splitIndex:]
		} else {
			// Send the remaining part of the message
			_, err := config.DiscordSession.ChannelMessageSend(config.GetErrorChannelID(), message)
			if err != nil {
				logger.Discord.Error("Error sending message to error channel: " + err.Error())
				return // Return whatever was sent before the error
			}
			break
		}
	}
}

// unsused (replaced with SendUntrackedMessageToErrorChannel) in 4.3, needed for having a restart button on the last exception message like in v2. Might remve this in the future, but for now let's keep it.
func sendMessageToErrorChannel(message string) []*discordgo.Message {
	if !config.GetIsDiscordEnabled() {
		return nil
	}
	if config.DiscordSession == nil {
		logger.Discord.Error("Discord Error: Discord is enabled but session is not initialized")
		return nil
	}

	maxMessageLength := 2000 // Discord's message character limit
	var sentMessages []*discordgo.Message

	// Function to split the message into chunks and send each one
	for len(message) > 0 {
		if len(message) > maxMessageLength {
			// Find a safe split point, for example, the last newline before the limit
			splitIndex := strings.LastIndex(message[:maxMessageLength], "\n")
			if splitIndex == -1 {
				splitIndex = maxMessageLength // No newline found, force split at max length
			}

			// Send the chunk
			sentMessage, err := config.DiscordSession.ChannelMessageSend(config.GetErrorChannelID(), message[:splitIndex])
			if err != nil {
				logger.Discord.Error("Error sending message to error channel: " + err.Error())
				return sentMessages // Return whatever was sent before the error
			}

			// Add sent message to the list
			sentMessages = append(sentMessages, sentMessage)

			// Remove the sent chunk from the message
			message = message[splitIndex:]
		} else {
			// Send the remaining part of the message
			sentMessage, err := config.DiscordSession.ChannelMessageSend(config.GetErrorChannelID(), message)
			if err != nil {
				logger.Discord.Error("Error sending message to error channel: " + err.Error())
				return sentMessages // Return whatever was sent before the error
			}

			// Add the final sent message to the list
			sentMessages = append(sentMessages, sentMessage)
			break
		}
	}

	return sentMessages
}

// This function is used to clear messages above the last N messages in a channel. If you call this with 5, it will clear all messages in the channel besides the most recent 5.
func clearMessagesAboveLastN(channelID string, keep int) {
	go func() {
		if !config.GetIsDiscordEnabled() {
			return
		}
		if config.DiscordSession == nil {
			logger.Discord.Error("Discord Error: Discord is enabled but session is not initialized")
			return
		}

		// Retrieve the last 100 messages in the channel (Discord API limit)
		messages, err := config.DiscordSession.ChannelMessages(channelID, 100, "", "", "")
		if err != nil {
			logger.Discord.Error("Error fetching messages from channel " + channelID + ": " + err.Error())
			return
		}

		// If there are more than 'keep' messages, delete the excess ones
		if len(messages) > keep {
			for _, message := range messages[keep:] {
				err := config.DiscordSession.ChannelMessageDelete(channelID, message.ID)
				if err != nil {
					logger.Discord.Error("Error deleting message " + message.ID + " in channel " + channelID + ": " + err.Error())
				}
			}
		}
	}()
}

// sendTemporaryMessage sends a message to the specified channel and deletes it after the given duration.
func sendTemporaryMessage(s *discordgo.Session, channelID, message string, duration time.Duration) {
	// Send the message
	msg, err := s.ChannelMessageSend(channelID, message)
	if err != nil {
		logger.Discord.Error("Error sending temporary message: " + err.Error())
		return
	}

	// Schedule deletion after the specified duration
	go func() {
		time.Sleep(duration)
		err := s.ChannelMessageDelete(channelID, msg.ID)
		if err != nil {
			logger.Discord.Error("Error deleting temporary message: " + err.Error())
		}
	}()
}
