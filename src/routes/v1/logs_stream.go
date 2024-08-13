package v1

import (
	"fmt"
	"github.com/dana-team/platform-backend/src/controllers"
	"github.com/dana-team/platform-backend/src/middleware"
	websocketpkg "github.com/dana-team/platform-backend/src/websocket"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	multicluster "github.com/oam-dev/cluster-gateway/pkg/apis/cluster/transport"
	"go.uber.org/zap"
	"io"
	"k8s.io/client-go/kubernetes"
	"net/http"
)

const (
	namespaceParam      = "namespace"
	cappNameParam       = "cappName"
	containerQueryParam = "container"
	podNameQueryParam   = "podName"
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

		logger, err := getLogger(c)
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
	client, err := getKubeClient(c)
	if err != nil {
		return nil, err
	}

	namespace := c.Param(namespaceParam)
	podName := c.Param(podNameQueryParam)
	containerName := c.Query(containerQueryParam)

	clusterName := c.Param(clusterNameParam)
	context := multicluster.WithMultiClusterContext(c.Request.Context(), clusterName)

	return controllers.FetchPodLogs(context, client, namespace, podName, containerName, logger)
}

// streamCappLogs streams logs for a specific Capp.
func streamCappLogs(c *gin.Context, logger *zap.Logger) (io.ReadCloser, error) {
	client, err := getKubeClient(c)
	if err != nil {
		return nil, err
	}

	namespace := c.Param(namespaceParam)
	cappName := c.Param(cappNameParam)
	containerName := c.DefaultQuery(containerQueryParam, cappName)
	podName := c.Query(podNameQueryParam)

	clusterName := c.Param(clusterNameParam)
	context := multicluster.WithMultiClusterContext(c.Request.Context(), clusterName)

	return controllers.FetchCappLogs(context, client, namespace, cappName, containerName, podName, logger)
}

// getKubeClient retrieves the Kubernetes client from the gin.Context.
func getKubeClient(c *gin.Context) (kubernetes.Interface, error) {
	kube, exists := c.Get(middleware.KubeClientCtxKey)
	if !exists {
		return nil, fmt.Errorf("kube client not found")
	}
	return kube.(kubernetes.Interface), nil
}

// getLogger retrieves the logger from the gin.Context.
func getLogger(c *gin.Context) (*zap.Logger, error) {
	logger, exists := c.Get(middleware.LoggerCtxKey)
	if !exists {
		return nil, fmt.Errorf("logger not found in context")
	}
	return logger.(*zap.Logger), nil
}
