package controllers

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
	fakeClient *fake.Clientset
	dynClient  runtimeClient.WithWatch
	logger     *zap.Logger
)

func TestMain(m *testing.M) {
	setup()
	m.Run()
}

func setup() {
	fakeClient = fake.NewSimpleClientset()
	dynClient = runtimeFake.NewClientBuilder().WithScheme(setupScheme()).Build()
	logger, _ = zap.NewProduction()

}

func createTestNamespace(name string) {
	namespace := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	_, err := fakeClient.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
}
func createTestSecret(secretName string, namespace string) {
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
		},
		Type: v1.SecretTypeOpaque,
		Data: map[string][]byte{
			"key-1": []byte("ZmFrZQ=="),
		},
	}
	_, err := fakeClient.CoreV1().Secrets(namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
}
func setupScheme() *runtime.Scheme {
	schema := scheme.Scheme
	_ = cappv1.AddToScheme(schema)
	return schema
}
