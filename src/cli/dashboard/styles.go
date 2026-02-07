package dashboard

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// SSUI Dashboard Style System
// A comprehensive styling module using lipgloss for a polished terminal UI

// Color Palette - SSUI Brand Colors

var (
	// Primary brand colors
	Purple      = lipgloss.Color("#7C3AED")
	PurpleLight = lipgloss.Color("#A78BFA")
	PurpleDark  = lipgloss.Color("#5B21B6")

	// Semantic colors
	Green       = lipgloss.Color("#10B981")
	GreenLight  = lipgloss.Color("#34D399")
	Red         = lipgloss.Color("#EF4444")
	RedLight    = lipgloss.Color("#F87171")
	Yellow      = lipgloss.Color("#F59E0B")
	YellowLight = lipgloss.Color("#FBBF24")
	Blue        = lipgloss.Color("#3B82F6")
	BlueLight   = lipgloss.Color("#60A5FA")
	Cyan        = lipgloss.Color("#06B6D4")
	CyanLight   = lipgloss.Color("#22D3EE")

	// Neutral colors
	White   = lipgloss.Color("#F9FAFB")
	Gray100 = lipgloss.Color("#F3F4F6")
	Gray200 = lipgloss.Color("#E5E7EB")
	Gray300 = lipgloss.Color("#D1D5DB")
	Gray400 = lipgloss.Color("#9CA3AF")
	Gray500 = lipgloss.Color("#6B7280")
	Gray600 = lipgloss.Color("#4B5563")
	Gray700 = lipgloss.Color("#374151")
	Gray800 = lipgloss.Color("#1F2937")
	Gray900 = lipgloss.Color("#111827")
	Black   = lipgloss.Color("#030712")
)

const (
	// Status indicators
	BulletFilled = "●"
	BulletEmpty  = "○"
	BulletHalf   = "◐"
	CheckMark    = "✓"
	CrossMark    = "✗"
	Lightning    = "⚡"
	Gear         = "⚙"
	Server       = "󰒋"
	Globe        = "🌐"
	Clock        = "⏱"
	Calendar     = "📅"
	Folder       = "📁"
	User         = "👤"
	Users        = "👥"
	Shield       = "🛡"
	Plug         = "🔌"
	Bot          = "🤖"
	Save         = "💾"
	Refresh      = "🔄"
	Play         = "▶"
	Stop         = "■"
	Pause        = "⏸"
	Rocket       = "🚀"
	Fire         = "🔥"
	Sparkle      = "✨"

	// Box drawing characters
	BoxTopLeft     = "╭"
	BoxTopRight    = "╮"
	BoxBottomLeft  = "╰"
	BoxBottomRight = "╯"
	BoxHorizontal  = "─"
	BoxVertical    = "│"
	BoxCross       = "┼"
	BoxTeeRight    = "├"
	BoxTeeLeft     = "┤"
	BoxTeeDown     = "┬"
	BoxTeeUp       = "┴"

	// Double line box
	BoxDoubleH = "═"
	BoxDoubleV = "║"

	// Progress bar characters
	ProgressFull  = "█"
	ProgressHigh  = "▓"
	ProgressMid   = "▒"
	ProgressLow   = "░"
	ProgressEmpty = "░"

	// Sparkline characters (for mini graphs)
	Spark1 = "▁"
	Spark2 = "▂"
	Spark3 = "▃"
	Spark4 = "▄"
	Spark5 = "▅"
	Spark6 = "▆"
	Spark7 = "▇"
	Spark8 = "█"

	// Arrows and pointers
	ArrowRight = "→"
	ArrowLeft  = "←"
	ArrowUp    = "↑"
	ArrowDown  = "↓"
	Pointer    = "▸"
	Diamond    = "◆"
	Triangle   = "▲"
)

// Spinner frames for animated loading indicators
var SpinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
var SpinnerDots = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
var SpinnerPulse = []string{"●", "◐", "○", "◑"}

// Base Styles

var (
	// Base text styles
	BaseStyle = lipgloss.NewStyle().
			Foreground(White)

	BoldStyle = lipgloss.NewStyle().
			Foreground(White).
			Bold(true)

	MutedStyle = lipgloss.NewStyle().
			Foreground(Gray500)

	DimStyle = lipgloss.NewStyle().
			Foreground(Gray600)
)

// Logo and Header Styles

var (
	// ASCII art logo style with gradient effect
	LogoStyle = lipgloss.NewStyle().
			Foreground(Purple).
			Bold(true)

	LogoAccentStyle = lipgloss.NewStyle().
			Foreground(PurpleLight)

	// Main header container
	HeaderContainerStyle = lipgloss.NewStyle().
				Background(Gray800).
				Padding(0, 2)

	// Title bar style
	TitleStyle = lipgloss.NewStyle().
			Foreground(White).
			Background(Purple).
			Bold(true).
			Padding(0, 2).
			MarginRight(1)

	// Version badge
	VersionBadgeStyle = lipgloss.NewStyle().
				Foreground(PurpleLight).
				Background(Gray700).
				Padding(0, 1)

	// Status pill styles (for header indicators)
	StatusPillOnline = lipgloss.NewStyle().
				Foreground(Gray900).
				Background(Green).
				Bold(true).
				Padding(0, 1)

	StatusPillOffline = lipgloss.NewStyle().
				Foreground(White).
				Background(Red).
				Bold(true).
				Padding(0, 1)

	StatusPillWarning = lipgloss.NewStyle().
				Foreground(Gray900).
				Background(Yellow).
				Bold(true).
				Padding(0, 1)
)

// Tab Bar Styles

var (
	TabBarStyle = lipgloss.NewStyle().
			Background(Gray800).
			Padding(0, 1)

	TabActiveStyle = lipgloss.NewStyle().
			Foreground(White).
			Background(Purple).
			Bold(true).
			Padding(0, 2)

	TabInactiveStyle = lipgloss.NewStyle().
				Foreground(Gray400).
				Background(Gray700).
				Padding(0, 2)

	TabSeparatorStyle = lipgloss.NewStyle().
				Foreground(Gray600)
)

// Panel Styles

var (
	// Main panel container with rounded border
	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Purple).
			Padding(1, 2)

	// Panel variants
	PanelActiveStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(PurpleLight).
				Padding(1, 2)

	PanelDimStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Gray600).
			Padding(1, 2)

	// Section header within panels
	SectionHeaderStyle = lipgloss.NewStyle().
				Foreground(PurpleLight).
				Bold(true).
				MarginBottom(1)

	// Divider line
	DividerStyle = lipgloss.NewStyle().
			Foreground(Gray600)
)

// Status Indicator Styles

var (
	// Online/Offline indicators
	OnlineStyle = lipgloss.NewStyle().
			Foreground(Green).
			Bold(true)

	OfflineStyle = lipgloss.NewStyle().
			Foreground(Red).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(Yellow).
			Bold(true)

	InfoStyle = lipgloss.NewStyle().
			Foreground(Blue)

	// Feature status badges
	FeatureEnabledStyle = lipgloss.NewStyle().
				Foreground(Gray900).
				Background(Green).
				Padding(0, 1)

	FeatureDisabledStyle = lipgloss.NewStyle().
				Foreground(Gray400).
				Background(Gray700).
				Padding(0, 1)

	FeatureActiveStyle = lipgloss.NewStyle().
				Foreground(Gray900).
				Background(Cyan).
				Padding(0, 1)
)

// Data Display Styles

var (
	// Key-value pair styles
	LabelStyle = lipgloss.NewStyle().
			Foreground(Gray400).
			Width(20)

	ValueStyle = lipgloss.NewStyle().
			Foreground(White)

	ValueHighlightStyle = lipgloss.NewStyle().
				Foreground(PurpleLight).
				Bold(true)

	ValueSuccessStyle = lipgloss.NewStyle().
				Foreground(Green)

	ValueErrorStyle = lipgloss.NewStyle().
			Foreground(Red)

	ValueWarningStyle = lipgloss.NewStyle().
				Foreground(Yellow)

	// Numeric value style
	NumberStyle = lipgloss.NewStyle().
			Foreground(Cyan)

	// ID/Code style (monospace look)
	CodeStyle = lipgloss.NewStyle().
			Foreground(Gray300).
			Italic(true)
)

// Player List Styles

var (
	PlayerRowStyle = lipgloss.NewStyle().
			Padding(0, 1)

	PlayerIndexStyle = lipgloss.NewStyle().
				Foreground(Gray500).
				Width(4)

	PlayerNameStyle = lipgloss.NewStyle().
			Foreground(GreenLight).
			Bold(true)

	PlayerIDStyle = lipgloss.NewStyle().
			Foreground(Gray500).
			Italic(true)

	PlayerEmptyStyle = lipgloss.NewStyle().
				Foreground(Gray500).
				Italic(true)
)

// Log Viewer Styles

var (
	LogContainerStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(Gray600).
				Padding(0, 1)

	LogTimestampStyle = lipgloss.NewStyle().
				Foreground(Gray500)

	LogLevelDebugStyle = lipgloss.NewStyle().
				Foreground(Gray400)

	LogLevelInfoStyle = lipgloss.NewStyle().
				Foreground(Blue)

	LogLevelWarnStyle = lipgloss.NewStyle().
				Foreground(Yellow)

	LogLevelErrorStyle = lipgloss.NewStyle().
				Foreground(Red).
				Bold(true)

	LogMessageStyle = lipgloss.NewStyle().
			Foreground(Gray300)

	LogScrollIndicatorStyle = lipgloss.NewStyle().
				Foreground(Gray500).
				Italic(true)
)

// Footer Styles

var (
	FooterStyle = lipgloss.NewStyle().
			Background(Gray800).
			Foreground(Gray400).
			Padding(0, 2)

	FooterBrandStyle = lipgloss.NewStyle().
				Foreground(Purple).
				Bold(true)

	FooterStatStyle = lipgloss.NewStyle().
			Foreground(Gray400)

	FooterSeparatorStyle = lipgloss.NewStyle().
				Foreground(Gray600)

	// Key hint styles
	KeyStyle = lipgloss.NewStyle().
			Foreground(Gray300).
			Background(Gray700).
			Padding(0, 1)

	KeyDescStyle = lipgloss.NewStyle().
			Foreground(Gray500)
)

// Help Styles

var (
	HelpStyle = lipgloss.NewStyle().
			Foreground(Gray500)

	HelpKeyStyle = lipgloss.NewStyle().
			Foreground(Gray300)

	HelpDescStyle = lipgloss.NewStyle().
			Foreground(Gray500)

	HelpSeparatorStyle = lipgloss.NewStyle().
				Foreground(Gray600)
)

// Utility Functions

// RenderProgressBar creates a visual progress bar
func RenderProgressBar(current, max, width int) string {
	if max == 0 {
		max = 1
	}
	if width < 5 {
		width = 10
	}

	ratio := float64(current) / float64(max)
	filled := min(int(ratio*float64(width)), width)

	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += ProgressFull
		} else {
			bar += ProgressEmpty
		}
	}

	// Color based on fill level
	var style lipgloss.Style
	switch {
	case ratio >= 0.9:
		style = lipgloss.NewStyle().Foreground(Red)
	case ratio >= 0.7:
		style = lipgloss.NewStyle().Foreground(Yellow)
	default:
		style = lipgloss.NewStyle().Foreground(Green)
	}

	return style.Render(bar)
}

// RenderStatusDot returns a colored status dot
func RenderStatusDot(online bool) string {
	if online {
		return OnlineStyle.Render(BulletFilled)
	}
	return OfflineStyle.Render(BulletEmpty)
}

// RenderFeatureBadge creates a badge for feature status
func RenderFeatureBadge(name string, enabled bool) string {
	if enabled {
		return FeatureEnabledStyle.Render(CheckMark + " " + name)
	}
	return FeatureDisabledStyle.Render(CrossMark + " " + name)
}

// RenderKeyValue renders a styled key-value pair
func RenderKeyValue(key, value string) string {
	return LabelStyle.Render(key+":") + " " + ValueStyle.Render(value)
}

// RenderKeyValueHighlight renders a key-value pair with highlighted value
func RenderKeyValueHighlight(key, value string) string {
	return LabelStyle.Render(key+":") + " " + ValueHighlightStyle.Render(value)
}

// RenderDivider creates a horizontal divider line
func RenderDivider(width int) string {
	if width < 1 {
		width = 40
	}
	var line strings.Builder
	for i := 0; i < width; i++ {
		line.WriteString(BoxHorizontal)
	}
	return DividerStyle.Render(line.String())
}

// RenderSectionTitle creates a styled section title
func RenderSectionTitle(title string) string {
	return SectionHeaderStyle.Render(Diamond + " " + title)
}

// Gradient applies a simple two-color gradient to text (line by line)
func Gradient(text string, from, to lipgloss.Color) string {
	// Simple implementation: alternate colors by line
	lines := splitLines(text)
	result := ""
	for i, line := range lines {
		var style lipgloss.Style
		if i%2 == 0 {
			style = lipgloss.NewStyle().Foreground(from)
		} else {
			style = lipgloss.NewStyle().Foreground(to)
		}
		result += style.Render(line)
		if i < len(lines)-1 {
			result += "\n"
		}
	}
	return result
}

// splitLines splits text into lines
func splitLines(text string) []string {
	var lines []string
	current := ""
	for _, r := range text {
		if r == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(r)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

// GetSpinnerFrame returns the current spinner frame based on tick count
func GetSpinnerFrame(tickCount int) string {
	return SpinnerFrames[tickCount%len(SpinnerFrames)]
}

// GetSpinnerDot returns the current dot spinner frame based on tick count
func GetSpinnerDot(tickCount int) string {
	return SpinnerDots[tickCount%len(SpinnerDots)]
}

// GetPulse returns the current pulse frame based on tick count
func GetPulse(tickCount int) string {
	return SpinnerPulse[tickCount%len(SpinnerPulse)]
}

// RenderSparkline creates a mini sparkline from values
func RenderSparkline(values []int, maxVal int) string {
	sparkChars := []string{Spark1, Spark2, Spark3, Spark4, Spark5, Spark6, Spark7, Spark8}
	if maxVal <= 0 {
		maxVal = 1
		for _, v := range values {
			if v > maxVal {
				maxVal = v
			}
		}
	}

	result := ""
	for _, v := range values {
		idx := max(int(float64(v)/float64(maxVal)*float64(len(sparkChars)-1)), 0)
		if idx >= len(sparkChars) {
			idx = len(sparkChars) - 1
		}
		result += sparkChars[idx]
	}
	return result
}

// RenderMiniBar creates a compact horizontal bar
func RenderMiniBar(value, max int, width int) string {
	if max <= 0 {
		max = 1
	}
	if width <= 0 {
		width = 10
	}

	filled := int(float64(value) / float64(max) * float64(width))
	if filled > width {
		filled = width
	}

	result := ""
	for i := 0; i < width; i++ {
		if i < filled {
			result += "▰"
		} else {
			result += "▱"
		}
	}
	return result
}

// RenderAnimatedDots creates animated dots based on tick count
func RenderAnimatedDots(tickCount int) string {
	dots := []string{"", ".", "..", "..."}
	return dots[tickCount%len(dots)]
}

// RenderBoxedText creates text surrounded by a box
func RenderBoxedText(text string, color lipgloss.Color) string {
	width := len(text) + 2
	top := BoxTopLeft + repeat(BoxHorizontal, width) + BoxTopRight
	bottom := BoxBottomLeft + repeat(BoxHorizontal, width) + BoxBottomRight
	middle := BoxVertical + " " + text + " " + BoxVertical

	style := lipgloss.NewStyle().Foreground(color)
	return style.Render(top + "\n" + middle + "\n" + bottom)
}

// repeat repeats a string n times
func repeat(s string, n int) string {
	result := ""
	for range n {
		result += s
	}
	return result
}

// RenderStatusIndicator creates an animated status indicator
func RenderStatusIndicator(online bool, tickCount int) string {
	if online {
		// Pulsing green dot when online
		pulse := GetPulse(tickCount)
		return lipgloss.NewStyle().Foreground(Green).Bold(true).Render(pulse)
	}
	// Static empty circle when offline
	return lipgloss.NewStyle().Foreground(Red).Render(BulletEmpty)
}

// RenderActivityIndicator shows activity with animated spinner
func RenderActivityIndicator(active bool, tickCount int, label string) string {
	if active {
		spinner := GetSpinnerFrame(tickCount)
		return lipgloss.NewStyle().Foreground(Cyan).Render(spinner) + " " + label
	}
	return MutedStyle.Render("○ " + label)
}
