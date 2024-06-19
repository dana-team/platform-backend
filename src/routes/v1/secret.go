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

// secretHandler handles the request of the client to the Kubernetes cluster.
func secretHandler(handler func(controller controllers.SecretController, c *gin.Context) (interface{}, error)) gin.HandlerFunc {
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

		secretController := controllers.NewSecretController(kubeClient, context, logger)
		result, err := handler(secretController, c)
		if err != nil {
			c.AbortWithStatusJSON(int(err.(*k8serrors.StatusError).ErrStatus.Code), gin.H{"error": "Operation failed", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// CreateSecret creates a new secret in a specific namespace.
func CreateSecret() gin.HandlerFunc {
	return func(c *gin.Context) {
		var uriRequest types.SecretNamespaceUriRequest
		if err := c.BindUri(&uriRequest); err != nil {
			return
		}
		var request types.CreateSecretRequest
		if err := c.BindJSON(&request); err != nil {
			return
		}

		secretHandler(func(controller controllers.SecretController, c *gin.Context) (interface{}, error) {
			return controller.CreateSecret(uriRequest.NamespaceName, request)
		})(c)
	}
}

// GetSecrets gets all secrets in a specific namespace.
func GetSecrets() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request types.SecretNamespaceUriRequest
		if err := c.BindUri(&request); err != nil {
			return
		}

		secretHandler(func(controller controllers.SecretController, c *gin.Context) (interface{}, error) {
			return controller.GetSecrets(request.NamespaceName)
		})(c)
	}
}

// GetSecret gets a specific secret from a specific namespace.
func GetSecret() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request types.SecretUriRequest
		if err := c.BindUri(&request); err != nil {
			return
		}

		secretHandler(func(controller controllers.SecretController, c *gin.Context) (interface{}, error) {
			return controller.GetSecret(request.NamespaceName, request.SecretName)
		})(c)
	}
}

// UpdateSecret updates a specific secret in a specific namespace.
func UpdateSecret() gin.HandlerFunc {
	return func(c *gin.Context) {
		var uriRequest types.SecretUriRequest
		if err := c.BindUri(&uriRequest); err != nil {
			return
		}
		var request types.UpdateSecretRequest
		if err := c.BindJSON(&request); err != nil {
			return
		}

		secretHandler(func(controller controllers.SecretController, c *gin.Context) (interface{}, error) {
			return controller.UpdateSecret(uriRequest.NamespaceName, uriRequest.SecretName, request)
		})(c)
	}
}

// DeleteSecret deletes a specific secret in a specific namespace.
func DeleteSecret() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request types.SecretUriRequest
		if err := c.BindUri(&request); err != nil {
			return
		}

		secretHandler(func(controller controllers.SecretController, c *gin.Context) (interface{}, error) {
			return controller.DeleteSecret(request.NamespaceName, request.SecretName)
		})(c)
	}
}
