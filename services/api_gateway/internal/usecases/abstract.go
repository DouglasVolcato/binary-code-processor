package usecases

import (
	"context"

	"github.com/douglasvolcato/binary-code-processor/api_gateway/internal/entities"
)

type TaskRepositoryInterface interface {
	GetTasks(ctx context.Context, limit int, offset int) ([]entities.Task, error)
}

type TaskProcessorInterface interface {
	SendTaskToProcess(ctx context.Context, messages []string) ([]entities.Task, error)
}
