package store

import (
	"context"
	"database/sql"
	"errors"
	"teado/internal/models"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type store struct {
	db      *sql.DB
	timeout time.Duration
}

func NewStore(path string, timeout time.Duration) (*store, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	const query = `
		CREATE TABLE IF NOT EXISTS tasks (
        	id INTEGER PRIMARY KEY AUTOINCREMENT,
        	title TEXT NOT NULL,
        	description TEXT,
        	is_done BOOLEAN,
			updated_at DATETIME NOT NULL
    	);`

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err = db.ExecContext(ctx, query)
	if err != nil {
		return nil, err
	}

	return &store{db, timeout}, nil
}

func (store store) Create(task *models.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), store.timeout)
	defer cancel()

	const query = "INSERT INTO tasks (title, description, is_done, updated_at) VALUES(?, ?, ?, ?)"
	task.UpdatedAt = time.Now()

	res, err := store.db.ExecContext(ctx, query, task.Title, task.Description, task.IsDone, task.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	task.Id = id

	return nil
}

func (store store) Read() (models.Tasks, error) {
	const query = "SELECT * FROM tasks"

	ctx, cancel := context.WithTimeout(context.Background(), store.timeout)
	defer cancel()

	rows, err := store.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks models.Tasks
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.Id, &t.Title, &t.Description, &t.IsDone, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (store store) Update(task *models.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), store.timeout)
	defer cancel()

	const query = `
        UPDATE tasks
        SET title = ?, description = ?, is_done = ?, updated_at = ?
        WHERE id = ?`

	task.UpdatedAt = time.Now()
	_, err := store.db.ExecContext(ctx, query, task.Title, task.Description, task.IsDone, task.UpdatedAt, task.Id)
	if err != nil {
		return err
	}

	return nil
}

func (store store) Delete(id int64) error {
	const query = "DELETE FROM tasks WHERE id = ?"

	ctx, cancel := context.WithTimeout(context.Background(), store.timeout)
	defer cancel()

	res, err := store.db.ExecContext(ctx, query, id)
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
