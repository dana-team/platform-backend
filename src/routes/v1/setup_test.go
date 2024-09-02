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

const (
	cluster = "test-cluster"
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
		c.Set(middleware.ClusterCtxKey, cluster)
		c.Next()
	})

	v1 := engine.Group("/v1")

	setupNamespaceRoutes(v1, nil, nil)
	setupClustersRoutes(v1, nil, nil)

	return engine
}

func setupScheme() *runtime.Scheme {
	schema := scheme.Scheme
	_ = cappv1alpha1.AddToScheme(schema)
	_ = dnsrecordv1alpha1.AddToScheme(schema)
	return schema
}
