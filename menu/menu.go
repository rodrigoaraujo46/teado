package menu

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	title    string
	options  []string
	cursor   int
	binds    binds
	help     help.Model
	showHelp bool
}

func NewMenu(title string, options ...string) Model {
	b := binds{
		up:     key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("↑/k", "up")),
		down:   key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("↓/j", "down")),
		choose: key.NewBinding(key.WithKeys("enter", " "), key.WithHelp("enter", "choose")),
	}

	return Model{binds: b, help: help.New(), title: title, options: options}
}

func (m Model) Init() tea.Cmd { return nil }

type GotOption struct {
	Index int
	Label string
}

func (m Model) getOption() tea.Cmd {
	return func() tea.Msg {
		return GotOption{
			Index: m.cursor,
			Label: m.options[m.cursor],
		}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.binds.up):
			m.cursor = max(m.cursor-1, 0)

		case key.Matches(msg, m.binds.down):
			m.cursor = min(m.cursor+1, len(m.options)-1)

		case key.Matches(msg, m.binds.choose):
			return m, m.getOption()
		}
	}

	return m, nil
}

func (m Model) View() string {
	titleS := lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	title := fmt.Sprintln(titleS.Render(m.title))

	options := ""
	for i, choice := range m.options {
		isSelected := (i == m.cursor)
		options += fmt.Sprintln(optionString(choice, isSelected))
	}

	return title + "\n" + options + "\n"
}

func (m Model) Binds() []key.Binding {
	return m.binds.slice()
}

func optionString(label string, isChecked bool) string {
	if isChecked {
		color := lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
		return color.Render(fmt.Sprintf("[x] %s", label))
	}

	return fmt.Sprintf("[ ] %s", label)
}
