package taskform

import (
	"errors"
	"teado/internal/models"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var backKey = key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "quit form"))

type Model struct {
	form   *huh.Form
	task   models.Task
	width  int
	height int
}

func New() *Model {
	return &Model{form: nil}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height-3
		if m.form != nil {
			m.form.WithWidth(m.width).WithHeight(m.height)
		}
		return *m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return *m, goBackCmd
		}
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		if m.form.State == huh.StateCompleted {
			m.task.Title = m.form.GetString("name")
			m.task.Description = m.form.GetString("description")
			cmd := newTaskCmd(m.task)

			return *m, cmd
		}

		m.form = f
	}

	return *m, cmd
}

func (m Model) View() string {
	help := m.form.Help().ShortHelpView(append(m.form.KeyBinds(), backKey))

	return lipgloss.JoinVertical(lipgloss.Top, m.form.View(), help)
}

func (m *Model) newForm(width, height int) *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("Title").
				Placeholder("Task Title").
				Validate(func(s string) error {
					if len(s) == 0 {
						return errors.New("Title Required")
					}
					return nil
				}),

			huh.NewText().
				Key("description").
				Title("Description").
				Placeholder("Task description"),

			huh.NewConfirm().
				Title("All done?").
				Validate(func(b bool) error {
					if !b {
						return errors.New("Well, finish up then")
					}
					return nil
				}),
		).Title(lipgloss.NewStyle().Padding(0, 0, 1, 0).Render("New Task")),
	).WithWidth(width).WithHeight(height).WithShowHelp(false)

	return form
}

func (m *Model) Start(isTaskDone bool) {
	m.task.IsDone = isTaskDone
	m.form = m.newForm(m.width, m.height)
	m.form.Init()
}

type NewTaskMsg struct {
	Task models.Task
}

func newTaskCmd(t models.Task) tea.Cmd {
	return func() tea.Msg {
		return NewTaskMsg{t}
	}
}

type GoBackMsg struct{}

func goBackCmd() tea.Msg {
	return GoBackMsg{}
}
