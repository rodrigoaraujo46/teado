package lists

import "github.com/charmbracelet/lipgloss"

type styles struct {
	focused lipgloss.Style
	column  lipgloss.Style
}

func defaultStyles() *styles {
	return &styles{
		focused: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2),
		column: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("241")).
			Padding(1, 2),
	}
}
