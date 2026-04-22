package entities

import (
	"errors"
	"strings"
)

type Event struct {
	ID     string
	Status string
}

func NewEvent(id string, status string) *Event {
	return &Event{
		ID:     id,
		Status: status,
	}
}

func (e *Event) Validate() error {
	if strings.TrimSpace(e.ID) == "" {
		return errors.New("invalid event ID")
	}
	if strings.TrimSpace(e.Status) == "" {
		return errors.New("invalid event status")
	}
	return nil
}
