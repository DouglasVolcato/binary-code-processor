package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/douglasvolcato/binary-code-processor/task_service/internal/infra/database"
	taskgrpc "github.com/douglasvolcato/binary-code-processor/task_service/internal/infra/grpc"
	"github.com/douglasvolcato/binary-code-processor/task_service/internal/infra/id"
	"github.com/douglasvolcato/binary-code-processor/task_service/internal/usecases"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

const defaultDatabaseURL = "postgres://postgres:postgres@localhost:5432/task_service?sslmode=disable"
const defaultPort = "50051"
const defaultMetricsPort = "8081"

func main() {
	loadDotEnv()

	databaseURL := getenv("DATABASE_URL", defaultDatabaseURL)
	port := getenv("GRPC_PORT", defaultPort)
	metricsPort := getenv("METRICS_PORT", defaultMetricsPort)

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

	go serveHTTPMetrics(metricsPort)

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

func serveHTTPMetrics(port string) {
	mux := http.NewServeMux()
	mux.Handle("/health", healthHandler())
	mux.Handle("/healthz", healthHandler())
	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("task service metrics listening on :%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
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

func healthHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
}
