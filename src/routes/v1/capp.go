package v1

import (
	"github.com/dana-team/platform-backend/src/customerrors"
	"github.com/dana-team/platform-backend/src/middleware"
	"github.com/dana-team/platform-backend/src/routes"
	"github.com/dana-team/platform-backend/src/utils/pagination"
	"net/http"

	"github.com/dana-team/platform-backend/src/controllers"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/gin-gonic/gin"
)

func cappHandler(handler func(controller controllers.CappController, c *gin.Context) (interface{}, error)) gin.HandlerFunc {
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
		cappController := controllers.NewCappController(kubeClient, context, logger)

		result, err := handler(cappController, c)
		if middleware.AddErrorToContext(c, err) {
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func GetCapps() gin.HandlerFunc {
	return func(c *gin.Context) {
		var cappUri types.CappNamespaceUri
		if err := c.BindUri(&cappUri); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		var cappQuery types.CappQuery
		if err := c.BindQuery(&cappQuery); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		limit, page, err := pagination.ExtractPaginationParamsFromCtx(c)
		if err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		cappHandler(func(controller controllers.CappController, c *gin.Context) (interface{}, error) {
			return controller.GetCapps(cappUri.NamespaceName, limit, page, cappQuery)
		})(c)
	}
}

func GetCapp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var cappUri types.CappUri
		if err := c.BindUri(&cappUri); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
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
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}
		var capp types.CreateCapp
		if err := c.BindJSON(&capp); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		cappHandler(func(controller controllers.CappController, c *gin.Context) (interface{}, error) {
			return controller.CreateCapp(cappUri.NamespaceName, capp)
		})(c)
	}
}

func UpdateCapp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var cappUri types.CappUri
		if err := c.BindUri(&cappUri); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}
		var capp types.UpdateCapp
		if err := c.BindJSON(&capp); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		cappHandler(func(controller controllers.CappController, c *gin.Context) (interface{}, error) {
			return controller.UpdateCapp(cappUri.NamespaceName, cappUri.CappName, capp)
		})(c)
	}
}

func EditCappState() gin.HandlerFunc {
	return func(c *gin.Context) {
		var cappUri types.CappUri
		if err := c.BindUri(&cappUri); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}
		var state types.CappState
		if err := c.BindJSON(&state); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		cappHandler(func(controller controllers.CappController, c *gin.Context) (interface{}, error) {
			return controller.EditCappState(cappUri.NamespaceName, cappUri.CappName, state.State)
		})(c)
	}
}

func GetCappState() gin.HandlerFunc {
	return func(c *gin.Context) {
		var cappUri types.CappUri
		if err := c.BindUri(&cappUri); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		cappHandler(func(controller controllers.CappController, c *gin.Context) (interface{}, error) {
			return controller.GetCappState(cappUri.NamespaceName, cappUri.CappName)
		})(c)
	}
}

func GetCappDNS() gin.HandlerFunc {
	return func(c *gin.Context) {
		var cappUri types.CappUri
		if err := c.BindUri(&cappUri); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		cappHandler(func(controller controllers.CappController, c *gin.Context) (interface{}, error) {
			return controller.GetCappDNS(cappUri.NamespaceName, cappUri.CappName)
		})(c)
	}
}

func DeleteCapp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var cappUri types.CappUri
		if err := c.BindUri(&cappUri); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		cappHandler(func(controller controllers.CappController, c *gin.Context) (interface{}, error) {
			return controller.DeleteCapp(cappUri.NamespaceName, cappUri.CappName)
		})(c)
	}
}
