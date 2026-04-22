package usecases

import (
	"errors"
	"testing"

	"github.com/douglasvolcato/binary-code-processor/websocket_service/internal/entities"
	"github.com/douglasvolcato/binary-code-processor/websocket_service/test"
	"github.com/stretchr/testify/assert"
)

type mockWebSocketClient struct {
	SendProcessedTasksToClientCalls int
	SendProcessedTasksToClientArgs  struct {
		Task entities.Task
	}
	SendProcessedTasksToClientFunc func(task entities.Task) error
}

func (m *mockWebSocketClient) SendProcessedTasksToClient(task entities.Task) error {
	m.SendProcessedTasksToClientCalls++
	m.SendProcessedTasksToClientArgs.Task = task
	if m.SendProcessedTasksToClientFunc != nil {
		return m.SendProcessedTasksToClientFunc(task)
	}
	return nil
}

func makeFakeTask() entities.Task {
	faker := test.FakeData{}
	return entities.Task{
		ID:         faker.ID(),
		BinaryCode: faker.Phrase(),
	}
}

func TestNewSendProcessedTasksUseCaseShouldCreateSendProcessedTasksUseCase(t *testing.T) {
	webSocket := &mockWebSocketClient{}
	sut := NewSendProcessedTasksUseCase(webSocket)

	assert.NotNil(t, sut)
	assert.Same(t, webSocket, sut.WebSocket)
}

func TestSendProcessedTasksExecuteShouldSendTaskToClient(t *testing.T) {
	task := makeFakeTask()
	webSocket := &mockWebSocketClient{
		SendProcessedTasksToClientFunc: func(task entities.Task) error {
			return nil
		},
	}
	sut := NewSendProcessedTasksUseCase(webSocket)

	output, err := sut.Execute(&SendProcessedTasksInput{Task: task})

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, 1, webSocket.SendProcessedTasksToClientCalls)
	assert.Equal(t, task, webSocket.SendProcessedTasksToClientArgs.Task)
}

func TestSendProcessedTasksExecuteShouldReturnErrorWhenWebSocketFails(t *testing.T) {
	task := makeFakeTask()
	expectedErr := errors.New("websocket failure")
	webSocket := &mockWebSocketClient{
		SendProcessedTasksToClientFunc: func(task entities.Task) error {
			return expectedErr
		},
	}
	sut := NewSendProcessedTasksUseCase(webSocket)

	output, err := sut.Execute(&SendProcessedTasksInput{Task: task})

	assert.Nil(t, output)
	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, 1, webSocket.SendProcessedTasksToClientCalls)
	assert.Equal(t, task, webSocket.SendProcessedTasksToClientArgs.Task)
}
