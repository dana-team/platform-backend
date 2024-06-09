package main

import (
	"github.com/dana-team/platform-backend/src/routes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
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

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Fatalf("Error syncing logger: %v", err)
		}
	}()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	nsGroup := r.Group("/namespaces")
	{
		nsGroup.GET("/", routes.ListNamespaces(client, logger))
		nsGroup.GET("/:name", routes.GetNamespace(client, logger))
		nsGroup.POST("/", routes.CreateNamespace(client, logger))
		nsGroup.DELETE("/:name", routes.DeleteNamespace(client, logger))
	}

	if err := r.Run(); err != nil {
		panic(err.Error())
	}
}
