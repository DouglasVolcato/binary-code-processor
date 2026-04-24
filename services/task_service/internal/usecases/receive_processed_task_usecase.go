package usecases

import (
	"context"
	"errors"
	"strings"
)

type ReceiveProcessedTaskUseCase struct {
	Repo TaskProcessorInterface
}

func NewReceiveProcessedTaskUseCase(repo TaskProcessorInterface) *ReceiveProcessedTaskUseCase {
	return &ReceiveProcessedTaskUseCase{
		Repo: repo,
	}
}

type ReceiveProcessedTaskInput struct {
	Ctx        context.Context
	ID         string
	BinaryCode string
}

type ReceiveProcessedTaskOutput struct {
	Success bool
}

func (u *ReceiveProcessedTaskUseCase) Execute(input *ReceiveProcessedTaskInput) (*ReceiveProcessedTaskOutput, error) {
	taskID := strings.TrimSpace(input.ID)
	binaryCode := strings.TrimSpace(input.BinaryCode)
	if taskID == "" {
		return nil, errors.New("task id is empty")
	}
	if binaryCode == "" {
		return nil, errors.New("binary code is empty")
	}

	ctx := input.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	_, err := u.Repo.FinishProcessing(ctx, FinishProcessingDTO{
		ID:         taskID,
		BinaryCode: binaryCode,
	})
	if err != nil {
		return nil, err
	}
	return &ReceiveProcessedTaskOutput{
		Success: true,
	}, nil
}
