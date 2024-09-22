package mocks

import (
	"context"
	"github.com/dana-team/platform-backend/src/utils/testutils"
	corev1 "k8s.io/api/core/v1"
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
func CreateTestCapp(dynClient runtimeClient.WithWatch, name, namespace, domain, site string, labels, annotations map[string]string) {
	cappRevision := PrepareCapp(name, namespace, domain, site, labels, annotations)
	err := dynClient.Create(context.TODO(), &cappRevision)
	if err != nil {
		panic(err)
	}
}

// CreateTestCappWithState creates a test Capp object with given state.
func CreateTestCappWithState(dynClient runtimeClient.WithWatch, name, namespace, state, site string, labels, annotations map[string]string) {
	cappRevision := PrepareCappWithKnativeObject(name, namespace, state, site, labels, annotations)
	err := dynClient.Create(context.TODO(), &cappRevision)
	if err != nil {
		panic(err)
	}
}

// CreateTestCappWithHostname creates a test Capp object with hostname.
func CreateTestCappWithHostname(dynClient runtimeClient.WithWatch, name, namespace, hostname, domain string, labels, annotations map[string]string) {
	capp := PrepareCappWithHostname(name, namespace, hostname, domain, labels, annotations)
	err := dynClient.Create(context.TODO(), &capp)
	if err != nil {
		panic(err)
	}
}

// CreateTestCappRevision creates a test CappRevision object.
func CreateTestCappRevision(dynClient runtimeClient.WithWatch, name, namespace, site string, labels, annotations map[string]string) {
	cappRevision := PrepareCappRevision(name, namespace, site, labels, annotations)
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

// CreateTestPod creates a test Pod object.
func CreateTestPod(fakeClient *fake.Clientset, namespace, name, cappName string, isMultipleContainers bool) {
	pod := PreparePod(namespace, name, cappName, isMultipleContainers)
	_, err := fakeClient.CoreV1().Pods(namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
}

// CreateTestCNAMERecord creates a test CNAME record
func CreateTestCNAMERecord(dynClient runtimeClient.WithWatch, name, cappName, cappNSName, hostname string, readyStatus, syncedStatus corev1.ConditionStatus) {
	record := prepareCNAMERecord(name, cappName, cappNSName, hostname, readyStatus, syncedStatus)
	err := dynClient.Create(context.TODO(), &record)
	if err != nil {
		panic(err)
	}
}

// CreateTestCNAMERecordWithoutConditions creates a test CNAME record without conditions
func CreateTestCNAMERecordWithoutConditions(dynClient runtimeClient.WithWatch, name, cappName, cappNSName, hostname string) {
	record := prepareBaseCNAMERecord(name, cappName, cappNSName, hostname)
	err := dynClient.Create(context.TODO(), &record)
	if err != nil {
		panic(err)
	}
}

// CreateTestServiceAccount creates a test service account
func CreateTestServiceAccount(fakeClient *fake.Clientset, namespace, name string, dockerCfgSecretName string) {
	serviceAccount := PrepareServiceAccount(name, namespace, dockerCfgSecretName)
	_, err := fakeClient.CoreV1().ServiceAccounts(namespace).Create(context.TODO(), serviceAccount, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
}

// CreateTestServiceAccountWithToken creates a test service account
func CreateTestServiceAccountWithToken(fakeClient *fake.Clientset, namespace, serviceAccountName, tokenSecretName, tokenValue, dockerCfgSecretName string) {
	tokenSecret := PrepareTokenSecret(tokenSecretName, namespace, tokenValue, serviceAccountName)
	createTestSecret(fakeClient, tokenSecret)

	dockerCfgSecret := PrepareDockerConfigSecret(dockerCfgSecretName, namespace, tokenSecretName)
	createTestSecret(fakeClient, dockerCfgSecret)

	CreateTestServiceAccount(fakeClient, namespace, serviceAccountName, dockerCfgSecretName)
}

// createTestSecret receives a secret and creates it.
func createTestSecret(fakeClient *fake.Clientset, secret corev1.Secret) {
	_, err := fakeClient.CoreV1().Secrets(secret.Namespace).Create(context.TODO(), &secret, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
}

// CreateTestPlacement creates a test Placement object.
func CreateTestPlacement(dynClient runtimeClient.WithWatch, name, namespace string, labels map[string]string) {
	placement := PreparePlacement(name, namespace, labels)
	err := dynClient.Create(context.TODO(), &placement)
	if err != nil {
		panic(err)
	}
}
