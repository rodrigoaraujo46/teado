package app

import tea "github.com/charmbracelet/bubbletea"

type view int8

const (
	tasks view = iota
)

type model struct {
	current view
	views   map[view]tea.Model
}

type TaskHandler interface {
	tea.Model
}

func New(th TaskHandler) *model {
	views := make(map[view]tea.Model)
	views[tasks] = th

	return &model{current: tasks, views: views}
}

func (m model) Init() tea.Cmd { return nil }

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.views[m.current], cmd = m.views[m.current].Update(msg)

	return m, cmd
}

func (m model) View() string {
	return m.views[m.current].View()
}
