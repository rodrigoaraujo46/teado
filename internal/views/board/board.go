package board

import (
	"slices"
	"strings"
	"teado/internal/messages"
	"teado/internal/models"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Store interface {
	Read() (models.Tasks, error)
	Update(*models.Task) error
	Delete(int64) error
}

type focus int

const (
	toDo focus = iota
	done
)

type board struct {
	current  focus
	lists    []list.Model
	store    Store
	keys     keys
	styles   styles
	fullHelp bool
	width    int
	height   int
}

func New(ts Store) *board {
	return &board{
		current: toDo,
		lists: []list.Model{
			newList("TO DO", ts),
			newList("DONE", ts),
		},
		store:  ts,
		keys:   *defaultKeys(),
		styles: *defaultStyles(),
	}
}

func (b board) Init() tea.Cmd {
	return b.ReadTasks
}

func (b board) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.width, b.height = msg.Width, msg.Height
		return b, nil

	case tea.KeyMsg:
		if b.lists[b.current].FilterState() == list.Filtering {
			break
		}
		switch {
		case key.Matches(msg, b.keys.insertItem):
			return b, createTask(b.current == done)

		case key.Matches(msg, b.keys.more):
			b.fullHelp = !b.fullHelp
			return b, nil

		case key.Matches(msg, b.keys.toggleList):
			b.current = (b.current + 1) % 2
			return b, nil
		}

	case messages.TaskCreated:
		focus := toDo
		if msg.Task.IsDone {
			focus = done
		}

		b.lists[focus].InsertItem(0, msg.Task)
		b.lists[focus].ResetFilter()
		b.lists[focus].ResetSelected()
		return b, nil

	case messages.TasksRead:
		b.updateTasks(msg.Tasks)
		return b, nil

	case messages.TaskUpdated:
		var cmd tea.Cmd
		if index := getTaskIndex(b.lists[toDo].Items(), msg.Task.Id); index != -1 {
			cmd = deleteSelectedItem(&b.lists[toDo])
		} else if index := getTaskIndex(b.lists[done].Items(), msg.Task.Id); index != -1 {
			cmd = deleteSelectedItem(&b.lists[done])
		} else {
			return b, nil
		}

		focus := toDo
		if msg.Task.IsDone {
			focus = done
		}

		b.lists[focus].InsertItem(0, msg.Task)
		return b, cmd

	case messages.TaskDeleted:
		var l *list.Model
		if msg.Task.IsDone {
			l = &b.lists[done]
		} else {
			l = &b.lists[toDo]
		}

		return b, deleteSelectedItem(l)
	}

	var cmd tea.Cmd
	b.lists[b.current], cmd = b.lists[b.current].Update(msg)
	return b, cmd
}

func (b board) View() string {
	views := make([]string, len(b.lists))
	b.calculateSizes()
	for i, l := range b.lists {
		view := l.View()
		if b.current == focus(i) {
			view = b.styles.focused.Render(view)
		} else {
			view = b.styles.column.Render(view)
		}
		views[i] = view
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Left, views...),
		b.helpView(),
	)
}

// I have to do this bullshit because charm is taking over a year to merge
// a pr that changes 3 lines for removeitem.
func deleteSelectedItem(l *list.Model) tea.Cmd {
	//Previous satate
	pVisible := l.VisibleItems()
	pFilter := l.FilterValue()
	wasFiltered := l.IsFiltered()
	i := l.Index()

	//Set Items Globally
	g := l.GlobalIndex()
	cmd := l.SetItems(slices.Delete(l.Items(), g, g+1))

	//Restore State
	if wasFiltered {
		l.SetFilterText(pFilter)
	}
	nLen := len(pVisible) - 1
	if i >= nLen {
		i = nLen - 1
	}
	if i < 0 {
		i = 0
		if wasFiltered {
			l.ResetFilter()
		}
	}
	l.Select(i)

	return cmd
}

func (b board) helpView() string {
	current := b.lists[b.current]

	/*
		I set width to full because we are rendering help outside the list,
		afterwards we set the size back so we render items properly in delegate.

		I should probably create my own help function instead of using list
		FullHelp/ShortHelp so that I can render help without doing this.
	*/

	current.SetWidth(b.width - 2)
	var view string
	if b.fullHelp {
		view = current.Help.FullHelpView(current.FullHelp())
	} else {
		view = current.Help.ShortHelpView(current.ShortHelp())
	}
	current.SetSize(b.width, b.height)

	return view
}

func (b board) helpLines() int {
	return strings.Count(b.helpView(), "\n") + 1
}

func (b *board) ReadTasks() tea.Msg {
	tasks, err := b.store.Read()
	if err != nil {
		return nil
	}

	return messages.TasksRead{Tasks: tasks}
}

func (b *board) calculateSizes() {
	helpHeight := b.helpLines()

	b.styles.focused = b.styles.focused.
		Width(b.width/2 - 2).
		MaxHeight(b.height - helpHeight)
	b.styles.column = b.styles.column.
		Width(b.width/2 - 2).
		MaxHeight(b.height - helpHeight)

	for i := range b.lists {
		b.lists[i].SetSize(b.width/2-2, b.height-helpHeight-5)
	}
}

func (b *board) updateTasks(tasks models.Tasks) {
	unfinished, finished := tasks.SortByMostRecent().SplitByIsDone()

	b.lists[toDo].SetItems(tasksToItems(unfinished))
	b.lists[done].SetItems(tasksToItems(finished))
}
