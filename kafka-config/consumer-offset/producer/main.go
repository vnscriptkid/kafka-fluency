package main

import (
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		panic(err)
	}
	defer p.Close()

	topic := "demo-topic"
	for i := 0; i < 10; i++ {
		message := fmt.Sprintf("Message %d", i)
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(message),
		}, nil)
		fmt.Printf("Sent: %s\n", message)
		time.Sleep(1 * time.Second)
	}

	p.Flush(15 * 1000)
}
