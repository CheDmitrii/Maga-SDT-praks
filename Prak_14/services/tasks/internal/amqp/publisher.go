package amqp

import (
	"context"
	"encoding/json"

	amqplib "github.com/rabbitmq/amqp091-go"
)

func Dial(url string) (*amqplib.Connection, error) {
	return amqplib.Dial(url)
}

func PublishJob(ch *amqplib.Channel, queue string, job any) error {
	body, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return ch.PublishWithContext(
		context.Background(),
		"",
		queue,
		false,
		false,
		amqplib.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqplib.Persistent,
			Body:         body,
		},
	)
}

func DeclareQueues(ch *amqplib.Channel) error {
	// DLQ first (must exist before main queue references it)
	_, err := ch.QueueDeclare(
		"task_jobs_dlq",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Main queue
	_, err = ch.QueueDeclare(
		"task_jobs",
		true,
		false,
		false,
		false,
		amqplib.Table{
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": "task_jobs_dlq",
		},
	)
	return err
}
