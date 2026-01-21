package ui

import (
	"fmt"
	"io"
	"os"
	"strings"

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
	if len(envs) == 0 {
		fmt.Fprintln(t.writer, "No environments found.")
		return
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
	fmt.Fprintf(t.writer, format, "PATH", "CVE", "APP")
	fmt.Fprintln(t.writer, strings.Repeat("-", pathWidth+cveWidth+appWidth+4))

	// Print rows
	for _, env := range envs {
		path := truncate(env.Path, pathWidth)
		cve := "-"
		if len(env.CVE) > 0 {
			cve = truncate(env.CVE[0], cveWidth)
		}
		app := truncate(env.App, appWidth)

		fmt.Fprintf(t.writer, format, path, cve, app)
	}

	fmt.Fprintf(t.writer, "\nTotal: %d environments\n", len(envs))
}

// PrintEnvironmentStatuses prints a list of environment statuses
func (t *Table) PrintEnvironmentStatuses(statuses []types.EnvironmentStatus) {
	if len(statuses) == 0 {
		fmt.Fprintln(t.writer, "No running environments.")
		return
	}

	// Print header
	fmt.Fprintf(t.writer, "%-35s  %-10s  %-30s\n", "PATH", "STATUS", "PORTS")
	fmt.Fprintln(t.writer, strings.Repeat("-", 80))

	for _, status := range statuses {
		statusStr := "stopped"
		if status.Running {
			statusStr = "running"
		}

		// Format ports
		ports := formatPorts(status.Containers)

		fmt.Fprintf(t.writer, "%-35s  %-10s  %-30s\n",
			truncate(status.Environment.Path, 35),
			statusStr,
			truncate(ports, 30),
		)
	}

	fmt.Fprintf(t.writer, "\nTotal: %d environments\n", len(statuses))
}

// PrintContainerStatuses prints detailed container statuses
func (t *Table) PrintContainerStatuses(containers []types.ContainerStatus) {
	if len(containers) == 0 {
		fmt.Fprintln(t.writer, "No containers found.")
		return
	}

	// Print header
	fmt.Fprintf(t.writer, "%-25s  %-30s  %-10s  %-30s\n", "NAME", "IMAGE", "STATE", "PORTS")
	fmt.Fprintln(t.writer, strings.Repeat("-", 100))

	for _, c := range containers {
		ports := formatContainerPorts(c.Ports)

		fmt.Fprintf(t.writer, "%-25s  %-30s  %-10s  %-30s\n",
			truncate(c.Name, 25),
			truncate(c.Image, 30),
			c.State,
			truncate(ports, 30),
		)
	}
}

// PrintEnvironmentInfo prints detailed environment information
func (t *Table) PrintEnvironmentInfo(info *types.EnvironmentInfo) {
	env := info.Environment

	fmt.Fprintln(t.writer, strings.Repeat("=", 60))
	fmt.Fprintf(t.writer, "Environment: %s\n", env.Path)
	fmt.Fprintln(t.writer, strings.Repeat("=", 60))

	fmt.Fprintf(t.writer, "Name:        %s\n", env.Name)
	if len(env.CVE) > 0 {
		fmt.Fprintf(t.writer, "CVE:         %s\n", strings.Join(env.CVE, ", "))
	}
	fmt.Fprintf(t.writer, "Application: %s\n", env.App)

	if len(env.Tags) > 0 {
		fmt.Fprintf(t.writer, "Tags:        %s\n", strings.Join(env.Tags, ", "))
	}

	fmt.Fprintf(t.writer, "Downloaded:  %v\n", info.Downloaded)
	if info.Downloaded {
		fmt.Fprintf(t.writer, "Local Path:  %s\n", info.LocalPath)
	}

	if info.Readme != "" {
		fmt.Fprintln(t.writer, strings.Repeat("-", 60))
		fmt.Fprintln(t.writer, "README:")
		fmt.Fprintln(t.writer, strings.Repeat("-", 60))
		fmt.Fprintln(t.writer, info.Readme)
	}
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
