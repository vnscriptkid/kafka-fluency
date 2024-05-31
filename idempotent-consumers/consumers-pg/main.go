package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/lib/pq"
)

// CREATE TABLE processed_messages (
//     message_id VARCHAR PRIMARY KEY
// );

func main() {
	db, err := sql.Open("postgres", "user=postgres dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

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

		// Try to insert the message ID into the database
		_, err = db.Exec("INSERT INTO processed_messages (message_id) VALUES ($1)", messageID)
		if err != nil {
			if err.Error() == "pq: duplicate key value violates unique constraint \"processed_messages_pkey\"" {
				// Message already processed
				fmt.Printf("Skipping already processed message: %s\n", string(msg.Value))
			} else {
				log.Fatal(err)
			}
		} else {
			// Message not processed, proceed with processing
			fmt.Printf("Processing message: %s\n", string(msg.Value))

			// Simulate message processing
			// ...

			// Commit the message offset
			c.CommitMessage(msg)
		}
	}
}
