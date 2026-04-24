package graphql

import "github.com/douglasvolcato/binary-code-processor/api_gateway/internal/usecases"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	GetTasksUseCase          *usecases.GetTasksUseCase
	SendTaskToProcessUseCase *usecases.SendTaskToProcessUseCase
}
