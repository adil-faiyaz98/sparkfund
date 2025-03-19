package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// Producer handles Kafka message publishing
type Producer struct {
	logger   *zap.Logger
	producer sarama.SyncProducer
	config   *sarama.Config
}

// NewProducer creates a new Kafka producer
func NewProducer(cfg *Config, logger *zap.Logger) (*Producer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Return.Errors = true
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Retry.Max = 5
	saramaConfig.Producer.Timeout = 10 * time.Second

	// Set up authentication if credentials are provided
	if cfg.Username != "" && cfg.Password != "" {
		saramaConfig.Net.SASL.Enable = true
		saramaConfig.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		saramaConfig.Net.SASL.User = cfg.Username
		saramaConfig.Net.SASL.Password = cfg.Password
	}

	producer, err := sarama.NewSyncProducer(cfg.Brokers, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %v", err)
	}

	return &Producer{
		logger:   logger,
		producer: producer,
		config:   saramaConfig,
	}, nil
}

// Publish sends a message to a Kafka topic
func (p *Producer) Publish(ctx context.Context, topic string, key string, value interface{}) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(jsonValue),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		p.logger.Error("Failed to send message",
			zap.String("topic", topic),
			zap.String("key", key),
			zap.Error(err))
		return fmt.Errorf("failed to send message: %v", err)
	}

	p.logger.Debug("Message sent successfully",
		zap.String("topic", topic),
		zap.String("key", key),
		zap.Int32("partition", partition),
		zap.Int64("offset", offset))

	return nil
}

// Close closes the producer
func (p *Producer) Close() error {
	if err := p.producer.Close(); err != nil {
		return fmt.Errorf("failed to close producer: %v", err)
	}
	return nil
}
