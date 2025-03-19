package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// Consumer handles Kafka message consumption
type Consumer struct {
	logger   *zap.Logger
	consumer sarama.ConsumerGroup
	config   *sarama.Config
}

// MessageHandler is a function type for handling consumed messages
type MessageHandler func(ctx context.Context, topic string, key string, value []byte) error

// consumerGroupHandler implements sarama.ConsumerGroupHandler
type consumerGroupHandler struct {
	logger  *zap.Logger
	handler MessageHandler
}

func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if err := h.handler(session.Context(), message.Topic, string(message.Key), message.Value); err != nil {
				h.logger.Error("Failed to handle message",
					zap.String("topic", message.Topic),
					zap.String("key", string(message.Key)),
					zap.Error(err))
			}
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(cfg *Config, logger *zap.Logger) (*Consumer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Return.Errors = true
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	saramaConfig.Consumer.Group.Session.Timeout = 20 * time.Second
	saramaConfig.Consumer.Group.Heartbeat.Interval = 6 * time.Second

	// Set up authentication if credentials are provided
	if cfg.Username != "" && cfg.Password != "" {
		saramaConfig.Net.SASL.Enable = true
		saramaConfig.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		saramaConfig.Net.SASL.User = cfg.Username
		saramaConfig.Net.SASL.Password = cfg.Password
	}

	consumer, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %v", err)
	}

	return &Consumer{
		logger:   logger,
		consumer: consumer,
		config:   saramaConfig,
	}, nil
}

// ConsumeMessages starts consuming messages from specified topics
func (c *Consumer) ConsumeMessages(ctx context.Context, topics []string, handler MessageHandler) error {
	consumerHandler := &consumerGroupHandler{
		logger:  c.logger,
		handler: handler,
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			err := c.consumer.Consume(ctx, topics, consumerHandler)
			if err != nil {
				c.logger.Error("Error from consumer",
					zap.Error(err))
			}
		}
	}
}

// Close closes the consumer
func (c *Consumer) Close() error {
	if err := c.consumer.Close(); err != nil {
		return fmt.Errorf("failed to close consumer: %v", err)
	}
	return nil
}
