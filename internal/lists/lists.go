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

type Model struct {
	current  focus
	lists    []list.Model
	keys     keys
	styles   styles
	fullHelp bool
	width    int
	height   int
}

func New() *Model {
	return &Model{
		current: toDo,
		lists: []list.Model{
			newList("TO DO"),
			newList("DONE"),
		},
		keys:   *defaultKeys(),
		styles: *defaultStyles(),
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.lists[m.current].FilterState() == list.Filtering {
			break
		}
		switch {
		case key.Matches(msg, m.keys.insertItem):
			return *m, addTaskCmd(m.current == done)

		case key.Matches(msg, m.keys.more):
			m.fullHelp = !m.fullHelp
			m.setSize(m.width, m.height)
			return *m, nil

		case key.Matches(msg, m.keys.toggle):
			if m.current == done {
				m.current = toDo
			} else {
				m.current = done
			}
			return *m, nil
		}

	case tea.WindowSizeMsg:
		m.setSize(msg.Width, msg.Height)
		return *m, nil

	}

	var cmd tea.Cmd
	m.lists[m.current], cmd = m.lists[m.current].Update(msg)

	return *m, cmd
}

func (m Model) View() string {
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

func (m *Model) InsertTask(task models.Task) {
	if task.IsDone {
		m.lists[done].InsertItem(0, task)
	} else {
		m.lists[toDo].InsertItem(0, task)
	}
}

func (m *Model) setSize(width, height int) {
	m.width, m.height = width, height
	helpHeight := m.helpLines()

	m.styles.focused = m.styles.focused.
		Width(m.width/2 - 2).
		MaxHeight(m.height - helpHeight)
	m.styles.column = m.styles.column.
		Width(m.width/2 - 2).
		MaxHeight(m.height - helpHeight)

	for i := range m.lists {
		m.lists[i].SetSize(m.width/2-2, m.height-helpHeight-5)
	}
}

func (m *Model) UpdateTasks(tasks models.Tasks) {
	m.lists[toDo].SetItems(tasksToItems(tasks.GetToDo()))
	m.lists[done].SetItems(tasksToItems(tasks.GetDone()))
}

func (m Model) helpView() string {
	current := m.lists[m.current]

	/*
		I set width to full because we are rendering help outside the list,
		afterwards we set the size back so we render items properly in delegate.

		I should probably create my own help function instead of using list
		FullHelp/ShortHelp so that I can render help without doing this.
	*/

	current.SetWidth(m.width - 2)
	var view string
	if m.fullHelp {
		view = current.Help.FullHelpView(current.FullHelp())
	} else {
		view = current.Help.ShortHelpView(current.ShortHelp())
	}
	current.SetSize(m.width, m.height)

	return view
}

func (m Model) helpLines() int {
	return strings.Count(m.helpView(), "\n") + 1
}
