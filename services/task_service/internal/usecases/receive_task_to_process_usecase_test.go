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

	// implement SetTaskAsProcessed so it satisfies TaskProcessorInterface if needed elsewhere
	SetTaskAsProcessedFunc func(taskID string) error
}

func (m *mockProcessorForReceive) MoveTaskToProcessing(dto CreateTaskDTO) (entities.Task, error) {
	m.MoveTaskToProcessingCalls++
	m.MoveTaskToProcessingArgs.DTO = dto
	if m.MoveTaskToProcessingFunc != nil {
		return m.MoveTaskToProcessingFunc(dto)
	}
	return entities.Task{}, nil
}

func (m *mockProcessorForReceive) SetTaskAsProcessed(taskID string) error {
	if m.SetTaskAsProcessedFunc != nil {
		return m.SetTaskAsProcessedFunc(taskID)
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
		BinaryData: faker.Binary(),
		CreatedAt:  faker.Date(),
		UpdatedAt:  faker.Date(),
	}
}

func TestNewReceiveTaskToProcessUseCase_ConstructsWithDeps(t *testing.T) {
	proc := &mockProcessorForReceive{}
	idGen := &mockIDGen{}
	sut := NewReceiveTaskToProcessUseCase(proc, idGen)

	assert.NotNil(t, sut)
	assert.Same(t, proc, sut.Repo)
	assert.Same(t, idGen, sut.IDGen)
}

func TestReceiveTaskToProcessExecute_Success(t *testing.T) {
	fakeTask := makeFakeTaskEntity()
	id := "generated-id-123"

	proc := &mockProcessorForReceive{
		MoveTaskToProcessingFunc: func(dto CreateTaskDTO) (entities.Task, error) {
			// return a task with the dto values
			t := fakeTask
			t.ID = dto.ID
			t.Message = dto.Message
			return t, nil
		},
	}
	idGen := &mockIDGen{GenerateIDFunc: func() string { return id }}

	sut := NewReceiveTaskToProcessUseCase(proc, idGen)

	input := &ReceiveTaskToProcessInput{Message: "hello"}
	output, err := sut.Execute(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.True(t, output.Success)
	assert.Equal(t, id, output.Task.ID)
	assert.Equal(t, "hello", output.Task.Message)
	assert.Equal(t, 1, proc.MoveTaskToProcessingCalls)
	assert.Equal(t, 1, idGen.GenerateIDCalls)
}

func TestReceiveTaskToProcessExecute_RepoError(t *testing.T) {
	expectedErr := errors.New("move failure")
	proc := &mockProcessorForReceive{
		MoveTaskToProcessingFunc: func(dto CreateTaskDTO) (entities.Task, error) {
			return entities.Task{}, expectedErr
		},
	}
	idGen := &mockIDGen{GenerateIDFunc: func() string { return "id" }}

	sut := NewReceiveTaskToProcessUseCase(proc, idGen)
	input := &ReceiveTaskToProcessInput{Message: "msg"}

	output, err := sut.Execute(input)

	assert.Nil(t, output)
	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, 1, proc.MoveTaskToProcessingCalls)
}
