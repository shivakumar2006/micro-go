package worker

import (
	"context"
	"fmt"
	"log/slog"
	"payment/internal/db/storage"
	"payment/internal/kafka"
	"payment/internal/models"
	"strconv"
)

type OutboxWorker struct {
	outbox storage.OutboxStorage
	kafka  *kafka.Producer
}

func NewOutboxWorker(outbox storage.OutboxStorage, kafka *kafka.Producer) *OutboxWorker {
	return &OutboxWorker{
		outbox: outbox,
		kafka:  kafka,
	}
}

func (w *OutboxWorker) ProcessPendingEvents(ctx context.Context) error {
	pendingEvents, err := w.outbox.GetPendingEvents(ctx, models.OutboxStatusPending)
	if err != nil {
		return fmt.Errorf("failed to get pending events : %w", err)
	}
	slog.Info(
		"pending outbox events fetched",
		slog.Int("count", len(pendingEvents)),
	)

	for _, event := range pendingEvents {
		key := strconv.FormatInt(event.AggregateID, 10)

		if err := w.kafka.PublishPayload(ctx, key, event.Payload); err != nil {
			slog.Error("failed to publish outbox event", slog.Int("event_id", event.ID), slog.String("error", err.Error()))
			continue
		}

		if err := w.outbox.MarkAsPublished(ctx, event.ID); err != nil {
			slog.Error("failed to mark outbox event as published", slog.Int("event_id", event.ID), slog.String("error", err.Error()))
			continue
		}
	}

	return nil
}
