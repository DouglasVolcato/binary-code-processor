package usecases

import (
	"errors"
	"testing"

	"github.com/douglasvolcato/binary-code-processor/event_publisher/internal/entities"
	"github.com/douglasvolcato/binary-code-processor/event_publisher/test"
	"github.com/stretchr/testify/assert"
)

type mockEventRepositoryForSendProcessed struct {
	GetUnprocessedEventsCalls int
	GetUnprocessedEventsArgs  struct {
		Limit  int
		Offset int
	}
	GetUnprocessedEventsFunc func(limit int, offset int) ([]entities.Event, error)

	GetProcessedEventsCalls int
	GetProcessedEventsArgs  struct {
		Limit  int
		Offset int
	}
	GetProcessedEventsFunc func(limit int, offset int) ([]entities.Event, error)

	DeleteEventByIDCalls int
	DeleteEventByIDArgs  struct {
		ID string
	}
	DeleteEventByIDFunc func(id string) error
}

func (m *mockEventRepositoryForSendProcessed) GetUnprocessedEvents(limit int, offset int) ([]entities.Event, error) {
	m.GetUnprocessedEventsCalls++
	m.GetUnprocessedEventsArgs.Limit = limit
	m.GetUnprocessedEventsArgs.Offset = offset
	if m.GetUnprocessedEventsFunc != nil {
		return m.GetUnprocessedEventsFunc(limit, offset)
	}
	return nil, nil
}

func (m *mockEventRepositoryForSendProcessed) GetProcessedEvents(limit int, offset int) ([]entities.Event, error) {
	m.GetProcessedEventsCalls++
	m.GetProcessedEventsArgs.Limit = limit
	m.GetProcessedEventsArgs.Offset = offset
	if m.GetProcessedEventsFunc != nil {
		return m.GetProcessedEventsFunc(limit, offset)
	}
	return nil, nil
}

func (m *mockEventRepositoryForSendProcessed) DeleteEventByID(id string) error {
	m.DeleteEventByIDCalls++
	m.DeleteEventByIDArgs.ID = id
	if m.DeleteEventByIDFunc != nil {
		return m.DeleteEventByIDFunc(id)
	}
	return nil
}

type mockEventProcessorForSendProcessed struct {
	SendFanoutEventCalls int
	SendFanoutEventArgs  struct {
		Event entities.Event
	}
	SendFanoutEventFunc func(event entities.Event) error

	SendEventToProcessCalls int
	SendEventToProcessFunc  func(event entities.Event) error
}

func (m *mockEventProcessorForSendProcessed) SendEventToProcess(event entities.Event) error {
	m.SendEventToProcessCalls++
	if m.SendEventToProcessFunc != nil {
		return m.SendEventToProcessFunc(event)
	}
	return nil
}

func (m *mockEventProcessorForSendProcessed) SendFanoutEvent(event entities.Event) error {
	m.SendFanoutEventCalls++
	m.SendFanoutEventArgs.Event = event
	if m.SendFanoutEventFunc != nil {
		return m.SendFanoutEventFunc(event)
	}
	return nil
}

type mockRemoteEventProcessorForSendProcessed struct {
	SendFanoutEventCalls int
	SendFanoutEventArgs  struct {
		Event entities.Event
	}
	SendFanoutEventFunc func(event entities.Event) error

	SendToQueueCalls int
	SendToQueueFunc  func(event entities.Event) error
}

func (m *mockRemoteEventProcessorForSendProcessed) SendToQueue(event entities.Event) error {
	m.SendToQueueCalls++
	if m.SendToQueueFunc != nil {
		return m.SendToQueueFunc(event)
	}
	return nil
}

func (m *mockRemoteEventProcessorForSendProcessed) SendFanoutEvent(event entities.Event) error {
	m.SendFanoutEventCalls++
	m.SendFanoutEventArgs.Event = event
	if m.SendFanoutEventFunc != nil {
		return m.SendFanoutEventFunc(event)
	}
	return nil
}

func makeFakeSendProcessedEvent() entities.Event {
	faker := test.FakeData{}
	return entities.Event{
		ID:     faker.ID(),
		Status: faker.Word(),
	}
}

func makeFakeSendProcessedEvents(count int) []entities.Event {
	events := make([]entities.Event, 0, count)
	for i := 0; i < count; i++ {
		events = append(events, makeFakeSendProcessedEvent())
	}
	return events
}

func TestNewSendProcessedEventsUseCaseShouldCreateSendProcessedEventsUseCase(t *testing.T) {
	repo := &mockEventRepositoryForSendProcessed{}
	processor := &mockEventProcessorForSendProcessed{}
	remoteProcessor := &mockRemoteEventProcessorForSendProcessed{}
	sut := NewSendProcessedEventsUseCase(repo, processor, remoteProcessor)

	assert.NotNil(t, sut)
	assert.Same(t, repo, sut.Repo)
	assert.Same(t, processor, sut.Processor)
	assert.Same(t, remoteProcessor, sut.RemoteProcessor)
}

func TestSendProcessedEventsExecuteShouldSendAllEvents(t *testing.T) {
	events := makeFakeSendProcessedEvents(2)
	repo := &mockEventRepositoryForSendProcessed{
		GetProcessedEventsFunc: func(limit int, offset int) ([]entities.Event, error) {
			return events, nil
		},
	}
	processor := &mockEventProcessorForSendProcessed{}
	remoteProcessor := &mockRemoteEventProcessorForSendProcessed{}
	sut := NewSendProcessedEventsUseCase(repo, processor, remoteProcessor)

	output, err := sut.Execute(&SendProcessedEventsInput{})

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, 1, repo.GetProcessedEventsCalls)
	assert.Equal(t, 100, repo.GetProcessedEventsArgs.Limit)
	assert.Equal(t, 0, repo.GetProcessedEventsArgs.Offset)
	assert.Equal(t, 0, repo.GetUnprocessedEventsCalls)
	assert.Equal(t, 2, remoteProcessor.SendFanoutEventCalls)
	assert.Equal(t, 2, processor.SendFanoutEventCalls)
	assert.Equal(t, events[1], remoteProcessor.SendFanoutEventArgs.Event)
	assert.Equal(t, events[1], processor.SendFanoutEventArgs.Event)
}

func TestSendProcessedEventsExecuteShouldReturnErrorWhenRepoFails(t *testing.T) {
	expectedErr := errors.New("repo failure")
	repo := &mockEventRepositoryForSendProcessed{
		GetProcessedEventsFunc: func(limit int, offset int) ([]entities.Event, error) {
			return nil, expectedErr
		},
	}
	processor := &mockEventProcessorForSendProcessed{}
	remoteProcessor := &mockRemoteEventProcessorForSendProcessed{}
	sut := NewSendProcessedEventsUseCase(repo, processor, remoteProcessor)

	output, err := sut.Execute(&SendProcessedEventsInput{})

	assert.Nil(t, output)
	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, 1, repo.GetProcessedEventsCalls)
	assert.Equal(t, 0, repo.GetUnprocessedEventsCalls)
	assert.Equal(t, 0, remoteProcessor.SendFanoutEventCalls)
	assert.Equal(t, 0, processor.SendFanoutEventCalls)
}

func TestSendProcessedEventsExecuteShouldReturnErrorWhenRemoteProcessorFails(t *testing.T) {
	events := makeFakeSendProcessedEvents(2)
	expectedErr := errors.New("remote failure")
	repo := &mockEventRepositoryForSendProcessed{
		GetProcessedEventsFunc: func(limit int, offset int) ([]entities.Event, error) {
			return events, nil
		},
	}
	processor := &mockEventProcessorForSendProcessed{}
	remoteProcessor := &mockRemoteEventProcessorForSendProcessed{
		SendFanoutEventFunc: func(event entities.Event) error {
			if event.ID == events[0].ID {
				return expectedErr
			}
			return nil
		},
	}
	sut := NewSendProcessedEventsUseCase(repo, processor, remoteProcessor)

	output, err := sut.Execute(&SendProcessedEventsInput{})

	assert.Nil(t, output)
	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, 1, repo.GetProcessedEventsCalls)
	assert.Equal(t, 0, repo.GetUnprocessedEventsCalls)
	assert.Equal(t, 1, remoteProcessor.SendFanoutEventCalls)
	assert.Equal(t, 0, processor.SendFanoutEventCalls)
	assert.Equal(t, events[0], remoteProcessor.SendFanoutEventArgs.Event)
}

func TestSendProcessedEventsExecuteShouldReturnErrorWhenProcessorFails(t *testing.T) {
	events := makeFakeSendProcessedEvents(2)
	expectedErr := errors.New("processor failure")
	repo := &mockEventRepositoryForSendProcessed{
		GetProcessedEventsFunc: func(limit int, offset int) ([]entities.Event, error) {
			return events, nil
		},
	}
	processor := &mockEventProcessorForSendProcessed{
		SendFanoutEventFunc: func(event entities.Event) error {
			if event.ID == events[1].ID {
				return expectedErr
			}
			return nil
		},
	}
	remoteProcessor := &mockRemoteEventProcessorForSendProcessed{}
	sut := NewSendProcessedEventsUseCase(repo, processor, remoteProcessor)

	output, err := sut.Execute(&SendProcessedEventsInput{})

	assert.Nil(t, output)
	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, 1, repo.GetProcessedEventsCalls)
	assert.Equal(t, 0, repo.GetUnprocessedEventsCalls)
	assert.Equal(t, 2, remoteProcessor.SendFanoutEventCalls)
	assert.Equal(t, 2, processor.SendFanoutEventCalls)
	assert.Equal(t, events[1], processor.SendFanoutEventArgs.Event)
}
