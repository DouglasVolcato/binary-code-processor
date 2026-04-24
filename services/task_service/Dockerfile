FROM golang:1.25.1-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags=prod -trimpath -ldflags="-s -w" -o /out/task_service cmd/main.go

FROM scratch

COPY --from=build /out/task_service /task_service
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 50051
ENTRYPOINT ["/task_service"]
