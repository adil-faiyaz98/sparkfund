package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Producer represents a Kafka producer
type Producer struct {
	producer sarama.SyncProducer
	tracer   trace.Tracer
}

// NewProducer creates a new Kafka producer
func NewProducer(brokers []string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &Producer{
		producer: producer,
		tracer:   otel.Tracer("kafka-producer"),
	}, nil
}

// SendMessage sends a message to a Kafka topic
func (p *Producer) SendMessage(ctx context.Context, topic string, key string, value []byte) error {
	ctx, span := p.tracer.Start(ctx, "SendMessage",
		trace.WithAttributes(
			attribute.String("topic", topic),
			attribute.String("key", key),
		))
	defer span.End()

	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.StringEncoder(key),
		Value:     sarama.ByteEncoder(value),
		Timestamp: time.Now(),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to send message: %w", err)
	}

	span.SetAttributes(
		attribute.Int64("partition", int64(partition)),
		attribute.Int64("offset", offset),
	)

	return nil
}

// Close closes the producer
func (p *Producer) Close() error {
	if err := p.producer.Close(); err != nil {
		return fmt.Errorf("failed to close producer: %w", err)
	}
	return nil
}
