package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/douglasvolcato/binary-code-processor/api_gateway/internal/entities"
	"github.com/douglasvolcato/binary-code-processor/api_gateway/test"
	"github.com/stretchr/testify/assert"
)

type mockTaskProcessor struct {
	SendTaskToProcessCalls int
	SendTaskToProcessArgs  struct {
		Messages []string
	}
	SendTaskToProcessFunc func(ctx context.Context, messages []string) ([]entities.Task, error)
}

func (m *mockTaskProcessor) SendTaskToProcess(ctx context.Context, messages []string) ([]entities.Task, error) {
	_ = ctx
	m.SendTaskToProcessCalls++
	m.SendTaskToProcessArgs.Messages = messages
	if m.SendTaskToProcessFunc != nil {
		return m.SendTaskToProcessFunc(ctx, messages)
	}
	return nil, nil
}

func makeFakeMessages(count int) []string {
	faker := test.FakeData{}
	msgs := make([]string, 0, count)
	for i := 0; i < count; i++ {
		msgs = append(msgs, faker.Phrase())
	}
	return msgs
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

func TestNewSendTaskToProcessUseCaseShouldCreateSendTaskToProcessUseCase(t *testing.T) {
	repo := &mockTaskProcessor{}
	sut := NewSendTaskToProcessUseCase(repo)

	assert.NotNil(t, sut)
	assert.Same(t, repo, sut.Repo)
}

func TestSendTaskToProcessExecuteShouldReturnTasks(t *testing.T) {
	expectedMessages := makeFakeMessages(3)
	expectedTasks := []entities.Task{
		makeFakeTaskEntity(),
		makeFakeTaskEntity(),
	}

	repo := &mockTaskProcessor{
		SendTaskToProcessFunc: func(ctx context.Context, messages []string) ([]entities.Task, error) {
			return expectedTasks, nil
		},
	}
	sut := NewSendTaskToProcessUseCase(repo)

	input := &SendTaskToProcessInput{
		Ctx:      context.Background(),
		Messages: expectedMessages,
	}

	output, err := sut.Execute(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.True(t, output.Success)
	assert.Equal(t, expectedTasks, output.Tasks)
	assert.Equal(t, 1, repo.SendTaskToProcessCalls)
	assert.Equal(t, expectedMessages, repo.SendTaskToProcessArgs.Messages)
}

func TestSendTaskToProcessExecuteShouldReturnErrorWhenRepoFails(t *testing.T) {
	expectedError := errors.New("repo failure")
	repo := &mockTaskProcessor{
		SendTaskToProcessFunc: func(ctx context.Context, messages []string) ([]entities.Task, error) {
			return nil, expectedError
		},
	}
	sut := NewSendTaskToProcessUseCase(repo)

	input := &SendTaskToProcessInput{
		Ctx:      context.Background(),
		Messages: makeFakeMessages(2),
	}

	output, err := sut.Execute(input)

	assert.Nil(t, output)
	assert.ErrorIs(t, err, expectedError)
	assert.Equal(t, 1, repo.SendTaskToProcessCalls)
}
