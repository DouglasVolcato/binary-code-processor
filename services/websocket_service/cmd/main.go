package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/douglasvolcato/binary-code-processor/websocket_service/internal/entities"
	"github.com/douglasvolcato/binary-code-processor/websocket_service/internal/infra/queue"
	ws "github.com/douglasvolcato/binary-code-processor/websocket_service/internal/infra/websocket"
	"github.com/douglasvolcato/binary-code-processor/websocket_service/internal/usecases"
)

const defaultRabbitURL = "amqp://guest:guest@localhost:5672/"
const defaultQueueName = "task.websocket"
const defaultExchangeName = "task.processed"
const defaultPort = "8082"

type eventMessage struct {
	ID         string `json:"id"`
	BinaryCode string `json:"binaryCode"`
}

func main() {
	rabbitURL := getenv("RABBITMQ_URL", defaultRabbitURL)
	queueName := getenv("WEBSOCKET_QUEUE", defaultQueueName)
	exchangeName := getenv("WEBSOCKET_EXCHANGE", defaultExchangeName)
	port := getenv("PORT", defaultPort)

	hub := ws.NewHub()
	usecase := usecases.NewSendProcessedTasksUseCase(hub)

	consumer, err := queue.NewConsumer(rabbitURL)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	if err := consumer.DeclareQueue(queueName); err != nil {
		log.Fatal(err)
	}
	if err := consumer.DeclareExchange(exchangeName, "fanout"); err != nil {
		log.Fatal(err)
	}
	if err := consumer.BindQueueToExchange(queueName, exchangeName); err != nil {
		log.Fatal(err)
	}
	if err := consumer.Consume(queueName, func(payload []byte) error {
		var msg eventMessage
		if err := json.Unmarshal(payload, &msg); err != nil {
			return err
		}
		_, err := usecase.Execute(&usecases.SendProcessedTasksInput{
			Task: entities.Task{
				ID:         msg.ID,
				BinaryCode: msg.BinaryCode,
			},
		})
		return err
	}); err != nil {
		log.Fatal(err)
	}

	http.Handle("/ws", hub)
	log.Printf("websocket service listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getenv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
