package app

import (
	"teado/menu"
	"teado/taskform"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type view int

const (
	mainMenu view = iota
	taskForm
)

type app struct {
	current  view
	mainMenu menu.Model
	taskForm taskform.Model
	keys     binds
	help     help.Model
}

func Start() tea.Model {
	keys := binds{
		help: key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
		quit: key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"), key.WithHelp("q", "quit")),
	}

	return app{
		current:  mainMenu,
		mainMenu: menu.NewMenu("Tea Do", "New Task", "Quit"),
		taskForm: taskform.NewModel(),
		keys:     keys,
		help:     help.New(),
	}
}

func (a app) Init() tea.Cmd {
	return a.taskForm.Init()
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if kMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(kMsg, a.keys.quit):
			return a, tea.Quit

		case key.Matches(kMsg, a.keys.help):
			a.help.ShowAll = !a.help.ShowAll
			return a, nil
		}
	}

	switch msg := msg.(type) {
	case menu.GotOption:
		if msg.Label == "New Task" {
			a.current = taskForm
			return a, nil
		}
		if msg.Label == "Quit" {
			return a, tea.Quit
		}
	}

	switch a.current {
	case mainMenu:
		newModel, cmd := a.mainMenu.Update(msg)
		a.mainMenu = newModel.(menu.Model)
		return a, cmd

	case taskForm:
		newModel, cmd := a.taskForm.Update(msg)
		a.taskForm = newModel.(taskform.Model)
		return a, cmd
	}

	return a, nil
}

func (a app) View() string {
	view := ""
	switch a.current {
	case mainMenu:
		view = a.mainMenu.View()
		a.keys.innerKeys = a.mainMenu.Binds()
	case taskForm:
		view = a.taskForm.View()
		a.keys.innerKeys = a.taskForm.Binds()
	}

	view += "\n" + a.help.View(a.keys)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2, 1, 2).
		Render(view)
}
