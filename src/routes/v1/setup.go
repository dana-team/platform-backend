package v1

import (
	"net/http"

	"github.com/dana-team/platform-backend/src/auth"
	"github.com/dana-team/platform-backend/src/middleware"
	"github.com/dana-team/platform-backend/src/utils/terminal_utils"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(engine *gin.Engine, tokenProvider auth.TokenProvider) {
	v1 := engine.Group("/v1")

	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	engine.GET("/ws/terminal", func(c *gin.Context) {
		handler := terminal_utils.CreateAttachHandler()
		handler.ServeHTTP(c.Writer, c.Request)
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
		secretsGroup.PUT("/:secretName", UpdateSecret())
		secretsGroup.DELETE("/:secretName", DeleteSecret())
	}

	cappGroup := namespacesGroup.Group("/:namespaceName/capps")
	cappGroup.Use(middleware.TokenAuthMiddleware(tokenProvider))
	{
		cappGroup.POST("/", CreateCapp())
		cappGroup.GET("/", GetCapps())
		cappGroup.GET("/:cappName", GetCapp())
		cappGroup.PUT("/:cappName", UpdateCapp())
		cappGroup.DELETE("/:cappName", DeleteCapp())
		cappGroup.POST("/:cappName/pods/:podName/containers/:containerName/terminal", StartTerminal())
	}

	cappRevisionGroup := namespacesGroup.Group("/:namespaceName/capprevisions")
	cappRevisionGroup.Use(middleware.TokenAuthMiddleware(tokenProvider))
	{
		cappRevisionGroup.GET("/", GetCappRevisions())
		cappRevisionGroup.GET("/:cappRevisionName", GetCappRevision())
	}
}
