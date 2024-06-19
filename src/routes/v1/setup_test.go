package v1_test

import (
	"context"
	"fmt"
	"testing"

	cappv1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	routev1 "github.com/dana-team/platform-backend/src/routes/v1"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes/scheme"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
	runtimeFake "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var (
	router    *gin.Engine
	client    *fake.Clientset
	dynClient runtimeClient.WithWatch
)

func TestMain(m *testing.M) {
	Setup()
	CreateTestSecret()
	CreateTestNamespace()
	CreateTestCapp()
	CreateTestCappRevision()
	m.Run()
}

func Setup() {
	client = fake.NewSimpleClientset()
	dynClient = runtimeFake.NewClientBuilder().WithScheme(setupScheme()).Build()
	logger, _ := zap.NewProduction()
	router = SetupRouter(logger)
}

func SetupRouter(logger *zap.Logger) *gin.Engine {
	engine := gin.Default()
	engine.Use(func(c *gin.Context) {
		c.Set("logger", logger)
		c.Set("kubeClient", client)
		c.Set("dynClient", dynClient)
		c.Next()
	})
	v1 := engine.Group("/v1")
	{
		namespacesGroup := v1.Group("/namespaces")
		{
			namespacesGroup.GET("", routev1.ListNamespaces())
			namespacesGroup.GET("/:namespaceName", routev1.GetNamespace())
			namespacesGroup.POST("/", routev1.CreateNamespace())
			namespacesGroup.DELETE("/:namespaceName", routev1.DeleteNamespace())

			secretsGroup := namespacesGroup.Group("/:namespaceName/secrets")
			{
				secretsGroup.POST("/", routev1.CreateSecret())
				secretsGroup.GET("", routev1.GetSecrets())
				secretsGroup.GET("/:secretName", routev1.GetSecret())
				secretsGroup.PUT("/:secretName", routev1.UpdateSecret())
				secretsGroup.DELETE("/:secretName", routev1.DeleteSecret())
			}

			cappGroup := namespacesGroup.Group("/:namespaceName/capps")
			{
				cappGroup.POST("/", routev1.CreateCapp())
				cappGroup.GET("", routev1.GetCapps())
				cappGroup.GET("/:cappName", routev1.GetCapp())
				cappGroup.PUT("/:cappName", routev1.UpdateCapp())
				cappGroup.DELETE("/:cappName", routev1.DeleteCapp())
			}

			cappRevisionGroup := namespacesGroup.Group("/:namespaceName/capprevisions")
			{
				cappRevisionGroup.GET("", routev1.GetCappRevisions())
				cappRevisionGroup.GET("/:cappRevisionName", routev1.GetCappRevision())
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

func CreateTestCapp() {
	capp := cappv1.Capp{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "test-capp",
			Namespace:   "test-namespace",
			Annotations: map[string]string{},
			Labels:      map[string]string{},
		},
		Spec:   cappv1.CappSpec{},
		Status: cappv1.CappStatus{},
	}
	err := dynClient.Create(context.TODO(), &capp)
	if err != nil {
		panic(err)
	}
}

func CreateTestCappRevision() {
	capp := cappv1.CappRevision{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "test-capprevision",
			Namespace:   "test-namespace",
			Annotations: map[string]string{},
			Labels:      map[string]string{},
		},
		Spec:   cappv1.CappRevisionSpec{},
		Status: cappv1.CappRevisionStatus{},
	}
	err := dynClient.Create(context.TODO(), &capp)
	if err != nil {
		panic(err)
	}
}

func setupScheme() *runtime.Scheme {
	schema := scheme.Scheme
	err := cappv1.AddToScheme(schema)
	fmt.Println(err)
	return schema
}
