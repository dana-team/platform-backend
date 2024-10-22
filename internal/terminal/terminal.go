package terminal

// The code in this file is largely copied from https://github.com/kubernetes/dashboard/blob/release/7.5.0/modules/api/pkg/handler/terminal.go

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
	"net/http"
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

// PtyHandler is what remotecommand expects from a pty.
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

// startProcess executes cmd in the container and connects it up with the ptyHandler (a session).
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

	exec, err := remotecommand.NewSPDYExecutor(cfg, http.MethodPost, req.URL())
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

// WaitForTerminal is called from capp controller as a goroutine.
// Waits for the WebSocket connection to be opened by the client the session to be Bound in handleTerminalSession.
// After the session is bound, it manages the session by opening the stream and closing it in case of an error or timeout.
func WaitForTerminal(k8sClient kubernetes.Interface, cfg *rest.Config, clusterName, namespaceName, podName, containerName, shell, sessionId string, logger *zap.Logger) {
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

// HandleTerminalSession is called for binding any new connections.
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
