package amqp

import (
	"context"
	"encoding/json"
	"log"

	amqplib "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	ch        *amqplib.Channel
	queueName string
}

func NewPublisher(conn *amqplib.Connection, queueName string) (*Publisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Publisher{ch: ch, queueName: queueName}, nil
}

// Publish sends any JSON-serialisable value to the queue.
// On failure the error is logged; the call is best-effort.
func (p *Publisher) Publish(ctx context.Context, v any) error {
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}

	err = p.ch.PublishWithContext(
		ctx,
		"",          // default exchange
		p.queueName, // routing key = queue name
		false,
		false,
		amqplib.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqplib.Persistent,
			Body:         body,
		},
	)
	if err != nil {
		log.Printf("publish error: %v", err)
		return err
	}

	log.Printf("published to %s: %s", p.queueName, body)
	return nil
}

func (p *Publisher) Close() { _ = p.ch.Close() }

// Dial opens a connection to RabbitMQ.
func Dial(url string) (*amqplib.Connection, error) {
	return amqplib.Dial(url)
}
