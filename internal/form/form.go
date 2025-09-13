package form

import (
	"errors"
	"teado/internal/messages"
	"teado/internal/models"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/rodrigoaraujo46/assert"
)

var backKey = key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "quit form"))

type Store interface {
	Create(*models.Task) error
	Update(*models.Task) error
}

type form struct {
	form       *huh.Form
	task       *models.Task
	isUpdating bool
	store      Store
	width      int
	height     int
}

func New(ts Store) *form {
	return &form{form: nil, task: &models.Task{}, store: ts}
}

func (f form) Init() tea.Cmd {
	return nil
}

func (m form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height-3
		if m.form != nil {
			m.form.WithWidth(m.width).WithHeight(m.height)
		}
		return m, nil

	case messages.CreateTask:
		m.startCreate(msg.Done)
		return m, nil

	case messages.UpdateTask:
		m.startUpdate(msg.Task)
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m, goBack
		}
	}

	var cmds []tea.Cmd
	f, cmd := m.form.Update(msg)
	cmds = append(cmds, cmd)

	_, ok := f.(*huh.Form)
	assert.Assert(ok, "form needs to be of type huh.form")
	if m.form.State == huh.StateCompleted {
		if m.isUpdating {
			return m, m.updateTask()
		}
		return m, m.createTask()
	}

	return m, tea.Batch(cmds...)
}

func (f form) View() string {
	help := f.form.Help().ShortHelpView(append(f.form.KeyBinds(), backKey))

	return lipgloss.JoinVertical(lipgloss.Top, f.form.View(), help)
}

func (f *form) newCreateForm(width, height int) *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Value(&f.task.Title).
				Title("Title").
				Placeholder("Task Title").
				Validate(func(s string) error {
					if len(s) == 0 {
						return errors.New("Title Required")
					}
					return nil
				}),

			huh.NewText().
				Value(&f.task.Description).
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

func (f *form) newUpdateForm(width, height int) *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Value(&f.task.Title).
				Title("Title").
				Placeholder("Task Title").
				Validate(func(s string) error {
					if len(s) == 0 {
						return errors.New("Title Required")
					}
					return nil
				}),

			huh.NewText().
				Value(&f.task.Description).
				Key("description").
				Title("Description").
				Placeholder("Task description"),

			huh.NewConfirm().
				Value(&f.task.IsDone).
				Key("isDone").
				Title("Is the task finished?").
				WithButtonAlignment(lipgloss.Left),

			huh.NewConfirm().
				Title("All done?").
				Validate(func(b bool) error {
					if !b {
						return errors.New("Well, finish up then")
					}
					return nil
				}),
		).Title(lipgloss.NewStyle().Padding(0, 0, 1, 0).Render("Edit Task")),
	).WithWidth(width).WithHeight(height).WithShowHelp(false)

	return form
}

func (f *form) startCreate(isDone bool) {
	f.task = &models.Task{IsDone: isDone}
	f.form = f.newCreateForm(f.width, f.height)
	f.isUpdating = false
	f.form.Init()
}

func (f *form) startUpdate(task models.Task) {
	f.task = &task
	f.form = f.newUpdateForm(f.width, f.height)
	f.isUpdating = true
	f.form.Init()
}

func (f form) createTask() tea.Cmd {
	return func() tea.Msg {
		f.store.Create(f.task)
		return messages.TaskCreated{Task: *f.task}
	}
}

func (f form) updateTask() tea.Cmd {
	return func() tea.Msg {
		if err := f.store.Update(f.task); err != nil {
			return nil
		}

		return messages.TaskUpdated{Task: *f.task}
	}
}

func goBack() tea.Msg {
	return messages.GoBack{}
}
