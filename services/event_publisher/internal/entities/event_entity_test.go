package entities

import (
	"testing"

	"github.com/douglasvolcato/binary-code-processor/event_publisher/test"
	"github.com/stretchr/testify/assert"
)

type testDataEventEntity struct {
	Task   Task
	Status string
}

func makeFakeDataEventEntity() *testDataEventEntity {
	faker := test.FakeData{}
	return &testDataEventEntity{
		Task: Task{
			ID:         faker.ID(),
			Message:    faker.Phrase(),
			BinaryCode: faker.Phrase(),
			CreatedAt:  faker.Date(),
			UpdatedAt:  faker.Date(),
		},
		Status: faker.Word(),
	}
}

func TestNewEventShouldCreateEvent(t *testing.T) {
	testData := makeFakeDataEventEntity()
	sut := NewEvent(testData.Task, testData.Status)

	assert.NotNil(t, sut)
	assert.Equal(t, testData.Task.ID, sut.Task.ID)
	assert.Equal(t, testData.Task.Message, sut.Task.Message)
	assert.Equal(t, testData.Task.BinaryCode, sut.Task.BinaryCode)
	assert.Equal(t, testData.Task.CreatedAt, sut.Task.CreatedAt)
	assert.Equal(t, testData.Task.UpdatedAt, sut.Task.UpdatedAt)
	assert.Equal(t, testData.Status, sut.Status)
}

func TestValidateShouldReturnErrorIfEventDataIsInvalid(t *testing.T) {
	testData := makeFakeDataEventEntity()
	tests := []struct {
		name  string
		event *Event
	}{
		{name: "invalid task", event: NewEvent(Task{}, testData.Status)},
		{name: "empty status", event: NewEvent(testData.Task, "")},
		{name: "blank status", event: NewEvent(testData.Task, " ")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Error(t, tt.event.Validate())
		})
	}
}

func TestValidateShouldReturnNilIfEventDataIsValid(t *testing.T) {
	testData := makeFakeDataEventEntity()
	assert.NoError(t, NewEvent(testData.Task, testData.Status).Validate())
}
