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

	SetTaskAsProcessedCalls int
	SetTaskAsProcessedFunc  func(taskID string) (entities.Task, error)
}

type mockOutboxForReceive struct {
	StoreUnprocessedEventCalls int
	StoreUnprocessedEventArgs  struct {
		Task entities.Task
	}
	StoreUnprocessedEventFunc func(task entities.Task) error

	StoreProcessedEventCalls int
	StoreProcessedEventFunc  func(task entities.Task) error
}

func (m *mockProcessorForReceive) MoveTaskToProcessing(dto CreateTaskDTO) (entities.Task, error) {
	m.MoveTaskToProcessingCalls++
	m.MoveTaskToProcessingArgs.DTO = dto
	if m.MoveTaskToProcessingFunc != nil {
		return m.MoveTaskToProcessingFunc(dto)
	}
	return entities.Task{}, nil
}

func (m *mockProcessorForReceive) SetTaskAsProcessed(taskID string) (entities.Task, error) {
	m.SetTaskAsProcessedCalls++
	if m.SetTaskAsProcessedFunc != nil {
		return m.SetTaskAsProcessedFunc(taskID)
	}
	return entities.Task{}, nil
}

func (m *mockOutboxForReceive) StoreUnprocessedEvent(task entities.Task) error {
	m.StoreUnprocessedEventCalls++
	m.StoreUnprocessedEventArgs.Task = task
	if m.StoreUnprocessedEventFunc != nil {
		return m.StoreUnprocessedEventFunc(task)
	}
	return nil
}

func (m *mockOutboxForReceive) StoreProcessedEvent(task entities.Task) error {
	m.StoreProcessedEventCalls++
	if m.StoreProcessedEventFunc != nil {
		return m.StoreProcessedEventFunc(task)
	}
	return nil
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
	outbox := &mockOutboxForReceive{}
	idGen := &mockIDGen{}
	sut := NewReceiveTaskToProcessUseCase(proc, outbox, idGen)

	assert.NotNil(t, sut)
	assert.Same(t, proc, sut.Repo)
	assert.Same(t, outbox, sut.Outbox)
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
	outbox := &mockOutboxForReceive{}
	idGen := &mockIDGen{GenerateIDFunc: func() string { return fakeTask.ID }}

	sut := NewReceiveTaskToProcessUseCase(proc, outbox, idGen)

	output, err := sut.Execute(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.True(t, output.Success)
	assert.Equal(t, fakeTask.ID, output.Task.ID)
	assert.Equal(t, input.Message, output.Task.Message)
	assert.Equal(t, 1, proc.MoveTaskToProcessingCalls)
	assert.Equal(t, fakeTask.ID, proc.MoveTaskToProcessingArgs.DTO.ID)
	assert.Equal(t, input.Message, proc.MoveTaskToProcessingArgs.DTO.Message)
	assert.Equal(t, 1, outbox.StoreUnprocessedEventCalls)
	assert.Equal(t, fakeTask.ID, outbox.StoreUnprocessedEventArgs.Task.ID)
	assert.Equal(t, input.Message, outbox.StoreUnprocessedEventArgs.Task.Message)
	assert.Equal(t, "", outbox.StoreUnprocessedEventArgs.Task.BinaryCode)
	assert.Equal(t, 1, idGen.GenerateIDCalls)
	assert.Equal(t, 0, proc.SetTaskAsProcessedCalls)
	assert.Equal(t, 0, outbox.StoreProcessedEventCalls)
}

func TestReceiveTaskToProcessExecuteShouldReturnErrorWhenRepoFails(t *testing.T) {
	faker := test.FakeData{}
	expectedErr := errors.New("move failure")
	proc := &mockProcessorForReceive{
		MoveTaskToProcessingFunc: func(dto CreateTaskDTO) (entities.Task, error) {
			return entities.Task{}, expectedErr
		},
	}
	outbox := &mockOutboxForReceive{}
	idGen := &mockIDGen{GenerateIDFunc: func() string { return faker.ID() }}

	sut := NewReceiveTaskToProcessUseCase(proc, outbox, idGen)
	input := &ReceiveTaskToProcessInput{Message: faker.Phrase()}

	output, err := sut.Execute(input)

	assert.Nil(t, output)
	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, 1, proc.MoveTaskToProcessingCalls)
	assert.Equal(t, input.Message, proc.MoveTaskToProcessingArgs.DTO.Message)
	assert.NotEmpty(t, proc.MoveTaskToProcessingArgs.DTO.ID)
	assert.Equal(t, 0, outbox.StoreUnprocessedEventCalls)
	assert.Equal(t, 0, outbox.StoreProcessedEventCalls)
	assert.Equal(t, 1, idGen.GenerateIDCalls)
	assert.Equal(t, 0, proc.SetTaskAsProcessedCalls)
}
