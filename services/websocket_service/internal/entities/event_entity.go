package entities

import (
	"strings"

	"github.com/stretchr/testify/assert"
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
		return assert.AnError
	}
	if err := e.Task.Validate(); err != nil {
		return err
	}
	return nil
}
