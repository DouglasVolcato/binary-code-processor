package usecases

import (
	"errors"
	"testing"

	"github.com/douglasvolcato/binary-code-processor/event_publisher/internal/entities"
	"github.com/douglasvolcato/binary-code-processor/event_publisher/test"
	"github.com/stretchr/testify/assert"
)

type mockEventRepositoryForProcessUnprocessed struct {
	GetUnprocessedEventsCalls int
	GetUnprocessedEventsArgs  struct {
		Limit  int
		Offset int
	}
	GetUnprocessedEventsFunc func(limit int, offset int) ([]entities.Event, error)
}

func (m *mockEventRepositoryForProcessUnprocessed) GetUnprocessedEvents(limit int, offset int) ([]entities.Event, error) {
	m.GetUnprocessedEventsCalls++
	m.GetUnprocessedEventsArgs.Limit = limit
	m.GetUnprocessedEventsArgs.Offset = offset
	if m.GetUnprocessedEventsFunc != nil {
		return m.GetUnprocessedEventsFunc(limit, offset)
	}
	return nil, nil
}

type mockEventProcessorForProcessUnprocessed struct {
	SendEventToProcessCalls int
	SendEventToProcessArgs  struct {
		Event entities.Event
	}
	SendEventToProcessFunc func(event entities.Event) error

	SendFanoutEventCalls int
	SendFanoutEventFunc  func(event entities.Event) error
}

func (m *mockEventProcessorForProcessUnprocessed) SendEventToProcess(event entities.Event) error {
	m.SendEventToProcessCalls++
	m.SendEventToProcessArgs.Event = event
	if m.SendEventToProcessFunc != nil {
		return m.SendEventToProcessFunc(event)
	}
	return nil
}

func (m *mockEventProcessorForProcessUnprocessed) SendFanoutEvent(event entities.Event) error {
	m.SendFanoutEventCalls++
	if m.SendFanoutEventFunc != nil {
		return m.SendFanoutEventFunc(event)
	}
	return nil
}

type mockRemoteEventProcessorForProcessUnprocessed struct {
	SendToQueueCalls int
	SendToQueueArgs  struct {
		Event entities.Event
	}
	SendToQueueFunc func(event entities.Event) error

	SendFanoutEventCalls int
	SendFanoutEventFunc  func(event entities.Event) error
}

func (m *mockRemoteEventProcessorForProcessUnprocessed) SendToQueue(event entities.Event) error {
	m.SendToQueueCalls++
	m.SendToQueueArgs.Event = event
	if m.SendToQueueFunc != nil {
		return m.SendToQueueFunc(event)
	}
	return nil
}

func (m *mockRemoteEventProcessorForProcessUnprocessed) SendFanoutEvent(event entities.Event) error {
	m.SendFanoutEventCalls++
	if m.SendFanoutEventFunc != nil {
		return m.SendFanoutEventFunc(event)
	}
	return nil
}

func makeFakeProcessUnprocessedEvent() entities.Event {
	faker := test.FakeData{}
	return entities.Event{
		ID:     faker.ID(),
		Status: faker.Word(),
	}
}

func makeFakeProcessUnprocessedEvents(count int) []entities.Event {
	events := make([]entities.Event, 0, count)
	for i := 0; i < count; i++ {
		events = append(events, makeFakeProcessUnprocessedEvent())
	}
	return events
}

func TestNewProcessUnprocessedEventsUseCaseShouldCreateProcessUnprocessedEventsUseCase(t *testing.T) {
	repo := &mockEventRepositoryForProcessUnprocessed{}
	processor := &mockEventProcessorForProcessUnprocessed{}
	remoteProcessor := &mockRemoteEventProcessorForProcessUnprocessed{}
	sut := NewProcessUnprocessedEventsUseCase(repo, processor, remoteProcessor)

	assert.NotNil(t, sut)
	assert.Same(t, repo, sut.Repo)
	assert.Same(t, processor, sut.Processor)
	assert.Same(t, remoteProcessor, sut.RemoteProcessor)
}

func TestProcessUnprocessedEventsExecuteShouldSendAllEventsToQueueAndProcess(t *testing.T) {
	events := makeFakeProcessUnprocessedEvents(2)
	repo := &mockEventRepositoryForProcessUnprocessed{
		GetUnprocessedEventsFunc: func(limit int, offset int) ([]entities.Event, error) {
			return events, nil
		},
	}
	processor := &mockEventProcessorForProcessUnprocessed{}
	remoteProcessor := &mockRemoteEventProcessorForProcessUnprocessed{}
	sut := NewProcessUnprocessedEventsUseCase(repo, processor, remoteProcessor)

	output, err := sut.Execute(&ProcessUnprocessedEventsInput{})

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, 1, repo.GetUnprocessedEventsCalls)
	assert.Equal(t, 100, repo.GetUnprocessedEventsArgs.Limit)
	assert.Equal(t, 0, repo.GetUnprocessedEventsArgs.Offset)
	assert.Equal(t, 2, remoteProcessor.SendToQueueCalls)
	assert.Equal(t, 2, processor.SendEventToProcessCalls)
	assert.Equal(t, events[1], remoteProcessor.SendToQueueArgs.Event)
	assert.Equal(t, events[1], processor.SendEventToProcessArgs.Event)
}

func TestProcessUnprocessedEventsExecuteShouldReturnErrorWhenRepoFails(t *testing.T) {
	expectedErr := errors.New("repo failure")
	repo := &mockEventRepositoryForProcessUnprocessed{
		GetUnprocessedEventsFunc: func(limit int, offset int) ([]entities.Event, error) {
			return nil, expectedErr
		},
	}
	processor := &mockEventProcessorForProcessUnprocessed{}
	remoteProcessor := &mockRemoteEventProcessorForProcessUnprocessed{}
	sut := NewProcessUnprocessedEventsUseCase(repo, processor, remoteProcessor)

	output, err := sut.Execute(&ProcessUnprocessedEventsInput{})

	assert.Nil(t, output)
	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, 1, repo.GetUnprocessedEventsCalls)
	assert.Equal(t, 0, remoteProcessor.SendToQueueCalls)
	assert.Equal(t, 0, processor.SendEventToProcessCalls)
}

func TestProcessUnprocessedEventsExecuteShouldReturnErrorWhenRemoteProcessorFails(t *testing.T) {
	events := makeFakeProcessUnprocessedEvents(2)
	expectedErr := errors.New("remote failure")
	repo := &mockEventRepositoryForProcessUnprocessed{
		GetUnprocessedEventsFunc: func(limit int, offset int) ([]entities.Event, error) {
			return events, nil
		},
	}
	processor := &mockEventProcessorForProcessUnprocessed{}
	remoteProcessor := &mockRemoteEventProcessorForProcessUnprocessed{
		SendToQueueFunc: func(event entities.Event) error {
			if event.ID == events[0].ID {
				return expectedErr
			}
			return nil
		},
	}
	sut := NewProcessUnprocessedEventsUseCase(repo, processor, remoteProcessor)

	output, err := sut.Execute(&ProcessUnprocessedEventsInput{})

	assert.Nil(t, output)
	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, 1, repo.GetUnprocessedEventsCalls)
	assert.Equal(t, 1, remoteProcessor.SendToQueueCalls)
	assert.Equal(t, 0, processor.SendEventToProcessCalls)
	assert.Equal(t, events[0], remoteProcessor.SendToQueueArgs.Event)
}

func TestProcessUnprocessedEventsExecuteShouldReturnErrorWhenProcessorFails(t *testing.T) {
	events := makeFakeProcessUnprocessedEvents(2)
	expectedErr := errors.New("processor failure")
	repo := &mockEventRepositoryForProcessUnprocessed{
		GetUnprocessedEventsFunc: func(limit int, offset int) ([]entities.Event, error) {
			return events, nil
		},
	}
	processor := &mockEventProcessorForProcessUnprocessed{
		SendEventToProcessFunc: func(event entities.Event) error {
			if event.ID == events[1].ID {
				return expectedErr
			}
			return nil
		},
	}
	remoteProcessor := &mockRemoteEventProcessorForProcessUnprocessed{}
	sut := NewProcessUnprocessedEventsUseCase(repo, processor, remoteProcessor)

	output, err := sut.Execute(&ProcessUnprocessedEventsInput{})

	assert.Nil(t, output)
	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, 1, repo.GetUnprocessedEventsCalls)
	assert.Equal(t, 2, remoteProcessor.SendToQueueCalls)
	assert.Equal(t, 2, processor.SendEventToProcessCalls)
	assert.Equal(t, events[1], processor.SendEventToProcessArgs.Event)
}
