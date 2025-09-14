package board

import (
	"fmt"
	"io"
	"strings"
	"teado/internal/messages"
	"teado/internal/models"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/rodrigoaraujo46/assert"
)

type delegateKeys struct {
	choose key.Binding
	toggle key.Binding
	remove key.Binding
}

func newDelegateKeys() *delegateKeys {
	return &delegateKeys{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "choose"),
		),
		toggle: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "toggle completion"),
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
		d.toggle,
		d.remove,
	}
}

func (d delegateKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.choose,
			d.toggle,
			d.remove,
		},
	}
}

type delegate struct {
	list.DefaultDelegate
	keys  delegateKeys
	store Store
}

func newDelegate(keys delegateKeys, store Store) *delegate {
	return &delegate{
		list.NewDefaultDelegate(),
		keys,
		store,
	}
}

func (d delegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		task, ok := m.SelectedItem().(models.Task)
		if !ok {
			return nil
		}
		switch {
		case key.Matches(msg, d.keys.choose):
			return updateTask(task)

		case key.Matches(msg, d.keys.toggle):
			task.IsDone = !task.IsDone
			return d.updateTask(task)

		case key.Matches(msg, d.keys.remove):
			return d.deleteTask(task)
		}
	}

	return nil
}

func (d delegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	const (
		ellipsis = "..."
	)

	if m.Width() <= 0 {
		return
	}

	task, ok := item.(models.Task)
	assert.Assert(ok, "Item should be task in item delegate")

	var matchedRunes []int

	title := task.Title
	desc := task.Description
	d.ShowDescription = len(desc) != 0

	s := &d.Styles

	textwidth := m.Width() - s.NormalTitle.GetPaddingLeft() - s.NormalTitle.GetPaddingRight() - len(ellipsis) - 1
	title = ansi.Truncate(title, textwidth, ellipsis)
	if d.ShowDescription {
		var lines []string
		for i, line := range strings.Split(desc, "\n") {
			if i >= d.Height()-1 {
				break
			}
			lines = append(lines, ansi.Truncate(line, textwidth, ellipsis))
		}
		desc = strings.Join(lines, "\n")
	}

	var (
		isSelected  = index == m.Index()
		emptyFilter = m.FilterState() == list.Filtering && m.FilterValue() == ""
		isFiltered  = m.FilterState() == list.Filtering || m.FilterState() == list.FilterApplied
	)

	if isFiltered && index < len(m.VisibleItems()) {
		matchedRunes = m.MatchesForItem(index)
	}

	if emptyFilter {
		title = s.DimmedTitle.Render(title)
		desc = s.DimmedDesc.Render(desc)
	} else if isSelected && m.FilterState() != list.Filtering {
		if isFiltered {
			unmatched := s.SelectedTitle.Inline(true)
			matched := unmatched.Inherit(s.FilterMatch)
			title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
		}
		title = s.SelectedTitle.Render(title)
		desc = s.SelectedDesc.Render(desc)
	} else {
		if isFiltered {
			unmatched := s.NormalTitle.Inline(true)
			matched := unmatched.Inherit(s.FilterMatch)
			title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
		}
		title = s.NormalTitle.Render(title)
		desc = s.NormalDesc.Render(desc)
	}

	if d.ShowDescription {
		fmt.Fprintf(w, "%s\n%s", title, desc)
		return
	}
	fmt.Fprintf(w, "%s", title)
}

func (d delegate) ShortHelp() []key.Binding {
	return []key.Binding{d.keys.choose, d.keys.toggle, d.keys.remove}
}

func (d delegate) FullHelp() [][]key.Binding {
	return [][]key.Binding{{d.keys.choose, d.keys.toggle, d.keys.remove}}
}

func (d delegate) updateTask(task models.Task) tea.Cmd {
	return func() tea.Msg {
		d.store.Update(&task)
		return messages.TaskUpdated{Task: task}
	}
}

func (d delegate) deleteTask(task models.Task) tea.Cmd {
	return func() tea.Msg {
		d.store.Delete(task.Id)
		return messages.TaskDeleted{Task: task}
	}
}
