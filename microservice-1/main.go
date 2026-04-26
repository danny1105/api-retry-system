package main

import (
	"net/http"
	"log"
	"time"
	"encoding/json"
	"bytes"
	"context"

	"github.com/segmentio/kafka-go"
)

type Message struct {
	ID string `json:"id"`
	Body string `json:"body"`
}

var messages = []Message{
	{ID: "1", Body: "message 1"},
	{ID: "2", Body: "message 2"},
	{ID: "3", Body: "message 3"},
}

const (
	kafkaBroker = "localhost:9092"
	kafkaTopic = "events-retry-topic"
	groupID = "microservice-1-group"
	apiURL = "http://localhost:8081/events"
)

func main() {
	log.Default().Println("Microservice 1 is running")

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBroker},
		Topic: kafkaTopic,
		GroupID: groupID,
		MinBytes: 1,	// 1B
		MaxBytes: 10e6, // 10MB
	})

	defer reader.Close()

	for {
		msg, err := reader.FetchMessage(context.Background())
		if err != nil {
			log.Default().Println("Error fetching message: ", err)
			continue
		}

		var m Message
		err = json.Unmarshal(msg.Value, &m)
		if err != nil {
			log.Default().Println("Failed to decode message: ", err)
			continue
		}

		log.Default().Println("Received message: ", m.ID)

		processWithRetry(m)

		err = reader.CommitMessages(context.Background(), msg)
		if err != nil {
			log.Default().Println("Failed to commit message: ", err)
		} else {
			log.Default().Println("Committed message: ", m.ID)
		}
		
	}
}

func processWithRetry(msg Message) {
	for {
		err := callPostAPI(msg)
		if err == nil {
			log.Println("Successfully processed message: ", msg.ID)
			return
		}

		log.Println("Retry in 10s for message:", msg.ID, "Error:", err)
		time.Sleep(10 * time.Second)
	}
}

func callPostAPI(msg Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	response, err := http.Post(apiURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return err
	}

	return nil
}
