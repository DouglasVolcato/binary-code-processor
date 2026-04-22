package usecases

import "github.com/douglasvolcato/binary-code-processor/event_publisher/internal/entities"

type EventRepositoryInterface interface {
	GetUnpublishedEvents(limit int, offset int) ([]entities.Event, error)
}

type EventProcessorInterface interface {
	SendEventToProcess(event entities.Event) error
}
