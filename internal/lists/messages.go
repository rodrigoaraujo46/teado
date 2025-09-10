package lists

import (
	"teado/internal/models"

	tea "github.com/charmbracelet/bubbletea"
)

type NewTaskMsg struct{}

func newTaskCmd() tea.Msg {
	return NewTaskMsg{}
}

type UpdateTasksMsg struct {
	Tasks models.Tasks
}
