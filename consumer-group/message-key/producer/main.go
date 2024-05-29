package main

import (
	"fmt"

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
			Key:            []byte(fmt.Sprintf("key-%d", i)),
		}, nil)
		fmt.Printf("Sent: %s\n", message)
	}

	p.Flush(15 * 1000)
}
