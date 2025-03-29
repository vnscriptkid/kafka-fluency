package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/IBM/sarama"
)

// Order represents a concrete order message.
type Order struct {
	OrderID  string  `json:"order_id"`
	Customer string  `json:"customer"`
	Amount   float64 `json:"amount"`
	Retry    int     `json:"retry"`
}

const (
	ordersTopic      = "orders"
	ordersRetryTopic = "orders_retry"
	ordersDLQTopic   = "orders_dlq"
	maxRetries       = 3
)

func main() {
	// Set up the Sarama configuration.
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Producer.Return.Successes = true
	// Use an appropriate Kafka version.
	config.Version = sarama.V2_1_0_0

	brokers := []string{"localhost:9092"}

	// Create a synchronous producer.
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Error creating producer: %v", err)
	}
	defer producer.Close()

	// Start the main consumer and the retry consumer in separate goroutines.
	go consumeOrders(brokers, ordersTopic, producer)
	go consumeRetryOrders(brokers, ordersRetryTopic, producer)

	// Start producing order messages.
	go produceOrders(producer, ordersTopic)

	// Wait for a termination signal (e.g., Ctrl+C).
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	<-signals
	fmt.Println("Shutting down gracefully...")
}

// produceOrders creates sample order messages and sends them to the orders topic.
func produceOrders(producer sarama.SyncProducer, topic string) {
	orderID := 1
	customers := []string{"Alice", "Bob", "Charlie", "Diana"}
	for {
		order := Order{
			OrderID:  fmt.Sprintf("%d", orderID),
			Customer: customers[rand.Intn(len(customers))],
			Amount:   float64(rand.Intn(500)) + 1, // Amount between 1 and 500.
			Retry:    0,
		}
		orderJSON, err := json.Marshal(order)
		if err != nil {
			log.Printf("Error marshaling order: %v", err)
			continue
		}
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(orderJSON),
		}
		_, _, err = producer.SendMessage(msg)
		if err != nil {
			log.Printf("Error sending order: %v", err)
		} else {
			log.Printf("Produced order: %s", string(orderJSON))
		}
		orderID++
		time.Sleep(1 * time.Second)
	}
}

// consumeOrders consumes order messages from the main orders topic.
func consumeOrders(brokers []string, topic string, producer sarama.SyncProducer) {
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		log.Fatalf("Error creating consumer: %v", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error starting partition consumer: %v", err)
	}
	defer partitionConsumer.Close()

	for msg := range partitionConsumer.Messages() {
		var order Order
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Printf("Error unmarshalling order: %v", err)
			continue
		}
		log.Printf("[Orders Consumer] Received Order: %+v", order)
		if err := processOrder(order); err != nil {
			if order.Retry >= maxRetries {
				log.Printf("[Orders Consumer] Order %s reached max retries. Sending to DLQ.", order.OrderID)
				sendToTopic(producer, ordersDLQTopic, msg.Value)
			} else {
				order.Retry++
				updatedOrder, err := json.Marshal(order)
				if err != nil {
					log.Printf("Error marshaling updated order: %v", err)
					continue
				}
				log.Printf("[Orders Consumer] Processing failed for Order %s. Sending to retry topic (Retry count: %d)", order.OrderID, order.Retry)
				sendToTopic(producer, ordersRetryTopic, updatedOrder)
			}
		} else {
			log.Printf("[Orders Consumer] Successfully processed Order %s", order.OrderID)
		}
	}
}

// consumeRetryOrders consumes order messages from the orders_retry topic.
func consumeRetryOrders(brokers []string, topic string, producer sarama.SyncProducer) {
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		log.Fatalf("Error creating retry consumer: %v", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error starting partition consumer for retry topic: %v", err)
	}
	defer partitionConsumer.Close()

	for msg := range partitionConsumer.Messages() {
		var order Order
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Printf("Error unmarshalling order: %v", err)
			continue
		}
		log.Printf("[Retry Consumer] Received Order: %+v", order)
		if err := processOrder(order); err != nil {
			if order.Retry >= maxRetries {
				log.Printf("[Retry Consumer] Order %s reached max retries. Sending to DLQ.", order.OrderID)
				sendToTopic(producer, ordersDLQTopic, msg.Value)
			} else {
				order.Retry++
				updatedOrder, err := json.Marshal(order)
				if err != nil {
					log.Printf("Error marshaling updated order: %v", err)
					continue
				}
				log.Printf("[Retry Consumer] Processing failed for Order %s. Retrying (Retry count: %d)", order.OrderID, order.Retry)
				// Optionally add a delay before requeueing.
				time.Sleep(100 * time.Millisecond)
				sendToTopic(producer, ordersRetryTopic, updatedOrder)
			}
		} else {
			log.Printf("[Retry Consumer] Successfully processed Order %s", order.OrderID)
		}
	}
}

// processOrder simulates processing an order. Here, a 50% chance of failure is simulated.
func processOrder(order Order) error {
	// Simulate random processing failure.
	if rand.Intn(2) == 0 {
		return fmt.Errorf("simulated order processing failure for order %s", order.OrderID)
	}
	return nil
}

// sendToTopic publishes a message with the given value to the specified topic.
func sendToTopic(producer sarama.SyncProducer, topic string, value []byte) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(value),
	}
	_, _, err := producer.SendMessage(msg)
	if err != nil {
		log.Printf("Error sending message to topic %s: %v", topic, err)
	}
}
