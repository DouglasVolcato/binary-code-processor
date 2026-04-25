package websocket

import (
	"net/http"
	"sync"

	"github.com/douglasvolcato/binary-code-processor/websocket_service/internal/entities"
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus"
)

type Hub struct {
	mu       sync.RWMutex
	clients  map[*websocket.Conn]*clientState
	upgrader websocket.Upgrader
}

var (
	websocketConnectionsActive = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "websocket_connections_active",
		Help: "Number of active WebSocket connections.",
	})
	registerMetricsOnce sync.Once
)

type clientState struct {
	conn      *websocket.Conn
	writeMu   sync.Mutex
	closeOnce sync.Once
}

func NewHub() *Hub {
	registerMetricsOnce.Do(func() {
		prometheus.MustRegister(websocketConnectionsActive)
	})

	return &Hub{
		clients: make(map[*websocket.Conn]*clientState),
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

	state := &clientState{conn: conn}

	h.mu.Lock()
	h.clients[conn] = state
	h.mu.Unlock()
	websocketConnectionsActive.Inc()

	go h.readLoop(state)
}

func (h *Hub) readLoop(state *clientState) {
	defer state.close(h)

	for {
		if _, _, err := state.conn.ReadMessage(); err != nil {
			return
		}
	}
}

func (h *Hub) SendProcessedTasksToClient(task entities.Task) error {
	if err := task.Validate(); err != nil {
		return err
	}

	h.mu.RLock()
	clients := make([]*clientState, 0, len(h.clients))
	for _, state := range h.clients {
		clients = append(clients, state)
	}
	h.mu.RUnlock()

	for _, state := range clients {
		state.writeMu.Lock()
		if err := state.conn.WriteJSON(taskPayload{
			ID:         task.ID,
			BinaryCode: task.BinaryCode,
		}); err != nil {
			state.writeMu.Unlock()
			state.close(h)
			continue
		}
		state.writeMu.Unlock()
	}
	return nil
}

type taskPayload struct {
	ID         string `json:"id"`
	BinaryCode string `json:"binaryCode"`
}

func (h *Hub) Close() {
	h.mu.RLock()
	clients := make([]*clientState, 0, len(h.clients))
	for _, state := range h.clients {
		clients = append(clients, state)
	}
	h.mu.RUnlock()

	for _, state := range clients {
		state.close(h)
	}
}

func (s *clientState) close(h *Hub) {
	s.closeOnce.Do(func() {
		h.mu.Lock()
		delete(h.clients, s.conn)
		h.mu.Unlock()
		websocketConnectionsActive.Dec()
		_ = s.conn.Close()
	})
}
