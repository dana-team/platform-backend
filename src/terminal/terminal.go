package terminal

// Copyright 2017 The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	multicluster "github.com/oam-dev/cluster-gateway/pkg/apis/cluster/transport"
	"go.uber.org/zap"
	"io"
	"log"
	"sync"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

const (
	END_OF_TRANSMISSION = "\u0004"
	terminalTimeOut     = 600
	resource            = "pods"
	operation           = "exec"
)

// PtyHandler is what remotecommand expects from a pty
type PtyHandler interface {
	io.Reader
	io.Writer
	remotecommand.TerminalSizeQueue
}

// TerminalMessage is the messaging protocol between the Client and the TerminalSession.
//
// OP      DIRECTION  FIELD(S) USED  DESCRIPTION
// ---------------------------------------------------------------------
// bind    fe->be     SessionID      Id sent back from TerminalResponse
// stdin   fe->be     Data           Keystrokes/paste buffer
// resize  fe->be     Rows, Cols     New terminal_utils size
// stdout  be->fe     Data           Output from the process
// toast   be->fe     Data           OOB message to be shown to the user
type TerminalMessage struct {
	Op, Data, SessionID string
	Rows, Cols          uint16
}

// TerminalSessions is a global variable that used to manage all terminal sessions
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

// Toast can be used to send the user any OOB messages
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

// SessionMap stores a map of all TerminalSession objects and a lock to avoid concurrent conflict
type SessionMap struct {
	Sessions map[string]TerminalSession
	Lock     sync.RWMutex
}

// Get return a given terminalSession by sessionId
func (sm *SessionMap) Get(sessionId string) TerminalSession {
	sm.Lock.RLock()
	defer sm.Lock.RUnlock()
	return sm.Sessions[sessionId]
}

// Set store a TerminalSession to SessionMap
func (sm *SessionMap) Set(sessionId string, session TerminalSession) {
	sm.Lock.Lock()
	defer sm.Lock.Unlock()
	sm.Sessions[sessionId] = session
}

// Close shuts down the WebSocket connection.
// Can happen if the process exits or if there is an error starting up the process
// For now the status code is unused and reason is shown to the user (unless "")
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

// startProcess executes cmd in the container and connects it up with the ptyHandler (a session)
func startProcess(k8sClient kubernetes.Interface, cfg *rest.Config, clusterName, namespaceName, podName, containerName string, cmd []string, ptyHandler PtyHandler) error {
	req := k8sClient.CoreV1().RESTClient().Post().
		Resource(resource).
		Name(podName).
		Namespace(namespaceName).
		SubResource(operation)
	req.VersionedParams(&v1.PodExecOptions{
		Container: containerName,
		Command:   cmd,
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(cfg, "POST", req.URL())
	if err != nil {
		return err
	}

	ctx := multicluster.WithMultiClusterContext(context.Background(), clusterName)
	err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:             ptyHandler,
		Stdout:            ptyHandler,
		Stderr:            ptyHandler,
		TerminalSizeQueue: ptyHandler,
		Tty:               true,
	})
	if err != nil {
		return err
	}

	return nil
}

// GenTerminalSessionId generates a random session ID string.
// This ID is used to identify the session when the client opens the WebSocket connection.
func GenTerminalSessionId() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	id := make([]byte, hex.EncodedLen(len(bytes)))
	hex.Encode(id, bytes)
	return string(id), nil
}

// isValidShell checks if the shell is an allowed one
func isValidShell(validShells []string, shell string) bool {
	for _, validShell := range validShells {
		if validShell == shell {
			return true
		}
	}
	return false
}

// WaitForTerminal is called from capp controller as a goroutine
// Waits for the WebSocket connection to be opened by the client the session to be Bound in handleTerminalSession
func WaitForTerminal(ctx context.Context, k8sClient kubernetes.Interface, cfg *rest.Config, clusterName, namespaceName, podName, containerName, shell, sessionId string, logger *zap.Logger) {
	select {
	case <-TerminalSessions.Get(sessionId).Bound:
		close(TerminalSessions.Get(sessionId).Bound)

		var err error
		validShells := []string{"bash", "sh", "powershell", "cmd"}

		if isValidShell(validShells, shell) {
			cmd := []string{shell}
			err = startProcess(k8sClient, cfg, clusterName, namespaceName, podName, containerName, cmd, TerminalSessions.Get(sessionId))
		} else {
			// No shell given or it was not valid: try some shells until one succeeds or all fail
			// FIXME: if the first shell fails then the first keyboard event is lost
			for _, testShell := range validShells {
				cmd := []string{testShell}
				if err = startProcess(k8sClient, cfg, clusterName, namespaceName, podName, containerName, cmd, TerminalSessions.Get(sessionId)); err == nil {
					break
				}
			}
		}

		if err != nil {
			logger.Error(fmt.Sprintf("coundn't start terminal for pod %s and container %s in namespace %s on cluster %s with err: %s",
				podName, containerName, namespaceName, clusterName, err.Error()))
			// status 2 - error
			TerminalSessions.Close(sessionId)
			return
		}

		// status 1 - process exited
		TerminalSessions.Close(sessionId)

	case <-time.After(terminalTimeOut * time.Second):
		// Close chan and delete session when sockjs connection was timeout
		close(TerminalSessions.Get(sessionId).Bound)
		delete(TerminalSessions.Sessions, sessionId)
		return
	}
}

// HandleTerminalSession is Called for binding any new connections.
func HandleTerminalSession(session *websocket.Conn) {
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
