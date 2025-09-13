package messages

import "teado/internal/models"

type (
	GoBack struct{}

	TasksRead struct{ Tasks models.Tasks }

	CreateTask struct{ Done bool }

	TaskCreated struct{ Task models.Task }

	UpdateTask struct{ Task models.Task }

	TaskUpdated struct{ Task models.Task }

	TaskDeleted struct{ Task models.Task }
)
