package discordbot

import "github.com/bwmarrin/discordgo"

// EmbedData represents the structure for creating embeds
type EmbedData struct {
	Title       string
	Description string
	Color       int
	Fields      []EmbedField
}

// EmbedField represents a single field in the embed
type EmbedField struct {
	Name   string
	Value  string
	Inline bool
}

// generateEmbed creates a Discord embed from EmbedData
func generateEmbed(data EmbedData) *discordgo.MessageEmbed {
	fields := make([]*discordgo.MessageEmbedField, len(data.Fields))
	for i, field := range data.Fields {
		fields[i] = &discordgo.MessageEmbedField{
			Name:   field.Name,
			Value:  field.Value,
			Inline: field.Inline,
		}
	}

	return &discordgo.MessageEmbed{
		Title:       data.Title,
		Description: data.Description,
		Color:       data.Color,
		Fields:      fields,
	}
}
