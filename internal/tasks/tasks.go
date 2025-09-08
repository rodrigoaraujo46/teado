package tasks

import (
	"teado/internal/models"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type keys struct {
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
	insertItem       key.Binding
}

func newKeyMaps() *keys {
	list.NewDefaultDelegate()
	return &keys{
		insertItem: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add item"),
		),
		toggleTitleBar: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "toggle title"),
		),
		toggleStatusBar: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "toggle status"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
	}
}

type TaskStore interface {
	Create(*models.Task) error
	Read() ([]models.Task, error)
	Delete(uint64) error
}

func tasksToItems(tasks []models.Task) []list.Item {
	items := make([]list.Item, len(tasks))
	for i, task := range tasks {
		items[i] = task
	}

	return items
}

type model struct {
	list         list.Model
	keys         keys
	delegateKeys delegateKeys
	taskStore    TaskStore
}

func New(taskStore TaskStore) (*model, error) {
	tasks, err := taskStore.Read()
	if err != nil {
		return nil, err
	}

	delegateKeys := newDelegateKeys()
	keys := newKeyMaps()

	delegate := newDelegate(*delegateKeys, taskStore)
	list := list.New(tasksToItems(tasks), delegate, 0, 0)
	list.Title = "Tasks"
	list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			keys.toggleSpinner,
			keys.insertItem,
			keys.toggleTitleBar,
			keys.toggleStatusBar,
			keys.togglePagination,
			keys.toggleHelpMenu,
		}
	}

	return &model{list, *keys, *delegateKeys, taskStore}, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.toggleSpinner):
			cmd := m.list.ToggleSpinner()
			return m, cmd

		case key.Matches(msg, m.keys.toggleTitleBar):
			v := !m.list.ShowTitle()
			m.list.SetShowTitle(v)
			m.list.SetShowFilter(v)
			m.list.SetFilteringEnabled(v)
			return m, nil

		case key.Matches(msg, m.keys.toggleStatusBar):
			m.list.SetShowStatusBar(!m.list.ShowStatusBar())
			return m, nil

		case key.Matches(msg, m.keys.togglePagination):
			m.list.SetShowPagination(!m.list.ShowPagination())
			return m, nil

		case key.Matches(msg, m.keys.toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil

		case key.Matches(msg, m.keys.insertItem):
			newItem := models.NewTask(0, "Goat", "MY MAN", false)
			err := m.taskStore.Create(newItem)
			if err != nil {
				return m, nil
			}

			insCmd := m.list.InsertItem(0, *newItem)
			statusCmd := m.list.NewStatusMessage("Added " + newItem.Title())
			return m, tea.Batch(insCmd, statusCmd)
		}
	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return m.list.View()
}
