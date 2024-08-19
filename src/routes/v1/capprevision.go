package v1

import (
	"github.com/dana-team/platform-backend/src/customerrors"
	"github.com/dana-team/platform-backend/src/routes"
	"github.com/dana-team/platform-backend/src/utils/pagination"
	"net/http"

	"github.com/dana-team/platform-backend/src/controllers"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/gin-gonic/gin"
)

func cappRevisionHandler(handler func(controller controllers.CappRevisionController, c *gin.Context) (interface{}, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		kubeClient, err := routes.GetDynClient(c)
		if routes.AddErrorToContext(c, err) {
			return
		}

		logger, err := routes.GetLogger(c)
		if routes.AddErrorToContext(c, err) {
			return
		}

		context := c.Request.Context()
		cappRevisionController := controllers.NewCappRevisionController(kubeClient, context, logger)

		result, err := handler(cappRevisionController, c)
		if routes.AddErrorToContext(c, err) {
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func GetCappRevisions() gin.HandlerFunc {
	return func(c *gin.Context) {
		var cappRevisionUri types.CappRevisionNamespaceUri
		if err := c.BindUri(&cappRevisionUri); err != nil {
			routes.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		var cappRevisionQuery types.CappRevisionQuery
		if err := c.BindQuery(&cappRevisionQuery); err != nil {
			routes.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		limit, page, err := pagination.ExtractPaginationParamsFromCtx(c)
		if err != nil {
			routes.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		cappRevisionHandler(func(controller controllers.CappRevisionController, c *gin.Context) (interface{}, error) {
			return controller.GetCappRevisions(cappRevisionUri.NamespaceName, limit, page, cappRevisionQuery)
		})(c)
	}
}

func GetCappRevision() gin.HandlerFunc {
	return func(c *gin.Context) {
		var cappUri types.CappRevisionUri
		if err := c.BindUri(&cappUri); err != nil {
			routes.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		cappRevisionHandler(func(controller controllers.CappRevisionController, c *gin.Context) (interface{}, error) {
			return controller.GetCappRevision(cappUri.NamespaceName, cappUri.CappRevisionName)
		})(c)
	}
}
