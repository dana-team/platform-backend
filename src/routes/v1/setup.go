package v1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"

	"github.com/dana-team/platform-backend/src/auth"
	"github.com/dana-team/platform-backend/src/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRoutes initializes the API routes for version 1.
func SetupRoutes(engine *gin.Engine, tokenProvider auth.TokenProvider, scheme *runtime.Scheme) {
	engine.Use(middleware.ErrorHandlingMiddleware())
	v1 := engine.Group("/v1")

	engine.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	setupAuthRoutes(v1, tokenProvider)
	setupNamespaceRoutes(v1, tokenProvider, scheme)
	setupClustersRoutes(v1, tokenProvider, scheme)
}

// setupAuthRoutes defines routes related to authentication.
func setupAuthRoutes(v1 *gin.RouterGroup, tokenProvider auth.TokenProvider) {
	authGroup := v1.Group("/login")
	{
		authGroup.POST("", Login(tokenProvider))
	}
}

// setupNamespaceRoutes defines routes related to namespaces and their resources.
func setupNamespaceRoutes(v1 *gin.RouterGroup, tokenProvider auth.TokenProvider, scheme *runtime.Scheme) {
	namespacesGroup := v1.Group("/namespaces")

	if tokenProvider != nil {
		namespacesGroup.Use(middleware.TokenAuthMiddleware(tokenProvider, scheme))
	}

	{
		getNamespaces := namespacesGroup.Group("")
		getNamespaces.Use(middleware.PaginationMiddleware())
		getNamespaces.GET("", GetNamespaces())

		namespacesGroup.GET("/:namespaceName", GetNamespace())
		namespacesGroup.POST("", CreateNamespace())
		namespacesGroup.DELETE("/:namespaceName", DeleteNamespace())
	}

	secretsGroup := namespacesGroup.Group("/:namespaceName/secrets")
	{
		getSecrets := secretsGroup.Group("")
		getSecrets.Use(middleware.PaginationMiddleware())
		getSecrets.GET("", GetSecrets())

		secretsGroup.POST("", CreateSecret())
		secretsGroup.GET("/:secretName", GetSecret())
		secretsGroup.PUT("/:secretName", UpdateSecret())
		secretsGroup.DELETE("/:secretName", DeleteSecret())
	}

	cappGroup := namespacesGroup.Group("/:namespaceName/capps")
	{
		getCapps := cappGroup.Group("")
		getCapps.Use(middleware.PaginationMiddleware())
		getCapps.GET("", GetCapps())

		cappGroup.POST("", CreateCapp())
		cappGroup.GET("/:cappName", GetCapp())
		cappGroup.PUT("/:cappName", UpdateCapp())
		cappGroup.PUT("/:cappName/state", EditCappState())
		cappGroup.GET("/:cappName/state", GetCappState())
		cappGroup.DELETE("/:cappName", DeleteCapp())

		getDns := cappGroup.Group("")
		getDns.Use(middleware.ClusterMiddleware())
		getDns.GET("/:cappName/dns", GetCappDNS())
	}

	cappRevisionGroup := namespacesGroup.Group("/:namespaceName/capps/:cappName/capprevisions")
	cappRevisionGroup.Use(middleware.ClusterMiddleware())
	{
		getCappRevisions := cappRevisionGroup.Group("")
		getCappRevisions.Use(middleware.PaginationMiddleware())
		getCappRevisions.GET("", GetCappRevisions())

		cappRevisionGroup.GET("/:cappRevisionName", GetCappRevision())
	}

	usersGroup := namespacesGroup.Group("/:namespaceName/users")
	{
		getUsers := usersGroup.Group("")
		getUsers.Use(middleware.PaginationMiddleware())
		getUsers.GET("", GetUsers())

		usersGroup.POST("", CreateUser())
		usersGroup.GET("/:userName", GetUser())
		usersGroup.PUT("/:userName", UpdateUser())
		usersGroup.DELETE("/:userName", DeleteUser())
	}

	configMapGroup := namespacesGroup.Group("/:namespaceName/configmaps")
	{
		configMapGroup.GET("/:configMapName", GetConfigMap())
	}

	containersGroup := namespacesGroup.Group("/:namespaceName/pods/:podName/containers")
	{
		containersGroup.GET("", GetPodsContainers())
	}

	podsGroup := namespacesGroup.Group("/:namespaceName/capps/:cappName/pods")
	podsGroup.Use(middleware.ClusterMiddleware())
	{
		podsGroup.Use(middleware.PaginationMiddleware()).GET("", GetPods())
	}

	logsGroup := namespacesGroup.Group("/:namespaceName")
	{
		logsGroup.GET("/pod/:podName/logs", GetPodLogs())
		logsGroup.GET("/capp/:cappName/logs", GetCappLogs()).Use(middleware.ClusterMiddleware())
	}

	serviceAccountsGroup := namespacesGroup.Group("/:namespaceName/serviceaccounts")
	{
		serviceAccountsGroup.GET("/:serviceAccountName/token", GetToken())
	}
}

// setupClustersRoutes defines routes related to clusters and their namespaces.
func setupClustersRoutes(v1 *gin.RouterGroup, tokenProvider auth.TokenProvider, scheme *runtime.Scheme) {
	clustersGroup := v1.Group("/clusters/:clusterName")
	clustersGroup.Use(middleware.ClusterMiddleware())

	if tokenProvider != nil {
		clustersGroup.Use(middleware.TokenAuthMiddleware(tokenProvider, scheme))
	}

	namespacesGroup := clustersGroup.Group("/namespaces")
	{
		logsGroup := namespacesGroup.Group("/:namespaceName")
		{
			logsGroup.GET("/pod/:podName/logs", GetPodLogs())
			logsGroup.GET("/capp/:cappName/logs", GetCappLogs())
		}

		cappRevisionGroup := namespacesGroup.Group("/:namespaceName/capprevisions")
		{
			getCappRevisions := cappRevisionGroup.Group("")
			getCappRevisions.Use(middleware.PaginationMiddleware())
			getCappRevisions.GET("", GetCappRevisions())

			cappRevisionGroup.GET("/:cappRevisionName", GetCappRevision())
		}

		containersGroup := namespacesGroup.Group("/:namespaceName/pods/:podName/containers")
		{
			containersGroup.GET("", GetPodsContainers())
		}

		podsGroup := namespacesGroup.Group("/:namespaceName/capps/:cappName/pods")
		{
			podsGroup.GET("", GetPods()).Use(middleware.PaginationMiddleware())
		}
	}
}
