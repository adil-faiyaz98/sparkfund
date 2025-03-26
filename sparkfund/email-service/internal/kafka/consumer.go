package kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap" // Import zap
)

// MessageHandler is a function type that processes Kafka messages
type MessageHandler func(context.Context, *sarama.ConsumerMessage) error

// Consumer represents a Kafka consumer
type Consumer struct {
	consumer        sarama.ConsumerGroup
	topics          []string
	handler         MessageHandler
	tracer          trace.Tracer
	wg              sync.WaitGroup
	ctx             context.Context
	cancel          context.CancelFunc
	consumerGroupID string      // Add consumer group ID
	logger          *zap.Logger // Add logger
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(brokers []string, consumerGroupID string, topics []string, handler MessageHandler, logger *zap.Logger) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = true // Enable error channel
	config.Version = sarama.V2_8_0_0     // Or your broker version

	consumerGroup, err := sarama.NewConsumerGroup(brokers, consumerGroupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Consumer{
		consumer:        consumerGroup,
		topics:          topics,
		handler:         handler,
		tracer:          otel.Tracer("kafka-consumer"),
		ctx:             ctx,
		cancel:          cancel,
		consumerGroupID: consumerGroupID,
		logger:          logger,
	}, nil
}

// Start starts consuming messages from Kafka topics
func (c *Consumer) Start() error {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		handler := consumerGroupHandler{
			consumer: c,
		}
		for {
			// `Consume` should be called inside a for loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			err := c.consumer.Consume(c.ctx, c.topics, handler)
			if err != nil {
				c.logger.Error("Error from consumer: %v", zap.Error(err))
			}
			// check if context was cancelled, signaling that the consumer should stop
			if c.ctx.Err() != nil {
				return
			}
			handler.Ready = make(chan bool)
		}
	}()

	return nil
}

// consumerGroupHandler implements the sarama.ConsumerGroupHandler interface
type consumerGroupHandler struct {
	consumer *Consumer
	Ready    chan bool
}

func (consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for message := range claim.Messages() {
		h.consumer.logger.Info("Message claimed: value = %s, timestamp = %v, topic = %s",
			zap.String("value", string(message.Value)),
			zap.Time("timestamp", message.Timestamp),
			zap.String("topic", message.Topic))

		carrier := propagation.MapCarrier{}
		for k, v := range message.Headers {
			carrier[string(k.Key)] = string(v.Value)
		}
		ctx := otel.GetTextMapPropagator().Extract(h.consumer.ctx, carrier)

		ctx, span := h.consumer.tracer.Start(ctx, "ProcessMessage",
			trace.WithAttributes(
				attribute.String("topic", message.Topic),
				attribute.Int64("partition", int64(message.Partition)),
				attribute.Int64("offset", message.Offset),
				attribute.String("key", string(message.Key)),
			))

		err := h.consumer.handler(ctx, message)
		if err != nil {
			span.RecordError(err)
			h.consumer.logger.Error("Failed to process message", zap.Error(err))
			// Decide if you want to commit the offset or not based on the error.
			// If you commit, the message will be considered processed and won't be retried.
			// If you don't commit, the message will be redelivered on the next rebalance.
		} else {
			session.MarkMessage(message, "") // Mark message as processed
		}

		span.End()
		session.Commit()
	}

	return nil
}

// Stop stops consuming messages and closes all connections
func (c *Consumer) Stop() error {
	c.cancel()
	c.wg.Wait()

	if err := c.consumer.Close(); err != nil {
		return fmt.Errorf("failed to close consumer: %w", err)
	}

	return nil
}
