package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/douglasvolcato/binary-code-processor/task_service/internal/entities"
	"github.com/douglasvolcato/binary-code-processor/task_service/test"
	"github.com/stretchr/testify/assert"
)

type mockRepoForGetByID struct {
	GetTaskByIDCalls int
	GetTaskByIDArgs  struct {
		TaskID string
	}
	GetTaskByIDFunc func(ctx context.Context, taskID string) (entities.Task, error)
	GetTasksCalls   int
}

func (m *mockRepoForGetByID) GetTasks(ctx context.Context, limit int, offset int) ([]entities.Task, error) {
	_ = ctx
	m.GetTasksCalls++
	return nil, nil
}

func (m *mockRepoForGetByID) GetTaskByID(ctx context.Context, taskID string) (entities.Task, error) {
	_ = ctx
	m.GetTaskByIDCalls++
	m.GetTaskByIDArgs.TaskID = taskID
	if m.GetTaskByIDFunc != nil {
		return m.GetTaskByIDFunc(ctx, taskID)
	}
	return entities.Task{}, nil
}

func makeFakeTask() entities.Task {
	faker := test.FakeData{}
	return entities.Task{
		ID:         faker.ID(),
		Message:    faker.Phrase(),
		BinaryCode: faker.Phrase(),
		CreatedAt:  faker.Date(),
		UpdatedAt:  faker.Date(),
	}
}

func TestNewGetTaskByIDUseCaseShouldCreateGetTaskByIDUseCase(t *testing.T) {
	repo := &mockRepoForGetByID{}
	sut := NewGetTaskByIDUseCase(repo)

	assert.NotNil(t, sut)
	assert.Same(t, repo, sut.Repo)
}

func TestGetTaskByIDExecuteShouldReturnTask(t *testing.T) {
	task := makeFakeTask()
	repo := &mockRepoForGetByID{
		GetTaskByIDFunc: func(ctx context.Context, taskID string) (entities.Task, error) {
			return task, nil
		},
	}
	sut := NewGetTaskByIDUseCase(repo)

	input := &GetTaskByIDInput{Ctx: context.Background(), ID: task.ID}
	output, err := sut.Execute(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, task, output.Task)
	assert.Equal(t, 1, repo.GetTaskByIDCalls)
	assert.Equal(t, task.ID, repo.GetTaskByIDArgs.TaskID)
	assert.Equal(t, 0, repo.GetTasksCalls)
}

func TestGetTaskByIDExecuteShouldReturnErrorWhenRepoFails(t *testing.T) {
	expectedErr := errors.New("not found")
	repo := &mockRepoForGetByID{
		GetTaskByIDFunc: func(ctx context.Context, taskID string) (entities.Task, error) {
			return entities.Task{}, expectedErr
		},
	}
	sut := NewGetTaskByIDUseCase(repo)

	input := &GetTaskByIDInput{Ctx: context.Background(), ID: "missing"}
	output, err := sut.Execute(input)

	assert.Nil(t, output)
	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, 1, repo.GetTaskByIDCalls)
	assert.Equal(t, 0, repo.GetTasksCalls)
}
