package grpc

import (
	"context"
	"encoding/json"

	"github.com/douglasvolcato/binary-code-processor/api_gateway/internal/entities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
)

type jsonCodec struct{}

func (jsonCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (jsonCodec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (jsonCodec) Name() string {
	return "json"
}

func init() {
	encoding.RegisterCodec(jsonCodec{})
}

type TaskRecord struct {
	ID         string `json:"id"`
	Message    string `json:"message"`
	BinaryCode string `json:"binaryCode"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

type ReceiveTaskRequest struct {
	Message string `json:"message"`
}

type ReceiveTaskResponse struct {
	Tasks []TaskRecord `json:"tasks"`
}

type TaskListRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type TaskListResponse struct {
	Tasks []TaskRecord `json:"tasks"`
}

type Client struct {
	conn *grpc.ClientConn
}

func NewClient(conn *grpc.ClientConn) *Client {
	return &Client{conn: conn}
}

func (c *Client) GetTasks(ctx context.Context, limit int, offset int) ([]entities.Task, error) {
	out := new(TaskListResponse)
	if err := c.conn.Invoke(
		ctx,
		"/task.TaskService/GetTasks",
		&TaskListRequest{Limit: limit, Offset: offset},
		out,
		grpc.ForceCodec(jsonCodec{}),
	); err != nil {
		return nil, err
	}
	return mapTasks(out.Tasks), nil
}

func (c *Client) SendTaskToProcess(ctx context.Context, messages []string) ([]entities.Task, error) {
	stream, err := c.conn.NewStream(
		ctx,
		&grpc.StreamDesc{StreamName: "ReceiveTaskToProcess", ClientStreams: true},
		"/task.TaskService/ReceiveTaskToProcess",
		grpc.ForceCodec(jsonCodec{}),
	)
	if err != nil {
		return nil, err
	}

	for _, message := range messages {
		if err := stream.SendMsg(&ReceiveTaskRequest{Message: message}); err != nil {
			return nil, err
		}
	}

	if err := stream.CloseSend(); err != nil {
		return nil, err
	}

	out := new(ReceiveTaskResponse)
	if err := stream.RecvMsg(out); err != nil {
		return nil, err
	}
	return mapTasks(out.Tasks), nil
}

func mapTasks(tasks []TaskRecord) []entities.Task {
	mapped := make([]entities.Task, 0, len(tasks))
	for _, task := range tasks {
		mapped = append(mapped, entities.Task{
			ID:         task.ID,
			Message:    task.Message,
			BinaryCode: task.BinaryCode,
			CreatedAt:  task.CreatedAt,
			UpdatedAt:  task.UpdatedAt,
		})
	}
	return mapped
}
