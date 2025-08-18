package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/IlyasYOY/tasks-assistant-tg-bot/internal/task"
)

type SQLStore struct {
	db *sql.DB
}

func NewSQLStore(db *sql.DB) (*SQLStore, error) {
	return &SQLStore{db: db}, nil
}

func (s *SQLStore) Get(userID int64) (task.TasksText, error) {
	var txt string
	if err := s.db.QueryRowContext(context.Background(), `
		select tasks from user_tasks
		where user_id = ?
	`, userID).Scan(&txt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No tasks stored for this user.
			return task.TasksText(""), nil
		}
		// Propagate other errors.
		return task.TasksText(
				"",
			), fmt.Errorf(
				"querying tasks for user %d: %w",
				userID,
				err,
			)
	}
	return task.TasksText(txt), nil
}

func (s *SQLStore) Set(userID int64, t task.TasksText) error {
	_, err := s.db.ExecContext(context.Background(), `
		insert into user_tasks (user_id, tasks) values (?, ?)
		on conflict (user_id) do update set tasks = excluded.tasks;
	 `, userID, string(t))
	if err != nil {
		return fmt.Errorf("storing tasks for user %d: %w", userID, err)
	}
	return nil
}

func (s *SQLStore) Close() error {
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("closing DB: %w", err)
	}
	return nil
}
