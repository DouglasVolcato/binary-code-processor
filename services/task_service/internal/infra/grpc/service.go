package grpc

import (
	"context"
	"encoding/json"
	"io"

	"github.com/douglasvolcato/binary-code-processor/task_service/internal/entities"
	"github.com/douglasvolcato/binary-code-processor/task_service/internal/usecases"
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

func NewJSONCodec() encoding.Codec {
	return jsonCodec{}
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

type ProcessedTaskRequest struct {
	ID         string `json:"id"`
	BinaryCode string `json:"binaryCode"`
}

type ProcessedTaskResponse struct {
	Success bool `json:"success"`
}

type TaskListRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type TaskListResponse struct {
	Tasks []TaskRecord `json:"tasks"`
}

type TaskByIDRequest struct {
	ID string `json:"id"`
}

type TaskByIDResponse struct {
	Task TaskRecord `json:"task"`
}

type Server struct {
	ReceiveTaskToProcessUseCase *usecases.ReceiveTaskToProcessUseCase
	ReceiveProcessedTaskUseCase *usecases.ReceiveProcessedTaskUseCase
	GetTasksUseCase             *usecases.GetTasksUseCase
	GetTaskByIDUseCase          *usecases.GetTaskByIDUseCase
}

type taskServiceHandler interface{}

func RegisterTaskAPIServer(s grpc.ServiceRegistrar, srv *Server) {
	s.RegisterService(&grpc.ServiceDesc{
		ServiceName: "task.TaskService",
		HandlerType: (*taskServiceHandler)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "GetTasks",
				Handler:    taskServiceGetTasksHandler,
			},
			{
				MethodName: "GetTaskByID",
				Handler:    taskServiceGetTaskByIDHandler,
			},
		},
		Streams: []grpc.StreamDesc{
			{
				StreamName:    "ReceiveTaskToProcess",
				Handler:       taskServiceReceiveTaskHandler,
				ClientStreams: true,
			},
			{
				StreamName:    "SendProcessedTask",
				Handler:       taskServiceSendProcessedTaskHandler,
				ClientStreams: true,
			},
		},
		Metadata: "task.proto",
	}, srv)
}

func taskServiceReceiveTaskHandler(srv interface{}, stream grpc.ServerStream) error {
	server := srv.(*Server)
	var tasks []TaskRecord

	for {
		req := new(ReceiveTaskRequest)
		if err := stream.RecvMsg(req); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		output, err := server.ReceiveTaskToProcessUseCase.Execute(&usecases.ReceiveTaskToProcessInput{
			Message: req.Message,
		})
		if err != nil {
			return err
		}

		tasks = append(tasks, mapTask(output.Task))
	}

	return stream.SendMsg(&ReceiveTaskResponse{Tasks: tasks})
}

func taskServiceSendProcessedTaskHandler(srv interface{}, stream grpc.ServerStream) error {
	server := srv.(*Server)
	success := true

	for {
		req := new(ProcessedTaskRequest)
		if err := stream.RecvMsg(req); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		_, err := server.ReceiveProcessedTaskUseCase.Execute(&usecases.ReceiveProcessedTaskInput{
			ID:         req.ID,
			BinaryCode: req.BinaryCode,
		})
		if err != nil {
			return err
		}
	}

	return stream.SendMsg(&ProcessedTaskResponse{Success: success})
}

func taskServiceGetTasksHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}

	if interceptor == nil {
		output, err := srv.(*Server).GetTasksUseCase.Execute(&usecases.GetTasksInput{
			Limit:  in.Limit,
			Offset: in.Offset,
		})
		if err != nil {
			return nil, err
		}
		return &TaskListResponse{Tasks: mapTasks(output.Tasks)}, nil
	}

	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/task.TaskService/GetTasks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		output, err := srv.(*Server).GetTasksUseCase.Execute(&usecases.GetTasksInput{
			Limit:  req.(*TaskListRequest).Limit,
			Offset: req.(*TaskListRequest).Offset,
		})
		if err != nil {
			return nil, err
		}
		return &TaskListResponse{Tasks: mapTasks(output.Tasks)}, nil
	}
	return interceptor(ctx, in, info, handler)
}

func taskServiceGetTaskByIDHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}

	if interceptor == nil {
		output, err := srv.(*Server).GetTaskByIDUseCase.Execute(&usecases.GetTaskByIDInput{
			ID: in.ID,
		})
		if err != nil {
			return nil, err
		}
		return &TaskByIDResponse{Task: mapTask(output.Task)}, nil
	}

	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/task.TaskService/GetTaskByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		output, err := srv.(*Server).GetTaskByIDUseCase.Execute(&usecases.GetTaskByIDInput{
			ID: req.(*TaskByIDRequest).ID,
		})
		if err != nil {
			return nil, err
		}
		return &TaskByIDResponse{Task: mapTask(output.Task)}, nil
	}
	return interceptor(ctx, in, info, handler)
}

func mapTask(task entities.Task) TaskRecord {
	return TaskRecord{
		ID:         task.ID,
		Message:    task.Message,
		BinaryCode: task.BinaryCode,
		CreatedAt:  task.CreatedAt,
		UpdatedAt:  task.UpdatedAt,
	}
}

func mapTasks(tasks []entities.Task) []TaskRecord {
	mapped := make([]TaskRecord, 0, len(tasks))
	for _, task := range tasks {
		mapped = append(mapped, mapTask(task))
	}
	return mapped
}
