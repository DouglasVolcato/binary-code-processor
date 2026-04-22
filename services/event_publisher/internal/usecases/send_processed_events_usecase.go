package usecases

type SendProcessedEventsUseCase struct {
	Repo            EventRepositoryInterface
	Processor       EventProcessorInterface
	RemoteProcessor RemoteEventProcessorInterface
}

func NewSendProcessedEventsUseCase(repo EventRepositoryInterface, processor EventProcessorInterface, remoteProcessor RemoteEventProcessorInterface) *SendProcessedEventsUseCase {
	return &SendProcessedEventsUseCase{
		Repo:            repo,
		Processor:       processor,
		RemoteProcessor: remoteProcessor,
	}
}

type SendProcessedEventsInput struct {
}

type SendProcessedEventsOutput struct {
}

func (u *SendProcessedEventsUseCase) Execute(input *SendProcessedEventsInput) (*SendProcessedEventsOutput, error) {
	events, err := u.Repo.GetProcessedEvents(100, 0)
	if err != nil {
		return nil, err
	}
	for _, event := range events {
		if err := u.RemoteProcessor.SendFanoutEvent(event); err != nil {
			return nil, err
		}
		if err := u.Processor.SendFanoutEvent(event); err != nil {
			return nil, err
		}
	}
	return &SendProcessedEventsOutput{}, nil
}
