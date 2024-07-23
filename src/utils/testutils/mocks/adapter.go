package mocks

import (
	"context"
	"github.com/dana-team/platform-backend/src/utils/testutils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
)

// CreateTestNamespace creates a test Namespace object.
func CreateTestNamespace(fakeClient *fake.Clientset, name string) {
	namespace := PrepareNamespace(name, map[string]string{})
	_, err := fakeClient.CoreV1().Namespaces().Create(context.TODO(), &namespace, metav1.CreateOptions{})

	if err != nil {
		panic(err)
	}
}

// CreateTestSecret creates a test Secret object.
func CreateTestSecret(fakeClient *fake.Clientset, name, namespace string) {
	secret := PrepareSecret(name, namespace, testutils.SecretDataKey, testutils.SecretDataValueEncoded)
	_, err := fakeClient.CoreV1().Secrets(namespace).Create(context.TODO(), &secret, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
}

// CreateTestCapp creates a test Capp object.
func CreateTestCapp(dynClient runtimeClient.WithWatch, name, namespace, domain string, labels, annotations map[string]string) {
	cappRevision := PrepareCapp(name, namespace, domain, labels, annotations)
	err := dynClient.Create(context.TODO(), &cappRevision)
	if err != nil {
		panic(err)
	}
}

// CreateTestCappWithHostname creates a test Capp object with a hostname.
func CreateTestCappWithHostname(dynClient runtimeClient.WithWatch, name, namespace string, labels, annotations map[string]string) {
	capp := PrepareCappWithHostname(name, namespace, labels, annotations)
	err := dynClient.Create(context.TODO(), &capp)
	if err != nil {
		panic(err)
	}
}

// CreateTestCappRevision creates a test CappRevision object.
func CreateTestCappRevision(dynClient runtimeClient.WithWatch, name, namespace string, labels, annotations map[string]string) {
	cappRevision := PrepareCappRevision(name, namespace, labels, annotations)
	err := dynClient.Create(context.TODO(), &cappRevision)
	if err != nil {
		panic(err)
	}
}

// CreateTestRoleBinding creates a test RoleBinding object.
func CreateTestRoleBinding(fakeClient *fake.Clientset, name, namespace, role string) {
	roleBinding := PrepareRoleBinding(name, namespace, role)

	_, err := fakeClient.RbacV1().RoleBindings(namespace).Create(context.TODO(), &roleBinding, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
}

// CreateTestConfigMap creates a test ConfigMap object.
func CreateTestConfigMap(fakeClient *fake.Clientset, name, namespace string) {
	configMap := PrepareConfigMap(name, namespace, map[string]string{testutils.ConfigMapDataKey: testutils.ConfigMapDataValue})
	_, err := fakeClient.CoreV1().ConfigMaps(namespace).Create(context.TODO(), &configMap, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
}
