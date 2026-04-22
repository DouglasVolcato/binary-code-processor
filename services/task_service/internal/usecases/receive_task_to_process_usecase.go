package usecases

import "github.com/douglasvolcato/binary-code-processor/task_service/internal/entities"

type ReceiveTaskToProcessUseCase struct {
	Repo     TaskProcessorInterface
	Outbox   TaskOutboxRepositoryInterface
	IDGen    IDGeneratorInterface
}

func NewReceiveTaskToProcessUseCase(repo TaskProcessorInterface, outbox TaskOutboxRepositoryInterface, idGen IDGeneratorInterface) *ReceiveTaskToProcessUseCase {
	return &ReceiveTaskToProcessUseCase{
		Repo:   repo,
		Outbox: outbox,
		IDGen:  idGen,
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
	createTaskDto := CreateTaskDTO{
		ID:      u.IDGen.GenerateID(),
		Message: input.Message,
	}
	task, err := u.Repo.MoveTaskToProcessing(createTaskDto)
	if err != nil {
		return nil, err
	}
	if err := u.Outbox.StoreUnprocessedEvent(task); err != nil {
		return nil, err
	}
	return &ReceiveTaskToProcessOutput{
		Success: true,
		Task:    task,
	}, nil
}
