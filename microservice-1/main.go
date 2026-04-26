package main

import (
	"net/http"
	"log"
	"sync"
	"time"
	"encoding/json"
	"bytes"
)

var mu sync.Mutex


type Message struct {
	ID string `json:"id"`
	Body string `json:"body"`
}

var messages = []Message{
	{ID: "1", Body: "message 1"},
	{ID: "2", Body: "message 2"},
	{ID: "3", Body: "message 3"},
}

var apiURL = "http://localhost:8081/events"

func main() {
	log.Default().Println("Microservice 1 is running")

	for _, msg := range messages {
		processWithRetry(msg)
	}

	log.Default().Println("All messages processed")
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
