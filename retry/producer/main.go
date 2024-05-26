package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Event struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type OrderCreatedPayload struct {
	OrderID string  `json:"order_id"`
	Amount  float64 `json:"amount"`
}

func main() {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		panic(err)
	}
	defer p.Close()

	topic := "ecommerce-events"

	events := []Event{
		{Type: "OrderCreated", Payload: OrderCreatedPayload{OrderID: "12340", Amount: 10}},
		{Type: "OrderCreated", Payload: OrderCreatedPayload{OrderID: "12341", Amount: 30}},
		{Type: "OrderCreated", Payload: OrderCreatedPayload{OrderID: "12345", Amount: 99.99}},
	}

	for _, event := range events {
		eventBytes, _ := json.Marshal(event)
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          eventBytes,
		}, nil)
		fmt.Printf("Sent: %s event\n", event.Type)
		time.Sleep(500 * time.Millisecond) // Small delay between messages
	}

	p.Flush(15 * 1000)
}
