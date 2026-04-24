package usecases

import (
	"context"
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
	GetTaskByIDFunc func(ctx context.Context, taskID string) (entities.Task, error)
}

func (m *mockRepo) GetTaskByID(ctx context.Context, taskID string) (entities.Task, error) {
	_ = ctx
	m.GetTaskByIDCalls++
	m.GetTaskByIDArgs.ID = taskID
	if m.GetTaskByIDFunc != nil {
		return m.GetTaskByIDFunc(ctx, taskID)
	}
	return entities.Task{}, nil
}

type mockProcessor struct {
	FinishProcessingCalls int
	FinishProcessingArgs  struct {
		DTO FinishProcessingDTO
	}
	FinishProcessingFunc func(ctx context.Context, dto FinishProcessingDTO) error
}

func (m *mockProcessor) FinishProcessing(ctx context.Context, dto FinishProcessingDTO) error {
	_ = ctx
	m.FinishProcessingCalls++
	m.FinishProcessingArgs.DTO = dto
	if m.FinishProcessingFunc != nil {
		return m.FinishProcessingFunc(ctx, dto)
	}
	return nil
}

func makeFakeTask(message string) entities.Task {
	faker := test.FakeData{}
	return entities.Task{
		ID:         faker.ID(),
		Message:    message,
		BinaryCode: faker.Phrase(),
		CreatedAt:  faker.Date(),
		UpdatedAt:  faker.Date(),
	}
}

func TestNewProcessTaskUseCaseShouldCreateProcessTaskUseCase(t *testing.T) {
	repo := &mockRepo{}
	proc := &mockProcessor{}
	sut := NewProcessTaskUseCase(repo, proc)

	assert.NotNil(t, sut)
	assert.Same(t, repo, sut.Repo)
	assert.Same(t, proc, sut.Processor)
}

func TestProcessTaskExecuteShouldReturnBinaryCode(t *testing.T) {
	task := makeFakeTask("AB")

	repo := &mockRepo{
		GetTaskByIDFunc: func(ctx context.Context, taskID string) (entities.Task, error) {
			task.ID = taskID
			return task, nil
		},
	}
	proc := &mockProcessor{
		FinishProcessingFunc: func(ctx context.Context, dto FinishProcessingDTO) error {
			return nil
		},
	}

	sut := NewProcessTaskUseCase(repo, proc)
	input := &ProcessTaskInput{Ctx: context.Background(), ID: "task-123"}
	out, err := sut.Execute(input)

	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.Equal(t, "task-123", out.ID)
	assert.Equal(t, "0100000101000010", out.BinaryCode)
	assert.Equal(t, 1, repo.GetTaskByIDCalls)
	assert.Equal(t, "task-123", repo.GetTaskByIDArgs.ID)
	assert.Equal(t, 1, proc.FinishProcessingCalls)
	assert.Equal(t, "task-123", proc.FinishProcessingArgs.DTO.ID)
	assert.Equal(t, "0100000101000010", proc.FinishProcessingArgs.DTO.BinaryCode)
}

func TestProcessTaskExecuteShouldReturnErrorWhenGetTaskFails(t *testing.T) {
	expectedErr := errors.New("not found")
	repo := &mockRepo{
		GetTaskByIDFunc: func(ctx context.Context, taskID string) (entities.Task, error) {
			return entities.Task{}, expectedErr
		},
	}
	proc := &mockProcessor{}

	sut := NewProcessTaskUseCase(repo, proc)
	input := &ProcessTaskInput{Ctx: context.Background(), ID: "task-404"}
	out, err := sut.Execute(input)

	assert.Nil(t, out)
	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, 1, repo.GetTaskByIDCalls)
	assert.Equal(t, 0, proc.FinishProcessingCalls)
}

func TestProcessTaskExecuteShouldReturnErrorWhenTaskMessageIsEmpty(t *testing.T) {
	task := makeFakeTask("")

	repo := &mockRepo{
		GetTaskByIDFunc: func(ctx context.Context, taskID string) (entities.Task, error) {
			task.ID = taskID
			return task, nil
		},
	}
	proc := &mockProcessor{}

	sut := NewProcessTaskUseCase(repo, proc)
	input := &ProcessTaskInput{Ctx: context.Background(), ID: "task-empty"}
	out, err := sut.Execute(input)

	assert.Nil(t, out)
	assert.Error(t, err)
	assert.Equal(t, 1, repo.GetTaskByIDCalls)
	assert.Equal(t, "task-empty", repo.GetTaskByIDArgs.ID)
	assert.Equal(t, 0, proc.FinishProcessingCalls)
}

func TestProcessTaskExecuteShouldReturnErrorWhenFinishProcessingFails(t *testing.T) {
	task := makeFakeTask("AB")

	expectedErr := errors.New("finish failure")
	repo := &mockRepo{
		GetTaskByIDFunc: func(ctx context.Context, taskID string) (entities.Task, error) {
			task.ID = taskID
			return task, nil
		},
	}
	proc := &mockProcessor{
		FinishProcessingFunc: func(ctx context.Context, dto FinishProcessingDTO) error {
			return expectedErr
		},
	}

	sut := NewProcessTaskUseCase(repo, proc)
	input := &ProcessTaskInput{Ctx: context.Background(), ID: "task-500"}
	out, err := sut.Execute(input)

	assert.Nil(t, out)
	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, 1, repo.GetTaskByIDCalls)
	assert.Equal(t, 1, proc.FinishProcessingCalls)
}
