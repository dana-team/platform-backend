package main

import (
	"fmt"
	"github.com/dana-team/platform-backend/src/models"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"net/http"
)

func main() {
	// Set up Kubernetes client
	config, err := clientcmd.BuildConfigFromFlags("", "/home/user/.kube/config")
	if err != nil {
		panic(err.Error())
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	nsGroup := r.Group("/namespaces")
	{
		nsGroup.GET("/", listNamespaces(client))
		nsGroup.GET("/:name", getNamespace(client))
		nsGroup.POST("/", createNamespace(client))
		nsGroup.DELETE("/:name", deleteNamespace(client))
	}

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func listNamespaces(client *kubernetes.Clientset) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespaces, err := client.CoreV1().Namespaces().List(c, metav1.ListOptions{})
		if err != nil {
			c.AbortWithStatusJSON(int(err.(*k8serrors.StatusError).ErrStatus.Code), gin.H{"error": "failed to list namespaces", "details": err.Error()})
			return
		}
		outputNamespaces := models.OutputNamespaces{}

		for _, namespace := range namespaces.Items {
			outputNamespaces.Namespaces = append(outputNamespaces.Namespaces, models.Namespace{Name: namespace.Name})
		}
		outputNamespaces.Count = len(outputNamespaces.Namespaces)

		c.JSON(http.StatusOK, outputNamespaces)
	}
}

func getNamespace(client *kubernetes.Clientset) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")

		outputNS := models.Namespace{}
		namespace, err := client.CoreV1().Namespaces().Get(c, name, metav1.GetOptions{})
		if err != nil {
			c.AbortWithStatusJSON(int(err.(*k8serrors.StatusError).ErrStatus.Code), gin.H{"error": "failed to get namespace", "details": err.Error()})
			return
		}
		outputNS.Name = namespace.Name
		c.JSON(http.StatusOK, outputNS)
	}
}

func createNamespace(client *kubernetes.Clientset) gin.HandlerFunc {
	return func(c *gin.Context) {
		var inputNS models.Namespace
		if err := c.BindJSON(&inputNS); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
			return
		}
		namespace := v1.Namespace{}
		namespace.Name = inputNS.Name
		createdNamespace, err := client.CoreV1().Namespaces().Create(c, &namespace, metav1.CreateOptions{})
		if err != nil {
			c.AbortWithStatusJSON(int(err.(*k8serrors.StatusError).ErrStatus.Code), gin.H{"error": "failed to create namespace", "details": err.Error()})
			return
		}
		outputNamespace := models.Namespace{Name: createdNamespace.Name}
		c.JSON(http.StatusCreated, outputNamespace)
	}
}

func deleteNamespace(client *kubernetes.Clientset) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		if err := client.CoreV1().Namespaces().Delete(c, name, metav1.DeleteOptions{}); err != nil {
			c.AbortWithStatusJSON(int(err.(*k8serrors.StatusError).ErrStatus.Code), gin.H{"error": "failed to delete namespace", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("deleted namespace successfully %s", name)})
	}
}
