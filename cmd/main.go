package main

import (
	"log"
	"net/http"

	"github.com/dana-team/platform-backend/src/routes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// TODO: Refactor main into smaller functions
func main() {
	// TODO: Use environment variable to get path
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
		log.Fatalf("Can't initialize zap logger: %v", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Fatalf("Error syncing logger: %v", err)
		}
	}()

	engine := gin.Default()
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	namespacesGroup := engine.Group("/namespaces")
	{
		namespacesGroup.GET("/", routes.ListNamespaces(client, logger))
		namespacesGroup.GET("/:name", routes.GetNamespace(client, logger))
		namespacesGroup.POST("/", routes.CreateNamespace(client, logger))
		namespacesGroup.DELETE("/:name", routes.DeleteNamespace(client, logger))
	}

	if err := engine.Run(); err != nil {
		panic(err.Error())
	}
}
