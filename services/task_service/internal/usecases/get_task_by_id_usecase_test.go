package usecases

import (
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
	GetTaskByIDFunc func(taskID string) (entities.Task, error)
}

func (m *mockRepoForGetByID) GetTasks(limit int, offset int) ([]entities.Task, error) {
	return nil, nil
}

func (m *mockRepoForGetByID) GetTaskByID(taskID string) (entities.Task, error) {
	m.GetTaskByIDCalls++
	m.GetTaskByIDArgs.TaskID = taskID
	if m.GetTaskByIDFunc != nil {
		return m.GetTaskByIDFunc(taskID)
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

func TestNewGetTaskByIDUseCase_Constructs(t *testing.T) {
	repo := &mockRepoForGetByID{}
	sut := NewGetTaskByIDUseCase(repo)

	assert.NotNil(t, sut)
	assert.Same(t, repo, sut.Repo)
}

func TestGetTaskByIDExecute_Success(t *testing.T) {
	task := makeFakeTask()
	repo := &mockRepoForGetByID{
		GetTaskByIDFunc: func(taskID string) (entities.Task, error) {
			return task, nil
		},
	}
	sut := NewGetTaskByIDUseCase(repo)

	input := &GetTaskByIDInput{ID: task.ID}
	output, err := sut.Execute(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, task, output.Task)
	assert.Equal(t, 1, repo.GetTaskByIDCalls)
	assert.Equal(t, task.ID, repo.GetTaskByIDArgs.TaskID)
}

func TestGetTaskByIDExecute_RepoError(t *testing.T) {
	expectedErr := errors.New("not found")
	repo := &mockRepoForGetByID{
		GetTaskByIDFunc: func(taskID string) (entities.Task, error) {
			return entities.Task{}, expectedErr
		},
	}
	sut := NewGetTaskByIDUseCase(repo)

	input := &GetTaskByIDInput{ID: "missing"}
	output, err := sut.Execute(input)

	assert.Nil(t, output)
	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, 1, repo.GetTaskByIDCalls)
}
