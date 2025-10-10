package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)
import "github.com/google/uuid"

type user struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(user{
			ID:   uuid.NewString(),
			Name: "Gopher",
		})
	})

	//maps.Collect("status", "ok", "time", "<RFC3339>" )
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(
			map[string]string{"status": "ok", "time": time.Now().UTC().Format(time.RFC3339)},
		)
	})

	addr := ":" + port
	log.Printf("Starting on %s ...", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
