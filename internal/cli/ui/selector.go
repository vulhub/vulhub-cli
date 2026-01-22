package ui

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/vulhub/vulhub-cli/pkg/types"
)

// Selector provides interactive selection functionality
type Selector struct{}

// NewSelector creates a new Selector
func NewSelector() *Selector {
	return &Selector{}
}

// SelectEnvironment prompts the user to select an environment from a list
func (s *Selector) SelectEnvironment(envs []types.Environment, prompt string) (*types.Environment, error) {
	if len(envs) == 0 {
		return nil, fmt.Errorf("no environments to select from")
	}

	if len(envs) == 1 {
		return &envs[0], nil
	}

	// Build options for huh.Select
	options := make([]huh.Option[int], len(envs))
	for i, env := range envs {
		cve := "-"
		if len(env.CVE) > 0 {
			cve = env.CVE[0]
		}
		label := fmt.Sprintf("%-30s %s", env.Path, cve)
		options[i] = huh.NewOption(label, i)
	}

	var selected int
	err := huh.NewSelect[int]().
		Title(prompt).
		Options(options...).
		Value(&selected).
		Run()

	if err != nil {
		if err == huh.ErrUserAborted {
			return nil, fmt.Errorf("selection cancelled")
		}
		return nil, fmt.Errorf("selection failed: %w", err)
	}

	return &envs[selected], nil
}

// Confirm prompts the user for confirmation
func (s *Selector) Confirm(prompt string, defaultYes bool) (bool, error) {
	var confirmed bool

	err := huh.NewConfirm().
		Title(prompt).
		Affirmative("Yes").
		Negative("No").
		Value(&confirmed).
		Run()

	if err != nil {
		if err == huh.ErrUserAborted {
			return false, nil
		}
		return false, fmt.Errorf("confirmation failed: %w", err)
	}

	return confirmed, nil
}

// PromptString prompts the user for a string input
func (s *Selector) PromptString(prompt string, defaultValue string) (string, error) {
	var input string

	inputField := huh.NewInput().
		Title(prompt).
		Value(&input)

	if defaultValue != "" {
		inputField.Placeholder(defaultValue)
	}

	err := inputField.Run()
	if err != nil {
		if err == huh.ErrUserAborted {
			return "", fmt.Errorf("input cancelled")
		}
		return "", fmt.Errorf("input failed: %w", err)
	}

	if input == "" {
		return defaultValue, nil
	}

	return input, nil
}
