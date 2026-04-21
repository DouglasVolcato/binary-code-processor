package usecases

import "github.com/douglasvolcato/binary-code-processer/api_gateway/internal/entities"

type TaskRepository interface {
	GetTasks(limit int, offset int) ([]entities.Task, error)
}

type GetTasksUseCase struct {
	Repo TaskRepository
}

func NewGetTasksUseCase(repo TaskRepository) *GetTasksUseCase {
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
	if input.Limit <= 0 || input.Limit > 20 {
		input.Limit = 20
	}
	tasks, err := u.Repo.GetTasks(input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}
	return &GetTasksOutput{
		Tasks: tasks,
	}, nil
}
