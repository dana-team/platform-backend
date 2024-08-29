package v1

import (
	"fmt"
	"github.com/dana-team/platform-backend/src/controllers"
	"github.com/dana-team/platform-backend/src/routes"
	websocketpkg "github.com/dana-team/platform-backend/src/websocket"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"io"
	"net/http"
)

const (
	namespaceParam      = "namespace"
	cappNameParam       = "cappName"
	containerQueryParam = "container"
	podNameQueryParam   = "cappName"
	previousQueryParam  = "previous"
	trueValue           = "true"
	trueValueCapital    = "True"
)

// GetPodLogs returns a handler function that fetches logs for a specified pod and container.
func GetPodLogs() gin.HandlerFunc {
	return createLogHandler(streamPodLogs, podNameQueryParam, "Pod")
}

// GetCappLogs returns a handler function that fetches logs for a specified Capp.
func GetCappLogs() gin.HandlerFunc {
	return createLogHandler(streamCappLogs, cappNameParam, "Capp")
}

// createLogHandler creates a gin.HandlerFunc for streaming logs using the provided stream function.
func createLogHandler(streamFunc func(*gin.Context, *zap.Logger) (io.ReadCloser, error), paramKey, logPrefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !websocket.IsWebSocketUpgrade(c.Request) {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		logger, err := routes.GetLogger(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error getting logger %s", err.Error())})
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
	client, err := routes.GetKubeClient(c)
	if err != nil {
		return nil, err
	}

	namespace := c.Param(namespaceParam)
	podName := c.Param(podNameQueryParam)
	containerName := c.Query(containerQueryParam)

	return controllers.FetchPodLogs(c.Request.Context(), client, namespace, podName, containerName, isPreviousLogsRequested(c), logger)
}

// streamCappLogs streams logs for a specific Capp.
func streamCappLogs(c *gin.Context, logger *zap.Logger) (io.ReadCloser, error) {
	client, err := routes.GetKubeClient(c)
	if err != nil {
		return nil, err
	}

	namespace := c.Param(namespaceParam)
	cappName := c.Param(cappNameParam)
	containerName := c.DefaultQuery(containerQueryParam, cappName)
	podName := c.Query(podNameQueryParam)

	return controllers.FetchCappLogs(c.Request.Context(), client, namespace, cappName, containerName, podName, isPreviousLogsRequested(c), logger)
}

// isPreviousLogsRequested returns true if the query parameter for previous logs is set to "true" or "True".
func isPreviousLogsRequested(c *gin.Context) bool {
	previous := c.Query(previousQueryParam)
	return previous == trueValue || previous == trueValueCapital
}
