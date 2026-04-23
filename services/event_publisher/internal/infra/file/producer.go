package file

import (
	"encoding/json"

	"github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
}

func NewProducer(url string) (*Producer, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Producer{
		conn:    conn,
		channel: ch,
	}, nil
}

func (p *Producer) DeclareQueue(name string) error {
	_, err := p.channel.QueueDeclare(
		name,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false,
		nil,
	)
	return err
}

func (p *Producer) Publish(queue string, payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return p.channel.Publish(
		"",
		queue,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (p *Producer) Close() {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
}
