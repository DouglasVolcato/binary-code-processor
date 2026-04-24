package processor

import "github.com/douglasvolcato/binary-code-processor/event_publisher/internal/entities"

type Processor struct{}

func NewProcessor() *Processor {
	return &Processor{}
}

func (p *Processor) SendEventToProcess(event entities.Event) error {
	return event.Validate()
}

func (p *Processor) SendFanoutEvent(event entities.Event) error {
	return event.Validate()
}

