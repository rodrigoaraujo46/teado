package models

import (
	"slices"
	"time"
)

type Tasks []Task

type Task struct {
	Id          int64
	Title       string
	Description string
	IsDone      bool
	UpdatedAt   time.Time
}

func NewTask(title, description string, isDone bool) *Task {
	return &Task{Title: title, Description: description, IsDone: isDone}
}

func (t Task) FilterValue() string { return t.Title }

func (tasks Tasks) SplitByIsDone() (todo Tasks, done Tasks) {
	for _, task := range tasks {
		if task.IsDone {
			done = append(done, task)
		} else {
			todo = append(todo, task)
		}
	}

	return todo, done
}

func (tasks *Tasks) SortByMostRecent() *Tasks {
	slices.SortFunc(*tasks, func(a, b Task) int {
		return b.UpdatedAt.Compare(a.UpdatedAt)
	})

	return tasks
}
