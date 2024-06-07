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

	producer, err := kafka.NewProducer(config)
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	defer producer.Close()

	topics := []string{"myTopic"}
	err = consumer.SubscribeTopics(topics, nil)
	if err != nil {
		fmt.Printf("Failed to subscribe to topics: %s\n", err)
		os.Exit(1)
	}

	retryQueue := make(chan *kafka.Message, 100)

	// Start retry handler goroutine
	go handleRetries(retryQueue, producer, consumer)

	fmt.Println("Waiting for messages...")
	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			go processMessage(consumer, msg, retryQueue)
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}

func processMessage(consumer *kafka.Consumer, msg *kafka.Message, retryQueue chan *kafka.Message) {
	if err := processMessageContent(msg); err != nil {
		retryQueue <- msg
	} else {
		_, err := consumer.CommitMessage(msg)
		if err != nil {
			fmt.Printf("Failed to commit message: %s\n", err)
		}
	}
}

func processMessageContent(msg *kafka.Message) error {
	fmt.Printf("Processing message: %s\n", string(msg.Value))
	if string(msg.Value) == "fail" {
		return fmt.Errorf("simulated processing failure")
	}
	return nil
}

func handleRetries(retryQueue chan *kafka.Message, producer *kafka.Producer, consumer *kafka.Consumer) {
	for msg := range retryQueue {
		go func(msg *kafka.Message) {
			if err := processMessageWithRetry(msg, 3, 10*time.Second); err != nil {
				fmt.Printf("Failed to process message after retries: %s. Sending to DLQ...\n", msg)
				sendToDLQ(producer, consumer, "myDLQ", msg)
			}
		}(msg)
	}
}

func processMessageWithRetry(msg *kafka.Message, retries int, waitTime time.Duration) error {
	for i := 0; i < retries; i++ {
		var err error
		if err := processMessageContent(msg); err == nil {
			return nil
		}
		fmt.Printf("Failed to process message: %s. Retrying %d/%d...\n", err, i+1, retries)
		time.Sleep(waitTime)
	}
	return fmt.Errorf("failed to process message after %d retries", retries)
}

func sendToDLQ(producer *kafka.Producer, consumer *kafka.Consumer, topic string, msg *kafka.Message) {
	dlqMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          msg.Value,
		Key:            msg.Key,
		Headers:        msg.Headers,
	}

	// Produce the message and handle delivery report asynchronously
	producer.Produce(dlqMsg, nil)

	// Wait for message delivery reports
	for e := range producer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				fmt.Printf("Failed to send message to DLQ: %s\n", ev.TopicPartition.Error)
			} else {
				fmt.Printf("Message sent to DLQ: %s\n", ev.TopicPartition)
				// Commit the message after it has been successfully sent to the DLQ
				if _, err := consumer.CommitMessage(msg); err != nil {
					fmt.Printf("Failed to commit message after sending to DLQ: %s\n", err)
				}
			}
			return
		}
	}
}
