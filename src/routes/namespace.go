package routes

import (
	"fmt"
	"net/http"

	"github.com/dana-team/platform-backend/src/controllers"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
)

func ListNamespaces(client *kubernetes.Clientset, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		context := c.Request.Context()
		namespaceController := controllers.New(client, context, logger)
		namespaces, err := namespaceController.GetNamespaces()
		if err != nil {
			c.AbortWithStatusJSON(int(err.(*k8serrors.StatusError).ErrStatus.Code), gin.H{"error": "Failed to list namespaces", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, namespaces)
	}
}

func GetNamespace(client *kubernetes.Clientset, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Validate parameter
		name := c.Param("name")

		context := c.Request.Context()
		namespaceController := controllers.New(client, context, logger)
		namespace, err := namespaceController.GetNamespace(name)
		if err != nil {
			c.AbortWithStatusJSON(int(err.(*k8serrors.StatusError).ErrStatus.Code), gin.H{"error": "Failed to fetch namespace", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, namespace)
	}
}

func CreateNamespace(client *kubernetes.Clientset, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var namespace types.Namespace

		if err := c.BindJSON(&namespace); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		context := c.Request.Context()
		namespaceController := controllers.New(client, context, logger)
		newNamespace, err := namespaceController.CreateNamespace(namespace.Name)
		if err != nil {
			c.AbortWithStatusJSON(int(err.(*k8serrors.StatusError).ErrStatus.Code), gin.H{"error": "Failed to create namespace", "details": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, newNamespace)
	}
}

func DeleteNamespace(client *kubernetes.Clientset, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Validate parameter
		name := c.Param("name")

		context := c.Request.Context()
		namespaceController := controllers.New(client, context, logger)
		if err := namespaceController.DeleteNamespace(name); err != nil {
			c.AbortWithStatusJSON(int(err.(*k8serrors.StatusError).ErrStatus.Code), gin.H{"error": "Failed to delete namespace", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Deleted namespace successfully %s", name)})
	}
}
