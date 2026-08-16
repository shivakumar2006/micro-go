package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &Producer{writer: writer}
}

func (p *Producer) Publish(ctx context.Context, event PaymentSuccessEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal the event : %w", err)
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(strconv.FormatInt(event.OrderID, 10)),
		Value: data,
	})

	if err != nil {
		return fmt.Errorf("failed to write message : %w", err)
	}

	return nil
}

func (p *Producer) PublishPayload(ctx context.Context, key string, payload []byte) error {
	err := p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: payload,
	})
	if err != nil {
		return fmt.Errorf("failed to write the message : %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
