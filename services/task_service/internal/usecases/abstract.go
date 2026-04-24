package usecases

import (
	"context"

	"github.com/douglasvolcato/binary-code-processor/task_service/internal/entities"
)

type IDGeneratorInterface interface {
	GenerateID() string
}

type TaskRepositoryInterface interface {
	GetTasks(ctx context.Context, limit int, offset int) ([]entities.Task, error)
	GetTaskByID(ctx context.Context, taskID string) (entities.Task, error)
}

type CreateTaskDTO struct {
	ID      string
	Message string
}

type TaskProcessorInterface interface {
	MoveTaskToProcessing(ctx context.Context, createTaskDto CreateTaskDTO) (entities.Task, error)
	FinishProcessing(ctx context.Context, dto FinishProcessingDTO) (entities.Task, error)
}

type FinishProcessingDTO struct {
	ID         string
	BinaryCode string
}
