package main

import (
	"os"
	"teado/internal/app"
	"teado/internal/store"
	"teado/internal/tasks"
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

	taskStore, err := store.NewTaskStore("./DB", time.Second)
	assert.NoError(err, "Error creating taskStore.")

	taskHandler, err := tasks.New(taskStore)
	assert.NoError(err, "Error creating taskHandler with provided store.")

	_, err = tea.NewProgram(app.New(taskHandler), tea.WithAltScreen()).Run()
	assert.NoError(err, "Error running program")
}
