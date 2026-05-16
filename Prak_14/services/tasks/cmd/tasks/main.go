package main

import (
	"log"
	"net/http"
	"os"

	amqpclient "Prak_14-producer/internal/amqp"
	taskhttp "Prak_14-producer/internal/http"
)

func main() {
	rabbitURL := os.Getenv("RABBIT_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	conn, err := amqpclient.Dial(rabbitURL)
	if err != nil {
		log.Fatal("rabbit connect error:", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("channel error:", err)
	}
	defer ch.Close()

	if err := amqpclient.DeclareQueues(ch); err != nil {
		log.Fatal("queue declare error:", err)
	}

	handler := taskhttp.NewHandler(ch)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/v1/jobs/process-task", handler.EnqueueJob)

	log.Println("tasks service started on :8082")
	if err := http.ListenAndServe(":8082", mux); err != nil {
		log.Fatal(err)
	}
}
