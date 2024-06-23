package v1_test

import (
	"context"
	"fmt"
	"github.com/dana-team/platform-backend/src/utils"
	"testing"

	rbacv1 "k8s.io/api/rbac/v1"

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

const (
	testNamespace = "test-namespace"
	testName      = "test"
)

func TestMain(m *testing.M) {
	setup()
	createTestSecret()
	createTestNamespace(testNamespace)
	setupCappRevisions()
	createTestCapp()
	createConfigMap()
	m.Run()
}

func setup() {
	client = fake.NewSimpleClientset()
	dynClient = runtimeFake.NewClientBuilder().WithScheme(setupScheme()).Build()
	logger, _ := zap.NewProduction()
	router = setupRouter(logger)
}

func setupRouter(logger *zap.Logger) *gin.Engine {
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

			usersGroup := namespacesGroup.Group("/:namespaceName/users")
			{
				usersGroup.POST("/", routev1.CreateUser())
				usersGroup.GET("/", routev1.GetUsers())
				usersGroup.GET("/:userName", routev1.GetUser())
				usersGroup.PUT("/:userName", routev1.UpdateUser())
				usersGroup.DELETE("/:userName", routev1.DeleteUser())
			}

			configMapGroup := namespacesGroup.Group("/:namespaceName/configmaps")
			{
				configMapGroup.GET("/:configMapName", routev1.GetConfigMap())
			}
		}
	}
	return engine
}

func createTestNamespace(name string) {
	namespace := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	_, err := client.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
}

func createTestSecret() {
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

func createTestCapp() {
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

// createRoleBinding creates a RoleBinding in the specified namespace
func createRoleBinding(namespace string, userName string) {
	roleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      userName,
			Namespace: namespace,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:     rbacv1.UserKind,
				Name:     userName,
				APIGroup: rbacv1.GroupName,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind: "ClusterRole", // Could be "Role" or "ClusterRole"
			Name: "capp-user-admin",
		},
	}

	_, err := client.RbacV1().RoleBindings(namespace).Create(context.TODO(), roleBinding, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
}

// createTestCappRevision creates a test CappRevision object
func createTestCappRevision(name, namespace string, labels, annotations map[string]string) {
	cappRevision := utils.GetBareCappRevision(name, namespace, labels, annotations)
	err := dynClient.Create(context.TODO(), &cappRevision)
	if err != nil {
		panic(err)
	}
}

func createConfigMap() {
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-configmap",
			Namespace: "test-namespace",
		},
		Data: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}
	_, err := client.CoreV1().ConfigMaps("test-namespace").Create(context.TODO(), configMap, metav1.CreateOptions{})
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
