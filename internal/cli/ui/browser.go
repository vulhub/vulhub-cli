package ui

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vulhub/vulhub-cli/pkg/types"
)

// OpenBrowser opens the specified URL in the default browser
func OpenBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default: // linux, freebsd, etc.
		cmd = exec.Command("xdg-open", url)
	}

	return cmd.Start()
}

var (
	browserTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("205")).
				MarginBottom(1)

	searchBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	browserHelpStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				MarginTop(1)

	counterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))
)

// browserModel is the bubbletea model for the environment browser
type browserModel struct {
	title        string
	textInput    textinput.Model
	environments []types.Environment
	filtered     []types.Environment
	cursor       int
	offset       int
	width        int
	height       int
	ready        bool
	quitting     bool
	selected     *types.Environment
}

// BrowseResult contains the result of browsing
type BrowseResult struct {
	Selected *types.Environment
	Quit     bool
}

// BrowseOptions contains options for the browser
type BrowseOptions struct {
	Title         string
	InitialSearch string
}

func newBrowserModel(envs []types.Environment, opts BrowseOptions) browserModel {
	ti := textinput.New()
	ti.Placeholder = "Type to search..."
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 40

	if opts.InitialSearch != "" {
		ti.SetValue(opts.InitialSearch)
	}

	title := opts.Title
	if title == "" {
		title = "Available Environments"
	}

	m := browserModel{
		title:        title,
		textInput:    ti,
		environments: envs,
		filtered:     envs,
		cursor:       0,
		offset:       0,
	}

	// Apply initial filter if search is provided
	if opts.InitialSearch != "" {
		m.filterEnvironments()
	}

	return m
}

func (m browserModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m browserModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			if len(m.filtered) > 0 && m.cursor < len(m.filtered) {
				m.selected = &m.filtered[m.cursor]
			}
			return m, tea.Quit

		case "up", "ctrl+p":
			if m.cursor > 0 {
				m.cursor--
				if m.cursor < m.offset {
					m.offset = m.cursor
				}
			}
			return m, nil

		case "down", "ctrl+n":
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
				visibleHeight := m.visibleHeight()
				if m.cursor >= m.offset+visibleHeight {
					m.offset = m.cursor - visibleHeight + 1
				}
			}
			return m, nil

		case "pgup":
			visibleHeight := m.visibleHeight()
			m.cursor -= visibleHeight
			if m.cursor < 0 {
				m.cursor = 0
			}
			if m.cursor < m.offset {
				m.offset = m.cursor
			}
			return m, nil

		case "pgdown":
			visibleHeight := m.visibleHeight()
			m.cursor += visibleHeight
			if m.cursor >= len(m.filtered) {
				m.cursor = len(m.filtered) - 1
			}
			if m.cursor < 0 {
				m.cursor = 0
			}
			if m.cursor >= m.offset+visibleHeight {
				m.offset = m.cursor - visibleHeight + 1
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil
	}

	// Handle text input
	prevValue := m.textInput.Value()
	m.textInput, cmd = m.textInput.Update(msg)

	// Filter environments if search text changed
	if m.textInput.Value() != prevValue {
		m.filterEnvironments()
	}

	return m, cmd
}

func (m *browserModel) filterEnvironments() {
	query := strings.ToLower(m.textInput.Value())
	if query == "" {
		m.filtered = m.environments
	} else {
		m.filtered = nil
		for _, env := range m.environments {
			// Search in path, name, app, CVE, and tags
			searchText := strings.ToLower(fmt.Sprintf("%s %s %s %s %s",
				env.Path,
				env.Name,
				env.App,
				strings.Join(env.CVE, " "),
				strings.Join(env.Tags, " "),
			))
			if strings.Contains(searchText, query) {
				m.filtered = append(m.filtered, env)
			}
		}
	}
	m.cursor = 0
	m.offset = 0
}

func (m browserModel) visibleHeight() int {
	// Account for title, search box, header, footer, and padding
	return m.height - 10
}

func (m browserModel) View() string {
	if !m.ready {
		return "\n  Loading..."
	}

	var b strings.Builder

	// Title
	b.WriteString(browserTitleStyle.Render(m.title))
	b.WriteString("\n")

	// Search box
	b.WriteString(searchBoxStyle.Render(m.textInput.View()))
	b.WriteString("\n\n")

	// Counter
	counter := counterStyle.Render(fmt.Sprintf("Showing %d of %d environments", len(m.filtered), len(m.environments)))
	b.WriteString(counter)
	b.WriteString("\n\n")

	// Calculate dynamic column widths based on terminal width
	// Minimum widths: PATH=25, CVE=15, APP=8, TAGS=10
	// Preferred: PATH=30, CVE=45, APP=10, TAGS=remaining
	availableWidth := m.width - 12 // account for padding and separators

	pathWidth := 30
	cveWidth := 45
	appWidth := 10
	tagsWidth := availableWidth - pathWidth - cveWidth - appWidth

	// Adjust if terminal is narrow
	if tagsWidth < 15 {
		tagsWidth = 15
		cveWidth = availableWidth - pathWidth - appWidth - tagsWidth
		if cveWidth < 30 {
			cveWidth = 30
			pathWidth = availableWidth - cveWidth - appWidth - tagsWidth
			if pathWidth < 25 {
				pathWidth = 25
			}
		}
	}

	// Table header
	headerFmt := fmt.Sprintf("  %%-%ds  %%-%ds  %%-%ds  %%s", pathWidth, cveWidth, appWidth)
	header := fmt.Sprintf(headerFmt, "PATH", "CVE", "APP", "TAGS")
	b.WriteString(dimStyle.Render(header))
	b.WriteString("\n")
	b.WriteString(dimStyle.Render("  " + strings.Repeat("─", m.width-4)))
	b.WriteString("\n")

	// Environment list
	visibleHeight := m.visibleHeight()
	if visibleHeight < 1 {
		visibleHeight = 1
	}

	end := m.offset + visibleHeight
	if end > len(m.filtered) {
		end = len(m.filtered)
	}

	if len(m.filtered) == 0 {
		b.WriteString(dimStyle.Render("  No environments match your search."))
		b.WriteString("\n")
	} else {
		lineFmt := fmt.Sprintf("  %%-%ds  %%-%ds  %%-%ds  %%s", pathWidth, cveWidth, appWidth)
		for i := m.offset; i < end; i++ {
			env := m.filtered[i]

			// Join all CVEs
			cve := "-"
			if len(env.CVE) > 0 {
				cve = strings.Join(env.CVE, ", ")
			}

			// Join all tags
			tags := "-"
			if len(env.Tags) > 0 {
				tags = strings.Join(env.Tags, ", ")
			}

			line := fmt.Sprintf(lineFmt,
				truncate(env.Path, pathWidth),
				truncate(cve, cveWidth),
				truncate(env.App, appWidth),
				truncate(tags, tagsWidth),
			)

			if i == m.cursor {
				b.WriteString(selectedStyle.Render("> " + line[2:]))
			} else {
				b.WriteString(normalStyle.Render(line))
			}
			b.WriteString("\n")
		}
	}

	// Scroll indicator
	if len(m.filtered) > visibleHeight {
		scrollInfo := dimStyle.Render(fmt.Sprintf("\n  [%d-%d of %d]",
			m.offset+1,
			min(m.offset+visibleHeight, len(m.filtered)),
			len(m.filtered),
		))
		b.WriteString(scrollInfo)
	}

	// Help
	help := browserHelpStyle.Render("↑/↓: navigate • Enter: select • PgUp/PgDn: page • Esc: quit")
	b.WriteString("\n")
	b.WriteString(help)

	return b.String()
}

// EnvironmentBrowser provides an interactive environment browser
type EnvironmentBrowser struct{}

// NewEnvironmentBrowser creates a new EnvironmentBrowser
func NewEnvironmentBrowser() *EnvironmentBrowser {
	return &EnvironmentBrowser{}
}

// Browse displays an interactive browser for selecting an environment
func (eb *EnvironmentBrowser) Browse(envs []types.Environment) (*BrowseResult, error) {
	return eb.BrowseWithOptions(envs, BrowseOptions{})
}

// BrowseWithOptions displays an interactive browser with custom options
func (eb *EnvironmentBrowser) BrowseWithOptions(envs []types.Environment, opts BrowseOptions) (*BrowseResult, error) {
	m := newBrowserModel(envs, opts)

	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
	)

	finalModel, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("browser error: %w", err)
	}

	result := finalModel.(browserModel)
	return &BrowseResult{
		Selected: result.selected,
		Quit:     result.quitting,
	}, nil
}
