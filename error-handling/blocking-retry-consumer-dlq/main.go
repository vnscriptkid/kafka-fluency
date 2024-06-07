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

	fmt.Println("Waiting for messages...")
	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			processMessageWithRetry(consumer, producer, msg, 3, 10*time.Second)
		} else {
			// Error is either due to poll timeout or a Kafka error
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}

func processMessageWithRetry(consumer *kafka.Consumer, producer *kafka.Producer, msg *kafka.Message, retries int, waitTime time.Duration) {
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

	// After exhausting retries, send the message to the DLQ
	fmt.Printf("Failed to process message after %d retries: %s. Sending to DLQ...\n", retries, msg)
	dlqTopic := "myDLQ"
	err := sendMessageToDLQ(producer, dlqTopic, msg)
	if err != nil {
		fmt.Printf("Failed to send message to DLQ: %s\n", err)
	}
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

func sendMessageToDLQ(producer *kafka.Producer, topic string, msg *kafka.Message) error {
	dlqMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          msg.Value,
		Key:            msg.Key,
		Headers:        msg.Headers,
	}

	deliveryChan := make(chan kafka.Event)
	err := producer.Produce(dlqMsg, deliveryChan)
	if err != nil {
		return err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}

	close(deliveryChan)
	return nil
}
