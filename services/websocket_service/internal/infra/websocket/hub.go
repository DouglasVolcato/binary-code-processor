package websocket

import (
	"net/http"
	"sync"

	"github.com/douglasvolcato/binary-code-processor/websocket_service/internal/entities"
	"github.com/gorilla/websocket"
)

type Hub struct {
	mu       sync.RWMutex
	writeMu  sync.Mutex
	clients  map[*websocket.Conn]struct{}
	upgrader websocket.Upgrader
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*websocket.Conn]struct{}),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	h.clients[conn] = struct{}{}
	h.mu.Unlock()

	go h.readLoop(conn)
}

func (h *Hub) readLoop(conn *websocket.Conn) {
	defer func() {
		h.mu.Lock()
		delete(h.clients, conn)
		h.mu.Unlock()
		_ = conn.Close()
	}()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
	}
}

func (h *Hub) SendProcessedTasksToClient(task entities.Task) error {
	if err := task.Validate(); err != nil {
		return err
	}

	h.mu.RLock()
	conns := make([]*websocket.Conn, 0, len(h.clients))
	for conn := range h.clients {
		conns = append(conns, conn)
	}
	h.mu.RUnlock()

	for _, conn := range conns {
		h.writeMu.Lock()
		if err := conn.WriteJSON(taskPayload{
			ID:         task.ID,
			BinaryCode: task.BinaryCode,
		}); err != nil {
			h.writeMu.Unlock()
			_ = conn.Close()
			h.mu.Lock()
			delete(h.clients, conn)
			h.mu.Unlock()
			continue
		}
		h.writeMu.Unlock()
	}
	return nil
}

type taskPayload struct {
	ID         string `json:"id"`
	BinaryCode string `json:"binaryCode"`
}

func (h *Hub) Close() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for conn := range h.clients {
		_ = conn.Close()
		delete(h.clients, conn)
	}
}
