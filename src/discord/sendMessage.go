package discord

import (
	"StationeersServerUI/src/config"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func SendMessageToControlChannel(message string) {
	if !config.IsDiscordEnabled {
		return
	}
	if config.DiscordSession == nil {
		fmt.Println("Discord Error: Discord is enabled but session is not initialized")
		return
	}
	//clearMessagesAboveLastN(config.ControlChannelID, 20)
	_, err := config.DiscordSession.ChannelMessageSend(config.ControlChannelID, message)
	if err != nil {
		fmt.Println("Error sending message to control channel:", err)
	} else {
		fmt.Println("Sent message to control channel:", message)
	}
}

func SendMessageToStatusChannel(message string) {
	if !config.IsDiscordEnabled {
		return
	}
	if config.DiscordSession == nil {
		fmt.Println("Discord Error: Discord is enabled but session is not initialized")
		return
	}
	//clearMessagesAboveLastN(config.StatusChannelID, 10)
	_, err := config.DiscordSession.ChannelMessageSend(config.StatusChannelID, message)
	if err != nil {
		fmt.Println("Error sending message to status channel:", err)
	} else {
		fmt.Println("Sent message to status channel:", message)
	}
}

func SendMessageToSavesChannel(message string) {
	if !config.IsDiscordEnabled {
		return
	}
	if config.DiscordSession == nil {
		fmt.Println("Discord Error: Discord is enabled but session is not initialized")
		return
	}
	//clearMessagesAboveLastN(config.SaveChannelID, 300)
	_, err := config.DiscordSession.ChannelMessageSend(config.SaveChannelID, message)
	if err != nil {
		fmt.Println("Error sending message to saves channel:", err)
	} else {
		fmt.Println("Sent message to saves channel:", message)
	}
}

// v4 OK, unsused in 4.3
func sendMessageToErrorChannel(message string) []*discordgo.Message {
	if !config.IsDiscordEnabled {
		return nil
	}
	if config.DiscordSession == nil {
		fmt.Println("Discord Error: Discord is enabled but session is not initialized")
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
			sentMessage, err := config.DiscordSession.ChannelMessageSend(config.ErrorChannelID, message[:splitIndex])
			if err != nil {
				fmt.Println("Error sending message to error channel:", err)
				return sentMessages // Return whatever was sent before the error
			}

			// Add sent message to the list
			sentMessages = append(sentMessages, sentMessage)

			// Remove the sent chunk from the message
			message = message[splitIndex:]
		} else {
			// Send the remaining part of the message
			sentMessage, err := config.DiscordSession.ChannelMessageSend(config.ErrorChannelID, message)
			if err != nil {
				fmt.Println("Error sending message to error channel:", err)
				return sentMessages // Return whatever was sent before the error
			}

			// Add the final sent message to the list
			sentMessages = append(sentMessages, sentMessage)
			break
		}
	}

	return sentMessages
}

func sendControlMessage() {
	messageContent := "Control Panel:\n\nReact with the following to perform actions:\n" +
		"▶️ Start the server\n\n" +
		"⏹️ Stop the server\n\n" +
		"♻️ Restart the server\n\n"

	msg, err := config.DiscordSession.ChannelMessageSend(config.ControlPanelChannelID, messageContent)
	if err != nil {
		fmt.Println("Error sending control message:", err)
		return
	}

	// Add reactions (acting as buttons) to the control message
	config.DiscordSession.MessageReactionAdd(config.ControlPanelChannelID, msg.ID, "▶️") // Start
	config.DiscordSession.MessageReactionAdd(config.ControlPanelChannelID, msg.ID, "⏹️") // Stop
	config.DiscordSession.MessageReactionAdd(config.ControlPanelChannelID, msg.ID, "♻️") // Restart
	config.ControlMessageID = msg.ID
	if config.Branch == "Prod" {
		clearMessagesAboveLastN(config.ControlPanelChannelID, 1)
	}
	if config.Branch != "Prod" {
		clearMessagesAboveLastN(config.ControlPanelChannelID, 15)
	}
}

// CLEAR MESSAGES
func clearMessagesAboveLastN(channelID string, keep int) {
	go func() {
		if !config.IsDiscordEnabled {
			return
		}
		if config.DiscordSession == nil {
			fmt.Println("Discord Error: Discord is enabled but session is not initialized")
			return
		}

		// Retrieve the last 100 messages in the channel (Discord API limit)
		messages, err := config.DiscordSession.ChannelMessages(channelID, 100, "", "", "")
		if err != nil {
			fmt.Printf("Error fetching messages from channel %s: %v\n", channelID, err)
			return
		}

		// If there are more than 'keep' messages, delete the excess ones
		if len(messages) > keep {
			for _, message := range messages[keep:] {
				err := config.DiscordSession.ChannelMessageDelete(channelID, message.ID)
				if err != nil {
					fmt.Printf("Error deleting message %s in channel %s: %v\n", message.ID, channelID, err)
				} else {
					fmt.Printf("Deleted message %s in channel %s\n", message.ID, channelID)
				}
			}
		}
	}()
}
