package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/douglasvolcato/binary-code-processor/task_service/internal/entities"
	"github.com/douglasvolcato/binary-code-processor/task_service/test"
	"github.com/stretchr/testify/assert"
)

type mockTaskRepository struct {
	GetTasksCalls int
	GetTasksArgs  struct {
		Limit  int
		Offset int
	}
	GetTasksFunc     func(ctx context.Context, limit int, offset int) ([]entities.Task, error)
	GetTaskByIDCalls int
	GetTaskByIDArgs  struct {
		TaskID string
	}
	GetTaskByIDFunc func(ctx context.Context, taskID string) (entities.Task, error)
}

func (m *mockTaskRepository) GetTasks(ctx context.Context, limit int, offset int) ([]entities.Task, error) {
	_ = ctx
	m.GetTasksCalls++
	m.GetTasksArgs.Limit = limit
	m.GetTasksArgs.Offset = offset
	if m.GetTasksFunc != nil {
		return m.GetTasksFunc(ctx, limit, offset)
	}
	return nil, nil
}

func (m *mockTaskRepository) GetTaskByID(ctx context.Context, taskID string) (entities.Task, error) {
	_ = ctx
	m.GetTaskByIDCalls++
	m.GetTaskByIDArgs.TaskID = taskID
	if m.GetTaskByIDFunc != nil {
		return m.GetTaskByIDFunc(ctx, taskID)
	}
	return entities.Task{}, nil
}

func makeFakeTasks(count int) []entities.Task {
	faker := test.FakeData{}
	tasks := make([]entities.Task, 0, count)
	for i := 0; i < count; i++ {
		tasks = append(tasks, entities.Task{
			ID:         faker.ID(),
			Message:    faker.Phrase(),
			BinaryCode: faker.Phrase(),
			CreatedAt:  faker.Date(),
			UpdatedAt:  faker.Date(),
		})
	}
	return tasks
}

func TestNewGetTasksUseCaseShouldCreateGetTasksUseCase(t *testing.T) {
	repo := &mockTaskRepository{}
	sut := NewGetTasksUseCase(repo)

	assert.NotNil(t, sut)
	assert.Same(t, repo, sut.Repo)
}

func TestGetTasksExecuteShouldReturnTasks(t *testing.T) {
	expectedTasks := makeFakeTasks(2)

	repo := &mockTaskRepository{
		GetTasksFunc: func(ctx context.Context, limit int, offset int) ([]entities.Task, error) {
			return expectedTasks, nil
		},
	}
	sut := NewGetTasksUseCase(repo)

	input := &GetTasksInput{Ctx: context.Background(), Limit: 10, Offset: 5}

	output, err := sut.Execute(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, expectedTasks, output.Tasks)
	assert.Equal(t, 1, repo.GetTasksCalls)
	assert.Equal(t, 10, repo.GetTasksArgs.Limit)
	assert.Equal(t, 5, repo.GetTasksArgs.Offset)
	assert.Equal(t, 0, repo.GetTaskByIDCalls)
}

func TestGetTasksExecuteShouldReturnErrorWhenRepoFails(t *testing.T) {
	expectedError := errors.New("repo failure")
	repo := &mockTaskRepository{
		GetTasksFunc: func(ctx context.Context, limit int, offset int) ([]entities.Task, error) {
			return nil, expectedError
		},
	}
	sut := NewGetTasksUseCase(repo)

	input := &GetTasksInput{Ctx: context.Background(), Limit: 5, Offset: 1}

	output, err := sut.Execute(input)

	assert.Nil(t, output)
	assert.ErrorIs(t, err, expectedError)
	assert.Equal(t, 1, repo.GetTasksCalls)
	assert.Equal(t, 5, repo.GetTasksArgs.Limit)
	assert.Equal(t, 1, repo.GetTasksArgs.Offset)
	assert.Equal(t, 0, repo.GetTaskByIDCalls)
}
