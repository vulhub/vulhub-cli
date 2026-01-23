package ui

import (
	"fmt"

	"github.com/charmbracelet/huh"
)

// StatusCard displays a status card with icon and message
type StatusCard struct {
	Icon    string
	Title   string
	Details []string
	Style   StatusStyle
}

type StatusStyle int

const (
	StatusSuccess StatusStyle = iota
	StatusWarning
	StatusError
	StatusInfo
)

// Render returns the rendered status card
func (s *StatusCard) Render() string {
	var iconStyle, titleStyle = SuccessStyle, SuccessStyle
	var boxStyle = SuccessBox()

	switch s.Style {
	case StatusWarning:
		iconStyle, titleStyle = WarningStyle, WarningStyle
		boxStyle = WarningBox()
	case StatusError:
		iconStyle, titleStyle = ErrorStyle, ErrorStyle
		boxStyle = ErrorBox()
	case StatusInfo:
		iconStyle, titleStyle = InfoStyle, InfoStyle
		boxStyle = InfoBox()
	}

	content := fmt.Sprintf("%s %s", iconStyle.Render(s.Icon), titleStyle.Render(s.Title))

	for _, detail := range s.Details {
		content += "\n\n" + MutedStyle.Render(detail)
	}

	return boxStyle.Render(content)
}

// Print prints the status card to stdout
func (s *StatusCard) Print() {
	fmt.Println(s.Render())
}

// AuthCodeCard displays the authorization code for OAuth device flow
type AuthCodeCard struct {
	Instruction string
	URL         string
	Code        string
}

// Render returns the rendered auth code card
func (a *AuthCodeCard) Render() string {
	content := fmt.Sprintf(
		"%s\n\n%s  %s\n\n%s  %s",
		MutedStyle.Render(a.Instruction),
		MutedStyle.Render("URL:"),
		URLStyle.Render(a.URL),
		MutedStyle.Render("Code:"),
		CodeStyle.Render(a.Code),
	)
	return InfoBox().Render(content)
}

// Print prints the auth code card to stdout
func (a *AuthCodeCard) Print() {
	fmt.Println(a.Render())
}

// MessageBox is a simple message box
type MessageBox struct {
	Title   string
	Message string
	Style   StatusStyle
}

// Render returns the rendered message box
func (m *MessageBox) Render() string {
	var titleStyle = SuccessStyle
	var boxStyle = SuccessBox()

	switch m.Style {
	case StatusWarning:
		titleStyle = WarningStyle
		boxStyle = WarningBox()
	case StatusError:
		titleStyle = ErrorStyle
		boxStyle = ErrorBox()
	case StatusInfo:
		titleStyle = InfoStyle
		boxStyle = InfoBox()
	}

	content := titleStyle.Render(m.Title)
	if m.Message != "" {
		content += "\n\n" + MutedStyle.Render(m.Message)
	}

	return boxStyle.Render(content)
}

// Print prints the message box to stdout
func (m *MessageBox) Print() {
	fmt.Println(m.Render())
}

// ConfirmPrompt shows a confirmation prompt
func ConfirmPrompt(title, description, yes, no string) (bool, error) {
	var confirmed bool

	confirm := huh.NewConfirm().
		Title(title).
		Affirmative(yes).
		Negative(no).
		Value(&confirmed)

	if description != "" {
		confirm.Description(description)
	}

	err := confirm.Run()
	if err != nil {
		if err == huh.ErrUserAborted {
			return false, nil
		}
		return false, err
	}

	return confirmed, nil
}

// PrintTitle prints a styled title
func PrintTitle(title string) {
	fmt.Println()
	fmt.Println(TitleStyle.Render(title))
}

// PrintMuted prints muted/secondary text
func PrintMuted(text string) {
	fmt.Println(MutedStyle.Render(text))
}

// PrintSuccess prints a success message with checkmark
func PrintSuccess(message string) {
	fmt.Println(SuccessStyle.Render("  ✓ " + message))
}

// PrintWarning prints a warning message
func PrintWarning(message string) {
	fmt.Println(WarningStyle.Render("  ✗ " + message))
}

// PrintInfo prints an info message
func PrintInfo(message string) {
	fmt.Println(MutedStyle.Render("  " + message))
}
