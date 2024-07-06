package v1

import (
	"context"
	cappv1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	"github.com/dana-team/platform-backend/src/routes/mocks"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	fakeclient "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	runtimeFake "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

var (
	router    *gin.Engine
	clientset *fakeclient.Clientset
	dynClient client.Client
)

const (
	testName          = "test"
	testNamespace     = testName + "-ns"
	nonExistentSuffix = "-non-existent"
)

const (
	labelSelectorKey = "labelSelector"
	labelKey         = "key"
	labelValue       = "value"
)

const (
	operationFailed = "Operation failed"
	invalidRequest  = "Invalid request"
	detailsKey      = "details"
	errorKey        = "error"
	messageKey      = "message"
)

const (
	metadata    = "metadata"
	labels      = "labels"
	annotations = "annotations"
	spec        = "spec"
	status      = "status"
	count       = "count"
	data        = "data"
	nameKey     = "name"
)

const (
	contentType     = "Content-Type"
	applicationJson = "application/json"
)

func TestMain(m *testing.M) {
	m.Run()
}

func setup() {
	clientset = fakeclient.NewSimpleClientset()
	dynClient = runtimeFake.NewClientBuilder().WithScheme(setupScheme()).Build()
	logger, _ := zap.NewProduction()
	router = setupRouter(logger)
}

func setupRouter(logger *zap.Logger) *gin.Engine {
	engine := gin.Default()
	engine.Use(func(c *gin.Context) {
		c.Set("logger", logger)
		c.Set("kubeClient", clientset)
		c.Set("dynClient", dynClient)
		c.Next()
	})
	v1 := engine.Group("/v1")
	{
		namespacesGroup := v1.Group("/namespaces")
		{
			namespacesGroup.GET("/", GetNamespaces())
			namespacesGroup.GET("/:namespaceName", GetNamespace())
			namespacesGroup.POST("/", CreateNamespace())
			namespacesGroup.DELETE("/:namespaceName", DeleteNamespace())

			secretsGroup := namespacesGroup.Group("/:namespaceName/secrets")
			{
				secretsGroup.POST("/", CreateSecret())
				secretsGroup.GET("/", GetSecrets())
				secretsGroup.GET("/:secretName", GetSecret())
				secretsGroup.PUT("/:secretName", UpdateSecret())
				secretsGroup.DELETE("/:secretName", DeleteSecret())
			}

			cappGroup := namespacesGroup.Group("/:namespaceName/capps")
			{
				cappGroup.POST("/", CreateCapp())
				cappGroup.GET("/", GetCapps())
				cappGroup.GET("/:cappName", GetCapp())
				cappGroup.PUT("/:cappName", UpdateCapp())
				cappGroup.DELETE("/:cappName", DeleteCapp())
			}

			cappRevisionGroup := namespacesGroup.Group("/:namespaceName/capprevisions")
			{
				cappRevisionGroup.GET("/", GetCappRevisions())
				cappRevisionGroup.GET("/:cappRevisionName", GetCappRevision())
			}

			usersGroup := namespacesGroup.Group("/:namespaceName/users")
			{
				usersGroup.POST("/", CreateUser())
				usersGroup.GET("/", GetUsers())
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

func createTestNamespace(name string) {
	namespace := mocks.PrepareNamespace(name)
	_, err := clientset.CoreV1().Namespaces().Create(context.TODO(), &namespace, metav1.CreateOptions{})

	if err != nil {
		panic(err)
	}
}

func setupScheme() *runtime.Scheme {
	schema := scheme.Scheme
	_ = cappv1.AddToScheme(schema)
	return schema
}
