package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

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
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Printf("%s received: %s\n", consumerId, string(msg.Value))
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			break
		}
	}
}

func main() {
	var wg sync.WaitGroup
	consumerGroupId1 := "demo-group-1"
	consumerGroupId2 := "demo-group-2"

	wg.Add(2)
	go consume(&wg, consumerGroupId1, "Consumer 1")
	go consume(&wg, consumerGroupId2, "Consumer 2")

	// Wait for termination signal
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	<-sigchan

	wg.Wait()
}
