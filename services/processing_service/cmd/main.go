package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	taskgrpc "github.com/douglasvolcato/binary-code-processor/processing_service/internal/infra/grpc"
	"github.com/douglasvolcato/binary-code-processor/processing_service/internal/queue"
	"github.com/douglasvolcato/binary-code-processor/processing_service/internal/usecases"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const defaultRabbitURL = "amqp://guest:guest@localhost:5672/"
const defaultTaskServiceAddr = "localhost:50051"
const defaultQueueName = "task.process"

type eventMessage struct {
	ID string `json:"id"`
}

func main() {
	loadDotEnv()

	rabbitURL := getenv("RABBITMQ_URL", defaultRabbitURL)
	taskServiceAddr := getenv("TASK_SERVICE_ADDR", defaultTaskServiceAddr)
	queueName := getenv("PROCESS_QUEUE", defaultQueueName)

	dialCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(dialCtx, taskServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := taskgrpc.NewClient(conn)
	usecase := usecases.NewProcessTaskUseCase(client, client)

	consumer, err := queue.NewConsumer(rabbitURL)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	if err := consumer.DeclareQueue(queueName); err != nil {
		log.Fatal(err)
	}

	if err := consumer.Consume(queueName, func(payload []byte) error {
		var msg eventMessage
		if err := json.Unmarshal(payload, &msg); err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		_, err := usecase.Execute(&usecases.ProcessTaskInput{
			Ctx: ctx,
			ID:  strings.TrimSpace(msg.ID),
		})
		return err
	}); err != nil {
		log.Fatal(err)
	}

	select {}
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
