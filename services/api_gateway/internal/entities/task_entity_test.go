package entities

import (
	"testing"

	faker "github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
)

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
