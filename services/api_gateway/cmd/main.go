package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/douglasvolcato/binary-code-processor/api_gateway/internal/infra/graphql"
	gatewaygrpc "github.com/douglasvolcato/binary-code-processor/api_gateway/internal/infra/grpc"
	"github.com/douglasvolcato/binary-code-processor/api_gateway/internal/infra/web"
	"github.com/douglasvolcato/binary-code-processor/api_gateway/internal/usecases"
	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const defaultPort = "8080"
const defaultTaskServiceAddr = "localhost:50051"

func main() {
	loadDotEnv()

	port := getenv("API_GATEWAY_PORT", defaultPort)
	taskServiceAddr := getenv("TASK_SERVICE_ADDR", defaultTaskServiceAddr)
	websocketURL := getenv("WEBSOCKET_URL", "")
	websocketPort := getenv("WEBSOCKET_PORT", "8082")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, taskServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := gatewaygrpc.NewClient(conn)
	resolver := &graphql.Resolver{
		GetTasksUseCase:          usecases.NewGetTasksUseCase(client),
		SendTaskToProcessUseCase: usecases.NewSendTaskToProcessUseCase(client),
	}

	gqlServer := handler.New(graphql.NewExecutableSchema(graphql.Config{Resolvers: resolver}))
	gqlServer.AddTransport(transport.Options{})
	gqlServer.AddTransport(transport.GET{})
	gqlServer.AddTransport(transport.POST{})
	gqlServer.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	gqlServer.Use(extension.Introspection{})
	gqlServer.Use(extension.AutomaticPersistedQuery{Cache: lru.New[string](100)})

	home, err := web.NewServer(websocketURL, websocketPort)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", home)
	mux.Handle("/query", gqlServer)
	mux.Handle("/healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	log.Printf("api gateway listening on :%s", port)
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
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		_ = os.Setenv(key, strings.Trim(value, `"'`))
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
