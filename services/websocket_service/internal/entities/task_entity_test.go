package entities

import (
	"testing"

	"github.com/douglasvolcato/binary-code-processor/websocket_service/test"
	"github.com/stretchr/testify/assert"
)

type testDataTaskEntity struct {
	ID         string
	Message    string
	BinaryCode string
	CreatedAt  string
	UpdatedAt  string
}

func makeFakeDataTaskEntity() *testDataTaskEntity {
	faker := test.FakeData{}
	return &testDataTaskEntity{
		ID:         faker.ID(),
		Message:    faker.Phrase(),
		BinaryCode: faker.Phrase(),
		CreatedAt:  faker.Date(),
		UpdatedAt:  faker.Date(),
	}
}

func TestNewTaskShouldCreateTask(t *testing.T) {
	testData := makeFakeDataTaskEntity()
	sut := NewTask(
		testData.ID,
		testData.Message,
		testData.BinaryCode,
		testData.CreatedAt,
		testData.UpdatedAt,
	)

	assert.NotNil(t, sut)
	assert.Equal(t, testData.ID, sut.ID)
	assert.Equal(t, testData.Message, sut.Message)
	assert.Equal(t, testData.BinaryCode, sut.BinaryCode)
	assert.Equal(t, testData.CreatedAt, sut.CreatedAt)
	assert.Equal(t, testData.UpdatedAt, sut.UpdatedAt)
}

func TestValidateShouldReturnErrorIfTaskDataIsInvalid(t *testing.T) {
	testData := makeFakeDataTaskEntity()
	tests := []struct {
		name string
		task *Task
	}{
		{name: "empty id", task: NewTask("", testData.Message, testData.BinaryCode, testData.CreatedAt, testData.UpdatedAt)},
		{name: "blank id", task: NewTask(" ", testData.Message, testData.BinaryCode, testData.CreatedAt, testData.UpdatedAt)},
		{name: "empty message", task: NewTask(testData.ID, "", testData.BinaryCode, testData.CreatedAt, testData.UpdatedAt)},
		{name: "blank message", task: NewTask(testData.ID, " ", testData.BinaryCode, testData.CreatedAt, testData.UpdatedAt)},
		{name: "empty binary code", task: NewTask(testData.ID, testData.Message, "", testData.CreatedAt, testData.UpdatedAt)},
		{name: "blank binary code", task: NewTask(testData.ID, testData.Message, " ", testData.CreatedAt, testData.UpdatedAt)},
		{name: "empty created at", task: NewTask(testData.ID, testData.Message, testData.BinaryCode, "", testData.UpdatedAt)},
		{name: "blank created at", task: NewTask(testData.ID, testData.Message, testData.BinaryCode, " ", testData.UpdatedAt)},
		{name: "empty updated at", task: NewTask(testData.ID, testData.Message, testData.BinaryCode, testData.CreatedAt, "")},
		{name: "blank updated at", task: NewTask(testData.ID, testData.Message, testData.BinaryCode, testData.CreatedAt, " ")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Error(t, tt.task.Validate())
		})
	}
}

func TestValidateShouldReturnNilIfTaskDataIsValid(t *testing.T) {
	testData := makeFakeDataTaskEntity()
	assert.NoError(t, NewTask(
		testData.ID,
		testData.Message,
		testData.BinaryCode,
		testData.CreatedAt,
		testData.UpdatedAt,
	).Validate())
}
