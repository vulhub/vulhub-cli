package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/vulhub/vulhub-cli/pkg/types"
)

// Selector provides interactive selection functionality
type Selector struct {
	reader *bufio.Reader
}

// NewSelector creates a new Selector
func NewSelector() *Selector {
	return &Selector{
		reader: bufio.NewReader(os.Stdin),
	}
}

// SelectEnvironment prompts the user to select an environment from a list
func (s *Selector) SelectEnvironment(envs []types.Environment, prompt string) (*types.Environment, error) {
	if len(envs) == 0 {
		return nil, fmt.Errorf("no environments to select from")
	}

	if len(envs) == 1 {
		return &envs[0], nil
	}

	// Display the list
	fmt.Println()
	fmt.Println(prompt)
	fmt.Println(strings.Repeat("-", 60))

	for i, env := range envs {
		cve := "-"
		if len(env.CVE) > 0 {
			cve = env.CVE[0]
		}
		fmt.Printf("  [%d] %-30s %s\n", i+1, env.Path, cve)
	}

	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("Enter number (1-%d) or 'q' to quit: ", len(envs))

	// Read user input
	input, err := s.reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read input: %w", err)
	}

	input = strings.TrimSpace(input)

	// Check for quit
	if strings.ToLower(input) == "q" {
		return nil, fmt.Errorf("selection cancelled")
	}

	// Parse number
	num, err := strconv.Atoi(input)
	if err != nil {
		return nil, fmt.Errorf("invalid input: please enter a number")
	}

	// Validate range
	if num < 1 || num > len(envs) {
		return nil, fmt.Errorf("invalid selection: please enter a number between 1 and %d", len(envs))
	}

	return &envs[num-1], nil
}

// Confirm prompts the user for confirmation
func (s *Selector) Confirm(prompt string, defaultYes bool) (bool, error) {
	var hint string
	if defaultYes {
		hint = "[Y/n]"
	} else {
		hint = "[y/N]"
	}

	fmt.Printf("%s %s: ", prompt, hint)

	input, err := s.reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read input: %w", err)
	}

	input = strings.TrimSpace(strings.ToLower(input))

	if input == "" {
		return defaultYes, nil
	}

	switch input {
	case "y", "yes":
		return true, nil
	case "n", "no":
		return false, nil
	default:
		return false, fmt.Errorf("invalid input: please enter 'y' or 'n'")
	}
}

// PromptString prompts the user for a string input
func (s *Selector) PromptString(prompt string, defaultValue string) (string, error) {
	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", prompt, defaultValue)
	} else {
		fmt.Printf("%s: ", prompt)
	}

	input, err := s.reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}

	input = strings.TrimSpace(input)

	if input == "" {
		return defaultValue, nil
	}

	return input, nil
}
