package usecases

import "github.com/douglasvolcato/binary-code-processor/task_service/internal/entities"

type GetTasksUseCase struct {
	Repo TaskRepositoryInterface
}

func NewGetTasksUseCase(repo TaskRepositoryInterface) *GetTasksUseCase {
	return &GetTasksUseCase{
		Repo: repo,
	}
}

type GetTasksInput struct {
	Limit  int
	Offset int
}

type GetTasksOutput struct {
	Tasks []entities.Task
}

func (u *GetTasksUseCase) Execute(input *GetTasksInput) (*GetTasksOutput, error) {
	tasks, err := u.Repo.GetTasks(input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}
	return &GetTasksOutput{
		Tasks: tasks,
	}, nil
}
