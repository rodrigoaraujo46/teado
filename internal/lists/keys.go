package lists

import "github.com/charmbracelet/bubbles/key"

type keys struct {
	insertItem key.Binding
	toggle     key.Binding
	more       key.Binding
}

func defaultKeys() *keys {
	return &keys{
		insertItem: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add item"),
		),
		toggle: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "switch list"),
		),
		more: key.NewBinding(
			key.WithKeys("?"),
		),
	}
}
