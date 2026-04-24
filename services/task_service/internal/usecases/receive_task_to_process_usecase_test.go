package usecases

import (
	"errors"
	"testing"

	"github.com/douglasvolcato/binary-code-processor/task_service/internal/entities"
	"github.com/douglasvolcato/binary-code-processor/task_service/test"
	"github.com/stretchr/testify/assert"
)

type mockProcessorForReceive struct {
	MoveTaskToProcessingCalls int
	MoveTaskToProcessingArgs  struct {
		DTO CreateTaskDTO
	}
	MoveTaskToProcessingFunc func(dto CreateTaskDTO) (entities.Task, error)

	FinishProcessingCalls int
	FinishProcessingFunc  func(dto FinishProcessingDTO) (entities.Task, error)
}

func (m *mockProcessorForReceive) MoveTaskToProcessing(dto CreateTaskDTO) (entities.Task, error) {
	m.MoveTaskToProcessingCalls++
	m.MoveTaskToProcessingArgs.DTO = dto
	if m.MoveTaskToProcessingFunc != nil {
		return m.MoveTaskToProcessingFunc(dto)
	}
	return entities.Task{}, nil
}

func (m *mockProcessorForReceive) FinishProcessing(dto FinishProcessingDTO) (entities.Task, error) {
	m.FinishProcessingCalls++
	if m.FinishProcessingFunc != nil {
		return m.FinishProcessingFunc(dto)
	}
	return entities.Task{}, nil
}

type mockIDGen struct {
	GenerateIDCalls int
	GenerateIDFunc  func() string
}

func (m *mockIDGen) GenerateID() string {
	m.GenerateIDCalls++
	if m.GenerateIDFunc != nil {
		return m.GenerateIDFunc()
	}
	return "fake-id"
}

func makeFakeTaskEntity() entities.Task {
	faker := test.FakeData{}
	return entities.Task{
		ID:         faker.ID(),
		Message:    faker.Phrase(),
		BinaryCode: "",
		CreatedAt:  faker.Date(),
		UpdatedAt:  faker.Date(),
	}
}

func TestNewReceiveTaskToProcessUseCaseShouldCreateReceiveTaskToProcessUseCase(t *testing.T) {
	proc := &mockProcessorForReceive{}
	idGen := &mockIDGen{}
	sut := NewReceiveTaskToProcessUseCase(proc, idGen)

	assert.NotNil(t, sut)
	assert.Same(t, proc, sut.Repo)
	assert.Same(t, idGen, sut.IDGen)
}

func TestReceiveTaskToProcessExecuteShouldReturnTask(t *testing.T) {
	faker := test.FakeData{}
	fakeTask := makeFakeTaskEntity()
	input := &ReceiveTaskToProcessInput{Message: faker.Phrase()}

	proc := &mockProcessorForReceive{
		MoveTaskToProcessingFunc: func(dto CreateTaskDTO) (entities.Task, error) {
			task := fakeTask
			task.ID = dto.ID
			task.Message = dto.Message
			return task, nil
		},
	}
	idGen := &mockIDGen{GenerateIDFunc: func() string { return fakeTask.ID }}

	sut := NewReceiveTaskToProcessUseCase(proc, idGen)

	output, err := sut.Execute(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.True(t, output.Success)
	assert.Equal(t, fakeTask.ID, output.Task.ID)
	assert.Equal(t, input.Message, output.Task.Message)
	assert.Equal(t, 1, proc.MoveTaskToProcessingCalls)
	assert.Equal(t, fakeTask.ID, proc.MoveTaskToProcessingArgs.DTO.ID)
	assert.Equal(t, input.Message, proc.MoveTaskToProcessingArgs.DTO.Message)
	assert.Equal(t, 1, idGen.GenerateIDCalls)
	assert.Equal(t, 0, proc.FinishProcessingCalls)
}

func TestReceiveTaskToProcessExecuteShouldReturnErrorWhenRepoFails(t *testing.T) {
	faker := test.FakeData{}
	expectedErr := errors.New("move failure")
	proc := &mockProcessorForReceive{
		MoveTaskToProcessingFunc: func(dto CreateTaskDTO) (entities.Task, error) {
			return entities.Task{}, expectedErr
		},
	}
	idGen := &mockIDGen{GenerateIDFunc: func() string { return faker.ID() }}

	sut := NewReceiveTaskToProcessUseCase(proc, idGen)
	input := &ReceiveTaskToProcessInput{Message: faker.Phrase()}

	output, err := sut.Execute(input)

	assert.Nil(t, output)
	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, 1, proc.MoveTaskToProcessingCalls)
	assert.Equal(t, input.Message, proc.MoveTaskToProcessingArgs.DTO.Message)
	assert.NotEmpty(t, proc.MoveTaskToProcessingArgs.DTO.ID)
	assert.Equal(t, 1, idGen.GenerateIDCalls)
	assert.Equal(t, 0, proc.FinishProcessingCalls)
}
