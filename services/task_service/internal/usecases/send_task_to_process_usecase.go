package usecases

import "github.com/douglasvolcato/binary-code-processor/task_service/internal/entities"

type ReceiveTaskToProcessUseCase struct {
	Repo TaskProcessorInterface
}

func NewReceiveTaskToProcessUseCase(repo TaskProcessorInterface) *ReceiveTaskToProcessUseCase {
	return &ReceiveTaskToProcessUseCase{
		Repo: repo,
	}
}

type ReceiveTaskToProcessInput struct {
	Task entities.Task
}

type ReceiveTaskToProcessOutput struct {
	Success bool
	Task    entities.Task
}

func (u *ReceiveTaskToProcessUseCase) Execute(input *ReceiveTaskToProcessInput) (*ReceiveTaskToProcessOutput, error) {
	task, err := u.Repo.ProcessTask(input.Task)
	if err != nil {
		return nil, err
	}
	return &ReceiveTaskToProcessOutput{
		Success: true,
		Task:    task,
	}, nil
}
