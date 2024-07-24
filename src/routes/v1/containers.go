package v1

import (
	"github.com/dana-team/platform-backend/src/controllers"
	"github.com/dana-team/platform-backend/src/middleware"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"net/http"
)

// containerHandler wraps a handler function with context setup for ContainerController.
func containerHandler(handler func(controller controllers.ContainerController, c *gin.Context) (interface{}, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		client, exists := c.Get(middleware.KubeClientCtxKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Kubernetes client not found"})
			return
		}

		ctxLogger, exists := c.Get(middleware.LoggerCtxKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Logger not found in context"})
			return
		}

		logger := ctxLogger.(*zap.Logger)
		kubeClient := client.(kubernetes.Interface)
		context := c.Request.Context()

		containerController := controllers.NewContainerController(kubeClient, context, logger)
		result, err := handler(containerController, c)
		if err != nil {
			c.AbortWithStatusJSON(int(err.(*k8serrors.StatusError).ErrStatus.Code), gin.H{"error": "Operation failed", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// GetContainers returns a Gin handler function for retrieving containers information.
func GetContainers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request types.ContainerRequestUri
		if err := c.BindUri(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		containerHandler(func(controller controllers.ContainerController, c *gin.Context) (interface{}, error) {
			return controller.GetContainers(request.NamespaceName, request.PodName)
		})(c)
	}
}
