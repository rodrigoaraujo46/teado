package lists

import (
	"teado/internal/models"

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

type NewTaskMsg struct {
	Task models.Task
}

type UpdateTasksMsg struct {
	Tasks models.Tasks
}
