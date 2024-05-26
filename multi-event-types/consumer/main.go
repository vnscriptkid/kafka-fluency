package main

import (
	"encoding/json"
	"fmt"
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

type OrderShippedPayload struct {
	OrderID string `json:"order_id"`
	Carrier string `json:"carrier"`
}

type OrderDeliveredPayload struct {
	OrderID string `json:"order_id"`
}

func handleOrderCreated(payload json.RawMessage) {
	var data OrderCreatedPayload
	if err := json.Unmarshal(payload, &data); err != nil {
		fmt.Printf("Failed to unmarshal OrderCreatedPayload: %v\n", err)
		return
	}
	fmt.Printf("Order Created: ID=%s, Amount=%.2f\n", data.OrderID, data.Amount)
}

func handleOrderShipped(payload json.RawMessage) {
	var data OrderShippedPayload
	if err := json.Unmarshal(payload, &data); err != nil {
		fmt.Printf("Failed to unmarshal OrderShippedPayload: %v\n", err)
		return
	}
	fmt.Printf("Order Shipped: ID=%s, Carrier=%s\n", data.OrderID, data.Carrier)
}

func handleOrderDelivered(payload json.RawMessage) {
	var data OrderDeliveredPayload
	if err := json.Unmarshal(payload, &data); err != nil {
		fmt.Printf("Failed to unmarshal OrderDeliveredPayload: %v\n", err)
		return
	}
	fmt.Printf("Order Delivered: ID=%s\n", data.OrderID)
}

func main() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "ecommerce-group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}
	defer c.Close()

	topic := "ecommerce-events"
	c.SubscribeTopics([]string{topic}, nil)

	run := true
	for run {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			var event Event
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				fmt.Printf("Failed to unmarshal event: %v\n", err)
				continue
			}

			switch event.Type {
			case "OrderCreated":
				handleOrderCreated(event.Payload)
			case "OrderShipped":
				handleOrderShipped(event.Payload)
			case "OrderDelivered":
				handleOrderDelivered(event.Payload)
			default:
				fmt.Printf("Unknown event type: %s\n", event.Type)
			}

			// Simulate variable processing time
			time.Sleep(time.Duration(100+msg.TopicPartition.Partition*100) * time.Millisecond)
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
