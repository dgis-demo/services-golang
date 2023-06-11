package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Point struct {
	Lat       float32 `json:"lat"`
	Lon       float32 `json:"lon"`
	Magnitude float32 `json:"magnitude"`
}

func createProducer() *kafka.Producer {
	producer, err := kafka.NewProducer(
		&kafka.ConfigMap{
			"bootstrap.servers": fmt.Sprintf(
				"%s:%s",
				KAFKA_HOST,
				KAFKA_PORT,
			),
			"client.id": "client_service",
			"acks":      "all",
		},
	)

	if err != nil {
		log.Printf("failed to create producer: %s", err)
	}
	return producer
}

func sendMessage(
	producer *kafka.Producer,
	topic string,
	message []byte,
) error {
	delivery_chan := make(chan kafka.Event, 10000)
	err := producer.Produce(
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: kafka.PartitionAny,
			},
			Value: message,
		},
		delivery_chan,
	)
	return err
}

func main() {
	producer := createProducer()

	for {
		point := &Point{
			Lat:       float32((rand.Intn(90) - rand.Intn(90))) + rand.Float32(),
			Lon:       float32((rand.Intn(180) - rand.Intn(180))) + rand.Float32(),
			Magnitude: float32(rand.Intn(10)) + rand.Float32(),
		}

		message, err := json.Marshal(point)
		if err != nil {
			log.Printf("failed to create a message: %s", err)
			break
		}

		err = sendMessage(producer, KAFKA_TOPIC, message)
		if err != nil {
			log.Printf("failed to send a message: %s", err)
			break
		}

		log.Printf("message has been sent, %v", point)
		time.Sleep(time.Second / RPS)
	}
}
