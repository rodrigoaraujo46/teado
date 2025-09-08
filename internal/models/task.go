package models

type Task struct {
	Id     uint64
	Name   string
	Info   string
	IsDone bool
}

func NewTask(id uint64, name, description string, isDone bool) *Task {
	return &Task{id, name, description, isDone}
}

func (t Task) FilterValue() string { return t.Name }

func (t Task) Title() string { return t.Name }

func (t Task) Description() string { return t.Info }
