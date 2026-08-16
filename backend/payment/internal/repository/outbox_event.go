package repository

import (
	"context"
	"database/sql"
	"fmt"
	"payment/internal/models"
)

type OutboxEventRepository struct {
	Db *sql.DB
}

func NewOutboxEventRepository(db *sql.DB) *OutboxEventRepository {
	return &OutboxEventRepository{
		Db: db,
	}
}

func (r *OutboxEventRepository) SaveEventTx(ctx context.Context, tx *sql.Tx, event *models.OutboxEvent) error {
	query := `INSERT INTO outbox_events (aggregate_type, aggregate_id, event_type, payload, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := tx.ExecContext(ctx, query, event.AggregateType, event.AggregateID, event.EventType, event.Payload, event.Status, event.CreatedAt, event.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to save outbox event: %w", err)
	}

	return nil
}

func (r *OutboxEventRepository) GetPendingEvents(ctx context.Context, status string) ([]models.OutboxEvent, error) {
	query := `SELECT id,aggregate_type,aggregate_id,event_type,payload,status,created_at,updated_at 
				FROM outbox_events 
				WHERE status = $1 
				ORDER BY created_at ASC
			`

	rows, err := r.Db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending status: %w", err)
	}

	defer rows.Close()

	var events []models.OutboxEvent
	for rows.Next() {
		var event models.OutboxEvent
		err := rows.Scan(&event.ID, &event.AggregateType, &event.AggregateID, &event.EventType, &event.Payload, &event.Status, &event.CreatedAt, &event.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows : %w", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating through rows : %w", err)
	}

	return events, nil
}

func (r *OutboxEventRepository) MarkAsPublished(ctx context.Context, eventID int) error {
	query := `UPDATE outbox_events
			SET 
				status = $1,
				updated_at = NOW()
			WHERE id = $2
	`

	_, err := r.Db.ExecContext(ctx, query, models.OutboxStatusPublished, eventID)
	if err != nil {
		return fmt.Errorf("failed to mark event as published: %w", err)
	}

	return nil
}
