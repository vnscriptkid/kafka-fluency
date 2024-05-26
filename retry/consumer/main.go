package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type OrderCreatedPayload struct {
	OrderID string  `json:"order_id"`
	Amount  float64 `json:"amount"`
}

const maxRetries = 3

var dlqTopic = "ecommerce-events-dlq"

func handleOrderCreated(payload json.RawMessage) error {
	var data OrderCreatedPayload
	if err := json.Unmarshal(payload, &data); err != nil {
		return fmt.Errorf("failed to unmarshal OrderCreatedPayload: %v", err)
	}
	fmt.Printf("Order Created: ID=%s, Amount=%.2f\n", data.OrderID, data.Amount)
	// Simulate processing error for demonstration
	if data.OrderID == "12345" {
		return fmt.Errorf("simulated processing error")
	}
	return nil
}

func processMessage(p *kafka.Producer, msg *kafka.Message) error {
	var event Event
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %v", err)
	}

	switch event.Type {
	case "OrderCreated":
		return handleOrderCreated(event.Payload)
	default:
		return fmt.Errorf("unknown event type: %s", event.Type)
	}
}

func sendToDLQ(p *kafka.Producer, msg *kafka.Message) {
	dlqMsg := kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &dlqTopic, Partition: kafka.PartitionAny},
		Value:          msg.Value,
		Key:            msg.Key,
	}
	err := p.Produce(&dlqMsg, nil)

	if err != nil {
		log.Fatalf("Failed to produce message to DLQ: %v", err)
		return
	}

	fmt.Printf("Sent to DLQ: %s\n", string(msg.Value))
}

func main() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "ecommerce-group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer c.Close()

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer p.Close()

	topic := "ecommerce-events"
	c.SubscribeTopics([]string{topic}, nil)

	run := true
	for run {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			for retry := 0; retry < maxRetries; retry++ {
				err = processMessage(p, msg)
				if err == nil {
					break
				}
				fmt.Printf("Processing failed: %v, retrying (%d/%d)\n", err, retry+1, maxRetries)
				time.Sleep(time.Duration(retry+1) * time.Second) // Exponential backoff
			}

			if err != nil {
				fmt.Printf("Processing failed after %d retries, sending to DLQ: %v\n", maxRetries, err)
				sendToDLQ(p, msg)
			}
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
