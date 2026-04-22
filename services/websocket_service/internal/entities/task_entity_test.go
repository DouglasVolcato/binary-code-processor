package entities

import (
	"testing"

	"github.com/douglasvolcato/binary-code-processor/websocket_service/test"
	"github.com/stretchr/testify/assert"
)

type testDataTaskEntity struct {
	ID         string
	BinaryCode string
}

func makeFakeDataTaskEntity() *testDataTaskEntity {
	faker := test.FakeData{}
	return &testDataTaskEntity{
		ID:         faker.ID(),
		BinaryCode: faker.Phrase(),
	}
}

func TestNewTaskShouldCreateTask(t *testing.T) {
	testData := makeFakeDataTaskEntity()
	sut := NewTask(testData.ID, testData.BinaryCode)

	assert.NotNil(t, sut)
	assert.Equal(t, testData.ID, sut.ID)
	assert.Equal(t, testData.BinaryCode, sut.BinaryCode)
}

func TestValidateShouldReturnErrorIfTaskDataIsInvalid(t *testing.T) {
	testData := makeFakeDataTaskEntity()
	tests := []struct {
		name string
		task *Task
	}{
		{
			name: "empty id",
			task: NewTask("", testData.BinaryCode),
		},
		{
			name: "blank id",
			task: NewTask(" ", testData.BinaryCode),
		},
		{
			name: "empty binary code",
			task: NewTask(testData.ID, ""),
		},
		{
			name: "blank binary code",
			task: NewTask(testData.ID, " "),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Error(t, tt.task.Validate())
		})
	}
}

func TestValidateShouldReturnNilIfTaskDataIsValid(t *testing.T) {
	testData := makeFakeDataTaskEntity()
	assert.NoError(t, NewTask(testData.ID, testData.BinaryCode).Validate())
}
