package v1

import (
	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	"github.com/dana-team/platform-backend/src/middleware"
	dnsrecordv1alpha1 "github.com/dana-team/provider-dns/apis/record/v1alpha1"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes/scheme"
	"os"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
	runtimeFake "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

var (
	router     *gin.Engine
	fakeClient *fake.Clientset
	dynClient  runtimeClient.WithWatch
	token      string
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func setup() {
	fakeClient = fake.NewSimpleClientset()
	dynClient = runtimeFake.NewClientBuilder().WithScheme(setupScheme()).Build()
	logger, _ := zap.NewProduction()
	router = setupRouter(logger)
}

func setupRouter(logger *zap.Logger) *gin.Engine {
	engine := gin.Default()
	engine.Use(middleware.ErrorHandlingMiddleware())

	engine.Use(func(c *gin.Context) {
		c.Set(middleware.LoggerCtxKey, logger)
		c.Set(middleware.KubeClientCtxKey, fakeClient)
		c.Set(middleware.DynamicClientCtxKey, dynClient)
		c.Set(middleware.TokenCtxKey, token)
		c.Next()
	})
	v1 := engine.Group("/v1")
	{
		namespacesGroup := v1.Group("/namespaces")
		{
			namespacesGroup.Use(middleware.PaginationMiddleware()).GET("", GetNamespaces())
			namespacesGroup.GET("/:namespaceName", GetNamespace())
			namespacesGroup.POST("", CreateNamespace())
			namespacesGroup.DELETE("/:namespaceName", DeleteNamespace())

			secretsGroup := namespacesGroup.Group("/:namespaceName/secrets")
			{
				secretsGroup.POST("", CreateSecret())
				secretsGroup.Use(middleware.PaginationMiddleware()).GET("", GetSecrets())
				secretsGroup.GET("/:secretName", GetSecret())
				secretsGroup.PUT("/:secretName", UpdateSecret())
				secretsGroup.DELETE("/:secretName", DeleteSecret())
			}

			cappGroup := namespacesGroup.Group("/:namespaceName/capps")
			{

				cappGroup.POST("", CreateCapp())
				cappGroup.Use(middleware.PaginationMiddleware()).GET("", GetCapps())
				cappGroup.GET("/:cappName", GetCapp())
				cappGroup.PUT("/:cappName", UpdateCapp())
				cappGroup.PUT("/:cappName/state", EditCappState())
				cappGroup.GET("/:cappName/state", GetCappState())
				cappGroup.GET("/:cappName/dns", GetCappDNS())
				cappGroup.DELETE("/:cappName", DeleteCapp())

			}

			cappRevisionGroup := namespacesGroup.Group("/:namespaceName/capprevisions")
			{

				cappRevisionGroup.Use(middleware.PaginationMiddleware()).GET("", GetCappRevisions())
				cappRevisionGroup.GET("/:cappRevisionName", GetCappRevision())
			}

			usersGroup := namespacesGroup.Group("/:namespaceName/users")
			{
				usersGroup.POST("", CreateUser())
				usersGroup.Use(middleware.PaginationMiddleware()).GET("", GetUsers())
				usersGroup.GET("/:userName", GetUser())
				usersGroup.PUT("/:userName", UpdateUser())
				usersGroup.DELETE("/:userName", DeleteUser())
			}

			logsGroup := v1.Group("/logs")
			{
				logsGroup.GET("/pod/:namespace/:cappName", GetPodLogs())
				logsGroup.GET("/capp/:namespace/:cappName", GetCappLogs())
			}

			configMapGroup := namespacesGroup.Group("/:namespaceName/configmaps")
			{
				configMapGroup.GET("/:configMapName", GetConfigMap())
			}

			containersGroup := namespacesGroup.Group("/:namespaceName")
			{
				containersGroup.GET("/pods/:podName/containers", GetContainers())
			}

			podsGroup := namespacesGroup.Group("/:namespaceName")
			{
				podsGroup.Use(middleware.PaginationMiddleware()).GET("/capps/:cappName/pods", GetPods())
			}

			serviceAccountsGroup := namespacesGroup.Group("/:namespaceName/serviceaccounts")
			{
				serviceAccountsGroup.GET("/:serviceAccountName/token", GetToken())
			}
		}
	}
	return engine
}

func setupScheme() *runtime.Scheme {
	schema := scheme.Scheme
	_ = cappv1alpha1.AddToScheme(schema)
	_ = dnsrecordv1alpha1.AddToScheme(schema)
	return schema
}
