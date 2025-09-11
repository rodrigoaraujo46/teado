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
	tasks     lists.Model
	form      taskform.Model
	taskStore TaskStore
}

func New(lists lists.Model, form taskform.Model, store TaskStore) *model {
	return &model{
		current:   tasks,
		tasks:     lists,
		form:      form,
		taskStore: store,
	}
}

func (m model) Init() tea.Cmd {
	return m.tasksLoadedCmd
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		_, cmd1 := m.tasks.Update(msg)
		_, cmd2 := m.form.Update(msg)

		return m, tea.Batch(cmd1, cmd2)

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}

	case taskform.GoBackMsg:
		m.current = tasks
		return m, nil

	case tasksLoadedMsg:
		m.tasks.UpdateTasks(msg.tasks)
		return m, nil

	case lists.AddTaskMsg:
		m.current = form
		m.form.Start(msg.Done)
		return m, nil

	case taskform.NewTaskMsg:
		if err := m.taskStore.Create(&msg.Task); err != nil {
			return m, nil
		}
		m.current = tasks
		m.tasks.InsertTask(msg.Task)
		return m, nil
	}

	var cmd tea.Cmd
	switch m.current {
	case tasks:
		_, cmd = m.tasks.Update(msg)
	case form:
		_, cmd = m.form.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	switch m.current {
	case tasks:
		return m.tasks.View()
	default:
		return m.form.View()
	}
}
