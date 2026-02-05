package dashboard

import (
	"fmt"
	"strings"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/core/loader"
	"github.com/charmbracelet/lipgloss"
)

// Config Items Definition

// buildConfigItems creates a list of editable config items grouped by section
func buildConfigItems() []ConfigItem {
	return []ConfigItem{
		// ─────────────────────────────────────────────────────────────────────
		// Basic Settings
		// ─────────────────────────────────────────────────────────────────────
		{Key: "ServerName", Label: "Server Name", Value: config.GetServerName(), Type: "string", Section: ConfigSectionBasic, Description: "Display name for your server"},
		{Key: "SaveName", Label: "Save Name", Value: config.GetSaveName(), Type: "string", Section: ConfigSectionBasic, Description: "World save file name"},
		{Key: "ServerMaxPlayers", Label: "Max Players", Value: config.GetServerMaxPlayers(), Type: "string", Section: ConfigSectionBasic, Description: "Maximum concurrent players"},
		{Key: "ServerPassword", Label: "Server Password", Value: config.GetServerPassword(), Type: "password", Section: ConfigSectionBasic, Description: "Password to join server (blank = none)"},
		{Key: "ServerVisible", Label: "Server Visible", Value: boolToStr(config.GetServerVisible()), Type: "bool", Section: ConfigSectionBasic, Description: "List server publicly"},
		{Key: "AutoSave", Label: "Auto-Save", Value: boolToStr(config.GetAutoSave()), Type: "bool", Section: ConfigSectionBasic, Description: "Enable automatic world saves"},
		{Key: "SaveInterval", Label: "Save Interval", Value: config.GetSaveInterval(), Type: "string", Section: ConfigSectionBasic, Description: "Time between saves (e.g., 300)"},

		// ─────────────────────────────────────────────────────────────────────
		// World Generation Settings
		// ─────────────────────────────────────────────────────────────────────
		{Key: "WorldID", Label: "World ID", Value: config.GetWorldID(), Type: "string", Section: ConfigSectionWorldGen, Description: "World seed identifier"},
		{Key: "Difficulty", Label: "Difficulty", Value: config.GetDifficulty(), Type: "string", Section: ConfigSectionWorldGen, Description: "Game difficulty setting"},
		{Key: "StartCondition", Label: "Start Condition", Value: config.GetStartCondition(), Type: "string", Section: ConfigSectionWorldGen, Description: "Initial spawn condition"},
		{Key: "StartLocation", Label: "Start Location", Value: config.GetStartLocation(), Type: "string", Section: ConfigSectionWorldGen, Description: "Starting world location"},

		// ─────────────────────────────────────────────────────────────────────
		// Network Settings
		// ─────────────────────────────────────────────────────────────────────
		{Key: "GamePort", Label: "Game Port", Value: config.GetGamePort(), Type: "string", Section: ConfigSectionNetwork, Description: "Main game connection port"},
		{Key: "UpdatePort", Label: "Update Port", Value: config.GetUpdatePort(), Type: "string", Section: ConfigSectionNetwork, Description: "Server query port"},
		{Key: "UPNPEnabled", Label: "UPnP Enabled", Value: boolToStr(config.GetUPNPEnabled()), Type: "bool", Section: ConfigSectionNetwork, Description: "Automatic port forwarding"},
		{Key: "LocalIpAddress", Label: "Local IP", Value: config.GetLocalIpAddress(), Type: "string", Section: ConfigSectionNetwork, Description: "Bind to specific IP (blank = auto)"},
		{Key: "StartLocalHost", Label: "Start LocalHost", Value: boolToStr(config.GetStartLocalHost()), Type: "bool", Section: ConfigSectionNetwork, Description: "Start in localhost mode"},
		{Key: "UseSteamP2P", Label: "Use Steam P2P", Value: boolToStr(config.GetUseSteamP2P()), Type: "bool", Section: ConfigSectionNetwork, Description: "Use Steam peer-to-peer networking"},

		// ─────────────────────────────────────────────────────────────────────
		// Advanced Settings
		// ─────────────────────────────────────────────────────────────────────
		{Key: "AutoStartServerOnStartup", Label: "Auto-Start Server", Value: boolToStr(config.GetAutoStartServerOnStartup()), Type: "bool", Section: ConfigSectionAdvanced, Description: "Start game server when SSUI starts"},
		{Key: "AutoRestartServerTimer", Label: "Auto-Restart Timer", Value: config.GetAutoRestartServerTimer(), Type: "string", Section: ConfigSectionAdvanced, Description: "Auto-restart interval (0 = disabled)"},
		{Key: "AutoPauseServer", Label: "Auto-Pause", Value: boolToStr(config.GetAutoPauseServer()), Type: "bool", Section: ConfigSectionAdvanced, Description: "Pause when no players"},
		{Key: "IsSSCMEnabled", Label: "SSCM/BepInEx", Value: boolToStr(config.GetIsSSCMEnabled()), Type: "bool", Section: ConfigSectionAdvanced, Description: "Enable mod support"},
		{Key: "Branch", Label: "Game Branch", Value: config.GetGameBranch(), Type: "string", Section: ConfigSectionAdvanced, Description: "Steam branch (public, beta, etc.)"},
		{Key: "LogClutterToConsole", Label: "Log Clutter", Value: boolToStr(config.GetLogClutterToConsole()), Type: "bool", Section: ConfigSectionAdvanced, Description: "Show verbose logs in console"},
	}
}

// Config Saving

// saveConfigItem saves a single config item using the appropriate setter
func saveConfigItem(item ConfigItem) error {
	switch item.Key {
	// Basic Settings
	case "ServerName":
		return config.SetServerName(item.Value)
	case "SaveName":
		return config.SetSaveName(item.Value)
	case "ServerMaxPlayers":
		return config.SetServerMaxPlayers(item.Value)
	case "ServerPassword":
		return config.SetServerPassword(item.Value)
	case "ServerVisible":
		return config.SetServerVisible(strToBool(item.Value))
	case "AutoSave":
		return config.SetAutoSave(strToBool(item.Value))
	case "SaveInterval":
		return config.SetSaveInterval(item.Value)

	// World Generation Settings
	case "WorldID":
		return config.SetWorldID(item.Value)
	case "Difficulty":
		return config.SetDifficulty(item.Value)
	case "StartCondition":
		return config.SetStartCondition(item.Value)
	case "StartLocation":
		return config.SetStartLocation(item.Value)

	// Network Settings
	case "GamePort":
		return config.SetGamePort(item.Value)
	case "UpdatePort":
		return config.SetUpdatePort(item.Value)
	case "UPNPEnabled":
		return config.SetUPNPEnabled(strToBool(item.Value))
	case "LocalIpAddress":
		return config.SetLocalIpAddress(item.Value)
	case "StartLocalHost":
		return config.SetStartLocalHost(strToBool(item.Value))
	case "UseSteamP2P":
		return config.SetUseSteamP2P(strToBool(item.Value))

	// Advanced Settings
	case "AutoStartServerOnStartup":
		return config.SetAutoStartServerOnStartup(strToBool(item.Value))
	case "AutoRestartServerTimer":
		return config.SetAutoRestartServerTimer(item.Value)
	case "AutoPauseServer":
		return config.SetAutoPauseServer(strToBool(item.Value))
	case "IsSSCMEnabled":
		return config.SetIsSSCMEnabled(strToBool(item.Value))
	case "Branch":
		return config.SetGameBranch(item.Value)
	case "LogClutterToConsole":
		return config.SetLogClutterToConsole(strToBool(item.Value))

	}
	return nil
}

// saveAllConfigChanges saves all config items and reloads the backend
func saveAllConfigChanges(items []ConfigItem) error {
	for _, item := range items {
		if err := saveConfigItem(item); err != nil {
			return fmt.Errorf("failed to save %s: %w", item.Key, err)
		}
	}
	loader.ReloadBackend()
	return nil
}

// Panel Rendering

func (m Model) renderConfigPanel() string {
	var lines []string

	// Header with save status
	headerLine := RenderSectionTitle("Configuration Editor")
	if m.configHasChanges {
		headerLine += "  " + lipgloss.NewStyle().Foreground(Yellow).Bold(true).Render("• Unsaved Changes, Ctrl+S to Save")
	}
	if m.configStatusMsg != "" && m.configStatusTick > 0 {
		headerLine += "  " + lipgloss.NewStyle().Foreground(Green).Render(m.configStatusMsg)
	}
	lines = append(lines, headerLine)
	lines = append(lines, "")

	lines = append(lines, MutedStyle.Render("  Use ↑↓ to navigate, Enter to edit, Space to toggle, Ctrl+S to save"))
	lines = append(lines, "")

	navIndex := 0
	for section := range ConfigSectionCount {
		sectionName := configSectionNames[section]
		isOpen := m.configSectionOpen[section]

		expandIcon := "▶"
		if isOpen {
			expandIcon = "▼"
		}

		sectionStyle := lipgloss.NewStyle().Bold(true).Foreground(PurpleLight)
		sectionHeader := sectionStyle.Render(expandIcon + " " + sectionName)

		if !m.configEditing && m.configSelectedIndex == navIndex && m.activePanel == PanelConfig {
			sectionHeader = lipgloss.NewStyle().
				Background(Purple).
				Foreground(White).
				Bold(true).
				Render(expandIcon + " " + sectionName)
		}

		lines = append(lines, sectionHeader)
		navIndex++

		if isOpen {
			sectionItems := m.getItemsForSection(section)
			for _, item := range sectionItems {
				line := m.renderConfigItem(item, navIndex)
				lines = append(lines, line)
				navIndex++
			}
		}
		lines = append(lines, "")
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return m.renderPanelContainer(content)
}

func (m Model) renderConfigItem(item ConfigItem, index int) string {
	isSelected := !m.configEditing && m.configSelectedIndex == index && m.activePanel == PanelConfig
	isEditing := m.configEditing && m.configSelectedIndex == index && m.activePanel == PanelConfig

	labelStyle := lipgloss.NewStyle().Width(24).Foreground(Gray400)
	label := labelStyle.Render("    " + item.Label + ":")

	var valueDisplay string

	if isEditing {
		editStyle := lipgloss.NewStyle().
			Background(Gray700).
			Foreground(White).
			Padding(0, 1)

		displayValue := m.configEditValue
		if len(displayValue) == 0 {
			displayValue = " "
		}
		valueDisplay = editStyle.Render(displayValue + "▌")
	} else if item.Type == "bool" {
		if strToBool(item.Value) {
			valueDisplay = lipgloss.NewStyle().Foreground(Green).Render(CheckMark + " true")
		} else {
			valueDisplay = lipgloss.NewStyle().Foreground(Red).Render(CrossMark + " false")
		}
	} else if item.Type == "password" && item.Value != "" {
		valueDisplay = MutedStyle.Render("••••••••")
	} else if item.Value == "" {
		valueDisplay = MutedStyle.Render("(not set)")
	} else {
		valueDisplay = ValueStyle.Render(item.Value)
	}

	if isSelected {
		label = lipgloss.NewStyle().
			Width(24).
			Background(Purple).
			Foreground(White).
			Render("  ▸ " + item.Label + ":")
	}

	return label + " " + valueDisplay
}

func (m Model) getItemsForSection(section ConfigSection) []ConfigItem {
	var items []ConfigItem
	for _, item := range m.configItems {
		if item.Section == section {
			items = append(items, item)
		}
	}
	return items
}

func (m Model) getTotalConfigItems() int {
	count := 0
	for section := range ConfigSectionCount {
		count++ // Section header
		if m.configSectionOpen[section] {
			count += len(m.getItemsForSection(section))
		}
	}
	return count
}

// getConfigItemAtIndex returns a config item at given nav index
func (m Model) getConfigItemAtIndex(index int) (*ConfigItem, bool, ConfigSection) {
	currentIndex := 0
	for section := range ConfigSectionCount {
		// Check if we're on the section header
		if currentIndex == index {
			return nil, true, section // It's a section header
		}
		currentIndex++

		if m.configSectionOpen[section] {
			items := m.getItemsForSection(section)
			for i := range items {
				if currentIndex == index {
					return &m.configItems[m.getGlobalItemIndex(section, i)], false, section
				}
				currentIndex++
			}
		}
	}
	return nil, false, 0
}

// getGlobalItemIndex converts section-local index to global configItems index
func (m Model) getGlobalItemIndex(section ConfigSection, localIndex int) int {
	globalIndex := 0
	for _, item := range m.configItems {
		if item.Section == section {
			if localIndex == 0 {
				return globalIndex
			}
			localIndex--
		}
		globalIndex++
	}
	return -1
}

// Helpers

func boolToStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func strToBool(s string) bool {
	return strings.ToLower(s) == "true"
}
