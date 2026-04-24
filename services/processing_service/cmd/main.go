package main

import (
	"encoding/json"
	"log"
	"os"

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
	rabbitURL := getenv("RABBITMQ_URL", defaultRabbitURL)
	taskServiceAddr := getenv("TASK_SERVICE_ADDR", defaultTaskServiceAddr)
	queueName := getenv("PROCESS_QUEUE", defaultQueueName)

	conn, err := grpc.Dial(taskServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
		_, err := usecase.Execute(&usecases.ProcessTaskInput{ID: msg.ID})
		return err
	}); err != nil {
		log.Fatal(err)
	}

	select {}
}

func getenv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
