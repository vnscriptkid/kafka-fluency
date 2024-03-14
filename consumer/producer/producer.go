package producer

import (
	"log"
	"time"

	"github.com/IBM/sarama"
)

func ProduceDispatchedOrder(brokerAddress []string, topic string, orderID string) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokerAddress, config)
	if err != nil {
		log.Fatalf("Failed to start Sarama producer: %v", err)
	}
	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.StringEncoder(orderID),
		Value:     sarama.StringEncoder("Order dispatched"),
		Timestamp: time.Now(),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
	log.Printf("Order %s dispatched, partition %d, offset %d", orderID, partition, offset)
}
