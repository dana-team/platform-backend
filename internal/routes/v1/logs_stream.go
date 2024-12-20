package v1

import (
	"fmt"
	"io"

	"github.com/dana-team/platform-backend/internal/controllers"
	"github.com/dana-team/platform-backend/internal/customerrors"
	"github.com/dana-team/platform-backend/internal/middleware"
	"github.com/dana-team/platform-backend/internal/routes"
	websocketpkg "github.com/dana-team/platform-backend/internal/websocket"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	namespaceParam      = "namespaceName"
	cappNameParam       = "cappName"
	containerQueryParam = "containerName"
	previousQueryParam  = "previous"
	trueValue           = "true"
	trueValueCapital    = "True"
	podNameParam        = "podName"
)

const (
	errWebsocketUpgrade = "failed to upgrade websocket connection"
)

// GetPodLogs returns a handler function that fetches logs for a specified pod and container.
func GetPodLogs() gin.HandlerFunc {
	return createLogHandler(streamPodLogs, podNameParam, "Pod")
}

// GetCappLogs returns a handler function that fetches logs for a specified Capp.
func GetCappLogs() gin.HandlerFunc {
	return createLogHandler(streamCappLogs, cappNameParam, "Capp")
}

// createLogHandler creates a gin.HandlerFunc for streaming logs using the provided stream function.
func createLogHandler(streamFunc func(*gin.Context, *zap.Logger) (io.ReadCloser, error), paramKey, logPrefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !websocket.IsWebSocketUpgrade(c.Request) {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(errWebsocketUpgrade))
			return
		}

		logger, err := middleware.GetLogger(c)
		if middleware.AddErrorToContext(c, err) {
			return
		}

		websocketClient := websocketpkg.NewWebSocket(nil)
		conn, err := websocketClient.Register(c)
		if err != nil {
			logger.Error(fmt.Sprintf("error streaming %q logs: %v", logPrefix, err.Error()))
			return
		}
		defer conn.Close()

		logStream, err := streamFunc(c, logger)
		if err != nil {
			logger.Debug(fmt.Sprintf("Error streaming %q logs: %v", logPrefix, err.Error()))
			websocketpkg.SendErrorMessage(conn, fmt.Sprintf("Error streaming %q logs: %v", logPrefix, err.Error()))
			return
		}
		defer logStream.Close()

		formatFunc := func(line string) string {
			return fmt.Sprintf("%v: %q line: %v", logPrefix, c.Param(paramKey), line)
		}

		websocketpkg.Stream(c, conn, logStream, formatFunc)
	}
}

// streamPodLogs streams logs for a specific pod and container.
func streamPodLogs(c *gin.Context, logger *zap.Logger) (io.ReadCloser, error) {
	client, err := middleware.GetKubeClient(c)
	if err != nil {
		return nil, err
	}

	namespace := c.Param(namespaceParam)
	podName := c.Param(podNameParam)
	containerName := c.Query(containerQueryParam)

	context := routes.GetContext(c)
	return controllers.FetchPodLogs(context, client, namespace, podName, containerName, isPreviousLogsRequested(c), logger)
}

// streamCappLogs streams logs for a specific Capp.
func streamCappLogs(c *gin.Context, logger *zap.Logger) (io.ReadCloser, error) {
	client, err := middleware.GetKubeClient(c)
	if err != nil {
		return nil, err
	}

	namespace := c.Param(namespaceParam)
	cappName := c.Param(cappNameParam)
	containerName := c.DefaultQuery(containerQueryParam, cappName)
	podName := c.Query(podNameParam)

	context := routes.GetContext(c)
	return controllers.FetchCappLogs(context, client, namespace, cappName, containerName, podName, isPreviousLogsRequested(c), logger)
}

// isPreviousLogsRequested returns true if the query parameter for previous logs is set to "true" or "True".
func isPreviousLogsRequested(c *gin.Context) bool {
	previous := c.Query(previousQueryParam)
	return previous == trueValue || previous == trueValueCapital
}
