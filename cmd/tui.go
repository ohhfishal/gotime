package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TuiCmd struct {
	// TODO: Have an instance of log and report cmd that this can then call on updates?
	Debug   bool          `help:"Enabled debug options"`
	Refresh time.Duration `default:"5s" help:"Duration between refreshing report."`
	Report  ReportCmd     `embed:""`
	model   Model         `kong:"-"`
}

func (cmd *TuiCmd) AfterApply() error {
	return cmd.Report.AfterApply()
}

type Model struct {
	state  tea.Model
	Report *ReportModel
	Log    *LogModel // TODO: Implement
	// Edit   *EditModel // TODO: Implement TUI to edit already existing entries
}

type LogModel struct {
	// TODO: Implement TUI to create a new entry
}

type ReportModel struct {
	content    string
	keys       keyMap
	help       help.Model
	inputStyle lipgloss.Style
	lastKey    string
	quitting   bool

	Refresh         func() (string, error)
	refreshDuration time.Duration
	lastTick        time.Time
}

func (cmd *TuiCmd) Run(log string) error {
	if cmd.Debug {
		f, err := tea.LogToFile("debug.log", "help")
		if err != nil {
			return fmt.Errorf("opening log file: %w", err)
		}
		defer f.Close() // nolint:errcheck
	}

	// Create the function to redraw the report
	// TODO: This doesn't *have* to be here?
	refresh := func() (string, error) {
		var stdout strings.Builder
		if err := cmd.Report.Run(&stdout, log); err != nil {
			return ``, fmt.Errorf(`getting report: %w`, err)
		}
		return stdout.String(), nil
	}

	content, err := refresh()
	if err != nil {
		return fmt.Errorf(`getting report: %w`, err)
	}

	cmd.model = Model{
		Report: &ReportModel{
			content:         content,
			keys:            keys,
			help:            help.New(),
			inputStyle:      lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")),
			Refresh:         refresh,
			refreshDuration: cmd.Refresh,
		},
	}
	cmd.model.state = cmd.model.Report

	if _, err := tea.NewProgram(&cmd.model).Run(); err != nil {
		return fmt.Errorf(`rending tui: %w`, err)
	}
	return nil
}

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	// TODO: Add more keybindings
	Up   key.Binding
	Down key.Binding
	Help key.Binding
	Quit key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},   // first column
		{k.Help, k.Quit}, // second column
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type TickMsg time.Time

func tickEvery(duration time.Duration) tea.Cmd {
	return tea.Every(duration, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// Generic Root Model
func (m *Model) Init() tea.Cmd {
	return m.Report.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.Report.Update(msg)
}

func (m *Model) View() string {
	return m.Report.View()
}

func (m *ReportModel) Init() tea.Cmd {
	return tickEvery(m.refreshDuration)
}

// Report Model

func (m *ReportModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: Update the data object (Probably after we are using a bubbles.Table)
	switch msg := msg.(type) {
	case TickMsg:
		m.lastTick = time.Now()
		content, err := m.Refresh()
		if err != nil {
			// TODO: Handle this more elegantly
			panic(fmt.Errorf(`refreshing: %w`, err))
		}
		m.content = content
		return m, tickEvery(m.refreshDuration)
	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can gracefully truncate
		// its view as needed.
		m.help.Width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Up):
			m.lastKey = "↑"
		case key.Matches(msg, m.keys.Down):
			m.lastKey = "↓"
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *ReportModel) View() string {
	if m.quitting {
		return ""
	}
	return fmt.Sprintf(
		"%s\n%s",
		m.content,
		m.help.View(m.keys),
	)
}

// Log Model
func (m *LogModel) Init() tea.Cmd {
	return nil
}

func (m *LogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *LogModel) View() string {
	return ``
}
