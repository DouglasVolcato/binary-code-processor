package entities

import (
	"errors"
	"strings"
)

type Task struct {
	ID         string
	BinaryCode string
}

func NewTask(id string, binaryCode string) *Task {
	return &Task{
		ID:         id,
		BinaryCode: binaryCode,
	}
}

func (t *Task) Validate() error {
	if strings.TrimSpace(t.ID) == "" {
		return errors.New("invalid task id")
	}
	if strings.TrimSpace(string(t.BinaryCode)) == "" {
		return errors.New("invalid task binary code")
	}
	return nil
}
