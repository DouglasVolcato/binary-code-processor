package usecases

import "github.com/douglasvolcato/binary-code-processor/api_gateway/internal/entities"

type SendTaskToProcessUseCase struct {
	Repo TaskProcessorInterface
}

func NewSendTaskToProcessUseCase(repo TaskProcessorInterface) *SendTaskToProcessUseCase {
	return &SendTaskToProcessUseCase{
		Repo: repo,
	}
}

type SendTaskToProcessInput struct {
	Messages []string
}

type SendTaskToProcessOutput struct {
	Success bool
	Tasks   []entities.Task
}

func (u *SendTaskToProcessUseCase) Execute(input *SendTaskToProcessInput) (*SendTaskToProcessOutput, error) {
	tasks, err := u.Repo.SendTaskToProcess(input.Messages)
	if err != nil {
		return nil, err
	}
	return &SendTaskToProcessOutput{
		Success: true,
		Tasks:   tasks,
	}, nil
}
