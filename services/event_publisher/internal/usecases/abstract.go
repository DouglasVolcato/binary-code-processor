package usecases

import "github.com/douglasvolcato/binary-code-processor/event_publisher/internal/entities"

type EventRepositoryInterface interface {
	GetUnprocessedEvents(limit int, offset int) ([]entities.Event, error)
}

type EventProcessorInterface interface {
	SendEventToProcess(event entities.Event) error
	SendFanoutEvent(event entities.Event) error
}

type RemoteEventProcessorInterface interface {
	SendToQueue(event entities.Event) error
	SendFanoutEvent(event entities.Event) error
}
