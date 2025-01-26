package sqlite

import (
	"context"
	"database/sql"

	"github.com/askerdev/unisync_bot/internal/domain"
)

type storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) (*storage, error) {
	return &storage{db}, ApplyMigrations(db)
}

func (s *storage) Insert(ctx context.Context, tasks []*domain.Task) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	const cleanup = `DELETE FROM tasks`

	_, err = tx.ExecContext(ctx, cleanup)
	if err != nil {
		tx.Rollback()
		return err
	}

	const query = `INSERT INTO tasks(chat_id, text, time_at) VALUES (?, ?, ?)`

	for _, t := range tasks {
		_, err := tx.ExecContext(ctx, query, t.ChatID, t.Text, t.TimeAt)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s *storage) Select(ctx context.Context) ([]*domain.Task, error) {
	result := []*domain.Task{}

	const query = `SELECT id, chat_id, text, time_at FROM tasks`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		task := &domain.Task{}
		if err := rows.Scan(
			&task.ID,
			&task.ChatID,
			&task.Text,
			&task.TimeAt,
		); err != nil {
			return nil, err
		}
		result = append(result, task)
	}

	return result, nil
}

func (s *storage) Delete(ctx context.Context, ID int) error {
	const query = `DELETE FROM tasks WHERE id = ?`
	_, err := s.db.ExecContext(ctx, query, ID)
	return err
}
