package v1

import (
	"testing"

	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes/scheme"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
	runtimeFake "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var (
	router     *gin.Engine
	fakeClient *fake.Clientset
	dynClient  runtimeClient.WithWatch
)

func TestMain(m *testing.M) {
	m.Run()
}

func setup() {
	fakeClient = fake.NewSimpleClientset()
	dynClient = runtimeFake.NewClientBuilder().WithScheme(setupScheme()).Build()
	logger, _ := zap.NewProduction()
	router = setupRouter(logger)
}

func setupRouter(logger *zap.Logger) *gin.Engine {
	engine := gin.Default()
	engine.Use(func(c *gin.Context) {
		c.Set("logger", logger)
		c.Set("kubeClient", fakeClient)
		c.Set("dynClient", dynClient)
		c.Next()
	})
	v1 := engine.Group("/v1")
	{
		namespacesGroup := v1.Group("/namespaces")
		{
			namespacesGroup.GET("", GetNamespaces())
			namespacesGroup.GET("/:namespaceName", GetNamespace())
			namespacesGroup.POST("", CreateNamespace())
			namespacesGroup.DELETE("/:namespaceName", DeleteNamespace())

			secretsGroup := namespacesGroup.Group("/:namespaceName/secrets")
			{
				secretsGroup.POST("", CreateSecret())
				secretsGroup.GET("", GetSecrets())
				secretsGroup.GET("/:secretName", GetSecret())
				secretsGroup.PUT("/:secretName", UpdateSecret())
				secretsGroup.DELETE("/:secretName", DeleteSecret())
			}

			cappGroup := namespacesGroup.Group("/:namespaceName/capps")
			{

				cappGroup.POST("", CreateCapp())
				cappGroup.GET("", GetCapps())
				cappGroup.GET("/:cappName", GetCapp())
				cappGroup.PUT("/:cappName", UpdateCapp())
				cappGroup.PUT("/:cappName/state", EditCappState())
				cappGroup.GET("/:cappName/state", GetCappState())
				cappGroup.DELETE("/:cappName", DeleteCapp())

			}

			cappRevisionGroup := namespacesGroup.Group("/:namespaceName/capprevisions")
			{

				cappRevisionGroup.GET("", GetCappRevisions())
				cappRevisionGroup.GET("/:cappRevisionName", GetCappRevision())
			}

			usersGroup := namespacesGroup.Group("/:namespaceName/users")
			{
				usersGroup.POST("", CreateUser())
				usersGroup.GET("", GetUsers())
				usersGroup.GET("/:userName", GetUser())
				usersGroup.PUT("/:userName", UpdateUser())
				usersGroup.DELETE("/:userName", DeleteUser())
			}

			configMapGroup := namespacesGroup.Group("/:namespaceName/configmaps")
			{
				configMapGroup.GET("/:configMapName", GetConfigMap())
			}
		}
	}
	return engine
}

func setupScheme() *runtime.Scheme {
	schema := scheme.Scheme
	_ = cappv1alpha1.AddToScheme(schema)
	return schema
}
