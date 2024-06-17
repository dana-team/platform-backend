package v1

import (
	"net/http"

	"github.com/dana-team/platform-backend/src/controllers"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	client "sigs.k8s.io/controller-runtime/pkg/client"
)

func containerAppRevisionHandler(handler func(controller controllers.ContainerAppRevisionController, c *gin.Context) (interface{}, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		dynClient, exists := c.Get("dynClient")
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
		kubeClient := dynClient.(client.Client)
		context := c.Request.Context()

		containerAppRevisionController, err := controllers.NewContainerAppRevisionController(kubeClient, context, logger)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create controller"})
			return
		}

		result, err := handler(containerAppRevisionController, c)
		if err != nil {
			c.AbortWithStatusJSON(int(err.(*k8serrors.StatusError).ErrStatus.Code), gin.H{"error": "Operation failed", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func GetContainerAppRevisions() gin.HandlerFunc {
	return func(c *gin.Context) {
		var cappUri types.CappRevisionNamespaceUri
		if err := c.BindUri(&cappUri); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}
		var cappRevisionQuery types.CappRevisionQuery
		if err := c.BindUri(&cappRevisionQuery); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		containerAppRevisionHandler(func(controller controllers.ContainerAppRevisionController, c *gin.Context) (interface{}, error) {
			return controller.GetContainerAppRevisions(cappUri.NamespaceName, cappRevisionQuery)
		})(c)
	}
}

func GetContainerAppRevision() gin.HandlerFunc {
	return func(c *gin.Context) {
		var cappUri types.CappRevisionUri
		if err := c.BindUri(&cappUri); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		containerAppRevisionHandler(func(controller controllers.ContainerAppRevisionController, c *gin.Context) (interface{}, error) {
			return controller.GetContainerAppRevision(cappUri.NamespaceName, cappUri.CappRevisionName)
		})(c)
	}
}
