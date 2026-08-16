package worker

import (
	"payment/internal/db/storage"
	"payment/internal/kafka"
)

type OutboxWorker struct {
	Repo  storage.Storage
	kafka kafka.Producer
}

func NewOutboxWorker(repo storage.Storage, kafka kafka.Producer) *OutboxWorker {
	return &OutboxWorker{
		Repo:  repo,
		kafka: kafka,
	}
}
