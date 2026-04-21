package entities

import (
	"strings"

	"github.com/stretchr/testify/assert"
)

type Task struct {
	ID         string
	Message    string
	BinaryData []byte
	CreatedAt  string
	UpdatedAt  string
}

func NewTask(id string, message string, binaryData []byte, createdAt string, updatedAt string) *Task {
	return &Task{
		ID:         id,
		Message:    message,
		BinaryData: binaryData,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
}

func (t *Task) Validate() error {
	if strings.TrimSpace(t.ID) == "" {
		return assert.AnError
	}
	if strings.TrimSpace(t.Message) == "" {
		return assert.AnError
	}
	if strings.TrimSpace(string(t.BinaryData)) == "" {
		return assert.AnError
	}
	if strings.TrimSpace(t.CreatedAt) == "" {
		return assert.AnError
	}
	if strings.TrimSpace(t.UpdatedAt) == "" {
		return assert.AnError
	}
	return nil
}
