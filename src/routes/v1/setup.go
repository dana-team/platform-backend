package v1

import (
	"github.com/dana-team/platform-backend/src/auth"
	"github.com/dana-team/platform-backend/src/middleware"
	"github.com/dana-team/platform-backend/src/routes/v1/operation"
	"github.com/dana-team/platform-backend/src/routes/v1/registry"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
)

func setupAPIRegistry(engine *gin.Engine) (huma.API, huma.Registry) {
	config := huma.DefaultConfig("platform-backend", "0.1.0")
	config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"bearer": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "jwt",
		},
		"basic": {
			Type:         "http",
			Scheme:       "basic",
			BearerFormat: "Basic Auth",
		},
	}

	registry := huma.NewMapRegistry("#/components/schemas/", registry.SchemaNamer)
	config.OpenAPI.Components.Schemas = registry
	api := humagin.New(engine, config)

	return api, registry
}

// SetupRoutes initializes the API routes for version 1.
func SetupRoutes(engine *gin.Engine, tokenProvider auth.TokenProvider, scheme *runtime.Scheme) {
	engine.Use(middleware.ErrorHandlingMiddleware())
	v1 := engine.Group("/v1")

	api, r := setupAPIRegistry(engine)

	engine.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	operation.AddHealthz(api)

	setupAuthRoutes(api, r, v1, tokenProvider)
	setupNamespaceRoutes(api, r, v1, tokenProvider, scheme)
	setupClustersRoutes(api, r, v1, tokenProvider, scheme)
}

// setupAuthRoutes defines routes related to authentication.
func setupAuthRoutes(api huma.API, r huma.Registry, v1 *gin.RouterGroup, tokenProvider auth.TokenProvider) {
	authGroup := v1.Group("/login")
	{
		authGroup.POST("", Login(tokenProvider))
		operation.AddLogin(api, r)
	}
}

// setupNamespaceRoutes defines routes related to namespaces and their resources.
func setupNamespaceRoutes(api huma.API, r huma.Registry, v1 *gin.RouterGroup, tokenProvider auth.TokenProvider, scheme *runtime.Scheme) {
	namespacesGroup := v1.Group("/namespaces")

	if tokenProvider != nil {
		namespacesGroup.Use(middleware.TokenAuthMiddleware(tokenProvider, scheme))
	}

	{
		getNamespaces := namespacesGroup.Group("")
		getNamespaces.Use(middleware.PaginationMiddleware())

		getNamespaces.GET("", GetNamespaces())
		operation.AddGetNamespaces(api, r)

		namespacesGroup.GET("/:namespaceName", GetNamespace())
		operation.AddGetNamespace(api, r)

		namespacesGroup.POST("", CreateNamespace())
		operation.AddCreateNamespace(api, r)

		namespacesGroup.DELETE("/:namespaceName", DeleteNamespace())
		operation.AddDeleteNamespace(api, r)
	}

	secretsGroup := namespacesGroup.Group("/:namespaceName/secrets")
	{
		getSecrets := secretsGroup.Group("")
		getSecrets.Use(middleware.PaginationMiddleware())
		getSecrets.GET("", GetSecrets())

		secretsGroup.POST("", CreateSecret())
		operation.AddCreateSecret(api, r)

		secretsGroup.GET("/:secretName", GetSecret())
		operation.AddGetSecret(api, r)

		secretsGroup.PUT("/:secretName", UpdateSecret())
		operation.AddUpdateSecret(api, r)

		secretsGroup.DELETE("/:secretName", DeleteSecret())
		operation.AddDeleteSecret(api, r)
	}

	cappGroup := namespacesGroup.Group("/:namespaceName/capps")
	{
		getCapps := cappGroup.Group("")
		getCapps.Use(middleware.PaginationMiddleware())

		getCapps.GET("", GetCapps())
		operation.AddGetCapps(api, r)

		cappGroup.POST("", CreateCapp())
		operation.AddCreateCapp(api, r)

		cappGroup.GET("/:cappName", GetCapp())
		operation.AddGetCapp(api, r)

		cappGroup.PUT("/:cappName", UpdateCapp())
		operation.AddUpdateCapp(api, r)

		cappGroup.PUT("/:cappName/state", EditCappState())
		operation.AddEditCappState(api, r)

		cappGroup.GET("/:cappName/state", GetCappState())
		operation.AddGetCappState(api, r)

		cappGroup.DELETE("/:cappName", DeleteCapp())
		operation.AddDeleteCapp(api, r)

		getDns := cappGroup.Group("")
		getDns.Use(middleware.ClusterMiddleware())
		getDns.GET("/:cappName/dns", GetCappDNS())
		operation.AddGetCappDNS(api, r)
	}

	cappRevisionGroup := namespacesGroup.Group("/:namespaceName/capps/:cappName/capprevisions")
	cappRevisionGroup.Use(middleware.ClusterMiddleware())
	{
		getCappRevisions := cappRevisionGroup.Group("")
		getCappRevisions.Use(middleware.PaginationMiddleware())

		getCappRevisions.GET("", GetCappRevisions())
		operation.AddGetCappRevisions(api, r)

		cappRevisionGroup.GET("/:cappRevisionName", GetCappRevision())
		operation.AddGetCappRevision(api, r)
	}

	usersGroup := namespacesGroup.Group("/:namespaceName/users")
	{
		getUsers := usersGroup.Group("")
		getUsers.Use(middleware.PaginationMiddleware())

		getUsers.GET("", GetUsers())
		operation.AddGetUsers(api, r)

		usersGroup.POST("", CreateUser())
		operation.AddCreateUser(api, r)

		usersGroup.GET("/:userName", GetUser())
		operation.AddGetUser(api, r)

		usersGroup.PUT("/:userName", UpdateUser())
		operation.AddUpdateUser(api, r)

		usersGroup.DELETE("/:userName", DeleteUser())
		operation.AddDeleteUser(api, r)

	}

	configMapGroup := namespacesGroup.Group("/:namespaceName/configmaps")
	{
		configMapGroup.GET("/:configMapName", GetConfigMap())
		operation.AddGetConfigMap(api, r)
	}

	containersGroup := namespacesGroup.Group("/:namespaceName/pods/:podName/containers")
	{
		containersGroup.GET("", GetPodsContainers())
		operation.AddGetContainers(api, r)
	}

	podsGroup := namespacesGroup.Group("/:namespaceName/capps/:cappName/pods")
	podsGroup.Use(middleware.ClusterMiddleware())
	{
		podsGroup.Use(middleware.PaginationMiddleware()).GET("", GetPods())
		operation.AddGetPods(api, r)

	}

	logsGroup := namespacesGroup.Group("/:namespaceName")
	{
		logsGroup.GET("/pod/:podName/logs", GetPodLogs())
		// TODO: fix
		operation.AddGetPodLogs(api, r)

		logsGroup.Use(middleware.ClusterMiddleware()).GET("/capp/:cappName/logs", GetCappLogs())
	}

	serviceAccountsGroup := namespacesGroup.Group("/:namespaceName/serviceaccounts")
	{
		serviceAccountsGroup.GET("/:serviceAccountName/token", GetToken())
		operation.AddGetToken(api, r)
	}
}

// setupClustersRoutes defines routes related to clusters and their namespaces.
func setupClustersRoutes(api huma.API, registry huma.Registry, v1 *gin.RouterGroup, tokenProvider auth.TokenProvider, scheme *runtime.Scheme) {
	clustersGroup := v1.Group("/clusters/:clusterName")

	if tokenProvider != nil {
		clustersGroup.Use(middleware.TokenAuthMiddleware(tokenProvider, scheme))
	}

	clustersGroup.Use(middleware.ClusterMiddleware())

	namespacesGroup := clustersGroup.Group("/namespaces")
	{
		logsGroup := namespacesGroup.Group("/:namespaceName")
		{
			// TODO
			logsGroup.GET("/pod/:podName/logs", GetPodLogs())
			// TODO
			logsGroup.GET("/capp/:cappName/logs", GetCappLogs())
		}

		cappRevisionGroup := namespacesGroup.Group("/:namespaceName/capprevisions")
		{
			getCappRevisions := cappRevisionGroup.Group("")
			getCappRevisions.Use(middleware.PaginationMiddleware())

			getCappRevisions.GET("", GetCappRevisions())
			operation.AddClusterGetCappRevisions(api, registry)

			cappRevisionGroup.GET("/:cappRevisionName", GetCappRevision())
			operation.AddClusterGetCappRevision(api, registry)
		}

		containersGroup := namespacesGroup.Group("/:namespaceName/pods/:podName/containers")
		{
			containersGroup.GET("", GetPodsContainers())
			operation.AddClusterGetContainers(api, registry)
		}

		podsGroup := namespacesGroup.Group("/:namespaceName/capps/:cappName/pods")
		{
			podsGroup.Use(middleware.PaginationMiddleware())
			podsGroup.GET("", GetPods())
			operation.AddClusterGetPods(api, registry)
		}
	}
}
