package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

type PurchaseEvent struct {
	MachineID      int32     `json:"machine_id"`
	BuyerID        int32     `json:"buyer_id"`
	DealExpiration time.Time `json:"deal_expiration"`
	PurchaseTime   time.Time `json:"purchase_time"`
	MachineName    string    `json:"machine_name"`
}

type Producer struct {
	writer *kafka.Writer
}

func New() (*Producer, error) {
	brokerAddress := os.Getenv("KAFKA_BROKER")
	if brokerAddress == "" {
		brokerAddress = "localhost:9092"
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokerAddress),
		Topic:        "machine-purchases",
		Balancer:     &kafka.LeastBytes{},
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	return &Producer{
		writer: writer,
	}, nil
}

func (p *Producer) SendPurchaseEvent(ctx context.Context, event PurchaseEvent) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal purchase event: %w", err)
	}

	message := kafka.Message{
		Key:   []byte(fmt.Sprintf("machine-%d", event.MachineID)),
		Value: eventBytes,
		Time:  time.Now(),
	}

	err = p.writer.WriteMessages(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to write message to kafka: %w", err)
	}

	log.Printf("Successfully sent purchase event to Kafka: MachineID=%d, BuyerID=%d",
		event.MachineID, event.BuyerID)

	return nil
}

func (p *Producer) Close() error {
	if p.writer != nil {
		return p.writer.Close()
	}
	return nil
}
