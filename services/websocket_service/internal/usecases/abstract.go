package usecases

import "github.com/douglasvolcato/binary-code-processor/websocket_service/internal/entities"

type WebSocketClient interface {
	SendProcessedTasksToClient(task entities.Task) error
}
