package entities

import (
	"errors"
	"strings"
)

type Event struct {
	Task   Task
	Status string
}

func NewEvent(task Task, status string) *Event {
	return &Event{
		Task:   task,
		Status: status,
	}
}

func (e *Event) Validate() error {
	if strings.TrimSpace(e.Status) == "" {
		return errors.New("invalid event status")
	}
	if err := e.Task.Validate(); err != nil {
		return err
	}
	return nil
}
