package usecases

import (
	"errors"
	"testing"

	"github.com/douglasvolcato/binary-code-processor/api_gateway/test"
	"github.com/stretchr/testify/assert"
)

type mockTaskProcessor struct {
	SendTaskToProcessCalls int
	SendTaskToProcessArgs  struct {
		Messages []string
	}
	SendTaskToProcessFunc func(messages []string) error
}

func (m *mockTaskProcessor) SendTaskToProcess(messages []string) error {
	m.SendTaskToProcessCalls++
	m.SendTaskToProcessArgs.Messages = messages
	if m.SendTaskToProcessFunc != nil {
		return m.SendTaskToProcessFunc(messages)
	}
	return nil
}

func makeFakeMessages(count int) []string {
	faker := test.FakeData{}
	msgs := make([]string, 0, count)
	for i := 0; i < count; i++ {
		msgs = append(msgs, faker.Phrase())
	}
	return msgs
}

func TestNewSendTaskToProcessUseCaseShouldReturnInstance(t *testing.T) {
	repo := &mockTaskProcessor{}
	sut := NewSendTaskToProcessUseCase(repo)

	assert.NotNil(t, sut)
	assert.Same(t, repo, sut.Repo)
}

func TestSendTaskExecuteShouldCallRepoWithInputValues(t *testing.T) {
	expectedMessages := makeFakeMessages(3)

	repo := &mockTaskProcessor{
		SendTaskToProcessFunc: func(messages []string) error {
			return nil
		},
	}
	sut := NewSendTaskToProcessUseCase(repo)

	input := &SendTaskToProcessInput{
		Messages: expectedMessages,
	}

	output, err := sut.Execute(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.True(t, output.Success)
	assert.Equal(t, 1, repo.SendTaskToProcessCalls)
	assert.Equal(t, expectedMessages, repo.SendTaskToProcessArgs.Messages)
}

func TestSendTaskExecuteShouldReturnErrorWhenRepoFails(t *testing.T) {
	expectedError := errors.New("repo failure")
	repo := &mockTaskProcessor{
		SendTaskToProcessFunc: func(messages []string) error {
			return expectedError
		},
	}
	sut := NewSendTaskToProcessUseCase(repo)

	input := &SendTaskToProcessInput{
		Messages: makeFakeMessages(2),
	}

	output, err := sut.Execute(input)

	assert.Nil(t, output)
	assert.ErrorIs(t, err, expectedError)
	assert.Equal(t, 1, repo.SendTaskToProcessCalls)
}
