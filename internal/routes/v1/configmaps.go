package v1

import (
	"github.com/dana-team/platform-backend/internal/customerrors"
	"github.com/dana-team/platform-backend/internal/middleware"
	"net/http"

	"github.com/dana-team/platform-backend/internal/controllers"
	"github.com/dana-team/platform-backend/internal/types"
	"github.com/gin-gonic/gin"
)

// configMapHandler handles the request of the client to the Kubernetes cluster.
func configMapHandler(handler func(controller controllers.ConfigMapController, c *gin.Context) (interface{}, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		kubeClient, err := middleware.GetKubeClient(c)
		if middleware.AddErrorToContext(c, err) {
			return
		}

		logger, err := middleware.GetLogger(c)
		if middleware.AddErrorToContext(c, err) {
			return
		}

		context := c.Request.Context()
		configMapController := controllers.NewConfigMapController(kubeClient, context, logger)

		result, err := handler(configMapController, c)
		if middleware.AddErrorToContext(c, err) {
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
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		configMapHandler(func(controller controllers.ConfigMapController, c *gin.Context) (interface{}, error) {
			return controller.GetConfigMap(uriRequest.NamespaceName, uriRequest.ConfigMapName)
		})(c)
	}
}
