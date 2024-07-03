package controllers_test

import (
	"context"
	"testing"

	cappv1 "github.com/dana-team/container-app-operator/api/v1alpha1"
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
	client    *fake.Clientset
	dynClient runtimeClient.WithWatch
	logger    *zap.Logger
)

const (
	testName     = "test"
	doesNotExist = "doesNotExist"
)

func TestMain(m *testing.M) {
	setup()
	setupCappRevisions()
	m.Run()
}

func setup() {
	client = fake.NewSimpleClientset()
	dynClient = runtimeFake.NewClientBuilder().WithScheme(setupScheme()).Build()
	logger, _ = zap.NewProduction()

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

func setupScheme() *runtime.Scheme {
	schema := scheme.Scheme
	_ = cappv1.AddToScheme(schema)
	return schema
}
