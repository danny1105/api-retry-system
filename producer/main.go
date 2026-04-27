package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

type Message struct {
	ID   string `json:"id"`
	Body string `json:"body"`
}

func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

func main() {
	broker := getEnv("KAFKA_BROKER", "localhost:9092")
	topic := "retry-topic"

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker},
		Topic:   topic,
	})

	defer writer.Close()

	messages := []Message{
		{ID: "1", Body: "message-1"},
		{ID: "2", Body: "message-2"},
		{ID: "3", Body: "message-3"},
	}

	log.Println("Producing messages to topic: ", topic)

	for _, m := range messages {
		data, err := json.Marshal(m)
		if err != nil {
			log.Println("marshal error: ", err)
			continue
		}

		err = writer.WriteMessages(context.Background(),
			kafka.Message{Value: data},
		)

		if err != nil {
			log.Println("failed to produce:", err)
		} else {
			log.Println("produced:", m.ID)
		}

		time.Sleep(1 * time.Second) // small gap for readability
	}

	log.Println("Done producing messages")
}