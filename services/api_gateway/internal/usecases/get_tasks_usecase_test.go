package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/douglasvolcato/binary-code-processor/api_gateway/internal/entities"
	"github.com/douglasvolcato/binary-code-processor/api_gateway/test"
	"github.com/stretchr/testify/assert"
)

type mockTaskRepository struct {
	GetTasksCalls int
	GetTasksArgs  struct {
		Limit  int
		Offset int
	}
	GetTasksFunc func(ctx context.Context, limit int, offset int) ([]entities.Task, error)
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

func makeFakeTasks(count int) []entities.Task {
	tasks := make([]entities.Task, 0, count)
	for i := 0; i < count; i++ {
		tasks = append(tasks, makeFakeTask())
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

	input := &GetTasksInput{
		Ctx:    context.Background(),
		Limit:  10,
		Offset: 5,
	}

	output, err := sut.Execute(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, expectedTasks, output.Tasks)
	assert.Equal(t, 1, repo.GetTasksCalls)
	assert.Equal(t, 10, repo.GetTasksArgs.Limit)
	assert.Equal(t, 5, repo.GetTasksArgs.Offset)
}

func TestGetTasksExecuteShouldUseDefaultLimitWhenInputLimitIsZero(t *testing.T) {
	expectedTasks := makeFakeTasks(1)

	repo := &mockTaskRepository{
		GetTasksFunc: func(ctx context.Context, limit int, offset int) ([]entities.Task, error) {
			return expectedTasks, nil
		},
	}
	sut := NewGetTasksUseCase(repo)

	input := &GetTasksInput{
		Ctx:    context.Background(),
		Limit:  0,
		Offset: 3,
	}

	output, err := sut.Execute(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, expectedTasks, output.Tasks)
	assert.Equal(t, 1, repo.GetTasksCalls)
	assert.Equal(t, 20, repo.GetTasksArgs.Limit)
	assert.Equal(t, 3, repo.GetTasksArgs.Offset)
	assert.Equal(t, 20, input.Limit)
}

func TestGetTasksExecuteShouldUseDefaultLimitWhenInputLimitIsGreaterThanMaximum(t *testing.T) {
	expectedTasks := makeFakeTasks(1)

	repo := &mockTaskRepository{
		GetTasksFunc: func(ctx context.Context, limit int, offset int) ([]entities.Task, error) {
			return expectedTasks, nil
		},
	}
	sut := NewGetTasksUseCase(repo)

	input := &GetTasksInput{
		Ctx:    context.Background(),
		Limit:  999,
		Offset: 7,
	}

	output, err := sut.Execute(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, expectedTasks, output.Tasks)
	assert.Equal(t, 1, repo.GetTasksCalls)
	assert.Equal(t, 20, repo.GetTasksArgs.Limit)
	assert.Equal(t, 7, repo.GetTasksArgs.Offset)
	assert.Equal(t, 20, input.Limit)
}

func TestGetTasksExecuteShouldReturnErrorWhenRepoFails(t *testing.T) {
	expectedError := errors.New("repo failure")
	repo := &mockTaskRepository{
		GetTasksFunc: func(ctx context.Context, limit int, offset int) ([]entities.Task, error) {
			return nil, expectedError
		},
	}
	sut := NewGetTasksUseCase(repo)

	input := &GetTasksInput{
		Ctx:    context.Background(),
		Limit:  5,
		Offset: 1,
	}

	output, err := sut.Execute(input)

	assert.Nil(t, output)
	assert.ErrorIs(t, err, expectedError)
	assert.Equal(t, 1, repo.GetTasksCalls)
	assert.Equal(t, 5, repo.GetTasksArgs.Limit)
	assert.Equal(t, 1, repo.GetTasksArgs.Offset)
}
