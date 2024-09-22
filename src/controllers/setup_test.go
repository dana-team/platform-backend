package controllers

import (
	"context"
	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	dnsrecordv1alpha1 "github.com/dana-team/provider-dns/apis/record/v1alpha1"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes/scheme"
	clusterv1beta1 "open-cluster-management.io/api/cluster/v1beta1"
	"os"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
	runtimeFake "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

var (
	fakeClient *fake.Clientset
	dynClient  runtimeClient.WithWatch
	logger     *zap.Logger
)

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

func setup() {
	fakeClient = fake.NewSimpleClientset()
	dynClient = runtimeFake.NewClientBuilder().WithScheme(setupScheme()).Build()
	logger, _ = zap.NewProduction()
}

func createTestNamespace(name string, labels map[string]string) {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
	}
	_, err := fakeClient.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
}

func createTestSecret(secretName string, namespace string, labels map[string]string) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
			Labels:    labels,
		},
		Type: corev1.SecretTypeOpaque,
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
	utilruntime.Must(cappv1alpha1.AddToScheme(schema))
	utilruntime.Must(dnsrecordv1alpha1.AddToScheme(schema))
	utilruntime.Must(clusterv1beta1.AddToScheme(schema))

	return schema
}
