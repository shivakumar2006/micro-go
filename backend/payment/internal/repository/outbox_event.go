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

func (r *OutboxEventRepository) SaveEvent(ctx context.Context, event *models.OutboxEvent) error {
	query := `INSERT INTO outbox_events (aggregate_type, aggregate_id, event_type, payload, created_at) 
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.Db.ExecContext(ctx, query, event.AggregateType, event.AggregateID, event.EventType, event.Payload, event.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to save outbox event: %w", err)
	}

	return nil
}
