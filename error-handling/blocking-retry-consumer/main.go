package main

import (
	"fmt"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	config := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		fmt.Printf("Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	defer consumer.Close()

	topics := []string{"myTopic"}
	err = consumer.SubscribeTopics(topics, nil)
	if err != nil {
		fmt.Printf("Failed to subscribe to topics: %s\n", err)
		os.Exit(1)
	}

	for {
		fmt.Printf("Waiting for messages...\n")
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			processMessageWithRetry(consumer, msg, 3, 10*time.Second)
		} else {
			// Error is either due to poll timeout or a Kafka error
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}

func processMessageWithRetry(consumer *kafka.Consumer, msg *kafka.Message, retries int, waitTime time.Duration) {
	for i := 0; i < retries; i++ {
		err := processMessage(msg)
		if err == nil {
			// Commit the message offset after successful processing
			_, err = consumer.CommitMessage(msg)
			if err != nil {
				fmt.Printf("Failed to commit message: %s\n", err)
			}
			return
		}

		fmt.Printf("Failed to process message: %s. Retrying %d/%d...\n", err, i+1, retries)
		time.Sleep(waitTime)
	}

	// After exhausting retries, log the failure and move on
	fmt.Printf("Failed to process message after %d retries: %s\n", retries, msg)
}

func processMessage(msg *kafka.Message) error {
	// Your message processing logic here
	fmt.Printf("Processing message: %s\n", string(msg.Value))

	// Simulate a failure for demonstration purposes
	if string(msg.Value) == "fail" {
		return fmt.Errorf("simulated processing failure")
	}

	return nil
}
