package v1

import (
	"github.com/dana-team/platform-backend/src/controllers"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"net/http"
)

// usersHandler handles the request of the client to the Kubernetes cluster.
func usersHandler(handler func(controller controllers.UserController, c *gin.Context) (interface{}, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		client, exists := c.Get("kubeClient")
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
		kubeClient := client.(kubernetes.Interface)
		context := c.Request.Context()

		userController := controllers.NewUserController(kubeClient, context, logger)
		result, err := handler(userController, c)
		if err != nil {
			c.AbortWithStatusJSON(int(err.(*k8serrors.StatusError).ErrStatus.Code), gin.H{"error": "Operation failed", "details": err.Error()})
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
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}
		var user types.User
		if err := c.BindJSON(&user); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
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
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}
		var userRole types.UpdateUserData
		if err := c.BindJSON(&userRole); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
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
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		usersHandler(func(controller controllers.UserController, c *gin.Context) (interface{}, error) {
			return controller.GetUsers(namespace.NamespaceName)
		})(c)
	}
}

// GetUser fetches a specific user from namespace.
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userIdentifier types.UserIdentifier
		if err := c.BindUri(&userIdentifier); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
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
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		usersHandler(func(controller controllers.UserController, c *gin.Context) (interface{}, error) {
			return controller.DeleteUser(userIdentifier)
		})(c)
	}
}
