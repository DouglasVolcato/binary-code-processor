package usecases

import (
	"errors"
	"fmt"
	"strings"
)

type ProcessTaskUseCase struct {
	Repo      TaskRepositoryInterface
	Processor TaskProcessorInterface
}

func NewProcessTaskUseCase(repo TaskRepositoryInterface, processor TaskProcessorInterface) *ProcessTaskUseCase {
	return &ProcessTaskUseCase{
		Repo:      repo,
		Processor: processor,
	}
}

type ProcessTaskInput struct {
	ID string
}

type ProcessTaskOutput struct {
	ID         string
	BinaryCode string
}

func (u *ProcessTaskUseCase) Execute(input *ProcessTaskInput) (*ProcessTaskOutput, error) {
	task, err := u.Repo.GetTaskByID(input.ID)
	if err != nil {
		return nil, err
	}
	binaryCode, err := u.convertMessageToBinaryCode(task.Message)
	if err != nil {
		return nil, err
	}
	err = u.Processor.FinishProcessing(FinishProcessingDTO{
		ID:         task.ID,
		BinaryCode: binaryCode,
	})
	if err != nil {
		return nil, err
	}
	return &ProcessTaskOutput{
		ID:         task.ID,
		BinaryCode: binaryCode,
	}, nil
}

func (u *ProcessTaskUseCase) convertMessageToBinaryCode(message string) (string, error) {
	if strings.TrimSpace(message) == "" {
		return "", errors.New("message is empty")
	}
	result := make([]string, len(message))

	for i, char := range message {
		result[i] = fmt.Sprintf("%08b", char)
	}
	return strings.Join(result, ""), nil
}
