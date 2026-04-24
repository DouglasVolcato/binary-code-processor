package usecases

import (
	"context"

	"github.com/douglasvolcato/binary-code-processor/processing_service/internal/entities"
)

type TaskRepositoryInterface interface {
	GetTaskByID(ctx context.Context, taskID string) (entities.Task, error)
}

type FinishProcessingDTO struct {
	ID         string
	BinaryCode string
}

type TaskProcessorInterface interface {
	FinishProcessing(ctx context.Context, dto FinishProcessingDTO) error
}
