package usecases

import "github.com/douglasvolcato/binary-code-processor/websocket_service/internal/entities"

type SendProcessedTasksUseCase struct {
	WebSocket WebSocketClient
}

func NewSendProcessedTasksUseCase(webSocket WebSocketClient) *SendProcessedTasksUseCase {
	return &SendProcessedTasksUseCase{
		WebSocket: webSocket,
	}
}

type SendProcessedTasksInput struct {
	Task entities.Task
}

type SendProcessedTasksOutput struct {
}

func (u *SendProcessedTasksUseCase) Execute(input *SendProcessedTasksInput) (*SendProcessedTasksOutput, error) {
	if err := u.WebSocket.SendProcessedTasksToClient(input.Task); err != nil {
		return nil, err
	}
	return &SendProcessedTasksOutput{}, nil
}
