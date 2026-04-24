package main

import (
	"log"
	"net"
	"os"

	"github.com/douglasvolcato/binary-code-processor/task_service/internal/infra/database"
	taskgrpc "github.com/douglasvolcato/binary-code-processor/task_service/internal/infra/grpc"
	"github.com/douglasvolcato/binary-code-processor/task_service/internal/infra/id"
	"github.com/douglasvolcato/binary-code-processor/task_service/internal/usecases"
	"google.golang.org/grpc"
)

const defaultDatabaseURL = "postgres://postgres:postgres@localhost:5432/task_service?sslmode=disable"
const defaultPort = "50051"

func main() {
	databaseURL := getenv("DATABASE_URL", defaultDatabaseURL)
	port := getenv("GRPC_PORT", defaultPort)

	db, err := database.NewDB(databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := database.NewRepository(db)
	idGen := id.NewGenerator()

	receiveTaskUseCase := usecases.NewReceiveTaskToProcessUseCase(repo, idGen)
	receiveProcessedUseCase := usecases.NewReceiveProcessedTaskUseCase(repo)
	getTasksUseCase := usecases.NewGetTasksUseCase(repo)
	getTaskByIDUseCase := usecases.NewGetTaskByIDUseCase(repo)

	server := grpc.NewServer(grpc.ForceServerCodec(taskgrpc.NewJSONCodec()))
	taskgrpc.RegisterTaskAPIServer(server, &taskgrpc.Server{
		ReceiveTaskToProcessUseCase: receiveTaskUseCase,
		ReceiveProcessedTaskUseCase: receiveProcessedUseCase,
		GetTasksUseCase:             getTasksUseCase,
		GetTaskByIDUseCase:          getTaskByIDUseCase,
	})

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("task service listening on :%s", port)
	log.Fatal(server.Serve(listener))
}

func getenv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
