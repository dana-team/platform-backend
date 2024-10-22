package v1

import (
	"github.com/dana-team/platform-backend/internal/customerrors"
	"github.com/dana-team/platform-backend/internal/middleware"

	"github.com/dana-team/platform-backend/internal/controllers"
	"github.com/dana-team/platform-backend/internal/types"
	"github.com/dana-team/platform-backend/internal/utils/pagination"
	"github.com/gin-gonic/gin"
	"net/http"
)

// usersHandler handles the request of the client to the Kubernetes cluster.
func usersHandler(handler func(controller controllers.UserController, c *gin.Context) (interface{}, error)) gin.HandlerFunc {
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
		userController := controllers.NewUserController(kubeClient, context, logger)

		result, err := handler(userController, c)
		if middleware.AddErrorToContext(c, err) {
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// CreateUser creates a specific user in a specific namespace.
func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var namespace types.NamespaceUri
		if err := c.BindUri(&namespace); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}
		var user types.User
		if err := c.BindJSON(&user); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		usersHandler(func(controller controllers.UserController, c *gin.Context) (interface{}, error) {
			return controller.AddUser(types.UserInput{Namespace: namespace.NamespaceName,
				User: types.User{Name: user.Name, Role: user.Role}})
		})(c)
	}
}

// UpdateUser updates a specific user in a specific namespace.
func UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userIdentifier types.UserIdentifier
		if err := c.BindUri(&userIdentifier); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}
		var userRole types.UpdateUserData
		if err := c.BindJSON(&userRole); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		usersHandler(func(controller controllers.UserController, c *gin.Context) (interface{}, error) {
			return controller.UpdateUser(types.UserInput{Namespace: userIdentifier.NamespaceName,
				User: types.User{Name: userIdentifier.UserName, Role: userRole.Role}})
		})(c)
	}
}

// GetUsers fetches all users from namespace.
func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var namespace types.NamespaceUri
		if err := c.BindUri(&namespace); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		limit, page, err := pagination.ExtractPaginationParamsFromCtx(c)
		if err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		usersHandler(func(controller controllers.UserController, c *gin.Context) (interface{}, error) {
			return controller.GetUsers(namespace.NamespaceName, limit, page)
		})(c)
	}
}

// GetUser fetches a specific user from namespace.
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userIdentifier types.UserIdentifier
		if err := c.BindUri(&userIdentifier); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		usersHandler(func(controller controllers.UserController, c *gin.Context) (interface{}, error) {
			return controller.GetUser(userIdentifier)
		})(c)
	}
}

// DeleteUser deletes a specific user from namespace.
func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userIdentifier types.UserIdentifier
		if err := c.BindUri(&userIdentifier); err != nil {
			middleware.AddErrorToContext(c, customerrors.NewValidationError(err.Error()))
			return
		}

		usersHandler(func(controller controllers.UserController, c *gin.Context) (interface{}, error) {
			return controller.DeleteUser(userIdentifier)
		})(c)
	}
}
