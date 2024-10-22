package v1

import (
	"github.com/dana-team/platform-backend/internal/customerrors"
	"github.com/dana-team/platform-backend/internal/middleware"
	"github.com/dana-team/platform-backend/internal/routes"
	"github.com/dana-team/platform-backend/internal/utils/pagination"
	"net/http"

	"github.com/dana-team/platform-backend/internal/controllers"
	"github.com/dana-team/platform-backend/internal/types"
	"github.com/gin-gonic/gin"
)

func cappRevisionHandler(handler func(controller controllers.CappRevisionController, c *gin.Context) (interface{}, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		kubeClient, err := middleware.GetDynClient(c)
		if middleware.AddErrorToContext(c, err) {
			return
		}

		logger, err := middleware.GetLogger(c)
		if middleware.AddErrorToContext(c, err) {
			return
		}

		context := routes.GetContext(c)
		cappRevisionController := controllers.NewCappRevisionController(kubeClient, context, logger)

		result, err := handler(cappRevisionController, c)
		if middleware.AddErrorToContext(c, err) {
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func GetCappRevisions() gin.HandlerFunc {
	return func(c *gin.Context) {
		var cappRevisionUri types.CappRevisionNamespaceUri
		if err := c.BindUri(&cappRevisionUri); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		limit, page, err := pagination.ExtractPaginationParamsFromCtx(c)
		if err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		cappRevisionHandler(func(controller controllers.CappRevisionController, c *gin.Context) (interface{}, error) {
			return controller.GetCappRevisions(cappRevisionUri.NamespaceName, limit, page, cappRevisionUri.CappName)
		})(c)
	}
}

func GetCappRevision() gin.HandlerFunc {
	return func(c *gin.Context) {
		var cappUri types.CappRevisionUri
		if err := c.BindUri(&cappUri); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		cappRevisionHandler(func(controller controllers.CappRevisionController, c *gin.Context) (interface{}, error) {
			return controller.GetCappRevision(cappUri.NamespaceName, cappUri.CappRevisionName)
		})(c)
	}
}
