package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var i = 0

// The consumer is supposed to simulate an error after processing 2 messages
func consume(wg *sync.WaitGroup, groupId, consumerId string) {
	defer wg.Done()
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          groupId,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}
	defer c.Close()

	c.SubscribeTopics([]string{"demo-topic"}, nil)

	for {
		if i == 2 {
			fmt.Printf("Consumer ERR: Simulating error\n")
			c.Close()
			break
		}

		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Printf("%s received: %s\n", consumerId, string(msg.Value))
			fmt.Printf("Processing...\n")
			i++

		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			break
		}
	}
}

func main() {
	var wg sync.WaitGroup
	consumerGroupId := "demo-group"

	wg.Add(1)
	go consume(&wg, consumerGroupId, "Consumer ERR")

	// Wait for termination signal
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	<-sigchan

	wg.Wait()
}
