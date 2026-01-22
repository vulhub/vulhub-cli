package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205"))

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

// pagerModel is the bubbletea model for the pager
type pagerModel struct {
	viewport    viewport.Model
	title       string
	ready       bool
	quitting    bool
}

func newPagerModel(title, content string) pagerModel {
	return pagerModel{
		title:   title,
		viewport: viewport.Model{},
	}
}

func (m pagerModel) Init() tea.Cmd {
	return nil
}

func (m pagerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		headerHeight := 3
		footerHeight := 2

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-headerHeight-footerHeight)
			m.viewport.YPosition = headerHeight
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - headerHeight - footerHeight
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m pagerModel) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	header := titleStyle.Render(m.title)
	scrollPercent := fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100)
	info := infoStyle.Render(fmt.Sprintf("Lines %d-%d of %d (%s)",
		m.viewport.YOffset+1,
		min(m.viewport.YOffset+m.viewport.Height, m.viewport.TotalLineCount()),
		m.viewport.TotalLineCount(),
		scrollPercent,
	))

	headerLine := fmt.Sprintf("%s  %s", header, info)
	footer := helpStyle.Render("↑/↓: scroll • PgUp/PgDn: page • q: quit")

	return fmt.Sprintf("%s\n%s\n%s\n%s",
		headerLine,
		strings.Repeat("─", m.viewport.Width),
		m.viewport.View(),
		footer,
	)
}

// Pager displays content in a scrollable viewport
type Pager struct{}

// NewPager creates a new Pager
func NewPager() *Pager {
	return &Pager{}
}

// Display shows the content in a pager with the given title
func (p *Pager) Display(title, content string) error {
	// Count lines to determine if pager is needed
	lines := strings.Count(content, "\n") + 1

	// If content is short, just print it directly
	if lines <= 20 {
		if title != "" {
			fmt.Println(titleStyle.Render(title))
			fmt.Println(strings.Repeat("─", 60))
		}
		fmt.Print(content)
		return nil
	}

	m := newPagerModel(title, content)

	prog := tea.NewProgram(
		m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// Set content after program initialization to get correct dimensions
	go func() {
		prog.Send(tea.WindowSizeMsg{})
	}()

	model, err := prog.Run()
	if err != nil {
		return fmt.Errorf("pager error: %w", err)
	}

	// Set content
	finalModel := model.(pagerModel)
	_ = finalModel

	return nil
}

// DisplayWithContent shows content in a pager, setting content on init
func (p *Pager) DisplayWithContent(title, content string) error {
	// Count lines to determine if pager is needed
	lines := strings.Count(content, "\n") + 1

	// If content is short, just print it directly
	if lines <= 20 {
		if title != "" {
			fmt.Println(titleStyle.Render(title))
			fmt.Println(strings.Repeat("─", 60))
		}
		fmt.Print(content)
		return nil
	}

	m := pagerModelWithContent{
		title:   title,
		content: content,
	}

	prog := tea.NewProgram(
		m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	_, err := prog.Run()
	if err != nil {
		return fmt.Errorf("pager error: %w", err)
	}

	return nil
}

// pagerModelWithContent is a pager model that stores content
type pagerModelWithContent struct {
	viewport viewport.Model
	title    string
	content  string
	ready    bool
	quitting bool
}

func (m pagerModelWithContent) Init() tea.Cmd {
	return nil
}

func (m pagerModelWithContent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		headerHeight := 3
		footerHeight := 2

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-headerHeight-footerHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.SetContent(m.content)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - headerHeight - footerHeight
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m pagerModelWithContent) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	header := titleStyle.Render(m.title)
	scrollPercent := fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100)
	info := infoStyle.Render(fmt.Sprintf("Lines %d-%d of %d (%s)",
		m.viewport.YOffset+1,
		min(m.viewport.YOffset+m.viewport.Height, m.viewport.TotalLineCount()),
		m.viewport.TotalLineCount(),
		scrollPercent,
	))

	headerLine := fmt.Sprintf("%s  %s", header, info)
	footer := helpStyle.Render("↑/↓: scroll • PgUp/PgDn: page • q: quit")

	return fmt.Sprintf("%s\n%s\n%s\n%s",
		headerLine,
		strings.Repeat("─", m.viewport.Width),
		m.viewport.View(),
		footer,
	)
}
