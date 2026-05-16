package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"Prak_14-consumer/internal/store"
	amqp "github.com/rabbitmq/amqp091-go"
)

const maxAttempts = 3

type TaskJob struct {
	Job       string `json:"job"`
	TaskID    string `json:"task_id"`
	Attempt   int    `json:"attempt"`
	MessageID string `json:"message_id"`
}

func processTask(job TaskJob) error {
	time.Sleep(1 * time.Second) // имитация тяжёлой работы

	// Искусственная ошибка для task_id == "t_fail"
	if job.TaskID == "t_fail" {
		return fmt.Errorf("simulated processing error for task_id=%s", job.TaskID)
	}
	return nil
}

func publishJob(ch *amqp.Channel, queue string, job TaskJob) error {
	body, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return ch.PublishWithContext(
		context.Background(), "", queue, false, false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}

func main() {
	rabbitURL := os.Getenv("RABBIT_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
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

	// Объявить DLQ первой
	if _, err := ch.QueueDeclare("task_jobs_dlq", true, false, false, false, nil); err != nil {
		log.Fatalf("dlq declare error: %v", err)
	}

	// Объявить основную очередь
	if _, err := ch.QueueDeclare("task_jobs", true, false, false, false,
		amqp.Table{
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": "task_jobs_dlq",
		}); err != nil {
		log.Fatalf("queue declare error: %v", err)
	}

	// prefetch = 1
	if err := ch.Qos(1, 0, false); err != nil {
		log.Fatalf("qos error: %v", err)
	}

	msgs, err := ch.Consume("task_jobs", "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("consume error: %v", err)
	}

	processed := store.NewProcessedStore()

	log.Println("worker started, waiting for jobs...")

	for d := range msgs {
		var job TaskJob
		if err := json.Unmarshal(d.Body, &job); err != nil {
			log.Printf("bad message: %v — nack", err)
			_ = d.Nack(false, false)
			continue
		}

		log.Printf("received job=%s task_id=%s attempt=%d message_id=%s",
			job.Job, job.TaskID, job.Attempt, job.MessageID)

		// Идемпотентная проверка
		if processed.Exists(job.MessageID) {
			log.Printf("duplicate message_id=%s — skip", job.MessageID)
			_ = d.Ack(false)
			continue
		}

		// Обработка
		if err := processTask(job); err != nil {
			log.Printf("process error (attempt %d/%d): %v", job.Attempt, maxAttempts, err)

			job.Attempt++
			if job.Attempt <= maxAttempts {
				// Повторная попытка
				if pubErr := publishJob(ch, "task_jobs", job); pubErr != nil {
					log.Printf("retry publish error: %v", pubErr)
				} else {
					log.Printf("retry scheduled: task_id=%s attempt=%d", job.TaskID, job.Attempt)
				}
			} else {
				// Превышен лимит — в DLQ
				if pubErr := publishJob(ch, "task_jobs_dlq", job); pubErr != nil {
					log.Printf("dlq publish error: %v", pubErr)
				} else {
					log.Printf("sent to DLQ: task_id=%s", job.TaskID)
				}
			}
			_ = d.Ack(false)
			continue
		}

		// Успех
		processed.MarkDone(job.MessageID)
		log.Printf("task processed successfully: task_id=%s", job.TaskID)
		_ = d.Ack(false)
	}
}
