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
		"client.id":         consumerId,
	})
	if err != nil {
		panic(err)
	}
	defer c.Close()

	c.SubscribeTopics([]string{"demo-topic"}, nil)

	fmt.Println("Consumer started")

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Printf("%s received: (%s) of key (%s)\n", consumerId, string(msg.Value), string(msg.Key))
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			break
		}
	}
}

const numOfConsumers = 3

func main() {
	var wg sync.WaitGroup
	consumerGroupId := "demo-group"

	for i := 0; i < numOfConsumers; i++ {
		wg.Add(1)
		go consume(&wg, consumerGroupId, fmt.Sprintf("Consumer-%d", i+1))
	}
	// Wait for termination signal
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	<-sigchan

	wg.Wait()
}
