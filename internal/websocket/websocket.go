package websocket

import (
	"github.com/dana-team/platform-backend/internal/middleware"
	"github.com/dana-team/platform-backend/internal/terminal"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"log"
	"net/http"
)

// WebSocketHandler defines an interface for registering a WebSocket connection.
type WebSocketHandler interface {
	Register(c *gin.Context) (*websocket.Conn, error)
}

// WebSocket struct holds the websocket.Upgrader.
type WebSocket struct {
	upgrader websocket.Upgrader
}

// ServeHTTP upgrades the HTTP connection to a WebSocket and handles the terminal session.
func (ws WebSocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error %s when upgrading connection to websocket", err)
		return
	}
	terminal.HandleTerminalSession(c)
}

// NewWebSocket creates a new WebSocket instance with the provided upgrader.
// If no upgrader is provided, it uses the DefaultUpgrader.
func NewWebSocket(upgrader *websocket.Upgrader) *WebSocket {
	if upgrader == nil {
		upgrader = DefaultUpgrader()
	}

	return &WebSocket{
		upgrader: *upgrader,
	}
}

// Register upgrades an HTTP connection to a WebSocket connection.
// It retrieves the token from the context and sets it in the WebSocket headers.
func (ws *WebSocket) Register(c *gin.Context) (*websocket.Conn, error) {
	token, exists := c.Get("token")
	if !exists {
		return nil, errors.New("Token not found")
	}

	h := http.Header{}
	h.Set(middleware.WebsocketTokenHeader, token.(string))

	conn, err := ws.upgrader.Upgrade(c.Writer, c.Request, h)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// CreateAttachHandler creates new web socket handler to handle terminal web socket
func CreateAttachHandler() http.Handler {
	handler := NewWebSocket(nil)
	return handler
}
