package usecases

type ProcessUnpublishedEventsUseCase struct {
	Repo      EventRepositoryInterface
	Processor EventProcessorInterface
}

func NewProcessUnpublishedEventsUseCase(repo EventRepositoryInterface, processor EventProcessorInterface) *ProcessUnpublishedEventsUseCase {
	return &ProcessUnpublishedEventsUseCase{
		Repo:      repo,
		Processor: processor,
	}
}

type ProcessUnpublishedEventsInput struct {
}

type ProcessUnpublishedEventsOutput struct {
}

func (u *ProcessUnpublishedEventsUseCase) Execute(input *ProcessUnpublishedEventsInput) (*ProcessUnpublishedEventsOutput, error) {
	events, err := u.Repo.GetUnpublishedEvents(100, 0)
	if err != nil {
		return nil, err
	}
	for _, event := range events {
		if err := u.Processor.SendEventToProcess(event); err != nil {
			return nil, err
		}
	}
	return &ProcessUnpublishedEventsOutput{}, nil
}
