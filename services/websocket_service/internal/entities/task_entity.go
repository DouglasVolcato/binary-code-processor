package entities

import (
	"strings"

	"github.com/stretchr/testify/assert"
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
		return assert.AnError
	}
	if strings.TrimSpace(string(t.BinaryCode)) == "" {
		return assert.AnError
	}
	return nil
}
