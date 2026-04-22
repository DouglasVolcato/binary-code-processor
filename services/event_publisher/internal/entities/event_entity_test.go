package entities

import (
	"testing"

	"github.com/douglasvolcato/binary-code-processor/event_publisher/test"
	"github.com/stretchr/testify/assert"
)

type testDataEventEntity struct {
	ID     string
	Status string
}

func makeFakeDataEventEntity() *testDataEventEntity {
	faker := test.FakeData{}
	return &testDataEventEntity{
		ID:     faker.ID(),
		Status: faker.Word(),
	}
}

func TestNewEventShouldCreateEvent(t *testing.T) {
	testData := makeFakeDataEventEntity()
	sut := NewEvent(testData.ID, testData.Status)

	assert.NotNil(t, sut)
	assert.Equal(t, testData.ID, sut.ID)
	assert.Equal(t, testData.Status, sut.Status)
}

func TestValidateShouldReturnErrorIfEventDataIsInvalid(t *testing.T) {
	testData := makeFakeDataEventEntity()
	tests := []struct {
		name  string
		event *Event
		err   string
	}{
		{
			name:  "empty id",
			event: NewEvent("", testData.Status),
			err:   "invalid event ID",
		},
		{
			name:  "blank id",
			event: NewEvent(" ", testData.Status),
			err:   "invalid event ID",
		},
		{
			name:  "empty status",
			event: NewEvent(testData.ID, ""),
			err:   "invalid event status",
		},
		{
			name:  "blank status",
			event: NewEvent(testData.ID, " "),
			err:   "invalid event status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.EqualError(t, tt.event.Validate(), tt.err)
		})
	}
}

func TestValidateShouldReturnNilIfEventDataIsValid(t *testing.T) {
	testData := makeFakeDataEventEntity()
	assert.NoError(t, NewEvent(testData.ID, testData.Status).Validate())
}
