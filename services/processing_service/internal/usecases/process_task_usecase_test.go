package usecases

import (
	"errors"
	"testing"

	"github.com/douglasvolcato/binary-code-processor/processing_service/internal/entities"
	"github.com/douglasvolcato/binary-code-processor/processing_service/test"
	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	GetTaskByIDCalls int
	GetTaskByIDArgs  struct {
		ID string
	}
	GetTaskByIDFunc func(taskID string) (entities.Task, error)
}

func (m *mockRepo) GetTaskByID(taskID string) (entities.Task, error) {
	m.GetTaskByIDCalls++
	m.GetTaskByIDArgs.ID = taskID
	if m.GetTaskByIDFunc != nil {
		return m.GetTaskByIDFunc(taskID)
	}
	return entities.Task{}, nil
}

type mockProcessor struct {
	FinishProcessingCalls int
	FinishProcessingArgs  struct {
		DTO FinishProcessingDTO
	}
	FinishProcessingFunc func(dto FinishProcessingDTO) error
}

func (m *mockProcessor) FinishProcessing(dto FinishProcessingDTO) error {
	m.FinishProcessingCalls++
	m.FinishProcessingArgs.DTO = dto
	if m.FinishProcessingFunc != nil {
		return m.FinishProcessingFunc(dto)
	}
	return nil
}

func makeFakeTask() entities.Task {
	faker := test.FakeData{}
	return entities.Task{
		ID:         faker.ID(),
		Message:    "A",
		BinaryCode: faker.Phrase(),
		CreatedAt:  faker.Date(),
		UpdatedAt:  faker.Date(),
	}
}

func TestNewProcessTaskUseCase_ShouldConstruct(t *testing.T) {
	repo := &mockRepo{}
	proc := &mockProcessor{}
	sut := NewProcessTaskUseCase(repo, proc)

	assert.NotNil(t, sut)
	assert.Same(t, repo, sut.Repo)
	assert.Same(t, proc, sut.Processor)
}

func TestExecute_ShouldReturnOutputOnSuccess(t *testing.T) {
	task := makeFakeTask()
	task.Message = "A" // ASCII 65 -> 01000001

	repo := &mockRepo{
		GetTaskByIDFunc: func(taskID string) (entities.Task, error) {
			task.ID = taskID
			return task, nil
		},
	}
	proc := &mockProcessor{
		FinishProcessingFunc: func(dto FinishProcessingDTO) error {
			return nil
		},
	}

	sut := NewProcessTaskUseCase(repo, proc)
	input := &ProcessTaskInput{ID: "task-123"}
	out, err := sut.Execute(input)

	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.Equal(t, "task-123", out.ID)
	assert.Equal(t, "01000001", out.BinaryCode)
	assert.Equal(t, 1, repo.GetTaskByIDCalls)
	assert.Equal(t, "task-123", repo.GetTaskByIDArgs.ID)
	assert.Equal(t, 1, proc.FinishProcessingCalls)
	assert.Equal(t, "task-123", proc.FinishProcessingArgs.DTO.ID)
	assert.Equal(t, "01000001", proc.FinishProcessingArgs.DTO.BinaryCode)
}

func TestExecute_ShouldReturnErrorWhenGetTaskFails(t *testing.T) {
	expectedErr := errors.New("not found")
	repo := &mockRepo{
		GetTaskByIDFunc: func(taskID string) (entities.Task, error) {
			return entities.Task{}, expectedErr
		},
	}
	proc := &mockProcessor{}

	sut := NewProcessTaskUseCase(repo, proc)
	input := &ProcessTaskInput{ID: "task-404"}
	out, err := sut.Execute(input)

	assert.Nil(t, out)
	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, 1, repo.GetTaskByIDCalls)
	assert.Equal(t, 0, proc.FinishProcessingCalls)
}

func TestExecute_ShouldReturnErrorWhenFinishProcessingFails(t *testing.T) {
	task := makeFakeTask()
	task.Message = "A"

	expectedErr := errors.New("finish failure")
	repo := &mockRepo{
		GetTaskByIDFunc: func(taskID string) (entities.Task, error) {
			task.ID = taskID
			return task, nil
		},
	}
	proc := &mockProcessor{
		FinishProcessingFunc: func(dto FinishProcessingDTO) error {
			return expectedErr
		},
	}

	sut := NewProcessTaskUseCase(repo, proc)
	input := &ProcessTaskInput{ID: "task-500"}
	out, err := sut.Execute(input)

	assert.Nil(t, out)
	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, 1, repo.GetTaskByIDCalls)
	assert.Equal(t, 1, proc.FinishProcessingCalls)
}
