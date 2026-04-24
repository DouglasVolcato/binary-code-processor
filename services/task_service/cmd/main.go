package main

import (
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/douglasvolcato/binary-code-processor/task_service/internal/infra/database"
	taskgrpc "github.com/douglasvolcato/binary-code-processor/task_service/internal/infra/grpc"
	"github.com/douglasvolcato/binary-code-processor/task_service/internal/infra/id"
	"github.com/douglasvolcato/binary-code-processor/task_service/internal/usecases"
	"google.golang.org/grpc"
)

const defaultDatabaseURL = "postgres://postgres:postgres@localhost:5432/task_service?sslmode=disable"
const defaultPort = "50051"

func main() {
	loadDotEnv()

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

func loadDotEnv() {
	for _, path := range []string{".env", filepath.Join("..", ".env"), filepath.Join("..", "..", ".env")} {
		if err := loadEnvFile(path); err == nil {
			return
		}
	}
}

func loadEnvFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		_ = os.Setenv(strings.TrimSpace(key), strings.Trim(strings.TrimSpace(value), `"'`))
	}

	return nil
}

func getenv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
