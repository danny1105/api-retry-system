package main

import ( 
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type Message struct {
	ID	  string    `json:"id"`
	Body  string `json:"body"`
}

var mu sync.Mutex

func eventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Default().Printf("Invalid request method: %s", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var msg Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		log.Default().Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if msg.ID == "" {
		log.Default().Println("Missing 'id' field in JSON")
		http.Error(w, "Missing 'id' field", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	log.Println("Processing Message: ", msg.ID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Processed"))
}


func main() {
	http.HandleFunc("/events", eventHandler)

	log.Default().Println("Microservice 2 is running on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}