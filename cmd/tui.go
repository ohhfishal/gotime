package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TuiCmd struct {
	// TODO: Have an instance of log and report cmd that this can then call on updates?
	Debug   bool          `help:"Enabled debug options"`
	Refresh time.Duration `default:"30s" help:"Duration between refreshing report."`
	Report  ReportCmd     `embed:""`
	model   Model         `kong:"-"`
}

func (cmd *TuiCmd) AfterApply() error {
	return cmd.Report.AfterApply()
}

type Model struct {
	Debug  bool
	Report tea.Model
	Log    *LogModel // TODO: Implement
	// Edit   *EditModel // TODO: Implement TUI to edit already existing entries
	state tea.Model
	// TODO: Enable this if Debug is true
	spinner spinner.Model
}

type LogModel struct {
	// TODO: Implement TUI to create a new entry
}

type ReportModel struct {
	Debug      bool
	content    string
	keys       keyMap
	help       help.Model
	inputStyle lipgloss.Style
	lastKey    string
	quitting   bool

	Refresh         func() (string, error)
	refreshDuration time.Duration
	lastTick        time.Time
	lastRefresh     time.Time
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
		Debug: cmd.Debug,
		Report: &ReportModel{
			Debug:           cmd.Debug,
			content:         content,
			keys:            keys,
			help:            help.New(),
			inputStyle:      lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")),
			Refresh:         refresh,
			refreshDuration: cmd.Refresh,
		},
	}
	cmd.model.state = cmd.model.Report

	s := spinner.New()
	s.Spinner = spinner.Ellipsis
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	cmd.model.spinner = s

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

// Generic Root Model
func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.Report.Init(),
		m.spinner.Tick,
	)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	var cmd tea.Cmd

	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)

	m.Report, cmd = m.Report.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(
		cmds...,
	)
}

func (m *Model) View() string {
	var views []string
	if m.Debug {
		views = append(views,
			fmt.Sprintf(" Debug: %s", m.spinner.View()),
		)
	}
	views = append(views, m.Report.View())
	return strings.Join(views, "\n")
}

func (m *ReportModel) Init() tea.Cmd {
	return nil
}

// Report Model

func (m *ReportModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: Update the data object (Probably after we are using a bubbles.Table)
	switch msg := msg.(type) {
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
	default:
		m.lastTick = time.Now()
		if time.Since(m.lastRefresh) > m.refreshDuration {
			content, err := m.Refresh()
			if err != nil {
				// TODO: Handle this more elegantly
				panic(fmt.Errorf(`refreshing: %w`, err))
			}
			m.content = content
			m.lastRefresh = time.Now()
		}
		return m, nil
	}
	return m, nil
}

func (m *ReportModel) View() string {
	if m.quitting {
		return ""
	}
	var view string
	if m.Debug {
		view += fmt.Sprintf(" Last Tick: %s\n Last Refresh: %s\n",
			m.lastTick.Format(time.DateTime),
			m.lastRefresh.Format(time.DateTime),
		)
	}

	return strings.Join(
		[]string{
			view,
			m.content,
			m.help.View(m.keys),
		}, "\n",
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
