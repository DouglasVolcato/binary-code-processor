package entities

import (
	"errors"
	"strings"
)

type Task struct {
	ID         string
	Message    string
	BinaryCode string
	CreatedAt  string
	UpdatedAt  string
}

func NewTask(id string, message string, binaryCode string, createdAt string, updatedAt string) *Task {
	return &Task{
		ID:         id,
		Message:    message,
		BinaryCode: binaryCode,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
}

func (t *Task) Validate() error {
	if strings.TrimSpace(t.ID) == "" {
		return errors.New("invalid task id")
	}
	if strings.TrimSpace(t.Message) == "" {
		return errors.New("invalid task message")
	}
	if strings.TrimSpace(t.CreatedAt) == "" {
		return errors.New("invalid task created at")
	}
	if strings.TrimSpace(t.UpdatedAt) == "" {
		return errors.New("invalid task updated at")
	}
	return nil
}
