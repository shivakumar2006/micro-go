package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})

	return &Consumer{reader: reader}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

func (c *Consumer) ReadMessage(ctx context.Context) (*PaymentSuccessEvent, error) {
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}

	var event PaymentSuccessEvent

	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return nil, err
	}

	return &event, nil
}

func (c *Consumer) FetchMessage(ctx context.Context) (kafka.Message, error) {
	msg, err := c.reader.FetchMessage(ctx)
	if err != nil {
		return kafka.Message{}, err
	}

	return msg, nil
}

func (c *Consumer) Start(ctx context.Context, handler func(PaymentSuccessEvent) error) error {
	for {
		fetch, err := c.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			slog.Error("kafka fetch failed", slog.String("error", err.Error()))
			continue
		}

		var event PaymentSuccessEvent

		if err := json.Unmarshal(fetch.Value, &event); err != nil {
			return fmt.Errorf("failed to unmarshal event : %w", err)
		}

		if err := handler(event); err != nil {
			continue
		}

		if err := c.reader.CommitMessages(ctx, fetch); err != nil {
			return fmt.Errorf("failed to commit message: %w", err)
		}
		slog.Info("messages committed successfully")
	}
}
