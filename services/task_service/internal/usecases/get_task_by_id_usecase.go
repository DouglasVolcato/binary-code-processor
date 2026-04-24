package usecases

import (
	"context"

	"github.com/douglasvolcato/binary-code-processor/task_service/internal/entities"
)

type GetTaskByIDUseCase struct {
	Repo TaskRepositoryInterface
}

func NewGetTaskByIDUseCase(repo TaskRepositoryInterface) *GetTaskByIDUseCase {
	return &GetTaskByIDUseCase{
		Repo: repo,
	}
}

type GetTaskByIDInput struct {
	Ctx context.Context
	ID  string
}

type GetTaskByIDOutput struct {
	Task entities.Task
}

func (u *GetTaskByIDUseCase) Execute(input *GetTaskByIDInput) (*GetTaskByIDOutput, error) {
	ctx := input.Ctx
	if ctx == nil {
		ctx = context.Background()
	}
	task, err := u.Repo.GetTaskByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	return &GetTaskByIDOutput{
		Task: task,
	}, nil
}
