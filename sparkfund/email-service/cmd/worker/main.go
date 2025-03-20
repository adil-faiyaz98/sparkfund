package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/sparkfund/email-service/internal/kafka"
	"go.uber.org/zap"
)

type emailMessage struct {
	LogID    string   `json:"log_id"`
	To       []string `json:"to"`
	Subject  string   `json:"subject"`
	Body     string   `json:"body"`
	Template string   `json:"template,omitempty"`
	Data     any      `json:"data,omitempty"`
}

func main() {
	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Initialize Kafka consumer
	consumer, err := kafka.NewConsumer([]string{"localhost:9092"}, []string{"email-queue"}, func(ctx context.Context, msg *sarama.ConsumerMessage) error {
		var emailMsg emailMessage
		if err := json.Unmarshal(msg.Value, &emailMsg); err != nil {
			logger.Error("Failed to unmarshal message",
				zap.Error(err))
			return err
		}

		// Process the email message
		// TODO: Implement email sending logic

		return nil
	})
	if err != nil {
		logger.Fatal("Failed to create Kafka consumer",
			zap.Error(err))
	}
	defer consumer.Stop()

	// Create context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Info("Received shutdown signal")
		cancel()
	}()

	// Start consuming messages
	if err := consumer.Start(); err != nil {
		logger.Fatal("Failed to start consumer",
			zap.Error(err))
	}

	<-ctx.Done()
	logger.Info("Shutting down worker")
}
