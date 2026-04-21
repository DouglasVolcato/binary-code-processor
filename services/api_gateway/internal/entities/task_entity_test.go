package entities

import (
	"strings"
	"testing"

	faker "github.com/go-faker/faker/v4"
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

func TestNewTaskShouldCreateTask(t *testing.T) {
	id := faker.UUIDDigit()
	message := faker.Sentence()
	binaryData := []byte(faker.Paragraph())
	createdAt := faker.Date()
	updatedAt := faker.Date()
	sut := NewTask(
		id,
		message,
		binaryData,
		createdAt,
		updatedAt,
	)
	assert.NotNil(t, sut)
	assert.Equal(t, id, sut.ID)
	assert.Equal(t, message, sut.Message)
	assert.Equal(t, binaryData, sut.BinaryData)
	assert.Equal(t, createdAt, sut.CreatedAt)
	assert.Equal(t, updatedAt, sut.UpdatedAt)
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

func TestValidateShouldReturnErrorIfDataIsInvalid(t *testing.T) {
	id := faker.UUIDDigit()
	message := faker.Sentence()
	binaryData := []byte(faker.Paragraph())
	createdAt := faker.Date()
	updatedAt := faker.Date()
	err := NewTask(
		"",
		message,
		binaryData,
		createdAt,
		updatedAt,
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		id,
		"",
		binaryData,
		createdAt,
		updatedAt,
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		id,
		message,
		([]byte)(""),
		createdAt,
		updatedAt,
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		id,
		message,
		binaryData,
		"",
		updatedAt,
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		id,
		message,
		binaryData,
		createdAt,
		"",
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		" ",
		message,
		binaryData,
		createdAt,
		updatedAt,
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		id,
		" ",
		binaryData,
		createdAt,
		updatedAt,
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		id,
		message,
		([]byte)(" "),
		createdAt,
		updatedAt,
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		id,
		message,
		binaryData,
		" ",
		updatedAt,
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		id,
		message,
		binaryData,
		createdAt,
		" ",
	).Validate()
	assert.Error(t, err)
}

// func TestValidateShouldValidateData(t *testing.T) {}
