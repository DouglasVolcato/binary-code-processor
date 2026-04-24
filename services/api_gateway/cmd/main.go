package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/douglasvolcato/binary-code-processor/api_gateway/internal/infra/graphql"
	taskgrpc "github.com/douglasvolcato/binary-code-processor/api_gateway/internal/infra/grpc"
	"github.com/douglasvolcato/binary-code-processor/api_gateway/internal/usecases"
	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const defaultPort = "8080"
const defaultTaskServiceAddr = "localhost:50051"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	taskServiceAddr := os.Getenv("TASK_SERVICE_ADDR")
	if taskServiceAddr == "" {
		taskServiceAddr = defaultTaskServiceAddr
	}

	conn, err := grpc.Dial(taskServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := taskgrpc.NewClient(conn)
	resolver := &graphql.Resolver{
		GetTasksUseCase:          usecases.NewGetTasksUseCase(client),
		SendTaskToProcessUseCase: usecases.NewSendTaskToProcessUseCase(client),
	}

	srv := handler.New(graphql.NewExecutableSchema(graphql.Config{Resolvers: resolver}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
