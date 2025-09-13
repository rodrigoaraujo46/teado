package board

import (
	"slices"
	"teado/internal/messages"
	"teado/internal/models"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rodrigoaraujo46/assert"
)

func tasksToItems(tasks models.Tasks) []list.Item {
	items := make([]list.Item, len(tasks))
	for i, task := range tasks {
		items[i] = task
	}
	return items
}

func newList(title string, store Store) list.Model {
	delegateKeys := newDelegateKeys()
	keys := defaultKeys()

	delegate := newDelegate(*delegateKeys, store)
	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = title
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{keys.insertItem}
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{keys.insertItem}
	}
	l.SetShowHelp(false)
	return l
}

func createTask(done bool) tea.Cmd {
	return func() tea.Msg {
		return messages.CreateTask{Done: done}
	}
}

func getTaskIndex(items []list.Item, taskId int64) int {
	index := slices.IndexFunc(items, func(i list.Item) bool {
		task, ok := i.(models.Task)
		assert.Assert(ok, "list.Item must be task in case TaskDeleted")
		if task.Id == taskId {
			return true
		}
		return false
	})

	return index
}

func updateTask(task models.Task) tea.Cmd {
	return func() tea.Msg {
		return messages.UpdateTask{Task: task}
	}
}
