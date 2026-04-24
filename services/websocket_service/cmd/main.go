package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	loadDotEnv()

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
				ID:         strings.TrimSpace(msg.ID),
				BinaryCode: strings.TrimSpace(msg.BinaryCode),
			},
		})
		return err
	}); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/ws", hub)
	mux.Handle("/healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("websocket service listening on :%s", port)
	log.Fatal(server.ListenAndServe())
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
