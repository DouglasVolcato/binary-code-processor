build:
	@$(MAKE) build-api_gateway
	@$(MAKE) build-event_publisher
	@$(MAKE) build-processing_service
	@$(MAKE) build-task_service
	@$(MAKE) build-websocket_service

docker-build:
	@$(MAKE) docker-build-api_gateway
	@$(MAKE) docker-build-event_publisher
	@$(MAKE) docker-build-processing_service
	@$(MAKE) docker-build-task_service
	@$(MAKE) docker-build-websocket_service

test:
	@$(MAKE) test-v-api_gateway
	@$(MAKE) test-v-event_publisher
	@$(MAKE) test-v-processing_service
	@$(MAKE) test-v-task_service
	@$(MAKE) test-v-websocket_service

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

k8s-up:
	k3d cluster create binary-code-processor -p "8080:30080@loadbalancer" -p "8082:30082@loadbalancer" -p "15672:31572@loadbalancer" -p "9090:30900@loadbalancer" -p "5432:30432@loadbalancer"
	k3d image import api_gateway:latest -c binary-code-processor
	k3d image import task_service:latest -c binary-code-processor
	k3d image import processing_service:latest -c binary-code-processor
	k3d image import websocket_service:latest -c binary-code-processor
	k3d image import event_publisher:latest -c binary-code-processor
	kubectl apply -f k8s/main.yaml
	kubectl get pods
	kubectl get svc

k8s-down:
	kubectl delete -f k8s/main.yaml
	k3d cluster delete binary-code-processor

build-api_gateway:
	cd services/api_gateway && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags=prod -trimpath -ldflags="-s -w" -o bin/api_gateway cmd/main.go

build-event_publisher:
	cd services/event_publisher && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags=prod -trimpath -ldflags="-s -w" -o bin/event_publisher cmd/main.go

build-processing_service:
	cd services/processing_service && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags=prod -trimpath -ldflags="-s -w" -o bin/processing_service cmd/main.go

build-task_service:
	cd services/task_service && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags=prod -trimpath -ldflags="-s -w" -o bin/task_service cmd/main.go

build-websocket_service:
	cd services/websocket_service && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags=prod -trimpath -ldflags="-s -w" -o bin/websocket_service cmd/main.go

docker-build-api_gateway:
	cd services/api_gateway && docker build -t api_gateway:latest .

docker-build-event_publisher:
	cd services/event_publisher && docker build -t event_publisher:latest .

docker-build-processing_service:
	cd services/processing_service && docker build -t processing_service:latest .

docker-build-task_service:
	cd services/task_service && docker build -t task_service:latest .

docker-build-websocket_service:
	cd services/websocket_service && docker build -t websocket_service:latest .

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
