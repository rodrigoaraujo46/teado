package app

import "github.com/charmbracelet/bubbles/key"

type binds struct {
	help      key.Binding
	quit      key.Binding
	innerKeys []key.Binding
}

func (k binds) ShortHelp() []key.Binding {
	return []key.Binding{k.help, k.quit}
}

func (k binds) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.innerKeys, {k.help, k.quit}}
}
