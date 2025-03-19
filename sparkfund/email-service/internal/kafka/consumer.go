package kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// MessageHandler is a function type that processes Kafka messages
type MessageHandler func(context.Context, *sarama.ConsumerMessage) error

// Consumer represents a Kafka consumer
type Consumer struct {
	consumer sarama.Consumer
	topics   []string
	handler  MessageHandler
	tracer   trace.Tracer
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(brokers []string, topics []string, handler MessageHandler) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Consumer{
		consumer: consumer,
		topics:   topics,
		handler:  handler,
		tracer:   otel.Tracer("kafka-consumer"),
		ctx:      ctx,
		cancel:   cancel,
	}, nil
}

// Start starts consuming messages from Kafka topics
func (c *Consumer) Start() error {
	for _, topic := range c.topics {
		partitions, err := c.consumer.Partitions(topic)
		if err != nil {
			return fmt.Errorf("failed to get partitions for topic %s: %w", topic, err)
		}

		for _, partition := range partitions {
			pc, err := c.consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
			if err != nil {
				return fmt.Errorf("failed to start consuming partition %d of topic %s: %w", partition, topic, err)
			}

			c.wg.Add(1)
			go c.consumePartition(pc)
		}
	}

	return nil
}

func (c *Consumer) consumePartition(pc sarama.PartitionConsumer) {
	defer c.wg.Done()
	defer pc.Close()

	for {
		select {
		case msg := <-pc.Messages():
			if msg == nil {
				return
			}

			ctx, span := c.tracer.Start(c.ctx, "ProcessMessage",
				trace.WithAttributes(
					attribute.String("topic", msg.Topic),
					attribute.Int64("partition", int64(msg.Partition)),
					attribute.Int64("offset", msg.Offset),
					attribute.String("key", string(msg.Key)),
				))

			if err := c.handler(ctx, msg); err != nil {
				span.RecordError(err)
			}

			span.End()

		case <-c.ctx.Done():
			return
		}
	}
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
