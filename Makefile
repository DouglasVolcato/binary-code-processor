build:
	@$(MAKE) build-api_gateway
	@$(MAKE) build-event_publisher
	@$(MAKE) build-processing_service
	@$(MAKE) build-task_service
	@$(MAKE) build-websocket_service

test:
	@$(MAKE) test-api_gateway
	@$(MAKE) test-event_publisher
	@$(MAKE) test-processing_service
	@$(MAKE) test-task_service
	@$(MAKE) test-websocket_service

test-cover:
	@$(MAKE) test-coverage-api_gateway
	@$(MAKE) test-coverage-event_publisher
	@$(MAKE) test-coverage-processing_service
	@$(MAKE) test-coverage-task_service
	@$(MAKE) test-coverage-websocket_service

test-bench:
	@$(MAKE) test-bench-api_gateway
	@$(MAKE) test-bench-event_publisher
	@$(MAKE) test-bench-processing_service
	@$(MAKE) test-bench-task_service
	@$(MAKE) test-bench-websocket_service

build-api_gateway:
	cd services/api_gateway && GOOS=linux GOARCH=amd64 go build -o api_gateway main.go

build-event_publisher:
	cd services/event_publisher && GOOS=linux GOARCH=amd64 go build -o event_publisher main.go

build-processing_service:
	cd services/processing_service && GOOS=linux GOARCH=amd64 go build -o processing_service main.go

build-task_service:
	cd services/task_service && GOOS=linux GOARCH=amd64 go build -o task_service main.go

build-websocket_service:
	cd services/websocket_service && GOOS=linux GOARCH=amd64 go build -o websocket_service main.go

test-v-api_gateway:
	cd services/api_gateway && go test -v ./...

test-v-event_publisher:
	cd services/event_publisher && go test -v ./...

test-v-processing_service:
	cd services/processing_service && go test -v ./...

test-v-task_service:
	cd services/task_service && go test -v ./...

test-v-websocket_service:
	cd services/websocket_service && go test -v ./...

test-coverage-api_gateway:
	cd services/api_gateway && go test -coverprofile=test/coverage.out ./... && go tool cover -html=test/coverage.out -o test/coverage.html

test-coverage-event_publisher:
	cd services/event_publisher && go test -coverprofile=test/coverage.out ./... && go tool cover -html=test/coverage.out -o test/coverage.html

test-coverage-processing_service:
	cd services/processing_service && go test -coverprofile=test/coverage.out ./... && go tool cover -html=test/coverage.out -o test/coverage.html

test-coverage-task_service:
	cd services/task_service && go test -coverprofile=test/coverage.out ./... && go tool cover -html=test/coverage.out -o test/coverage.html

test-coverage-websocket_service:
	cd services/websocket_service && go test -coverprofile=test/coverage.out ./... && go tool cover -html=test/coverage.out -o test/coverage.html

test-bench-api_gateway:
	cd services/api_gateway && go test -bench=. -benchmem ./...

test-bench-event_publisher:
	cd services/event_publisher && go test -bench=. -benchmem ./...

test-bench-processing_service:
	cd services/processing_service && go test -bench=. -benchmem ./...

test-bench-task_service:
	cd services/task_service && go test -bench=. -benchmem ./...

test-bench-websocket_service:
	cd services/websocket_service && go test -bench=. -benchmem ./...
