package usecases

import (
	"errors"
	"testing"

	"github.com/douglasvolcato/binary-code-processor/event_publisher/internal/entities"
	"github.com/douglasvolcato/binary-code-processor/event_publisher/test"
	"github.com/stretchr/testify/assert"
)

type mockEventRepository struct {
	GetUnpublishedEventsCalls int
	GetUnpublishedEventsArgs  struct {
		Limit  int
		Offset int
	}
	GetUnpublishedEventsFunc func(limit int, offset int) ([]entities.Event, error)
}

func (m *mockEventRepository) GetUnpublishedEvents(limit int, offset int) ([]entities.Event, error) {
	m.GetUnpublishedEventsCalls++
	m.GetUnpublishedEventsArgs.Limit = limit
	m.GetUnpublishedEventsArgs.Offset = offset
	if m.GetUnpublishedEventsFunc != nil {
		return m.GetUnpublishedEventsFunc(limit, offset)
	}
	return nil, nil
}

type mockEventProcessor struct {
	SendEventToProcessCalls int
	SendEventToProcessArgs  struct {
		Event entities.Event
	}
	SendEventToProcessFunc func(event entities.Event) error
}

func (m *mockEventProcessor) SendEventToProcess(event entities.Event) error {
	m.SendEventToProcessCalls++
	m.SendEventToProcessArgs.Event = event
	if m.SendEventToProcessFunc != nil {
		return m.SendEventToProcessFunc(event)
	}
	return nil
}

func makeFakeEvent() entities.Event {
	faker := test.FakeData{}
	return entities.Event{
		ID:     faker.ID(),
		Status: faker.Word(),
	}
}

func makeFakeEvents(count int) []entities.Event {
	events := make([]entities.Event, 0, count)
	for i := 0; i < count; i++ {
		events = append(events, makeFakeEvent())
	}
	return events
}

func TestNewProcessUnpublishedEventsUseCaseShouldCreateProcessUnpublishedEventsUseCase(t *testing.T) {
	repo := &mockEventRepository{}
	processor := &mockEventProcessor{}
	sut := NewProcessUnpublishedEventsUseCase(repo, processor)

	assert.NotNil(t, sut)
	assert.Same(t, repo, sut.Repo)
	assert.Same(t, processor, sut.Processor)
}

func TestProcessUnpublishedEventsExecuteShouldProcessAllEvents(t *testing.T) {
	events := makeFakeEvents(2)
	repo := &mockEventRepository{
		GetUnpublishedEventsFunc: func(limit int, offset int) ([]entities.Event, error) {
			return events, nil
		},
	}
	processor := &mockEventProcessor{}
	sut := NewProcessUnpublishedEventsUseCase(repo, processor)

	output, err := sut.Execute(&ProcessUnpublishedEventsInput{})

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, 1, repo.GetUnpublishedEventsCalls)
	assert.Equal(t, 100, repo.GetUnpublishedEventsArgs.Limit)
	assert.Equal(t, 0, repo.GetUnpublishedEventsArgs.Offset)
	assert.Equal(t, 2, processor.SendEventToProcessCalls)
	assert.Equal(t, events[1], processor.SendEventToProcessArgs.Event)
}

func TestProcessUnpublishedEventsExecuteShouldReturnErrorWhenRepoFails(t *testing.T) {
	expectedErr := errors.New("repo failure")
	repo := &mockEventRepository{
		GetUnpublishedEventsFunc: func(limit int, offset int) ([]entities.Event, error) {
			return nil, expectedErr
		},
	}
	processor := &mockEventProcessor{}
	sut := NewProcessUnpublishedEventsUseCase(repo, processor)

	output, err := sut.Execute(&ProcessUnpublishedEventsInput{})

	assert.Nil(t, output)
	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, 1, repo.GetUnpublishedEventsCalls)
	assert.Equal(t, 0, processor.SendEventToProcessCalls)
}

func TestProcessUnpublishedEventsExecuteShouldReturnErrorWhenProcessorFails(t *testing.T) {
	events := makeFakeEvents(2)
	expectedErr := errors.New("processor failure")
	repo := &mockEventRepository{
		GetUnpublishedEventsFunc: func(limit int, offset int) ([]entities.Event, error) {
			return events, nil
		},
	}
	processor := &mockEventProcessor{
		SendEventToProcessFunc: func(event entities.Event) error {
			if event.ID == events[1].ID {
				return expectedErr
			}
			return nil
		},
	}
	sut := NewProcessUnpublishedEventsUseCase(repo, processor)

	output, err := sut.Execute(&ProcessUnpublishedEventsInput{})

	assert.Nil(t, output)
	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, 1, repo.GetUnpublishedEventsCalls)
	assert.Equal(t, 2, processor.SendEventToProcessCalls)
	assert.Equal(t, events[1], processor.SendEventToProcessArgs.Event)
}
