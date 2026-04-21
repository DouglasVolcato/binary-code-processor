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

///////////////////////////

type TestData struct {
	ID         string
	Message    string
	BinaryData []byte
	CreatedAt  string
	UpdatedAt  string
}

func MakeTestData() TestData {
	return TestData{
		ID:         faker.UUIDDigit(),
		Message:    faker.Sentence(),
		BinaryData: []byte(faker.Paragraph()),
		CreatedAt:  faker.Date(),
		UpdatedAt:  faker.Date(),
	}
}

var testData = MakeTestData()

func TestNewTaskShouldCreateTask(t *testing.T) {
	sut := NewTask(
		testData.ID,
		testData.Message,
		testData.BinaryData,
		testData.CreatedAt,
		testData.UpdatedAt,
	)
	assert.NotNil(t, sut)
	assert.Equal(t, testData.ID, sut.ID)
	assert.Equal(t, testData.Message, sut.Message)
	assert.Equal(t, testData.BinaryData, sut.BinaryData)
	assert.Equal(t, testData.CreatedAt, sut.CreatedAt)
	assert.Equal(t, testData.UpdatedAt, sut.UpdatedAt)
}

func TestValidateShouldReturnErrorIfDataIsInvalid(t *testing.T) {
	err := NewTask(
		"",
		testData.Message,
		testData.BinaryData,
		testData.CreatedAt,
		testData.UpdatedAt,
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		testData.ID,
		"",
		testData.BinaryData,
		testData.CreatedAt,
		testData.UpdatedAt,
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		testData.ID,
		testData.Message,
		([]byte)(""),
		testData.CreatedAt,
		testData.UpdatedAt,
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		testData.ID,
		testData.Message,
		testData.BinaryData,
		"",
		testData.UpdatedAt,
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		testData.ID,
		testData.Message,
		testData.BinaryData,
		testData.CreatedAt,
		"",
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		" ",
		testData.Message,
		testData.BinaryData,
		testData.CreatedAt,
		testData.UpdatedAt,
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		testData.ID,
		" ",
		testData.BinaryData,
		testData.CreatedAt,
		testData.UpdatedAt,
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		testData.ID,
		testData.Message,
		([]byte)(" "),
		testData.CreatedAt,
		testData.UpdatedAt,
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		testData.ID,
		testData.Message,
		testData.BinaryData,
		" ",
		testData.UpdatedAt,
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		testData.ID,
		testData.Message,
		testData.BinaryData,
		testData.CreatedAt,
		" ",
	).Validate()
	assert.Error(t, err)
}

func TestValidateShouldReturnNilIfDataIsValid(t *testing.T) {
	err := NewTask(
		testData.ID,
		testData.Message,
		testData.BinaryData,
		testData.CreatedAt,
		testData.UpdatedAt,
	).Validate()
	assert.NoError(t, err)
}
