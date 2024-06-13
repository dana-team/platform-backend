package v1_test

import (
	"context"
	"testing"

	routev1 "github.com/dana-team/platform-backend/src/routes/v1"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

var (
	router *gin.Engine
	client *fake.Clientset
)

func TestMain(m *testing.M) {
	Setup()
	CreateTestSecret()
	CreateTestNamespace()
	m.Run()
}

func Setup() {
	client = fake.NewSimpleClientset()
	logger, _ := zap.NewProduction()
	router = SetupRouter(logger)
}

func SetupRouter(logger *zap.Logger) *gin.Engine {
	engine := gin.Default()
	engine.Use(func(c *gin.Context) {
		c.Set("logger", logger)
		c.Set("kubeClient", client)
		c.Next()
	})
	v1 := engine.Group("/v1")
	{
		namespacesGroup := v1.Group("/namespaces")
		{
			namespacesGroup.GET("/", routev1.ListNamespaces())
			namespacesGroup.GET("/:namespaceName", routev1.GetNamespace())
			namespacesGroup.POST("/", routev1.CreateNamespace())
			namespacesGroup.DELETE("/:namespaceName", routev1.DeleteNamespace())

			secretsGroup := namespacesGroup.Group("/:namespaceName/secrets")
			{
				secretsGroup.POST("/", routev1.CreateSecret())
				secretsGroup.GET("/", routev1.GetSecrets())
				secretsGroup.GET("/:secretName", routev1.GetSecret())
				secretsGroup.PATCH("/:secretName", routev1.PatchSecret())
				secretsGroup.DELETE("/:secretName", routev1.DeleteSecret())
			}
		}
	}
	return engine
}

func CreateTestNamespace() {
	namespace := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-namespace",
		},
	}
	_, err := client.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
}

func CreateTestSecret() {
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-secret",
			Namespace: "default",
		},
		Type: v1.SecretTypeOpaque,
		Data: map[string][]byte{
			"key1": []byte("ZmFrZQ=="),
		},
	}
	_, err := client.CoreV1().Secrets("default").Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
}
