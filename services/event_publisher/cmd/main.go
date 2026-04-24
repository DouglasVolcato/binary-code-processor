package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/douglasvolcato/binary-code-processor/event_publisher/internal/infra/database"
	"github.com/douglasvolcato/binary-code-processor/event_publisher/internal/infra/processor"
	"github.com/douglasvolcato/binary-code-processor/event_publisher/internal/infra/rabbitmq"
	"github.com/douglasvolcato/binary-code-processor/event_publisher/internal/usecases"
)

const defaultDatabaseURL = "postgres://postgres:postgres@localhost:5432/task_service?sslmode=disable"
const defaultRabbitURL = "amqp://guest:guest@localhost:5672/"

func main() {
	loadDotEnv()

	databaseURL := getenv("DATABASE_URL", defaultDatabaseURL)
	rabbitURL := getenv("RABBITMQ_URL", defaultRabbitURL)

	db, err := database.NewDB(databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := database.NewRepository(db)
	localProcessor := processor.NewProcessor()
	remoteProcessor, err := rabbitmq.NewBroker(rabbitURL)
	if err != nil {
		log.Fatal(err)
	}
	defer remoteProcessor.Close()

	unprocessedUseCase := usecases.NewProcessUnprocessedEventsUseCase(repo, localProcessor, remoteProcessor)
	processedUseCase := usecases.NewSendProcessedEventsUseCase(repo, localProcessor, remoteProcessor)

	interval := time.Second * 5
	if raw := os.Getenv("POLL_INTERVAL_SECONDS"); raw != "" {
		if parsed, parseErr := time.ParseDuration(raw); parseErr == nil {
			interval = parsed
		} else if parsed, parseErr := time.ParseDuration(raw + "s"); parseErr == nil {
			interval = parsed
		}
	}
	if interval <= 0 {
		interval = 5 * time.Second
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		if _, err := unprocessedUseCase.Execute(&usecases.ProcessUnprocessedEventsInput{}); err != nil {
			log.Println("process unprocessed events:", err)
		}
		if _, err := processedUseCase.Execute(&usecases.SendProcessedEventsInput{}); err != nil {
			log.Println("send processed events:", err)
		}
		<-ticker.C
	}
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
