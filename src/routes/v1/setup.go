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
		namespacesGroup.GET("/:name", GetNamespace())
		namespacesGroup.POST("/", CreateNamespace())
		namespacesGroup.DELETE("/:name", DeleteNamespace())
	}

	secretsGroup := namespacesGroup.Group("/:namespace/secrets")
	secretsGroup.Use(middleware.TokenAuthMiddleware(tokenProvider))
	{
		secretsGroup.POST("/", CreateSecret())
		secretsGroup.GET("/", GetSecrets())
		secretsGroup.GET("/:name", GetSecret())
		secretsGroup.PATCH("/:name", PatchSecret())
		secretsGroup.DELETE("/:name", DeleteSecret())
	}
}
