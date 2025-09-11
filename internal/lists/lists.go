package lists

import (
	"strings"
	"teado/internal/models"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type focus int

const (
	toDo focus = iota
	done
)

type model struct {
	current  focus
	lists    []list.Model
	keys     keys
	styles   styles
	fullHelp bool
	width    int
	height   int
}

func New() *model {
	return &model{
		current: toDo,
		lists: []list.Model{
			newList("TO DO"),
			newList("DONE"),
		},
		keys:   *defaultKeys(),
		styles: *defaultStyles(),
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.lists[m.current].FilterState() == list.Filtering {
			break
		}
		switch {
		case key.Matches(msg, m.keys.insertItem):
			return m, addTaskCmd(m.current == done)

		case key.Matches(msg, m.keys.more):
			m.fullHelp = !m.fullHelp
			m.setSize(m.width, m.height)
			return m, nil

		case key.Matches(msg, m.keys.toggle):
			if m.current == done {
				m.current = toDo
			} else {
				m.current = done
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.setSize(msg.Width, msg.Height)
		return m, nil

	case UpdateTasksMsg:
		m.updateLists(msg.Tasks)
		return m, nil

	case NewTaskMsg:
		task := msg.Task
		if task.IsDone {
			m.lists[done].InsertItem(0, task)
		} else {
			m.lists[toDo].InsertItem(0, task)
		}

		return m, nil
	}

	var cmd tea.Cmd
	m.lists[m.current], cmd = m.lists[m.current].Update(msg)

	return m, cmd
}

func (m model) View() string {
	views := make([]string, len(m.lists))
	for i, l := range m.lists {
		view := l.View()
		if m.current == focus(i) {
			view = m.styles.focused.Render(view)
		} else {
			view = m.styles.column.Render(view)
		}
		views[i] = view
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Left, views...),
		m.helpView(),
	)
}

func (m *model) setSize(width, height int) {
	m.width, m.height = width, height
	helpHeight := m.helpSize()

	m.styles.focused = m.styles.focused.
		Width(m.width/2 - 2).
		MaxHeight(m.height - helpHeight)
	m.styles.column = m.styles.column.
		Width(m.width/2 - 2).
		MaxHeight(m.height - helpHeight)

	for i := range m.lists {
		m.lists[i].SetSize(m.width, m.height-helpHeight-5)
	}
}

func (m *model) updateLists(tasks models.Tasks) {
	m.lists[toDo].SetItems(tasksToItems(tasks.GetToDo()))
	m.lists[done].SetItems(tasksToItems(tasks.GetDone()))
}

func (m model) helpView() string {
	current := m.lists[m.current]
	if m.fullHelp {
		return current.Help.FullHelpView(current.FullHelp())
	}
	return current.Help.ShortHelpView(current.ShortHelp())
}

func (m model) helpSize() int {
	return strings.Count(m.helpView(), "\n") + 1
}
