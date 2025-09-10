package taskform

import (
	"teado/internal/models"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type model struct {
	form *huh.Form
}

func New() *model {
	return &model{
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("name").
					Title("Name"),

				huh.NewText().
					Key("description").
					Title("Description"),
			).Title("New Task"),
		),
	}
}

func (m model) Init() tea.Cmd {
	return m.form.Init()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		if m.form.State == huh.StateCompleted {
			name := m.form.GetString("name")
			description := m.form.GetString("description")

			task := models.NewTask(0, name, description, false)
			return m, newTaskCmd(*task)
		}

		m.form = f
	}

	return m, cmd
}

func (m model) View() string {
	return m.form.View()
}

type NewTaskMsg struct {
	Task models.Task
}

func newTaskCmd(t models.Task) tea.Cmd {
	return func() tea.Msg {
		return NewTaskMsg{t}
	}
}
