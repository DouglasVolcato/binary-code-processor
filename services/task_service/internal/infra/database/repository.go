package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/douglasvolcato/binary-code-processor/task_service/internal/entities"
	"github.com/douglasvolcato/binary-code-processor/task_service/internal/usecases"
)

const (
	tasksTable       = "tasks"
	outboxTable      = "task_outbox_events"
	timeLayout       = "2006-01-02 15:04:05"
	eventUnprocessed = "unprocessed"
	eventProcessed   = "processed"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetTasks(ctx context.Context, limit int, offset int) ([]entities.Task, error) {
	rows, err := r.DB.QueryContext(
		ctx,
		fmt.Sprintf(`
			SELECT id, message, binary_code, created_at, updated_at
			FROM %s
			ORDER BY created_at ASC, id ASC
			LIMIT $1 OFFSET $2
		`, tasksTable),
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []entities.Task
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *Repository) GetTaskByID(ctx context.Context, taskID string) (entities.Task, error) {
	row := r.DB.QueryRowContext(
		ctx,
		fmt.Sprintf(`
			SELECT id, message, binary_code, created_at, updated_at
			FROM %s
			WHERE id = $1
		`, tasksTable),
		taskID,
	)
	return scanTaskRow(row)
}

func (r *Repository) MoveTaskToProcessing(ctx context.Context, createTaskDto usecases.CreateTaskDTO) (entities.Task, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return entities.Task{}, err
	}

	now := time.Now().UTC()
	nowString := now.Format(timeLayout)
	task := entities.Task{
		ID:         createTaskDto.ID,
		Message:    createTaskDto.Message,
		BinaryCode: "",
		CreatedAt:  nowString,
		UpdatedAt:  nowString,
	}

	if err := insertTask(tx, task, now); err != nil {
		_ = tx.Rollback()
		return entities.Task{}, err
	}

	if err := insertEvent(tx, task, eventUnprocessed, now); err != nil {
		_ = tx.Rollback()
		return entities.Task{}, err
	}

	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return entities.Task{}, err
	}

	return task, nil
}

func (r *Repository) FinishProcessing(ctx context.Context, dto usecases.FinishProcessingDTO) (entities.Task, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return entities.Task{}, err
	}

	now := time.Now().UTC()
	row := tx.QueryRow(
		fmt.Sprintf(`
			UPDATE %s
			SET binary_code = $1, updated_at = $2
			WHERE id = $3
			RETURNING id, message, binary_code, created_at, updated_at
		`, tasksTable),
		dto.BinaryCode,
		now,
		dto.ID,
	)

	task, err := scanTaskRow(row)
	if err != nil {
		_ = tx.Rollback()
		return entities.Task{}, err
	}

	if err := insertEvent(tx, task, eventProcessed, now); err != nil {
		_ = tx.Rollback()
		return entities.Task{}, err
	}

	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return entities.Task{}, err
	}

	return task, nil
}

func insertTask(tx *sql.Tx, task entities.Task, now time.Time) error {
	_, err := tx.Exec(
		fmt.Sprintf(`
			INSERT INTO %s (id, message, binary_code, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5)
		`, tasksTable),
		task.ID,
		task.Message,
		task.BinaryCode,
		now,
		now,
	)
	return err
}

func insertEvent(tx *sql.Tx, task entities.Task, status string, now time.Time) error {
	_, err := tx.Exec(
		fmt.Sprintf(`
			INSERT INTO %s (task_id, status, message, binary_code, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, outboxTable),
		task.ID,
		status,
		task.Message,
		task.BinaryCode,
		now,
		now,
	)
	return err
}

func scanTask(rows *sql.Rows) (entities.Task, error) {
	var task entities.Task
	var createdAt time.Time
	var updatedAt time.Time
	if err := rows.Scan(&task.ID, &task.Message, &task.BinaryCode, &createdAt, &updatedAt); err != nil {
		return entities.Task{}, err
	}
	task.CreatedAt = createdAt.Format(timeLayout)
	task.UpdatedAt = updatedAt.Format(timeLayout)
	return task, nil
}

func scanTaskRow(row *sql.Row) (entities.Task, error) {
	var task entities.Task
	var createdAt time.Time
	var updatedAt time.Time
	if err := row.Scan(&task.ID, &task.Message, &task.BinaryCode, &createdAt, &updatedAt); err != nil {
		return entities.Task{}, err
	}
	task.CreatedAt = createdAt.Format(timeLayout)
	task.UpdatedAt = updatedAt.Format(timeLayout)
	return task, nil
}
