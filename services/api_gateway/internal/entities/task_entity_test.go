package entities

import (
	"testing"

	"github.com/douglasvolcato/task-monitor/api_gateway/test"
	"github.com/stretchr/testify/assert"
)

type testDataTaskEntity struct {
	ID         string
	Message    string
	BinaryData []byte
	CreatedAt  string
	UpdatedAt  string
}

func makeFakeDataTaskEntity() *testDataTaskEntity {
	faker := test.FakeData{}
	return &testDataTaskEntity{
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

func TestValidateShouldReturnErrorIfTaskDataIsInvalid(t *testing.T) {
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

func TestValidateShouldReturnNilIfTaskDataIsValid(t *testing.T) {
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
