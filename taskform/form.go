package taskform

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

const maxWidth = 40

var (
	red    = lipgloss.AdaptiveColor{Light: "#FE5F86", Dark: "#FE5F86"}
	indigo = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	green  = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
)

type Styles struct {
	Base,
	HeaderText,
	Status,
	StatusHeader,
	Highlight,
	ErrorHeaderText,
	Help lipgloss.Style
}

func NewStyles(lg *lipgloss.Renderer) *Styles {
	s := Styles{}
	s.Base = lg.NewStyle().
		Padding(1, 4, 0, 1)
	s.HeaderText = lg.NewStyle().
		Foreground(indigo).
		Bold(true).
		Padding(0, 1, 0, 2)
	s.Status = lg.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(indigo).
		PaddingLeft(1).
		MarginTop(1)
	s.StatusHeader = lg.NewStyle().
		Foreground(green).
		Bold(true)
	s.Highlight = lg.NewStyle().
		Foreground(lipgloss.Color("212"))
	s.ErrorHeaderText = s.HeaderText.
		Foreground(red)
	s.Help = lg.NewStyle().
		Foreground(lipgloss.Color("240"))
	return &s
}

type Model struct {
	lg     *lipgloss.Renderer
	styles *Styles
	form   *huh.Form
	width  int
}

func NewModel() Model {
	m := Model{width: maxWidth}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("Name").
				Validate(func(s string) error {
					if len(s) <= 0 {
						return fmt.Errorf("Name Required")
					}
					return nil
				}),

			huh.NewText().
				Key("description").
				Title("Description").Editor(),

			huh.NewConfirm().
				Title("All done?").
				Validate(func(v bool) error {
					if !v {
						return fmt.Errorf("Welp, finish up then")
					}
					return nil
				}).
				Affirmative("Yep").
				Negative("Nay"),
		),
	).
		WithWidth(maxWidth).
		WithShowHelp(false).
		WithShowErrors(false)
	return m
}

func (m Model) Binds() []key.Binding {
	return m.form.KeyBinds()
}

func (m Model) Init() tea.Cmd {
	return m.form.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	if m.form.State == huh.StateCompleted {
		cmds = append(cmds, func() tea.Msg { return "back" })
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	s := m.styles

	v := strings.TrimSuffix(m.form.View(), "\n\n")
	form := m.lg.NewStyle().Margin(1, 0).Render(v)

	errors := m.form.Errors()
	header := m.appBoundaryView("New Task")
	if len(errors) > 0 {
		header = m.appErrorBoundaryView(m.errorView())
	}

	return s.Base.Render(header + "\n" + form)
}

func (m Model) errorView() string {
	var s string
	for _, err := range m.form.Errors() {
		s += "| " + err.Error() + " "
	}
	return s
}

func (m Model) appBoundaryView(text string) string {
	return m.styles.HeaderText.Render(text)
}

func (m Model) appErrorBoundaryView(text string) string {
	return m.styles.ErrorHeaderText.Render(text)
}
