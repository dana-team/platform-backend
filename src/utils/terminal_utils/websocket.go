package terminal_utils

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type WebSocketHandler struct {
	upgrader websocket.Upgrader
}

func NewWebSocketHandler() WebSocketHandler {
	webSocketHandler := WebSocketHandler{
		upgrader: websocket.Upgrader{},
	}
	return webSocketHandler
}

func (wsh WebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error %s when upgrading connection to websocket", err)
		return
	}
	handleTerminalSession(c)
}

// handleTerminalSession is Called for binding any new connections.
func handleTerminalSession(session *websocket.Conn) {
	var (
		buf             []byte
		err             error
		msg             TerminalMessage
		terminalSession TerminalSession
		mt              int
	)

	if mt, buf, err = session.ReadMessage(); err != nil {
		log.Printf("handleTerminalSession: can't Recv: %v", err)
		return
	}
	if mt == websocket.BinaryMessage {
		log.Printf("handleTerminalSession: server doesn't support binary messages")
	}

	if err = json.Unmarshal(buf, &msg); err != nil {
		log.Printf("handleTerminalSession: can't UnMarshal (%v): %s", err, buf)
		return
	}

	if msg.Op != "bind" {
		log.Printf("handleTerminalSession: expected 'bind' message, got: %s", buf)
		return
	}

	if terminalSession = TerminalSessions.Get(msg.SessionID); terminalSession.Id == "" {
		log.Printf("handleTerminalSession: can't find session '%s'", msg.SessionID)
		return
	}

	terminalSession.WebSocket = session
	TerminalSessions.Set(msg.SessionID, terminalSession)
	terminalSession.Bound <- nil
}
