package main

import (
	"os"
	"teado/internal/store"
	"teado/internal/views"
	"teado/internal/views/board"
	"teado/internal/views/form"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rodrigoaraujo46/assert"
)

func main() {
	f, err := tea.LogToFile("debug.log", "")
	assert.NoError(err, "Couldn't open debug.log")

	defer func(f *os.File) {
		assert.NoError(f.Close(), "File has already been closed")
	}(f)

	store, err := store.NewStore("./DB", time.Second)
	assert.NoError(err, "Error creating taskStore.")

	board := board.New(store)
	assert.NoError(err, "Error creating taskList with provided store.")

	form := form.New(store)

	_, err = tea.NewProgram(views.New(*board, *form), tea.WithAltScreen()).Run()
	assert.NoError(err, "Error running program")
}
