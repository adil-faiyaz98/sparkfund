package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Producer represents a Kafka producer
type Producer struct {
	producer sarama.AsyncProducer
	tracer   trace.Tracer
	topic    string
	logger   *zap.Logger
}

// NewProducer creates a new Kafka producer
func NewProducer(brokers []string, topic string, logger *zap.Logger) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Or Local if you want faster but less durable
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true                  // Required to check for success
	config.Producer.Compression = sarama.CompressionGZIP     //Example
	config.Producer.Flush.Frequency = 500 * time.Millisecond // Example - how often to flush messages

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	// Handle errors on the Errors channel
	go func() {
		for err := range producer.Errors() {
			logger.Error("Failed to write message to Kafka", zap.Error(err.Err))
		}
	}()

	return &Producer{
		producer: producer,
		tracer:   otel.Tracer("kafka-producer"),
		topic:    topic,
		logger:   logger,
	}, nil
}

// SendMessage sends a message to a Kafka topic
func (p *Producer) SendMessage(ctx context.Context, key string, value []byte) error {
	ctx, span := p.tracer.Start(ctx, "SendMessage",
		trace.WithAttributes(
			attribute.String("topic", p.topic),
			attribute.String("key", key),
		))
	defer span.End()

	msg := &sarama.ProducerMessage{
		Topic:     p.topic,
		Key:       sarama.StringEncoder(key),
		Value:     sarama.ByteEncoder(value),
		Timestamp: time.Now(),
	}

	// Propagate tracing context
	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	for k, v := range carrier {
		msg.Headers = append(msg.Headers, sarama.RecordHeader{Key: []byte(k), Value: []byte(v)})
	}

	p.producer.Input() <- msg // Send message to the input channel

	return nil
}

// Close closes the producer
func (p *Producer) Close() error {
	p.logger.Info("Closing Kafka producer...")
	defer p.logger.Info("Kafka producer closed.")

	if err := p.producer.Close(); err != nil {
		return fmt.Errorf("failed to close producer: %w", err)
	}
	return nil
}
