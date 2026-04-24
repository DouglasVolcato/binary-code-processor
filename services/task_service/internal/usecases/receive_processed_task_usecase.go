package usecases

import (
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

	_, err := u.Repo.FinishProcessing(FinishProcessingDTO{
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
