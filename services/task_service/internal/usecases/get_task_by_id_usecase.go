package usecases

import "github.com/douglasvolcato/binary-code-processor/task_service/internal/entities"

type GetTaskByIDUseCase struct {
	Repo TaskRepositoryInterface
}

func NewGetTaskByIDUseCase(repo TaskRepositoryInterface) *GetTaskByIDUseCase {
	return &GetTaskByIDUseCase{
		Repo: repo,
	}
}

type GetTaskByIDInput struct {
	ID string
}

type GetTaskByIDOutput struct {
	Task entities.Task
}

func (u *GetTaskByIDUseCase) Execute(input *GetTaskByIDInput) (*GetTaskByIDOutput, error) {
	task, err := u.Repo.GetTaskByID(input.ID)
	if err != nil {
		return nil, err
	}
	return &GetTaskByIDOutput{
		Task: task,
	}, nil
}
