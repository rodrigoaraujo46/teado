package lists

import (
	"teado/internal/models"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

func tasksToItems(tasks models.Tasks) []list.Item {
	items := make([]list.Item, len(tasks))
	for i, task := range tasks {
		items[i] = task
	}
	return items
}

func newList(title string) list.Model {
	delegateKeys := newDelegateKeys()
	keys := defaultKeys()

	delegate := newDelegate(*delegateKeys)
	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = title
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{keys.insertItem, keys.toggle}
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{keys.insertItem, keys.toggle}
	}
	l.SetShowHelp(false)
	return l
}
