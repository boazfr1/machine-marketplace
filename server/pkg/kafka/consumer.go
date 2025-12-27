package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type (
	Consumer struct {
		reader *kafka.Reader
	}
)

func NewConsumer(address, groupID string) (*Consumer, error) {
	brokerAddress := address
	if brokerAddress == "" {
		brokerAddress = "localhost:9092"
	}

	if groupID == "" {
		groupID = "machine-purchases-consumer-group"
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{brokerAddress},
		Topic:          "machine-purchases",
		GroupID:        groupID,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: 1 * time.Second,
		StartOffset:    kafka.LastOffset,
	})

	return &Consumer{
		reader: reader,
	}, nil
}

func (c *Consumer) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return c.reader.ReadMessage(ctx)
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}