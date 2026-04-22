package usecases

import "github.com/douglasvolcato/binary-code-processor/event_publisher/internal/entities"

type EventRepositoryInterface interface {
	GetUnprocessedEvents(limit int, offset int) ([]entities.Event, error)
	GetProcessedEvents(limit int, offset int) ([]entities.Event, error)
}

type EventProcessorInterface interface {
	// SendEventToProcess hands the queued event to the local processing pipeline.
	SendEventToProcess(event entities.Event) error
	// SendFanoutEvent hands the processed event to the local fanout pipeline.
	SendFanoutEvent(event entities.Event) error
}

type RemoteEventProcessorInterface interface {
	// SendToQueue publishes the event to RabbitMQ's queue flow.
	SendToQueue(event entities.Event) error
	// SendFanoutEvent publishes the event to RabbitMQ's fanout flow.
	SendFanoutEvent(event entities.Event) error
}
