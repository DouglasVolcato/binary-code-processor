package usecases

type ProcessUnprocessedEventsUseCase struct {
	Repo            EventRepositoryInterface
	Processor       EventProcessorInterface
	RemoteProcessor RemoteEventProcessorInterface
}

func NewProcessUnprocessedEventsUseCase(repo EventRepositoryInterface, processor EventProcessorInterface, remoteProcessor RemoteEventProcessorInterface) *ProcessUnprocessedEventsUseCase {
	return &ProcessUnprocessedEventsUseCase{
		Repo:            repo,
		Processor:       processor,
		RemoteProcessor: remoteProcessor,
	}
}

type ProcessUnprocessedEventsInput struct {
}

type ProcessUnprocessedEventsOutput struct {
}

func (u *ProcessUnprocessedEventsUseCase) Execute(input *ProcessUnprocessedEventsInput) (*ProcessUnprocessedEventsOutput, error) {
	events, err := u.Repo.GetUnprocessedEvents(100, 0)
	if err != nil {
		return nil, err
	}
	for _, event := range events {
		if err := u.RemoteProcessor.SendToQueue(event); err != nil {
			return nil, err
		}
		if err := u.Processor.SendEventToProcess(event); err != nil {
			return nil, err
		}
		if err := u.Repo.DeleteEventByID(event.ID); err != nil {
			return nil, err
		}
	}
	return &ProcessUnprocessedEventsOutput{}, nil
}
