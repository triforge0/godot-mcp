package bridge

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

var (
	ErrNotConnected = errors.New("godot plugin is not connected")
	ErrTimeout      = errors.New("godot plugin request timed out")
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type EventHandler func(eventType string, data json.RawMessage)

// WebSocket implements JSON-RPC transport to the Godot plugin.
type WebSocket struct {
	addr    string
	mu      sync.RWMutex
	conn    *websocket.Conn
	pending map[string]chan Response
	nextID  atomic.Uint64
	onEvent EventHandler
}

func NewWebSocket(addr string) *WebSocket {
	return &WebSocket{
		addr:    addr,
		pending: make(map[string]chan Response),
	}
}

func (s *WebSocket) Addr() string  { return s.addr }
func (s *WebSocket) Connected() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.conn != nil
}

func (s *WebSocket) SetEventHandler(handler EventHandler) {
	s.mu.Lock()
	s.onEvent = handler
	s.mu.Unlock()
}

func (s *WebSocket) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/ws", s.handleWebSocket)

	server := &http.Server{Addr: s.addr, Handler: mux}
	go func() {
		<-ctx.Done()
		_ = server.Shutdown(context.Background())
	}()

	slog.Info("godot bridge listening", "addr", s.addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *WebSocket) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	status := "disconnected"
	if s.Connected() {
		status = "connected"
	}
	_ = json.NewEncoder(w).Encode(map[string]string{"status": status})
}

func (s *WebSocket) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("websocket upgrade failed", "error", err)
		return
	}

	s.mu.Lock()
	if s.conn != nil {
		_ = s.conn.Close()
	}
	s.conn = conn
	s.mu.Unlock()

	slog.Info("godot plugin connected", "remote", r.RemoteAddr)
	defer func() {
		s.mu.Lock()
		if s.conn == conn {
			s.conn = nil
		}
		s.mu.Unlock()
		_ = conn.Close()
		slog.Info("godot plugin disconnected")
	}()

	conn.SetReadLimit(16 << 20)
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return
		}
		s.handleIncoming(message)
	}
}

func (s *WebSocket) handleIncoming(message []byte) {
	var probe struct {
		ID     string `json:"id"`
		Method string `json:"method"`
	}
	if err := json.Unmarshal(message, &probe); err != nil {
		return
	}
	if probe.ID == "" && probe.Method == EventMethod {
		s.handleNotification(message)
		return
	}
	s.dispatchResponse(message)
}

func (s *WebSocket) handleNotification(message []byte) {
	var note Notification
	if err := json.Unmarshal(message, &note); err != nil {
		return
	}
	var params EventParams
	if len(note.Params) > 0 {
		_ = json.Unmarshal(note.Params, &params)
	}
	s.mu.RLock()
	handler := s.onEvent
	s.mu.RUnlock()
	if handler != nil && params.Type != "" {
		handler(params.Type, params.Data)
	}
}

func (s *WebSocket) dispatchResponse(message []byte) {
	var resp Response
	if err := json.Unmarshal(message, &resp); err != nil || resp.ID == "" {
		return
	}
	s.mu.RLock()
	ch, ok := s.pending[resp.ID]
	s.mu.RUnlock()
	if !ok {
		return
	}
	select {
	case ch <- resp:
	default:
	}
}

func (s *WebSocket) Call(ctx context.Context, method string, params any) (json.RawMessage, error) {
	s.mu.RLock()
	conn := s.conn
	s.mu.RUnlock()
	if conn == nil {
		return nil, ErrNotConnected
	}

	var rawParams json.RawMessage
	if params != nil {
		b, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("marshal params: %w", err)
		}
		rawParams = b
	}

	id := fmt.Sprintf("%d", s.nextID.Add(1))
	req := Request{JSONRPC: JSONRPCVersion, ID: id, Method: method, Params: rawParams}

	ch := make(chan Response, 1)
	s.mu.Lock()
	s.pending[id] = ch
	s.mu.Unlock()
	defer func() {
		s.mu.Lock()
		delete(s.pending, id)
		s.mu.Unlock()
	}()

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	s.mu.RLock()
	writeConn := s.conn
	s.mu.RUnlock()
	if writeConn == nil {
		return nil, ErrNotConnected
	}
	if err := writeConn.WriteMessage(websocket.TextMessage, payload); err != nil {
		return nil, err
	}

	timeout := 30 * time.Second
	if deadline, ok := ctx.Deadline(); ok {
		timeout = time.Until(deadline)
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case resp := <-ch:
		if resp.Error != nil {
			return nil, fmt.Errorf("plugin error %d: %s", resp.Error.Code, resp.Error.Message)
		}
		return resp.Result, nil
	case <-time.After(timeout):
		return nil, ErrTimeout
	}
}
