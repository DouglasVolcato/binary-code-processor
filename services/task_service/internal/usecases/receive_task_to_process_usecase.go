package usecases

import (
	"errors"
	"strings"

	"github.com/douglasvolcato/binary-code-processor/task_service/internal/entities"
)

type ReceiveTaskToProcessUseCase struct {
	Repo  TaskProcessorInterface
	IDGen IDGeneratorInterface
}

func NewReceiveTaskToProcessUseCase(repo TaskProcessorInterface, idGen IDGeneratorInterface) *ReceiveTaskToProcessUseCase {
	return &ReceiveTaskToProcessUseCase{
		Repo:  repo,
		IDGen: idGen,
	}
}

type ReceiveTaskToProcessInput struct {
	Message string
}

type ReceiveTaskToProcessOutput struct {
	Success bool
	Task    entities.Task
}

func (u *ReceiveTaskToProcessUseCase) Execute(input *ReceiveTaskToProcessInput) (*ReceiveTaskToProcessOutput, error) {
	message := strings.TrimSpace(input.Message)
	if message == "" {
		return nil, errors.New("message is empty")
	}

	createTaskDto := CreateTaskDTO{
		ID:      u.IDGen.GenerateID(),
		Message: message,
	}
	task, err := u.Repo.MoveTaskToProcessing(createTaskDto)
	if err != nil {
		return nil, err
	}
	return &ReceiveTaskToProcessOutput{
		Success: true,
		Task:    task,
	}, nil
}
