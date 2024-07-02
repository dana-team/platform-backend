package v1

import (
	"net/http"

	"github.com/dana-team/platform-backend/src/controllers"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
)

// configMapHandler handles the request of the client to the Kubernetes cluster.
func configMapHandler(handler func(controller controllers.ConfigMapController, c *gin.Context) (interface{}, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		client, exists := c.Get("kubeClient")
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Kubernetes client not found"})
			return
		}

		ctxLogger, exists := c.Get("logger")
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Logger not found in context"})
			return
		}

		logger := ctxLogger.(*zap.Logger)
		kubeClient := client.(kubernetes.Interface)
		context := c.Request.Context()

		configMapController := controllers.NewConfigMapController(kubeClient, context, logger)
		result, err := handler(configMapController, c)
		if err != nil {
			c.AbortWithStatusJSON(int(err.(*k8serrors.StatusError).ErrStatus.Code), gin.H{"error": "Operation failed", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// GetConfigMap gets a specific config map from the specified namespace.
func GetConfigMap() gin.HandlerFunc {
	return func(c *gin.Context) {
		var uriRequest types.ConfigMapUri
		if err := c.BindUri(&uriRequest); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		configMapHandler(func(controller controllers.ConfigMapController, c *gin.Context) (interface{}, error) {
			return controller.GetConfigMap(uriRequest.NamespaceName, uriRequest.ConfigMapName)
		})(c)
	}
}
