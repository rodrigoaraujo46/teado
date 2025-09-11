package main

import (
	"os"
	"teado/internal/app"
	"teado/internal/form"
	"teado/internal/lists"
	"teado/internal/store"
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

	lists := lists.New()
	assert.NoError(err, "Error creating taskList with provided store.")

	taskForm := taskform.New()

	_, err = tea.NewProgram(app.New(*lists, *taskForm, taskStore), tea.WithAltScreen()).Run()
	assert.NoError(err, "Error running program")
}
