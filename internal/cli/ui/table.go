package ui

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/vulhub/vulhub-cli/pkg/types"
)

// Table provides table formatting functionality
type Table struct {
	writer io.Writer
}

// NewTable creates a new Table that writes to stdout
func NewTable() *Table {
	return &Table{
		writer: os.Stdout,
	}
}

// NewTableWithWriter creates a new Table that writes to the specified writer
func NewTableWithWriter(w io.Writer) *Table {
	return &Table{
		writer: w,
	}
}

// PrintEnvironments prints a list of environments in table format
func (t *Table) PrintEnvironments(envs []types.Environment) {
	fmt.Fprint(t.writer, t.FormatEnvironments(envs))
}

// FormatEnvironments returns a formatted list of environments as a string
func (t *Table) FormatEnvironments(envs []types.Environment) string {
	var buf bytes.Buffer

	if len(envs) == 0 {
		fmt.Fprintln(&buf, "No environments found.")
		return buf.String()
	}

	// Calculate column widths
	pathWidth := 30
	cveWidth := 20
	appWidth := 15

	for _, env := range envs {
		if len(env.Path) > pathWidth {
			pathWidth = len(env.Path)
		}
		if len(env.CVE) > 0 && len(env.CVE[0]) > cveWidth {
			cveWidth = len(env.CVE[0])
		}
		if len(env.App) > appWidth {
			appWidth = len(env.App)
		}
	}

	// Limit widths
	if pathWidth > 50 {
		pathWidth = 50
	}
	if cveWidth > 25 {
		cveWidth = 25
	}
	if appWidth > 20 {
		appWidth = 20
	}

	// Print header
	format := fmt.Sprintf("%%-%ds  %%-%ds  %%-%ds\n", pathWidth, cveWidth, appWidth)
	fmt.Fprintf(&buf, format, "PATH", "CVE", "APP")
	fmt.Fprintln(&buf, strings.Repeat("-", pathWidth+cveWidth+appWidth+4))

	// Print rows
	for _, env := range envs {
		path := truncate(env.Path, pathWidth)
		cve := "-"
		if len(env.CVE) > 0 {
			cve = truncate(env.CVE[0], cveWidth)
		}
		app := truncate(env.App, appWidth)

		fmt.Fprintf(&buf, format, path, cve, app)
	}

	fmt.Fprintf(&buf, "\nTotal: %d environments\n", len(envs))

	return buf.String()
}

// PrintEnvironmentStatuses prints a list of environment statuses
func (t *Table) PrintEnvironmentStatuses(statuses []types.EnvironmentStatus) {
	if len(statuses) == 0 {
		fmt.Fprintln(t.writer, "No downloaded environments.")
		return
	}

	const (
		pathWidth   = 35
		statusWidth = 14
		portsWidth  = 40
		gapWidth    = 2
	)

	header := lipgloss.JoinHorizontal(lipgloss.Top,
		TableHeaderStyle.Width(pathWidth).PaddingRight(gapWidth).Render("PATH"),
		TableHeaderStyle.Width(statusWidth).PaddingRight(gapWidth).Render("STATUS"),
		TableHeaderStyle.Width(portsWidth).Render("PORTS"),
	)
	separatorWidth := pathWidth + statusWidth + portsWidth + (gapWidth * 2)
	separator := TableBorderStyle.Render(strings.Repeat("─", separatorWidth))

	pathStyle := PathStyle.Width(pathWidth).PaddingRight(gapWidth)
	statusStyle := lipgloss.NewStyle().Width(statusWidth).PaddingRight(gapWidth)
	portsStyle := PortStyle.Width(portsWidth)

	fmt.Fprintln(t.writer, header)
	fmt.Fprintln(t.writer, separator)

	for _, status := range statuses {
		statusStr := StatusStoppedStyle.Render("○ stopped")
		if status.Running {
			statusStr = StatusRunningStyle.Render("● running")
		}

		// Format ports
		ports := formatPorts(status.Containers)

		row := lipgloss.JoinHorizontal(lipgloss.Top,
			pathStyle.Render(truncate(status.Environment.Path, pathWidth)),
			statusStyle.Render(statusStr),
			portsStyle.Render(truncate(ports, portsWidth)),
		)
		fmt.Fprintln(t.writer, row)
	}

	totalLine := MutedStyle.Render(fmt.Sprintf("Total: %d environments", len(statuses)))
	fmt.Fprintf(t.writer, "\n%s\n", totalLine)
}

// PrintContainerStatuses prints detailed container statuses
func (t *Table) PrintContainerStatuses(containers []types.ContainerStatus) {
	if len(containers) == 0 {
		fmt.Fprintln(t.writer, "No containers found.")
		return
	}

	const (
		nameWidth  = 25
		imageWidth = 30
		stateWidth = 10
		portsWidth = 30
		gapWidth   = 2
	)

	header := lipgloss.JoinHorizontal(lipgloss.Top,
		TableHeaderStyle.Width(nameWidth).PaddingRight(gapWidth).Render("NAME"),
		TableHeaderStyle.Width(imageWidth).PaddingRight(gapWidth).Render("IMAGE"),
		TableHeaderStyle.Width(stateWidth).PaddingRight(gapWidth).Render("STATE"),
		TableHeaderStyle.Width(portsWidth).Render("PORTS"),
	)
	separatorWidth := nameWidth + imageWidth + stateWidth + portsWidth + (gapWidth * 3)
	separator := TableBorderStyle.Render(strings.Repeat("─", separatorWidth))

	nameStyle := PathStyle.Width(nameWidth).PaddingRight(gapWidth)
	imageStyle := lipgloss.NewStyle().Width(imageWidth).PaddingRight(gapWidth)
	stateStyle := lipgloss.NewStyle().Width(stateWidth).PaddingRight(gapWidth)
	portsStyle := PortStyle.Width(portsWidth)

	fmt.Fprintln(t.writer, header)
	fmt.Fprintln(t.writer, separator)

	for _, c := range containers {
		ports := formatContainerPorts(c.Ports)
		stateText := c.State
		switch strings.ToLower(c.State) {
		case "running":
			stateText = StatusRunningStyle.Render(c.State)
		case "exited", "stopped":
			stateText = StatusStoppedStyle.Render(c.State)
		}

		row := lipgloss.JoinHorizontal(lipgloss.Top,
			nameStyle.Render(truncate(c.Name, nameWidth)),
			imageStyle.Render(truncate(c.Image, imageWidth)),
			stateStyle.Render(stateText),
			portsStyle.Render(truncate(ports, portsWidth)),
		)
		fmt.Fprintln(t.writer, row)
	}
}

// PrintEnvironmentInfo prints detailed environment information
func (t *Table) PrintEnvironmentInfo(info *types.EnvironmentInfo) {
	fmt.Fprint(t.writer, t.FormatEnvironmentInfo(info))
}

// FormatEnvironmentInfo returns formatted environment information as a string
func (t *Table) FormatEnvironmentInfo(info *types.EnvironmentInfo) string {
	var buf bytes.Buffer
	env := info.Environment

	fmt.Fprintln(&buf, strings.Repeat("=", 60))
	fmt.Fprintf(&buf, "Environment: %s\n", env.Path)
	fmt.Fprintln(&buf, strings.Repeat("=", 60))

	fmt.Fprintf(&buf, "Name:        %s\n", env.Name)
	if len(env.CVE) > 0 {
		fmt.Fprintf(&buf, "CVE:         %s\n", strings.Join(env.CVE, ", "))
	}
	fmt.Fprintf(&buf, "Application: %s\n", env.App)

	if len(env.Tags) > 0 {
		fmt.Fprintf(&buf, "Tags:        %s\n", strings.Join(env.Tags, ", "))
	}

	fmt.Fprintf(&buf, "Downloaded:  %v\n", info.Downloaded)
	if info.Downloaded {
		fmt.Fprintf(&buf, "Local Path:  %s\n", info.LocalPath)
	}

	if info.Readme != "" {
		fmt.Fprintln(&buf, strings.Repeat("-", 60))
		fmt.Fprintln(&buf, "README:")
		fmt.Fprintln(&buf, strings.Repeat("-", 60))
		fmt.Fprintln(&buf, info.Readme)
	}

	return buf.String()
}

// PrintSuccess prints a success message
func (t *Table) PrintSuccess(message string) {
	fmt.Fprintf(t.writer, "[OK] %s\n", message)
}

// PrintError prints an error message
func (t *Table) PrintError(message string) {
	fmt.Fprintf(t.writer, "[ERROR] %s\n", message)
}

// PrintWarning prints a warning message
func (t *Table) PrintWarning(message string) {
	fmt.Fprintf(t.writer, "[WARN] %s\n", message)
}

// PrintInfo prints an info message
func (t *Table) PrintInfo(message string) {
	fmt.Fprintf(t.writer, "[INFO] %s\n", message)
}

// Helper functions

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

func formatPorts(containers []types.ContainerStatus) string {
	var ports []string
	for _, c := range containers {
		for _, p := range c.Ports {
			if p.HostPort != "" && p.HostPort != "0" {
				ports = append(ports, fmt.Sprintf("%s:%s->%s/%s",
					p.HostIP, p.HostPort, p.ContainerPort, p.Protocol))
			}
		}
	}
	if len(ports) == 0 {
		return "-"
	}
	return strings.Join(ports, ", ")
}

func formatContainerPorts(ports []types.PortMapping) string {
	var result []string
	for _, p := range ports {
		if p.HostPort != "" && p.HostPort != "0" {
			result = append(result, fmt.Sprintf("%s:%s->%s/%s",
				p.HostIP, p.HostPort, p.ContainerPort, p.Protocol))
		}
	}
	if len(result) == 0 {
		return "-"
	}
	return strings.Join(result, ", ")
}
