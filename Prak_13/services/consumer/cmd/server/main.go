package main

import (
	"encoding/json"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type TaskEvent struct {
	Event     string `json:"event"`
	TaskID    string `json:"task_id"`
	TS        string `json:"ts"`
	RequestID string `json:"request_id,omitempty"`
	Producer  string `json:"producer,omitempty"`
}

func main() {
	rabbitURL := os.Getenv("RABBIT_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	queueName := os.Getenv("QUEUE_NAME")
	if queueName == "" {
		queueName = "task_events"
	}

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("rabbit connect error: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("channel error: %v", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		log.Fatalf("queue declare error: %v", err)
	}

	// prefetch = 1: не брать следующее сообщение до подтверждения текущего
	if err := ch.Qos(1, 0, false); err != nil {
		log.Fatalf("qos error: %v", err)
	}

	msgs, err := ch.Consume(
		queueName,
		"",    // consumer tag
		false, // auto-ack = false (ручное подтверждение)
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		log.Fatalf("consume error: %v", err)
	}

	log.Println("consumer started, waiting for messages from queue:", queueName)

	for d := range msgs {
		var ev TaskEvent
		if err := json.Unmarshal(d.Body, &ev); err != nil {
			log.Printf("bad message body: %v — sending nack", err)
			_ = d.Nack(false, false) // не перекладывать обратно
			continue
		}

		log.Printf("received event=%s task_id=%s ts=%s request_id=%s",
			ev.Event, ev.TaskID, ev.TS, ev.RequestID)

		// Подтверждаем успешную обработку
		if err := d.Ack(false); err != nil {
			log.Printf("ack error: %v", err)
		}
	}
}
