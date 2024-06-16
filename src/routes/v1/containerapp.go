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

func containerAppHandler(handler func(controller controllers.ContainerAppController, c *gin.Context) (interface{}, error)) gin.HandlerFunc {
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

		containerAppController, err := controllers.NewContainerAppController(kubeClient, context, logger)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create controller"})
			return
		}

		result, err := handler(containerAppController, c)
		if err != nil {
			c.AbortWithStatusJSON(int(err.(*k8serrors.StatusError).ErrStatus.Code), gin.H{"error": "Operation failed", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func GetContainerApps() gin.HandlerFunc {
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

		containerAppHandler(func(controller controllers.ContainerAppController, c *gin.Context) (interface{}, error) {
			return controller.GetContainerApps(cappUri.NamespaceName, cappQuery)
		})(c)
	}
}

func GetContainerApp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var cappUri types.CappUri
		if err := c.BindUri(&cappUri); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		containerAppHandler(func(controller controllers.ContainerAppController, c *gin.Context) (interface{}, error) {
			return controller.GetContainerApp(cappUri.NamespaceName, cappUri.CappName)
		})(c)
	}
}

func CreateContainerApp() gin.HandlerFunc {
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

		containerAppHandler(func(controller controllers.ContainerAppController, c *gin.Context) (interface{}, error) {
			return controller.CreateContainerApp(cappUri.NamespaceName, capp)
		})(c)
	}
}

func PatchContainerApp() gin.HandlerFunc {
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

		containerAppHandler(func(controller controllers.ContainerAppController, c *gin.Context) (interface{}, error) {
			return controller.PatchContainerApp(cappUri.NamespaceName, cappUri.CappName, capp)
		})(c)
	}
}

func DeleteContainerApp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var cappUri types.CappUri
		if err := c.BindUri(&cappUri); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		containerAppHandler(func(controller controllers.ContainerAppController, c *gin.Context) (interface{}, error) {
			return controller.DeleteContainerApp(cappUri.NamespaceName, cappUri.CappName)
		})(c)
	}
}
