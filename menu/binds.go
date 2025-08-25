package menu

import "github.com/charmbracelet/bubbles/key"

type binds struct {
	up     key.Binding
	down   key.Binding
	choose key.Binding
}

func (b binds) slice() []key.Binding {
	return []key.Binding{b.up, b.down, b.choose}
}
