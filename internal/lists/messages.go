package lists

import (
	tea "github.com/charmbracelet/bubbletea"
)

type AddTaskMsg struct {
	Done bool
}

func addTaskCmd(done bool) tea.Cmd {
	return func() tea.Msg {
		return AddTaskMsg{done}
	}
}
