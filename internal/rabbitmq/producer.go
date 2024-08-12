package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

type Producer struct {
	Channel *amqp.Channel
	Queue   string
}

func NewProducer(url string, queueName string) (*Producer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %v", err)
	}

	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %v", err)
	}

	return &Producer{Channel: ch, Queue: queueName}, nil
}

func (p *Producer) PublishEvent(event string) error {
	err := p.Channel.Publish(
		"",
		p.Queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)

	if err != nil {
		log.Printf("Failed to publish a message: %s", event)
	} else {
		log.Printf("Event published: %s", event)
	}

	return err
}
