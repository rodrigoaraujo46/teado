package tasks

import (
	"teado/internal/models"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rodrigoaraujo46/assert"
)

type delegateKeys struct {
	choose key.Binding
	remove key.Binding
}

func newDelegateKeys() *delegateKeys {
	return &delegateKeys{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "choose"),
		),
		remove: key.NewBinding(
			key.WithKeys("x", "backspace"),
			key.WithHelp("x", "delete"),
		),
	}
}

func (d delegateKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		d.choose,
		d.remove,
	}
}

func (d delegateKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.choose,
			d.remove,
		},
	}
}

type delegate struct {
	list.DefaultDelegate
	store TaskStore
	keys  delegateKeys
}

func newDelegate(keys delegateKeys, ts TaskStore) *delegate {
	d := &delegate{
		DefaultDelegate: list.NewDefaultDelegate(),
		keys:            keys,
		store:           ts,
	}

	return d
}

func (d *delegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	var title string
	if i, ok := m.SelectedItem().(models.Task); ok {
		title = i.Title()
	} else {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, d.keys.choose):
			return m.NewStatusMessage("You chose " + title)

		case key.Matches(msg, d.keys.remove):
			task, ok := m.SelectedItem().(models.Task)
			assert.Assert(ok, "Item needs to be a models.Task")

			err := d.store.Delete(task.Id)
			if err != nil {
				return nil
			}

			index := m.Index()
			m.RemoveItem(index)

			return m.NewStatusMessage("Deleted " + title)
		}
	}

	return nil
}

func (d delegate) ShortHelp() []key.Binding {
	return []key.Binding{d.keys.choose, d.keys.remove}
}

func (d delegate) FullHelp() [][]key.Binding {
	return [][]key.Binding{{d.keys.choose, d.keys.remove}}
}
