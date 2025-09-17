package ws

import (
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/go-hclog"
)

func WebsocketsCheckMiddleware(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

type WebSocketHub struct {
	clients map[string]*websocket.Conn // In future, it should be uuid.UUID
	mu      sync.Mutex
	logger  hclog.Logger
}

func NewWebSocketHub(logger hclog.Logger) *WebSocketHub {
	return &WebSocketHub{
		clients: make(map[string]*websocket.Conn),
		logger:  logger,
	}
}

func (w *WebSocketHub) Push(payload interface{}) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	for clientID, conn := range w.clients {
		// WriteJSON will encode the location and send it
		// to the client
		err := conn.WriteJSON(payload)
		if err != nil {
			// Handle errors gracefully, such as removing disconnected clients
			conn.Close()
			delete(w.clients, clientID)
			return err
		}
	}

	return nil
}

func (w *WebSocketHub) Register(clientId string, conn *websocket.Conn) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.clients[clientId] = conn
}

func (w *WebSocketHub) DeRegister(clientID string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if conn, exists := w.clients[clientID]; exists {
		conn.Close()
		delete(w.clients, clientID)
	}
}
