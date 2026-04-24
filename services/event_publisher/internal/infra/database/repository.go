package database

import (
	"database/sql"
	"fmt"

	"github.com/douglasvolcato/binary-code-processor/event_publisher/internal/entities"
)

const outboxTable = "task_outbox_events"

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetUnprocessedEvents(limit int, offset int) ([]entities.Event, error) {
	return r.getEventsByStatus("unprocessed", limit, offset)
}

func (r *Repository) GetProcessedEvents(limit int, offset int) ([]entities.Event, error) {
	return r.getEventsByStatus("processed", limit, offset)
}

func (r *Repository) getEventsByStatus(status string, limit int, offset int) ([]entities.Event, error) {
	rows, err := r.DB.Query(
		fmt.Sprintf(`
			SELECT task_id, status, binary_code
			FROM %s
			WHERE status = $1
			ORDER BY created_at ASC
			LIMIT $2 OFFSET $3
		`, outboxTable),
		status,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []entities.Event
	for rows.Next() {
		var event entities.Event
		if err := rows.Scan(&event.ID, &event.Status, &event.BinaryCode); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func (r *Repository) DeleteEventByID(id string) error {
	_, err := r.DB.Exec(
		fmt.Sprintf(`
			DELETE FROM %s
			WHERE task_id = $1
		`, outboxTable),
		id,
	)
	return err
}
