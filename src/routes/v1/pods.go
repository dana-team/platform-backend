package v1

import (
	"github.com/dana-team/platform-backend/src/customerrors"

	"github.com/dana-team/platform-backend/src/controllers"
	"github.com/dana-team/platform-backend/src/routes"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/dana-team/platform-backend/src/utils/pagination"
	"github.com/gin-gonic/gin"
	"net/http"
)

// podHandler wraps a handler function with context setup for PodController.
func podHandler(handler func(controller controllers.PodController, c *gin.Context) (interface{}, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		kubeClient, err := routes.GetKubeClient(c)
		if routes.AddErrorToContext(c, err) {
			return
		}

		logger, err := routes.GetLogger(c)
		if routes.AddErrorToContext(c, err) {
			return
		}

		context := c.Request.Context()
		podController := controllers.NewPodController(kubeClient, context, logger)

		result, err := handler(podController, c)
		if routes.AddErrorToContext(c, err) {
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// GetPods returns a Gin handler function for retrieving pods of a specific capp.
func GetPods() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request types.PodRequestUri
		if err := c.BindUri(&request); err != nil {
			routes.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		limit, page, err := pagination.ExtractPaginationParamsFromCtx(c)
		if err != nil {
			routes.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		podHandler(func(controller controllers.PodController, c *gin.Context) (interface{}, error) {
			return controller.GetPods(request.NamespaceName, request.CappName, limit, page)
		})(c)
	}
}
