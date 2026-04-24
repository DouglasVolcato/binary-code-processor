package grpc

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/douglasvolcato/binary-code-processor/processing_service/internal/entities"
	"github.com/douglasvolcato/binary-code-processor/processing_service/internal/usecases"
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

type TaskByIDRequest struct {
	ID string `json:"id"`
}

type TaskByIDResponse struct {
	Task TaskRecord `json:"task"`
}

type ProcessedTaskRequest struct {
	ID         string `json:"id"`
	BinaryCode string `json:"binaryCode"`
}

type ProcessedTaskResponse struct {
	Success bool `json:"success"`
}

type Client struct {
	conn *grpc.ClientConn
}

func NewClient(conn *grpc.ClientConn) *Client {
	return &Client{conn: conn}
}

func (c *Client) GetTaskByID(taskID string) (entities.Task, error) {
	out := new(TaskByIDResponse)
	if err := c.conn.Invoke(
		context.Background(),
		"/task.TaskService/GetTaskByID",
		&TaskByIDRequest{ID: taskID},
		out,
		grpc.ForceCodec(jsonCodec{}),
	); err != nil {
		return entities.Task{}, err
	}

	return entities.Task{
		ID:         out.Task.ID,
		Message:    out.Task.Message,
		BinaryCode: out.Task.BinaryCode,
		CreatedAt:  out.Task.CreatedAt,
		UpdatedAt:  out.Task.UpdatedAt,
	}, nil
}

func (c *Client) FinishProcessing(dto usecases.FinishProcessingDTO) error {
	stream, err := c.conn.NewStream(
		context.Background(),
		&grpc.StreamDesc{StreamName: "SendProcessedTask", ClientStreams: true},
		"/task.TaskService/SendProcessedTask",
		grpc.ForceCodec(jsonCodec{}),
	)
	if err != nil {
		return err
	}

	if err := stream.SendMsg(&ProcessedTaskRequest{
		ID:         dto.ID,
		BinaryCode: dto.BinaryCode,
	}); err != nil {
		return err
	}

	if err := stream.CloseSend(); err != nil {
		return err
	}

	out := new(ProcessedTaskResponse)
	if err := stream.RecvMsg(out); err != nil {
		return err
	}
	if !out.Success {
		return errors.New("task service rejected processed task")
	}
	return nil
}
