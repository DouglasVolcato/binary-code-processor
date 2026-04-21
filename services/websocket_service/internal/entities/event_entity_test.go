package entities

import (
	"testing"

	"github.com/douglasvolcato/binary-code-processer/websocket_service/test"
	"github.com/stretchr/testify/assert"
)

type testDataEventEntity struct {
	Task   Task
	Status string
}

func makeFakeDataEventEntity(invalidTask bool) *testDataEventEntity {
	faker := test.FakeData{}
	id := faker.ID()
	if invalidTask {
		id = ""
	}
	return &testDataEventEntity{
		Task: Task{
			ID:         id,
			Message:    faker.Phrase(),
			BinaryData: faker.Binary(),
			CreatedAt:  faker.Date(),
			UpdatedAt:  faker.Date(),
		},
		Status: faker.Word(),
	}
}

func TestNewEventShouldCreateEvent(t *testing.T) {
	testData := makeFakeDataEventEntity(false)
	sut := NewEvent(testData.Task, testData.Status)
	assert.NotNil(t, sut)
	assert.Equal(t, testData.Task.ID, sut.Task.ID)
	assert.Equal(t, testData.Task.Message, sut.Task.Message)
	assert.Equal(t, testData.Task.BinaryData, sut.Task.BinaryData)
	assert.Equal(t, testData.Task.CreatedAt, sut.Task.CreatedAt)
	assert.Equal(t, testData.Task.UpdatedAt, sut.Task.UpdatedAt)
	assert.Equal(t, testData.Status, sut.Status)
}

func TestValidateShouldReturnErrorIfEventDataIsInvalid(t *testing.T) {
	testData := makeFakeDataEventEntity(true)
	err := NewEvent(testData.Task, testData.Status).Validate()
	assert.Error(t, err)
	err = NewEvent(testData.Task, "").Validate()
	assert.Error(t, err)
}

func TestValidateShouldReturnNilIfEventDataIsValid(t *testing.T) {
	testData := makeFakeDataEventEntity(false)
	err := NewEvent(testData.Task, testData.Status).Validate()
	assert.NoError(t, err)
}
