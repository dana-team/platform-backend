package controllers

import (
	"context"
	"fmt"
	"github.com/dana-team/platform-backend/src/terminal"
	"github.com/dana-team/platform-backend/src/types"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

// HandleStartTerminal Handles execute shell API call
func HandleStartTerminal(ctx context.Context, clientSet kubernetes.Interface, config *rest.Config, clusterName, namespaceName, podName, containerName, shell string, logger *zap.Logger) (types.StartTerminalResponse, error) {

	sessionID, err := terminal.GenTerminalSessionId()
	if err != nil {
		logger.Error(fmt.Sprintf("coundn't generate terminal_utils session for pod %s and container %s in namespace %s on cluster %s with err: %s",
			podName, containerName, namespaceName, clusterName, err.Error()))
	}

	terminal.TerminalSessions.Set(sessionID, terminal.TerminalSession{
		Id:       sessionID,
		Bound:    make(chan error),
		SizeChan: make(chan remotecommand.TerminalSize),
	})

	go terminal.WaitForTerminal(ctx, clientSet, config, clusterName, namespaceName, podName, containerName, shell, sessionID, logger)
	return types.StartTerminalResponse{ID: sessionID}, nil
}
