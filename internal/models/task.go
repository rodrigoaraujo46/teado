package models

type Tasks []Task

type Task struct {
	Id          uint64
	Title       string
	Description string
	IsDone      bool
}

func NewTask(id uint64, name, description string, isDone bool) *Task {
	return &Task{id, name, description, isDone}
}

func (t Task) FilterValue() string { return t.Title }

func (tasks Tasks) GetToDo() Tasks {
	notDone := make(Tasks, 0)
	for _, task := range tasks {
		if !task.IsDone {
			notDone = append(notDone, task)
		}
	}

	return notDone
}

func (tasks Tasks) GetDone() Tasks {
	done := make(Tasks, 0)
	for _, task := range tasks {
		if task.IsDone {
			done = append(done, task)
		}
	}

	return done
}
