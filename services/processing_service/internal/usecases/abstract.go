package usecases

import "github.com/douglasvolcato/binary-code-processor/processing_service/internal/entities"

type TaskRepositoryInterface interface {
	GetTaskByID(taskID string) (entities.Task, error)
}

type FinishProcessingDTO struct {
	ID         string
	BinaryCode string
}

type TaskProcessorInterface interface {
	FinishProcessing(dto FinishProcessingDTO) error
}
