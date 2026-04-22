package usecases

import (
	"errors"
	"testing"

	"github.com/douglasvolcato/binary-code-processor/task_service/internal/entities"
	"github.com/douglasvolcato/binary-code-processor/task_service/test"
	"github.com/stretchr/testify/assert"
)

type mockProcessorForReceiveProcessed struct {
	SetTaskAsProcessedCalls int
	SetTaskAsProcessedArgs  struct {
		TaskID string
	}
	SetTaskAsProcessedFunc    func(taskID string) (entities.Task, error)
	MoveTaskToProcessingCalls int
}

type mockOutboxForReceiveProcessed struct {
	StoreUnprocessedEventCalls int
	StoreUnprocessedEventFunc   func(task entities.Task) error

	StoreProcessedEventCalls int
	StoreProcessedEventArgs  struct {
		Task entities.Task
	}
	StoreProcessedEventFunc func(task entities.Task) error
}

func (m *mockProcessorForReceiveProcessed) MoveTaskToProcessing(dto CreateTaskDTO) (entities.Task, error) {
	m.MoveTaskToProcessingCalls++
	return entities.Task{}, nil
}

func (m *mockProcessorForReceiveProcessed) SetTaskAsProcessed(taskID string) (entities.Task, error) {
	m.SetTaskAsProcessedCalls++
	m.SetTaskAsProcessedArgs.TaskID = taskID
	if m.SetTaskAsProcessedFunc != nil {
		return m.SetTaskAsProcessedFunc(taskID)
	}
	return entities.Task{}, nil
}

func (m *mockOutboxForReceiveProcessed) StoreUnprocessedEvent(task entities.Task) error {
	m.StoreUnprocessedEventCalls++
	if m.StoreUnprocessedEventFunc != nil {
		return m.StoreUnprocessedEventFunc(task)
	}
	return nil
}

func (m *mockOutboxForReceiveProcessed) StoreProcessedEvent(task entities.Task) error {
	m.StoreProcessedEventCalls++
	m.StoreProcessedEventArgs.Task = task
	if m.StoreProcessedEventFunc != nil {
		return m.StoreProcessedEventFunc(task)
	}
	return nil
}

func TestNewReceiveProcessedTaskUseCaseShouldCreateReceiveProcessedTaskUseCase(t *testing.T) {
	repo := &mockProcessorForReceiveProcessed{}
	outbox := &mockOutboxForReceiveProcessed{}
	sut := NewReceiveProcessedTaskUseCase(repo, outbox)

	assert.NotNil(t, sut)
	assert.Same(t, repo, sut.Repo)
	assert.Same(t, outbox, sut.Outbox)
}

func TestReceiveProcessedTaskExecuteShouldReturnSuccess(t *testing.T) {
	faker := test.FakeData{}
	task := entities.Task{
		ID:         faker.ID(),
		Message:    faker.Phrase(),
		BinaryCode: faker.Phrase(),
		CreatedAt:  faker.Date(),
		UpdatedAt:  faker.Date(),
	}
	repo := &mockProcessorForReceiveProcessed{
		SetTaskAsProcessedFunc: func(taskID string) (entities.Task, error) {
			task.ID = taskID
			return task, nil
		},
	}
	outbox := &mockOutboxForReceiveProcessed{}
	sut := NewReceiveProcessedTaskUseCase(repo, outbox)
	input := &ReceiveProcessedTaskInput{ID: faker.ID()}
	output, err := sut.Execute(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.True(t, output.Success)
	assert.Equal(t, 1, repo.SetTaskAsProcessedCalls)
	assert.Equal(t, input.ID, repo.SetTaskAsProcessedArgs.TaskID)
	assert.Equal(t, 1, outbox.StoreProcessedEventCalls)
	assert.Equal(t, input.ID, outbox.StoreProcessedEventArgs.Task.ID)
	assert.Equal(t, 0, repo.MoveTaskToProcessingCalls)
	assert.Equal(t, 0, outbox.StoreUnprocessedEventCalls)
}

func TestReceiveProcessedTaskExecuteShouldReturnErrorWhenRepoFails(t *testing.T) {
	faker := test.FakeData{}
	expectedErr := errors.New("set processed failure")
	repo := &mockProcessorForReceiveProcessed{
		SetTaskAsProcessedFunc: func(taskID string) (entities.Task, error) {
			return entities.Task{}, expectedErr
		},
	}
	outbox := &mockOutboxForReceiveProcessed{}
	sut := NewReceiveProcessedTaskUseCase(repo, outbox)

	input := &ReceiveProcessedTaskInput{ID: faker.ID()}
	output, err := sut.Execute(input)

	assert.Nil(t, output)
	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, 1, repo.SetTaskAsProcessedCalls)
	assert.Equal(t, 0, repo.MoveTaskToProcessingCalls)
	assert.Equal(t, 0, outbox.StoreProcessedEventCalls)
	assert.Equal(t, 0, outbox.StoreUnprocessedEventCalls)
}
