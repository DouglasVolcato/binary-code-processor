package usecases

import "github.com/douglasvolcato/binary-code-processor/task_service/internal/entities"

type IDGeneratorInterface interface {
	GenerateID() string
}

type TaskRepositoryInterface interface {
	GetTasks(limit int, offset int) ([]entities.Task, error)
	GetTaskByID(taskID string) (entities.Task, error)
}

type CreateTaskDTO struct {
	ID      string
	Message string
}

type TaskProcessorInterface interface {
	MoveTaskToProcessing(createTaskDto CreateTaskDTO) (entities.Task, error)
	FinishProcessing(dto FinishProcessingDTO) (entities.Task, error)
}

type FinishProcessingDTO struct {
	ID         string
	BinaryCode string
}
