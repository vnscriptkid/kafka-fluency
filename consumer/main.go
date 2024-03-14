package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/vnscriptkid/kafka-fluency/consumer/deserializers"
	"github.com/vnscriptkid/kafka-fluency/consumer/processors"
)

type consumerGroupHandler struct {
	dispatchService processors.DispatchService
}

func NewConsumerGroupHandler() consumerGroupHandler {
	return consumerGroupHandler{
		dispatchService: processors.NewDispatchService(),
	}
}

func (consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(msg.Value), msg.Timestamp, msg.Topic)

		// Deserialize
		oCreated, err := deserializers.DeserializeOrderCreated(msg.Value)

		if err != nil {
			log.Printf("Error deserializing order created: %v", err)
			// Ignore for now
			// Could potentially log to a dead-letter queue
			continue
		}

		// Process
		h.dispatchService.Process(oCreated)

		sess.MarkMessage(msg, "")
		// Manually commit offsets
		sess.Commit()

	}
	return nil
}

func (s *svc) consumePlacedOrders(brokerAddress []string, topic, groupID string) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0 // Specify appropriate Kafka version
	config.Consumer.Return.Errors = true
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Offsets.AutoCommit.Enable = false
	// config.Consumer.Offsets.AutoCommit.Enable = true
	// config.Consumer.Offsets.AutoCommit.Interval = time.Minute * 1

	consumerGroup, err := sarama.NewConsumerGroup(brokerAddress, groupID, config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Trap SIGINT and SIGTERM to trigger a graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	go func() {
		for {
			if err := consumerGroup.Consume(ctx, []string{topic}, s.consumerHandler); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	<-signals // Wait for signal
	log.Println("Terminating: signal caught")
}

type syncProducer struct {
	client sarama.SyncProducer
}

type svc struct {
	consumerHandler consumerGroupHandler
	producer        syncProducer
}

func main() {
	brokerAddress := []string{"localhost:9092"}
	topic := "order.created"
	groupID := "dispatch.order.created.group"

	s := svc{
		consumerHandler: NewConsumerGroupHandler(),
	}
	s.consumePlacedOrders(brokerAddress, topic, groupID)
}

func NewProducer(brokerAddress []string, topic string, orderID string) syncProducer {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokerAddress, config)
	if err != nil {
		log.Fatalf("Failed to start Sarama producer: %v", err)
	}
	// defer producer.Close()
	return syncProducer{
		client: producer,
	}

}
