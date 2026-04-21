package entities

import (
	"testing"

	"github.com/douglasvolcato/task-monitor/websocket_service/test"
	"github.com/stretchr/testify/assert"
)

type testDataType struct {
	ID         string
	Message    string
	BinaryData []byte
	CreatedAt  string
	UpdatedAt  string
}

func makeFakeDataTaskEntity() *testDataType {
	faker := test.FakeData{}
	return &testDataType{
		ID:         faker.ID(),
		Message:    faker.Phrase(),
		BinaryData: faker.Binary(),
		CreatedAt:  faker.Date(),
		UpdatedAt:  faker.Date(),
	}
}

func TestNewTaskShouldCreateTask(t *testing.T) {
	testData := makeFakeDataTaskEntity()
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
	testData := makeFakeDataTaskEntity()
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
	testData := makeFakeDataTaskEntity()
	err := NewTask(
		testData.ID,
		testData.Message,
		testData.BinaryData,
		testData.CreatedAt,
		testData.UpdatedAt,
	).Validate()
	assert.NoError(t, err)
}
