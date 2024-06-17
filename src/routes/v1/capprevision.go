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

func cappRevisionHandler(handler func(controller controllers.CappRevisionController, c *gin.Context) (interface{}, error)) gin.HandlerFunc {
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

		cappRevisionController, err := controllers.NewCappRevisionController(kubeClient, context, logger)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create controller"})
			return
		}

		result, err := handler(cappRevisionController, c)
		if err != nil {
			c.AbortWithStatusJSON(int(err.(*k8serrors.StatusError).ErrStatus.Code), gin.H{"error": "Operation failed", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func GetCappRevisions() gin.HandlerFunc {
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

		cappRevisionHandler(func(controller controllers.CappRevisionController, c *gin.Context) (interface{}, error) {
			return controller.GetCappRevisions(cappUri.NamespaceName, cappRevisionQuery)
		})(c)
	}
}

func GetCappRevision() gin.HandlerFunc {
	return func(c *gin.Context) {
		var cappUri types.CappRevisionUri
		if err := c.BindUri(&cappUri); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		cappRevisionHandler(func(controller controllers.CappRevisionController, c *gin.Context) (interface{}, error) {
			return controller.GetCappRevision(cappUri.NamespaceName, cappUri.CappRevisionName)
		})(c)
	}
}
