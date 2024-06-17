package v1

import (
	"net/http"

	"github.com/dana-team/platform-backend/src/auth"
	"github.com/dana-team/platform-backend/src/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(engine *gin.Engine, tokenProvider auth.TokenProvider) {
	v1 := engine.Group("/v1")

	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	authGroup := v1.Group("/login")
	{
		authGroup.POST("/", Login(tokenProvider))
	}

	namespacesGroup := v1.Group("/namespaces")
	namespacesGroup.Use(middleware.TokenAuthMiddleware(tokenProvider))
	{
		namespacesGroup.GET("/", ListNamespaces())
		namespacesGroup.GET("/:namespaceName", GetNamespace())
		namespacesGroup.POST("/", CreateNamespace())
		namespacesGroup.DELETE("/:namespaceName", DeleteNamespace())
	}

	secretsGroup := namespacesGroup.Group("/:namespaceName/secrets")
	secretsGroup.Use(middleware.TokenAuthMiddleware(tokenProvider))
	{
		secretsGroup.POST("/", CreateSecret())
		secretsGroup.GET("/", GetSecrets())
		secretsGroup.GET("/:secretName", GetSecret())
		secretsGroup.PATCH("/:secretName", PatchSecret())
		secretsGroup.DELETE("/:secretName", DeleteSecret())
	}

	containerAppGroup := namespacesGroup.Group("/:namespaceName/capps")
	containerAppGroup.Use(middleware.TokenAuthMiddleware(tokenProvider))
	{
		containerAppGroup.POST("/", CreateContainerApp())
		containerAppGroup.GET("/", GetContainerApps())
		containerAppGroup.GET("/:cappName", GetContainerApp())
		containerAppGroup.PATCH("/:cappName", PatchContainerApp())
		containerAppGroup.DELETE("/:cappName", DeleteContainerApp())
	}

	containerAppRevisionGroup := namespacesGroup.Group("/:namespaceName/capprevisions")
	containerAppRevisionGroup.Use(middleware.TokenAuthMiddleware(tokenProvider))
	{
		containerAppRevisionGroup.GET("/", GetContainerAppRevisions())
		containerAppRevisionGroup.GET("/:cappRevisionName", GetContainerAppRevision())
	}
}
