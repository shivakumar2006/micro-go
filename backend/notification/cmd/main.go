package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"notification/internal/client"
	"notification/internal/config"
	"notification/internal/kafka"
	"notification/internal/service"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.LoadConfig()
	if cfg == nil {
		log.Fatal("Config not loaded")
	}

	slog.Info("config loaded successfully")

	// client
	emailClient := client.NewEmailClient(cfg.Brevo.SMTPHost, cfg.Brevo.SMTPPort, cfg.Brevo.SMTPUser, cfg.Brevo.SMTPPassword, cfg.Brevo.SenderEmail, cfg.Brevo.SenderName)

	notifyservice := service.NewNotificationService(emailClient)

	// server
	server := &http.Server{
		Addr:         cfg.Server.Addr,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	kafkaCtx, kafkaCancel := context.WithCancel(context.Background())
	defer kafkaCancel()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		backoff := 1 * time.Second
		maxBackoff := 10 * time.Second

		for {
			if kafkaCtx.Err() != nil {
				return
			}

			consumer := kafka.NewConsumer([]string{cfg.Kafka.Addr}, cfg.Kafka.Topic, cfg.Kafka.GroupID)

			slog.Info("notification consumer started", slog.String("topic", cfg.Kafka.Topic), slog.String("groupId", cfg.Kafka.GroupID))

			err := consumer.Start(kafkaCtx, func(event kafka.PaymentSuccessEvent) error {
				slog.Info("payment success event received", slog.Any("event", event))

				return notifyservice.HandlePaymentSuccess(kafkaCtx, event)
			})

			consumer.Close()

			if kafkaCtx.Err() != nil {
				return
			}

			slog.Error("kafka consumer stopped", slog.String("error", err.Error()))

			timer := time.NewTimer(backoff)

			select {
			case <-timer.C:
			case <-kafkaCtx.Done():
				timer.Stop()
				return
			}

			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}
	}()

	go func() {
		slog.Info("notification http server listening", slog.String("address", cfg.Server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start notification service")
			log.Fatal("failed to start notification service", "error", err)
		}
	}()

	<-quit
	kafkaCancel()

	slog.Info("inventory service is shutting down gracefully")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("inventory service forced to shutdown", "error", err)
		log.Fatal("inventory service forced to shutdown", "error", err)
	}

	slog.Info("notification service stopped")
}
