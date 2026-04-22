# binary-code-processor

Distributed system that processes messages and turns them into binary code. It was made using Golang, TDD, Clean Architecture and following scalability principles.

## Technologies:
- Golang: Code language
- WebSocket: Used to get real time statuses
- RabbitMQ: Used to queue processes
- GraphQL: Used to simplify the api interaction
- gRPC: Used to simplify the microservices interaction
- Postgres: Database
- Kubernetes: Used to manage the containers and scale services

## Architecture:
<img src="docs/diagram.png">

### 1 - Client
- Fetches tasks from the API Gateway through GraphQL
- Sends task messages through GraphQL mutations
- Connects to the WebSocket Service for realtime updates

### 2 - API Gateway
- Receives GraphQL queries and mutations from the client
- Validates pagination and message payloads
- Lists tasks with `limit` and `offset`
- Sends `messages[]` to the Task Service
- Does not persist domain data
- Does not handle WebSocket connections

### 3 - Task Service
- Owns the Postgres database as the source of truth
- Creates tasks from incoming messages
- Returns tasks by ID and paginated lists
- Marks tasks as processed when the Processing Service finishes
- Receives only the task `ID` for the processed-task flow

### 4 - Processing Service
- Receives only the task `ID`
- Loads the task from the Task Service
- Converts `message` into `binaryCode`
- Returns `ID + BinaryCode`
- Sends the processed result back through its processor boundary so the task can be marked as processed

### 5 - Event Publisher
- Reads events by status from the outbox
- Reads unprocessed events and sends them to the queue flow
- Reads processed events and sends them to the fanout flow
- Keeps the queue and fanout reads separated
- Bridges event states to RabbitMQ and downstream consumers

### 6 - WebSocket Service
- Receives only `Task{ID, BinaryCode}`
- Sends realtime updates to connected clients
- Does not depend on message text, timestamps, or other task fields

## End-to-end Flow Example

This is the exact path a batch of client messages follows until the result reaches the client again through WebSocket.

### Contract shapes used in the flow

- `api_gateway/internal/entities.Task`: `ID`, `Message`, `BinaryCode`, `CreatedAt`, `UpdatedAt`
- `task_service/internal/entities.Task`: `ID`, `Message`, `BinaryCode`, `CreatedAt`, `UpdatedAt`
- `processing_service/internal/entities.Task`: `ID`, `BinaryCode`
- `websocket_service/internal/entities.Task`: `ID`, `BinaryCode`

### Example input from the client

```text
messages = ["hello", "world"]
```

### 1 - API Gateway

Use case:
`SendTaskToProcessUseCase`

Input type:
`SendTaskToProcessInput`

Input example:
```go
&SendTaskToProcessInput{
  Messages: []string{"hello", "world"},
}
```

Method used:
`Execute(input *SendTaskToProcessInput)`

Boundary method called:
`TaskProcessorInterface.SendTaskToProcess(messages []string)`

Output type:
`SendTaskToProcessOutput`

Output example:
```go
&SendTaskToProcessOutput{
  Success: true,
  Tasks: []entities.Task{
    {
      ID:         "task-1",
      Message:    "hello",
      BinaryCode: "1111111111111111",
      CreatedAt:  "2026-04-22T10:00:00Z",
      UpdatedAt:  "2026-04-22T10:00:00Z",
    },
    {
      ID:         "task-2",
      Message:    "world",
      BinaryCode: "1111111111111111",
      CreatedAt:  "2026-04-22T10:00:01Z",
      UpdatedAt:  "2026-04-22T10:00:01Z",
    },
  },
}
```

### 2 - Task Service, create flow

Use case:
`ReceiveTaskToProcessUseCase`

Input type:
`ReceiveTaskToProcessInput`

Input example:
```go
&ReceiveTaskToProcessInput{
  Message: "hello",
}
```

Method used:
`Execute(input *ReceiveTaskToProcessInput)`

Boundary methods called:
`IDGeneratorInterface.GenerateID()`
`TaskProcessorInterface.MoveTaskToProcessing(createTaskDto CreateTaskDTO)`

Output type:
`ReceiveTaskToProcessOutput`

Output example:
```go
&ReceiveTaskToProcessOutput{
  Success: true,
  Task: entities.Task{
    ID:         "task-1",
    Message:    "hello",
    BinaryCode: "1111111111111111",
    CreatedAt:  "2026-04-22T10:00:00Z",
    UpdatedAt:  "2026-04-22T10:00:00Z",
  },
}
```

The same flow repeats for `"world"`, producing a second task.

### 3 - Processing Service

Use case:
`ProcessTaskUseCase`

Input type:
`ProcessTaskInput`

Input example:
```go
&ProcessTaskInput{
  ID: "task-1",
}
```

Method used:
`Execute(input *ProcessTaskInput)`

Boundary methods called:
`TaskRepositoryInterface.GetTaskByID(taskID string)`
`TaskProcessorInterface.FinishProcessing(dto FinishProcessingDTO)`

Output type:
`ProcessTaskOutput`

Output example:
```go
&ProcessTaskOutput{
  ID: "task-1",
  BinaryCode: "0110100001100101011011000110110001101111",
}
```

### 4 - Task Service, processed flow

Use case:
`ReceiveProcessedTaskUseCase`

Input type:
`ReceiveProcessedTaskInput`

Input example:
```go
&ReceiveProcessedTaskInput{
  ID: "task-1",
}
```

Method used:
`Execute(input *ReceiveProcessedTaskInput)`

Boundary method called:
`TaskProcessorInterface.SetTaskAsProcessed(taskID string)`

Output type:
`ReceiveProcessedTaskOutput`

Output example:
```go
&ReceiveProcessedTaskOutput{
  Success: true,
}
```

### 5 - WebSocket Service

Use case:
`SendProcessedTasksUseCase`

Input type:
`SendProcessedTasksInput`

Input example:
```go
&SendProcessedTasksInput{
  Task: entities.Task{
    ID: "task-1",
    BinaryCode: "0110100001100101011011000110110001101111",
  },
}
```

Method used:
`Execute(input *SendProcessedTasksInput)`

Boundary method called:
`WebSocketClient.SendProcessedTasksToClient(task entities.Task)`

Output type:
`SendProcessedTasksOutput`

Output example:
```go
&SendProcessedTasksOutput{}
```

### Result

For each message in the original batch, the system follows this chain:

`SendTaskToProcessUseCase -> ReceiveTaskToProcessUseCase -> ProcessTaskUseCase -> ReceiveProcessedTaskUseCase -> SendProcessedTasksUseCase`

The client gets the final update through WebSocket as a `Task` with `ID` and `BinaryCode`.

## Author:
<a href="https://github.com/DouglasVolcato?tab=repositories">Douglas Volcato</a>
