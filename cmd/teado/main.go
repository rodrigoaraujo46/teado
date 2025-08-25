package main

import (
	"fmt"
	"os"
	"teado/app"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(app.Start())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
