package usecases

import (
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
	GetTasksFunc func(limit int, offset int) ([]entities.Task, error)
}

func (m *mockTaskRepository) GetTasks(limit int, offset int) ([]entities.Task, error) {
	m.GetTasksCalls++
	m.GetTasksArgs.Limit = limit
	m.GetTasksArgs.Offset = offset
	if m.GetTasksFunc != nil {
		return m.GetTasksFunc(limit, offset)
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
	for range count {
		tasks = append(tasks, makeFakeTask())
	}
	return tasks
}

func TestNewGetTasksUseCaseShouldReturnInstance(t *testing.T) {
	repo := &mockTaskRepository{}
	sut := NewGetTasksUseCase(repo)

	assert.NotNil(t, sut)
	assert.Same(t, repo, sut.Repo)
}

func TestExecuteShouldCallRepoWithInputValues(t *testing.T) {
	expectedTasks := makeFakeTasks(2)

	repo := &mockTaskRepository{
		GetTasksFunc: func(limit int, offset int) ([]entities.Task, error) {
			return expectedTasks, nil
		},
	}
	sut := NewGetTasksUseCase(repo)

	input := &GetTasksInput{
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

func TestExecuteShouldUseDefaultLimitWhenInputLimitIsZero(t *testing.T) {
	expectedTasks := makeFakeTasks(1)

	repo := &mockTaskRepository{
		GetTasksFunc: func(limit int, offset int) ([]entities.Task, error) {
			return expectedTasks, nil
		},
	}
	sut := NewGetTasksUseCase(repo)

	input := &GetTasksInput{
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
}

func TestExecuteShouldUseDefaultLimitWhenInputLimitIsGreaterThanMaximum(t *testing.T) {
	expectedTasks := makeFakeTasks(1)

	repo := &mockTaskRepository{
		GetTasksFunc: func(limit int, offset int) ([]entities.Task, error) {
			return expectedTasks, nil
		},
	}
	sut := NewGetTasksUseCase(repo)

	input := &GetTasksInput{
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
}

func TestExecuteShouldReturnErrorWhenRepoFails(t *testing.T) {
	expectedError := errors.New("repo failure")
	repo := &mockTaskRepository{
		GetTasksFunc: func(limit int, offset int) ([]entities.Task, error) {
			return nil, expectedError
		},
	}
	sut := NewGetTasksUseCase(repo)

	input := &GetTasksInput{
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
