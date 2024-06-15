package v1

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/dana-team/platform-backend/src/controllers"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/gin-gonic/gin"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
)

func namespaceHandler(handler func(controller controllers.NamespaceController, c *gin.Context) (interface{}, error)) gin.HandlerFunc {
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

		namespaceController := controllers.NewNamespaceController(kubeClient, context, logger)
		result, err := handler(namespaceController, c)
		if err != nil {
			c.AbortWithStatusJSON(int(err.(*k8serrors.StatusError).ErrStatus.Code), gin.H{"error": "Operation failed", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func ListNamespaces() gin.HandlerFunc {
	return namespaceHandler(func(controller controllers.NamespaceController, c *gin.Context) (interface{}, error) {
		return controller.GetNamespaces()
	})
}

func GetNamespace() gin.HandlerFunc {
	return func(c *gin.Context) {
		var namespaceUri types.NamespaceUri
		if err := c.BindUri(&namespaceUri); err != nil {
			return
		}

		namespaceHandler(func(controller controllers.NamespaceController, c *gin.Context) (interface{}, error) {
			return controller.GetNamespace(namespaceUri.NamespaceName)
		})(c)
	}
}

func CreateNamespace() gin.HandlerFunc {
	return func(c *gin.Context) {
		var namespace types.Namespace
		if err := c.BindJSON(&namespace); err != nil {
			return
		}

		namespaceHandler(func(controller controllers.NamespaceController, c *gin.Context) (interface{}, error) {
			return controller.CreateNamespace(namespace.Name)
		})(c)
	}
}

func DeleteNamespace() gin.HandlerFunc {
	return func(c *gin.Context) {
		var namespaceUri types.NamespaceUri
		if err := c.BindUri(&namespaceUri); err != nil {
			return
		}

		namespaceHandler(func(controller controllers.NamespaceController, c *gin.Context) (interface{}, error) {
			name := namespaceUri.NamespaceName
			message := fmt.Sprintf("Deleted namespace successfully %s", name)
			return gin.H{"message": message}, controller.DeleteNamespace(name)
		})(c)
	}
}
