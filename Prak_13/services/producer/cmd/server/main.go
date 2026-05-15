package main

import (
	"log"
	"net/http"
	"os"

	amqpclient "Prak_13-producer/internal/amqp"
	taskhttp "Prak_13-producer/internal/http"
)

func main() {
	rabbitURL := os.Getenv("RABBIT_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	queueName := os.Getenv("QUEUE_NAME")
	if queueName == "" {
		queueName = "task_events"
	}

	conn, err := amqpclient.Dial(rabbitURL)
	if err != nil {
		log.Fatal("rabbit connect error:", err)
	}
	defer conn.Close()

	publisher, err := amqpclient.NewPublisher(conn, queueName)
	if err != nil {
		log.Fatal("publisher error:", err)
	}
	defer publisher.Close()

	handler := taskhttp.NewHandler(publisher)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/v1/tasks", handler.CreateTask)

	log.Println("producer service started on :8082")
	if err := http.ListenAndServe(":8082", mux); err != nil {
		log.Fatal(err)
	}
}
