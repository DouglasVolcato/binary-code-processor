package usecases

type ReceiveProcessedTaskUseCase struct {
	Repo   TaskProcessorInterface
	Outbox TaskOutboxRepositoryInterface
}

func NewReceiveProcessedTaskUseCase(repo TaskProcessorInterface, outbox TaskOutboxRepositoryInterface) *ReceiveProcessedTaskUseCase {
	return &ReceiveProcessedTaskUseCase{
		Repo:   repo,
		Outbox: outbox,
	}
}

type ReceiveProcessedTaskInput struct {
	ID string
}

type ReceiveProcessedTaskOutput struct {
	Success bool
}

func (u *ReceiveProcessedTaskUseCase) Execute(input *ReceiveProcessedTaskInput) (*ReceiveProcessedTaskOutput, error) {
	task, err := u.Repo.SetTaskAsProcessed(input.ID)
	if err != nil {
		return nil, err
	}
	if err := u.Outbox.StoreProcessedEvent(task); err != nil {
		return nil, err
	}
	return &ReceiveProcessedTaskOutput{
		Success: true,
	}, nil
}
