package controllers

import (
	"fmt"
	"github.com/dana-team/platform-backend/src/terminal"
	"github.com/dana-team/platform-backend/src/types"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

// HandleStartTerminal Handles execute shell API call
func HandleStartTerminal(clientSet kubernetes.Interface, config *rest.Config, namespaceName, podName, containerName, shell string, logger *zap.Logger) (types.StartTerminalResponse, error) {
	// TODO: verify input - check if the given capp is the owner of the given pod

	sessionID, err := terminal.GenTerminalSessionId()
	if err != nil {
		logger.Error(fmt.Sprintf("coundn't generate terminal_utils session for %s capp with pod %s and container %s with err: %s",
			namespaceName, podName, containerName, err.Error()))
	}

	terminal.TerminalSessions.Set(sessionID, terminal.TerminalSession{
		Id:       sessionID,
		Bound:    make(chan error),
		SizeChan: make(chan remotecommand.TerminalSize),
	})

	go terminal.WaitForTerminal(clientSet, config, namespaceName, podName, containerName, shell, sessionID)
	return types.StartTerminalResponse{ID: sessionID}, nil
}
