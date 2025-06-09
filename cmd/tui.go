package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TuiCmd struct {
	// TODO: Have an instance of log and report cmd that this can then call on updates?
}

func (cmd *TuiCmd) AfterApply() error {
	return nil
}

func (cmd *TuiCmd) Run(log string) error {
	// TODO: Move this to an option or remove
	if os.Getenv("HELP_DEBUG") != "" {
		f, err := tea.LogToFile("debug.log", "help")
		if err != nil {
			return fmt.Errorf("opening log file: %w", err)
		}
		defer f.Close() // nolint:errcheck
	}

	// TODO: Have this in the CMD
	// TODO: Also should probably use entries.GetAll and use a bubbles table instead
	helper := ReportCmd{
		DurationFormat: "default",
		Output: "default",
	}
	// TODO: Call in our own AfterApply
	if err := helper.AfterApply(); err != nil {
		return fmt.Errorf(`applying: %w`, err)
	}

	// TODO: Call in the update method or somewhere else
	var stdout strings.Builder
	if err := helper.Run(&stdout, log); err != nil {
		return fmt.Errorf(`getting report: %w`, err)
	}

	if _, err := tea.NewProgram(newModel(stdout.String())).Run(); err != nil {
		return fmt.Errorf(`rending tui: %w`, err)
	}
	return nil
}

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	// TODO: Add more keybindings
	Up    key.Binding
	Down  key.Binding
	Help  key.Binding
	Quit  key.Binding
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
		{k.Up, k.Down}, // first column
		{k.Help, k.Quit},                // second column
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

type model struct {
	content string
	keys       keyMap
	help       help.Model
	inputStyle lipgloss.Style
	lastKey    string
	quitting   bool
}

func newModel(content string) model {
	return model{
		content: content,
		keys:       keys,
		help:       help.New(),
		inputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	return "\n" + m.content + "\n" + m.help.View(m.keys)
}
