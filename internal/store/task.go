package store

import (
	"context"
	"database/sql"
	"errors"
	"teado/internal/models"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type taskStore struct {
	db      *sql.DB
	timeout time.Duration
}

func NewTaskStore(path string, timeout time.Duration) (*taskStore, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	const query = `
		CREATE TABLE IF NOT EXISTS tasks (
        	id INTEGER PRIMARY KEY AUTOINCREMENT,
        	name TEXT,
        	description TEXT,
        	is_done BOOLEAN
    	);`

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err = db.ExecContext(ctx, query)
	if err != nil {
		return nil, err
	}

	return &taskStore{db, timeout}, nil
}

func (ts taskStore) Create(task *models.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), ts.timeout)
	defer cancel()

	const query = "INSERT INTO tasks (name, description, is_done) VALUES(?, ?, ?)"

	res, err := ts.db.ExecContext(ctx, query, task.Title, task.Description, task.IsDone)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	task.Id = uint64(id)

	return nil
}

func (ts taskStore) Read() (models.Tasks, error) {
	const query = "SELECT * FROM tasks"

	ctx, cancel := context.WithTimeout(context.Background(), ts.timeout)
	defer cancel()

	rows, err := ts.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks models.Tasks
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.Id, &t.Title, &t.Description, &t.IsDone); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (ts taskStore) Delete(id uint64) error {
	const query = "DELETE FROM tasks WHERE id = ?"

	ctx, cancel := context.WithTimeout(context.Background(), ts.timeout)
	defer cancel()

	res, err := ts.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected < 1 {
		return errors.New("0 tasks deleted")
	}

	return nil
}
