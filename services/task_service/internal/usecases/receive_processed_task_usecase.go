package usecases

type ReceiveProcessedTaskUseCase struct {
	Repo TaskProcessorInterface
}

func NewReceiveProcessedTaskUseCase(repo TaskProcessorInterface) *ReceiveProcessedTaskUseCase {
	return &ReceiveProcessedTaskUseCase{
		Repo: repo,
	}
}

type ReceiveProcessedTaskInput struct {
	ID string
}

type ReceiveProcessedTaskOutput struct {
	Success bool
}

func (u *ReceiveProcessedTaskUseCase) Execute(input *ReceiveProcessedTaskInput) (*ReceiveProcessedTaskOutput, error) {
	err := u.Repo.SetTaskAsProcessed(input.ID)
	if err != nil {
		return nil, err
	}
	return &ReceiveProcessedTaskOutput{
		Success: true,
	}, nil
}
