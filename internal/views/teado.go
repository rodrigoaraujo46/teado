package views

import (
	"teado/internal/messages"

	tea "github.com/charmbracelet/bubbletea"
)

type view uint8

const (
	boardView = iota
	formView
)

type root struct {
	current view
	board   tea.Model
	form    tea.Model
}

func New(board tea.Model, form tea.Model) *root {
	return &root{
		current: boardView,
		board:   board,
		form:    form,
	}
}

func (r root) Init() tea.Cmd {
	return r.board.Init()
}

func (r *root) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		r.board, cmd = r.board.Update(msg)
		cmds = append(cmds, cmd)
		r.form, cmd = r.form.Update(msg)
		cmds = append(cmds, cmd)

		return r, tea.Batch(cmds...)

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return r, tea.Quit
		}

	case messages.CreateTask:
		r.current = formView

	case messages.TaskCreated:
		r.current = boardView

	case messages.UpdateTask:
		r.current = formView

	case messages.TaskUpdated:
		r.current = boardView

	case messages.GoBack:
		r.current = boardView
	}

	switch r.current {
	case boardView:
		r.board, cmd = r.board.Update(msg)
	case formView:
		r.form, cmd = r.form.Update(msg)
	}

	return r, cmd
}

func (r root) View() string {
	switch r.current {
	case boardView:
		return r.board.View()
	default:
		return r.form.View()
	}
}
