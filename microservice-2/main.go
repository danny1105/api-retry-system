package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type Message struct {
	ID   string `json:"id"`
	Body string `json:"body"`
}

var (
	mu        sync.Mutex
	processed = make(map[string]bool)
)

func eventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("Invalid request method: %s", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var msg Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if msg.ID == "" {
		log.Println("Missing 'id' field in JSON")
		http.Error(w, "Missing 'id' field", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if processed[msg.ID] {
		log.Println("Duplicate message received:", msg.ID)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Duplicate ignored"))
		return
	}

	log.Println("Processing message:", msg.ID)

	processed[msg.ID] = true

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Processed"))
}

func main() {
	http.HandleFunc("/events", eventHandler)

	log.Println("Microservice 2 is running on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}