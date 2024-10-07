package terminal

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"k8s.io/client-go/tools/remotecommand"
	"log"
	"sync"
)

// The code in this file is largely copied from https://github.com/kubernetes/dashboard/blob/release/7.5.0/modules/api/pkg/handler/terminal.go

// TerminalSessions is a global variable that is used to manage all terminal sessions.
var TerminalSessions = SessionMap{Sessions: make(map[string]TerminalSession)}

// TerminalSession implements PtyHandler (using a WebSocket connection)
// This is the object that is used to communicate with container terminal.
type TerminalSession struct {
	Id        string
	Bound     chan error
	WebSocket *websocket.Conn
	SizeChan  chan remotecommand.TerminalSize
}

// Next handles pty->process resize events
// Called in a loop from remotecommand as long as the process is running
func (t TerminalSession) Next() *remotecommand.TerminalSize {
	size := <-t.SizeChan
	if size.Height == 0 && size.Width == 0 {
		return nil
	}
	return &size
}

// Read handles pty->process messages (stdin, resize)
// Called in a loop from remotecommand as long as the process is running
func (t TerminalSession) Read(p []byte) (int, error) {
	mt, message, err := t.WebSocket.ReadMessage()
	if err != nil {
		// Send terminated signal to process to avoid resource leak
		return copy(p, END_OF_TRANSMISSION), err
	}
	if mt == websocket.BinaryMessage {
		return copy(p, END_OF_TRANSMISSION), fmt.Errorf("server doesn't support binary messages")
	}

	var msg TerminalMessage
	if err = json.Unmarshal(message, &msg); err != nil {
		return copy(p, END_OF_TRANSMISSION), err
	}

	switch msg.Op {
	case "stdin":
		return copy(p, msg.Data), nil
	case "resize":
		t.SizeChan <- remotecommand.TerminalSize{Width: msg.Cols, Height: msg.Rows}
		return 0, nil
	default:
		return copy(p, END_OF_TRANSMISSION), fmt.Errorf("unknown message type '%s'", msg.Op)
	}
}

// Write handles process->pty stdout
// Called from remotecommand whenever there is any output
func (t TerminalSession) Write(p []byte) (int, error) {
	msg, err := json.Marshal(TerminalMessage{
		Op:   "stdout",
		Data: string(p),
	})
	if err != nil {
		return 0, err
	}

	if err = t.WebSocket.WriteMessage(websocket.TextMessage, msg); err != nil {
		return 0, err
	}
	return len(p), nil
}

// Toast can be used to send the user any OOB messages.
// (OOB -messages are non-standard messages or notifications sent outside the main communication channel.)
// hterm puts these in the center of the terminal_utils
func (t TerminalSession) Toast(p string) error {
	msg, err := json.Marshal(TerminalMessage{
		Op:   "toast",
		Data: p,
	})
	if err != nil {
		return err
	}

	if err = t.WebSocket.WriteMessage(websocket.TextMessage, msg); err != nil {
		return err
	}
	return nil
}

// SessionMap stores a map of all TerminalSession objects and a lock to avoid concurrent conflict.
type SessionMap struct {
	Sessions map[string]TerminalSession
	Lock     sync.RWMutex
}

// Get returns a given terminalSession by sessionId.
func (sm *SessionMap) Get(sessionId string) TerminalSession {
	sm.Lock.RLock()
	defer sm.Lock.RUnlock()
	return sm.Sessions[sessionId]
}

// Set stores a TerminalSession to the SessionMap.
func (sm *SessionMap) Set(sessionId string, session TerminalSession) {
	sm.Lock.Lock()
	defer sm.Lock.Unlock()
	sm.Sessions[sessionId] = session
}

// Close shuts down the WebSocket connection.
// Can happen if the process exits or if there is an error starting up the process
func (sm *SessionMap) Close(sessionId string) {
	sm.Lock.Lock()
	defer sm.Lock.Unlock()
	ses := sm.Sessions[sessionId]
	err := ses.WebSocket.Close()
	if err != nil {
		log.Println(err)
	}
	close(ses.SizeChan)
	delete(sm.Sessions, sessionId)
}
