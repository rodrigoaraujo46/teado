package app

import (
	"teado/internal/form"
	"teado/internal/lists"
	"teado/internal/models"

	tea "github.com/charmbracelet/bubbletea"
)

type TaskStore interface {
	Create(*models.Task) error
	Read() (models.Tasks, error)
	Delete(uint64) error
}

type tasksLoadedMsg struct {
	tasks models.Tasks
}

func (m model) tasksLoadedCmd() tea.Msg {
	tasks, err := m.taskStore.Read()
	if err != nil {
		return nil
	}

	return tasksLoadedMsg{tasks}
}

type view uint8

const (
	tasks view = iota
	form
)

type model struct {
	current   view
	views     map[view]tea.Model
	taskStore TaskStore
}

func New(tasklist tea.Model, taskform tea.Model, taskstore TaskStore) *model {
	views := make(map[view]tea.Model)
	views[tasks] = tasklist
	views[form] = taskform

	return &model{
		current:   tasks,
		views:     views,
		taskStore: taskstore,
	}
}

func (m model) Init() tea.Cmd { return m.tasksLoadedCmd }

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}

	case tasksLoadedMsg:
		var cmd tea.Cmd
		m.views[tasks], cmd = m.views[tasks].Update(lists.UpdateTasksMsg{Tasks: msg.tasks})
		return m, cmd

	case lists.NewTaskMsg:
		m.current = form
		m.views[form] = taskform.New()
		return m, m.views[form].Init()

	case taskform.NewTaskMsg:
		if err := m.taskStore.Create(&msg.Task); err != nil {
			return m, nil
		}
		m.current = tasks
		return m, nil
	}

	var cmd tea.Cmd
	m.views[m.current], cmd = m.views[m.current].Update(msg)

	return m, cmd
}

func (m model) View() string {
	return m.views[m.current].View()
}
