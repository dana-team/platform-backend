package routes

import (
	"fmt"
	"github.com/dana-team/platform-backend/src/controllers"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"net/http"
)

func ListNamespaces(client *kubernetes.Clientset, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		namespaceController := controllers.NewNSController(client, logger)
		namespaces, err := namespaceController.GetNamespaces()
		if err != nil {
			c.AbortWithStatusJSON(int(err.(*k8serrors.StatusError).ErrStatus.Code), gin.H{"error": "failed to list namespaces", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, namespaces)
	}
}

func GetNamespace(client *kubernetes.Clientset, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")

		namespaceController := controllers.NewNSController(client, logger)
		namespace, err := namespaceController.GetNamespace(name)
		if err != nil {
			c.AbortWithStatusJSON(int(err.(*k8serrors.StatusError).ErrStatus.Code), gin.H{"error": "failed to fetch namespace", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, namespace)
	}
}

func CreateNamespace(client *kubernetes.Clientset, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var inputNS types.Namespace
		if err := c.BindJSON(&inputNS); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
			return
		}
		namespaceController := controllers.NewNSController(client, logger)
		createdNamespace, err := namespaceController.CreateNamespace(inputNS.Name)
		if err != nil {
			c.AbortWithStatusJSON(int(err.(*k8serrors.StatusError).ErrStatus.Code), gin.H{"error": "failed to create namespace", "details": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, createdNamespace)
	}
}

func DeleteNamespace(client *kubernetes.Clientset, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		namespaceController := controllers.NewNSController(client, logger)
		if err := namespaceController.DeleteNamespace(name); err != nil {
			c.AbortWithStatusJSON(int(err.(*k8serrors.StatusError).ErrStatus.Code), gin.H{"error": "failed to delete namespace", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("deleted namespace successfully %s", name)})
	}
}
