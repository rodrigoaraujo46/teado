package taskform

import (
	"errors"
	"teado/internal/models"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	form *huh.Form
	task models.Task
}

func New() *model {
	return &model{
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("name").
					Title("Name").
					Validate(func(s string) error {
						if len(s) == 0 {
							return errors.New("Title Required")
						}
						return nil
					}),

				huh.NewText().
					Key("description").
					Title("Description"),

				huh.NewConfirm().
					Title("All done?").
					Validate(func(b bool) error {
						if !b {
							return errors.New("Well, finish up then")
						}
						return nil
					}),
			).Title(lipgloss.NewStyle().Padding(0, 0, 1, 0).Render("New Task")),
		),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case NewTaskFormMsg:
		m.task.IsDone = msg.Done
		return m, m.form.Init()
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		if m.form.State == huh.StateCompleted {
			m.task.Title = m.form.GetString("name")
			m.task.Description = m.form.GetString("description")

			return m, newTaskCmd(m.task)
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

type NewTaskFormMsg struct {
	Done bool
}
