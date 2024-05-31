package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "my-group",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		log.Fatal(err)
	}

	c.SubscribeTopics([]string{"my-topic"}, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err != nil {
			log.Printf("Consumer error: %v (%v)\n", err, msg)
			continue
		}

		messageID := string(msg.Key)

		// Use SETNX to ensure only one consumer processes the message
		set, err := rdb.SetNX(ctx, messageID, "processing", 10*time.Minute).Result()
		if err != nil {
			log.Fatal(err)
		}

		if set {
			// Message is not processed yet
			fmt.Printf("Processing message: %s\n", string(msg.Value))

			// Simulate message processing
			// ...

			// Mark message as processed by setting the key with an indefinite expiration
			err = rdb.Set(ctx, messageID, "processed", 0).Err()
			if err != nil {
				log.Fatal(err)
			}

			// Commit the message offset
			_, err = c.CommitMessage(msg)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			// Message is being processed by another consumer or already processed
			fmt.Printf("Skipping already processed message: %s\n", string(msg.Value))
		}
	}
}
