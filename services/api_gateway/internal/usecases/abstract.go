package usecases

import "github.com/douglasvolcato/binary-code-processer/api_gateway/internal/entities"

type TaskRepositoryInterface interface {
	GetTasks(limit int, offset int) ([]entities.Task, error)
}

type TaskProcessorInterface interface {
	SendTaskToProcess(messages []string) error
}
