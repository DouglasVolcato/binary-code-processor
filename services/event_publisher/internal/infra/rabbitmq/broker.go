package rabbitmq

import (
	"github.com/douglasvolcato/binary-code-processor/event_publisher/internal/entities"
	queuefile "github.com/douglasvolcato/binary-code-processor/event_publisher/internal/infra/queue"
)

const (
	defaultQueueName    = "task.process"
	defaultExchangeName = "task.processed"
)

type Broker struct {
	producer      *queuefile.Producer
	queueName     string
	exchangeName  string
}

func NewBroker(url string) (*Broker, error) {
	producer, err := queuefile.NewProducer(url)
	if err != nil {
		return nil, err
	}

	if err := producer.DeclareQueue(defaultQueueName); err != nil {
		producer.Close()
		return nil, err
	}

	if err := producer.DeclareExchange(defaultExchangeName, "fanout"); err != nil {
		producer.Close()
		return nil, err
	}

	return &Broker{
		producer:     producer,
		queueName:    defaultQueueName,
		exchangeName: defaultExchangeName,
	}, nil
}

func (b *Broker) SendToQueue(event entities.Event) error {
	return b.producer.Publish(b.queueName, eventPayload{
		ID:         event.ID,
		Status:     event.Status,
		BinaryCode: event.BinaryCode,
	})
}

func (b *Broker) SendFanoutEvent(event entities.Event) error {
	return b.producer.PublishToExchange(b.exchangeName, "", eventPayload{
		ID:         event.ID,
		Status:     event.Status,
		BinaryCode: event.BinaryCode,
	})
}

func (b *Broker) Close() {
	b.producer.Close()
}

type eventPayload struct {
	ID         string `json:"id"`
	Status     string `json:"status"`
	BinaryCode string `json:"binaryCode"`
}
