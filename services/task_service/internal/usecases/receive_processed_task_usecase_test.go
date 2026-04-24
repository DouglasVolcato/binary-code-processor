package usecases

import (
	"errors"
	"testing"

	"github.com/douglasvolcato/binary-code-processor/task_service/internal/entities"
	"github.com/douglasvolcato/binary-code-processor/task_service/test"
	"github.com/stretchr/testify/assert"
)

type mockProcessorForReceiveProcessed struct {
	MoveTaskToProcessingCalls int
	MoveTaskToProcessingFunc   func(dto CreateTaskDTO) (entities.Task, error)

	FinishProcessingCalls int
	FinishProcessingArgs  struct {
		DTO FinishProcessingDTO
	}
	FinishProcessingFunc func(dto FinishProcessingDTO) (entities.Task, error)
}

func (m *mockProcessorForReceiveProcessed) MoveTaskToProcessing(dto CreateTaskDTO) (entities.Task, error) {
	m.MoveTaskToProcessingCalls++
	if m.MoveTaskToProcessingFunc != nil {
		return m.MoveTaskToProcessingFunc(dto)
	}
	return entities.Task{}, nil
}

func (m *mockProcessorForReceiveProcessed) FinishProcessing(dto FinishProcessingDTO) (entities.Task, error) {
	m.FinishProcessingCalls++
	m.FinishProcessingArgs.DTO = dto
	if m.FinishProcessingFunc != nil {
		return m.FinishProcessingFunc(dto)
	}
	return entities.Task{}, nil
}

func TestNewReceiveProcessedTaskUseCaseShouldCreateReceiveProcessedTaskUseCase(t *testing.T) {
	repo := &mockProcessorForReceiveProcessed{}
	sut := NewReceiveProcessedTaskUseCase(repo)

	assert.NotNil(t, sut)
	assert.Same(t, repo, sut.Repo)
}

func TestReceiveProcessedTaskExecuteShouldReturnSuccess(t *testing.T) {
	faker := test.FakeData{}
	task := entities.Task{
		ID:         faker.ID(),
		Message:    faker.Phrase(),
		BinaryCode: faker.Phrase(),
		CreatedAt:  faker.Date(),
		UpdatedAt:  faker.Date(),
	}
	repo := &mockProcessorForReceiveProcessed{
		FinishProcessingFunc: func(dto FinishProcessingDTO) (entities.Task, error) {
			task.ID = dto.ID
			task.BinaryCode = dto.BinaryCode
			return task, nil
		},
	}
	sut := NewReceiveProcessedTaskUseCase(repo)
	input := &ReceiveProcessedTaskInput{ID: faker.ID(), BinaryCode: task.BinaryCode}
	output, err := sut.Execute(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.True(t, output.Success)
	assert.Equal(t, 1, repo.FinishProcessingCalls)
	assert.Equal(t, input.ID, repo.FinishProcessingArgs.DTO.ID)
	assert.Equal(t, input.BinaryCode, repo.FinishProcessingArgs.DTO.BinaryCode)
	assert.Equal(t, 0, repo.MoveTaskToProcessingCalls)
}

func TestReceiveProcessedTaskExecuteShouldReturnErrorWhenRepoFails(t *testing.T) {
	faker := test.FakeData{}
	expectedErr := errors.New("set processed failure")
	repo := &mockProcessorForReceiveProcessed{
		FinishProcessingFunc: func(dto FinishProcessingDTO) (entities.Task, error) {
			return entities.Task{}, expectedErr
		},
	}
	sut := NewReceiveProcessedTaskUseCase(repo)

	input := &ReceiveProcessedTaskInput{ID: faker.ID(), BinaryCode: faker.Phrase()}
	output, err := sut.Execute(input)

	assert.Nil(t, output)
	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, 1, repo.FinishProcessingCalls)
	assert.Equal(t, 0, repo.MoveTaskToProcessingCalls)
}
