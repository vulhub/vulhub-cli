package ui

import "github.com/charmbracelet/lipgloss"

// Color palette
const (
	ColorSuccess = lipgloss.Color("42")  // Green
	ColorWarning = lipgloss.Color("214") // Orange
	ColorError   = lipgloss.Color("167") // Red
	ColorInfo    = lipgloss.Color("39")  // Blue
	ColorMuted   = lipgloss.Color("245") // Gray
	ColorAccent  = lipgloss.Color("205") // Pink
	ColorCode    = lipgloss.Color("212") // Light pink
	ColorCodeBg  = lipgloss.Color("236") // Dark gray
	ColorBorder  = lipgloss.Color("240") // Border gray
)

// Text styles
var (
	// TitleStyle is for main titles/headers
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAccent).
			MarginBottom(1)

	// SuccessStyle is for success messages
	SuccessStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSuccess)

	// WarningStyle is for warning messages
	WarningStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorWarning)

	// ErrorStyle is for error messages
	ErrorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorError)

	// InfoStyle is for informational text
	InfoStyle = lipgloss.NewStyle().
			Foreground(ColorInfo)

	// MutedStyle is for secondary/muted text
	MutedStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	// CodeStyle is for inline code or important values
	CodeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorCode).
			Background(ColorCodeBg).
			Padding(0, 2)

	// URLStyle is for URLs
	URLStyle = lipgloss.NewStyle().
			Foreground(ColorInfo).
			Underline(true)

	// TableHeaderStyle is for table headers
	TableHeaderStyle = lipgloss.NewStyle().
				Foreground(ColorMuted).
				Bold(true)

	// TableBorderStyle is for table separators
	TableBorderStyle = lipgloss.NewStyle().
				Foreground(ColorBorder)

	// StatusRunningStyle is for running statuses
	StatusRunningStyle = lipgloss.NewStyle().
				Foreground(ColorSuccess).
				Bold(true)

	// StatusStoppedStyle is for stopped statuses
	StatusStoppedStyle = lipgloss.NewStyle().
				Foreground(ColorError).
				Bold(true)

	// PathStyle is for environment paths
	PathStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255"))

	// PortStyle is for port lists
	PortStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)
)

// Box styles
var (
	// BoxStyle is the base box style with rounded border
	BoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Padding(1, 2).
		MarginTop(1).
		MarginBottom(1)
)

// BoxWithColor returns a box style with a specific border color
func BoxWithColor(color lipgloss.Color) lipgloss.Style {
	return BoxStyle.BorderForeground(color)
}

// SuccessBox returns a box with success-colored border
func SuccessBox() lipgloss.Style {
	return BoxWithColor(ColorSuccess)
}

// WarningBox returns a box with warning-colored border
func WarningBox() lipgloss.Style {
	return BoxWithColor(ColorWarning)
}

// ErrorBox returns a box with error-colored border
func ErrorBox() lipgloss.Style {
	return BoxWithColor(ColorError)
}

// InfoBox returns a box with info-colored border
func InfoBox() lipgloss.Style {
	return BoxWithColor(ColorInfo)
}
