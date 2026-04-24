package usecases

import (
	"context"

	"github.com/douglasvolcato/binary-code-processor/api_gateway/internal/entities"
)

type GetTasksUseCase struct {
	Repo TaskRepositoryInterface
}

func NewGetTasksUseCase(repo TaskRepositoryInterface) *GetTasksUseCase {
	return &GetTasksUseCase{
		Repo: repo,
	}
}

type GetTasksInput struct {
	Ctx    context.Context
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
	ctx := input.Ctx
	if ctx == nil {
		ctx = context.Background()
	}
	tasks, err := u.Repo.GetTasks(ctx, input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}
	return &GetTasksOutput{
		Tasks: tasks,
	}, nil
}
