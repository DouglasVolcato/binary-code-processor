package usecases

import (
	"errors"
	"testing"

	"github.com/douglasvolcato/binary-code-processor/task_service/internal/entities"
	"github.com/douglasvolcato/binary-code-processor/task_service/test"
	"github.com/stretchr/testify/assert"
)

type mockTaskProcessor struct {
	ProcessTaskCalls int
	ProcessTaskArgs  struct {
		Task entities.Task
	}
	ProcessTaskFunc func(task entities.Task) (entities.Task, error)
}

func (m *mockTaskProcessor) ProcessTask(task entities.Task) (entities.Task, error) {
	m.ProcessTaskCalls++
	m.ProcessTaskArgs.Task = task
	if m.ProcessTaskFunc != nil {
		return m.ProcessTaskFunc(task)
	}
	return task, nil
}

func makeFakeTask() entities.Task {
	faker := test.FakeData{}
	return entities.Task{
		ID:         faker.ID(),
		Message:    faker.Phrase(),
		BinaryData: faker.Binary(),
		CreatedAt:  faker.Date(),
		UpdatedAt:  faker.Date(),
	}
}

func TestNewReceiveTaskToProcessUseCaseShouldReturnInstance(t *testing.T) {
	repo := &mockTaskProcessor{}
	sut := NewReceiveTaskToProcessUseCase(repo)

	assert.NotNil(t, sut)
	assert.Same(t, repo, sut.Repo)
}

func TestReceiveTaskExecuteShouldCallRepoAndReturnTask(t *testing.T) {
	inputTask := makeFakeTask()

	processedTask := inputTask
	processedTask.Message = processedTask.Message + " - processed"

	repo := &mockTaskProcessor{
		ProcessTaskFunc: func(task entities.Task) (entities.Task, error) {
			return processedTask, nil
		},
	}
	sut := NewReceiveTaskToProcessUseCase(repo)

	input := &ReceiveTaskToProcessInput{
		Task: inputTask,
	}

	output, err := sut.Execute(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.True(t, output.Success)
	assert.Equal(t, processedTask, output.Task)
	assert.Equal(t, 1, repo.ProcessTaskCalls)
	assert.Equal(t, inputTask, repo.ProcessTaskArgs.Task)
}

func TestReceiveTaskExecuteShouldReturnErrorWhenRepoFails(t *testing.T) {
	expectedError := errors.New("process failure")
	repo := &mockTaskProcessor{
		ProcessTaskFunc: func(task entities.Task) (entities.Task, error) {
			return entities.Task{}, expectedError
		},
	}
	sut := NewReceiveTaskToProcessUseCase(repo)

	input := &ReceiveTaskToProcessInput{
		Task: makeFakeTask(),
	}

	output, err := sut.Execute(input)

	assert.Nil(t, output)
	assert.ErrorIs(t, err, expectedError)
	assert.Equal(t, 1, repo.ProcessTaskCalls)
}
