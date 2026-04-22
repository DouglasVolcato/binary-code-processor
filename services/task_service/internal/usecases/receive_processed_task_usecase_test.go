package usecases

import (
	"errors"
	"testing"

	"github.com/douglasvolcato/binary-code-processor/task_service/internal/entities"
	"github.com/stretchr/testify/assert"
)

type mockProcessorForReceiveProcessed struct {
	SetTaskAsProcessedCalls int
	SetTaskAsProcessedArgs  struct {
		TaskID string
	}
	SetTaskAsProcessedFunc func(taskID string) error
}

func (m *mockProcessorForReceiveProcessed) MoveTaskToProcessing(dto CreateTaskDTO) (entities.Task, error) {
	// not used in these tests
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

func TestNewReceiveProcessedTaskUseCase_Constructs(t *testing.T) {
	repo := &mockProcessorForReceiveProcessed{}
	sut := NewReceiveProcessedTaskUseCase(repo)

	assert.NotNil(t, sut)
	assert.Same(t, repo, sut.Repo)
}

func TestReceiveProcessedTaskExecute_Success(t *testing.T) {
	repo := &mockProcessorForReceiveProcessed{
		SetTaskAsProcessedFunc: func(taskID string) error { return nil },
	}
	sut := NewReceiveProcessedTaskUseCase(repo)

	input := &ReceiveProcessedTaskInput{ID: "task-1"}
	output, err := sut.Execute(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.True(t, output.Success)
	assert.Equal(t, 1, repo.SetTaskAsProcessedCalls)
	assert.Equal(t, "task-1", repo.SetTaskAsProcessedArgs.TaskID)
}

func TestReceiveProcessedTaskExecute_RepoError(t *testing.T) {
	expectedErr := errors.New("set processed failure")
	repo := &mockProcessorForReceiveProcessed{
		SetTaskAsProcessedFunc: func(taskID string) error { return expectedErr },
	}
	sut := NewReceiveProcessedTaskUseCase(repo)

	input := &ReceiveProcessedTaskInput{ID: "task-2"}
	output, err := sut.Execute(input)

	assert.Nil(t, output)
	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, 1, repo.SetTaskAsProcessedCalls)
}
