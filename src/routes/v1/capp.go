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

func cappHandler(handler func(controller controllers.CappController, c *gin.Context) (interface{}, error)) gin.HandlerFunc {
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

		cappController, err := controllers.NewCappController(kubeClient, context, logger)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create controller"})
			return
		}

		result, err := handler(cappController, c)
		if err != nil {
			c.AbortWithStatusJSON(int(err.(*k8serrors.StatusError).ErrStatus.Code), gin.H{"error": "Operation failed", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func GetCapps() gin.HandlerFunc {
	return func(c *gin.Context) {
		var cappUri types.CappNamespaceUri
		if err := c.BindUri(&cappUri); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}
		var cappQuery types.CappQuery
		if err := c.BindUri(&cappQuery); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		cappHandler(func(controller controllers.CappController, c *gin.Context) (interface{}, error) {
			return controller.GetCapps(cappUri.NamespaceName, cappQuery)
		})(c)
	}
}

func GetCapp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var cappUri types.CappUri
		if err := c.BindUri(&cappUri); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		cappHandler(func(controller controllers.CappController, c *gin.Context) (interface{}, error) {
			return controller.GetCapp(cappUri.NamespaceName, cappUri.CappName)
		})(c)
	}
}

func CreateCapp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var cappUri types.CappNamespaceUri
		if err := c.BindUri(&cappUri); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}
		var capp types.Capp
		if err := c.BindJSON(&capp); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		cappHandler(func(controller controllers.CappController, c *gin.Context) (interface{}, error) {
			return controller.CreateCapp(cappUri.NamespaceName, capp)
		})(c)
	}
}

func PatchCapp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var cappUri types.CappUri
		if err := c.BindUri(&cappUri); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}
		var capp types.Capp
		if err := c.BindJSON(&capp); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		cappHandler(func(controller controllers.CappController, c *gin.Context) (interface{}, error) {
			return controller.PatchCapp(cappUri.NamespaceName, cappUri.CappName, capp)
		})(c)
	}
}

func DeleteCapp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var cappUri types.CappUri
		if err := c.BindUri(&cappUri); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		cappHandler(func(controller controllers.CappController, c *gin.Context) (interface{}, error) {
			return controller.DeleteCapp(cappUri.NamespaceName, cappUri.CappName)
		})(c)
	}
}
