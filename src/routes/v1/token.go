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

// tokenHandler handles the request of the client to the Kubernetes cluster.
func tokenHandler(handler func(controller controllers.TokenController, c *gin.Context) (interface{}, error)) gin.HandlerFunc {
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

		tokenController := controllers.NewTokenController(kubeClient, context, logger)
		result, err := handler(tokenController, c)
		if err != nil {
			c.AbortWithStatusJSON(int(err.(*k8serrors.StatusError).ErrStatus.Code), gin.H{"error": "Operation failed", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// CreateServiceAccount creates a new service account in the specified namespace.
func CreateServiceAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var uriRequest types.NamespaceUri
		if err := c.BindUri(&uriRequest); err != nil {
			return
		}
		var request types.ServiceAccount
		if err := c.BindJSON(&request); err != nil {
			return
		}

		tokenHandler(func(controller controllers.TokenController, c *gin.Context) (interface{}, error) {
			return controller.CreateServiceAccount(uriRequest.NamespaceName, request)
		})(c)
	}
}

// GetToken gets a new token from the specified service account.
func GetToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var uriRequest types.NamespaceUri
		if err := c.BindUri(&uriRequest); err != nil {
			return
		}
		var request types.ServiceAccount
		if err := c.BindJSON(&request); err != nil {
			return
		}

		tokenHandler(func(controller controllers.TokenController, c *gin.Context) (interface{}, error) {
			return controller.GetToken(uriRequest.NamespaceName, request)
		})(c)
	}
}
