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
	SetTaskAsProcessedFunc    func(taskID string) error
	MoveTaskToProcessingCalls int
}

func (m *mockProcessorForReceiveProcessed) MoveTaskToProcessing(dto CreateTaskDTO) (entities.Task, error) {
	m.MoveTaskToProcessingCalls++
	return entities.Task{}, nil
}

func (m *mockProcessorForReceiveProcessed) SetTaskAsProcessed(taskID string) error {
	m.SetTaskAsProcessedCalls++
	m.SetTaskAsProcessedArgs.TaskID = taskID
	if m.SetTaskAsProcessedFunc != nil {
		return m.SetTaskAsProcessedFunc(taskID)
	}
	return nil
}

func TestNewReceiveProcessedTaskUseCaseShouldCreateReceiveProcessedTaskUseCase(t *testing.T) {
	repo := &mockProcessorForReceiveProcessed{}
	sut := NewReceiveProcessedTaskUseCase(repo)

	assert.NotNil(t, sut)
	assert.Same(t, repo, sut.Repo)
}

func TestReceiveProcessedTaskExecuteShouldReturnSuccess(t *testing.T) {
	faker := test.FakeData{}
	repo := &mockProcessorForReceiveProcessed{
		SetTaskAsProcessedFunc: func(taskID string) error { return nil },
	}
	sut := NewReceiveProcessedTaskUseCase(repo)
	input := &ReceiveProcessedTaskInput{ID: faker.ID()}
	output, err := sut.Execute(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.True(t, output.Success)
	assert.Equal(t, 1, repo.SetTaskAsProcessedCalls)
	assert.Equal(t, input.ID, repo.SetTaskAsProcessedArgs.TaskID)
	assert.Equal(t, 0, repo.MoveTaskToProcessingCalls)
}

func TestReceiveProcessedTaskExecuteShouldReturnErrorWhenRepoFails(t *testing.T) {
	faker := test.FakeData{}
	expectedErr := errors.New("set processed failure")
	repo := &mockProcessorForReceiveProcessed{
		SetTaskAsProcessedFunc: func(taskID string) error { return expectedErr },
	}
	sut := NewReceiveProcessedTaskUseCase(repo)

	input := &ReceiveProcessedTaskInput{ID: faker.ID()}
	output, err := sut.Execute(input)

	assert.Nil(t, output)
	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, 1, repo.SetTaskAsProcessedCalls)
	assert.Equal(t, 0, repo.MoveTaskToProcessingCalls)
}
