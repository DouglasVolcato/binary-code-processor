package usecases

import "github.com/douglasvolcato/binary-code-processor/task_service/internal/entities"

type TaskRepositoryInterface interface {
	GetTasks(limit int, offset int) ([]entities.Task, error)
}

type TaskProcessorInterface interface {
	ProcessTask(task entities.Task) (entities.Task, error)
}
