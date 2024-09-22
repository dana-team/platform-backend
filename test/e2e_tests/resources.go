package e2e_tests

import (
	"context"
	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	"github.com/dana-team/platform-backend/src/utils/testutils"
	"github.com/dana-team/platform-backend/src/utils/testutils/mocks"
	rcsv1alpha1 "github.com/dana-team/rcs-ocm-deployer/api/v1alpha1"
	multicluster "github.com/oam-dev/cluster-gateway/pkg/apis/cluster/transport"
	. "github.com/onsi/gomega"
	configv1 "github.com/openshift/api/config/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"math/rand"
	clusterv1beta1 "open-cluster-management.io/api/cluster/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

const (
	CappImageName = "ghcr.io/dana-team/capp-gin-app:v0.2.0"
	charset       = "abcdefghijklmnopqrstuvwxyz0123456789"
	randStrLength = 10
)

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// generateRandomString returns a random string of the specified length using characters from the charset.
func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// doesResourceExist checks if a given Kubernetes object exists in the cluster.
func doesResourceExist(k8sClient client.Client, obj client.Object) bool {
	copyObject := obj.DeepCopyObject().(client.Object)
	key := client.ObjectKeyFromObject(copyObject)
	err := k8sClient.Get(context.Background(), key, copyObject)
	Expect(err).To(SatisfyAny(BeNil(), WithTransform(errors.IsNotFound, BeTrue())))
	return !errors.IsNotFound(err)
}

// getResource fetches an existing resource and returns an instance of it.
func getResource(k8sClient client.Client, obj client.Object, name, namespace string) {
	Expect(k8sClient.Get(context.Background(), client.ObjectKey{Name: name, Namespace: namespace}, obj))
}

// getClusterResource fetches an existing Cluster resource and returns an instance of it.
func getClusterResource(k8sClient client.Client, obj client.Object, name string) {
	Expect(k8sClient.Get(context.Background(), client.ObjectKey{Name: name}, obj))
}

// createResource creates a new resource.
func createResource(k8sClient client.Client, object client.Object) {
	Expect(k8sClient.Create(context.Background(), object)).To(SatisfyAny(BeNil(), WithTransform(errors.IsAlreadyExists, BeTrue())))
}

// deleteResource deletes a resource.
func deleteResource(k8sClient client.Client, object client.Object) {
	Expect(k8sClient.Delete(context.Background(), object)).To(SatisfyAny(BeNil(), WithTransform(errors.IsNotFound, BeTrue())))
	Eventually(func() bool {
		return doesResourceExist(k8sClient, object)
	}, testutils.Timeout, testutils.Interval).ShouldNot(BeTrue())
}

// generateName generates a new name by combining the given baseName
// with a randomly generated string of a specified length.
func generateName(baseName string) string {
	randString := generateRandomString(randStrLength)
	return baseName + "-" + randString
}

// createCapp creates a new Capp instance with a unique name and returns it.
func createCapp(k8sClient client.Client, capp *cappv1alpha1.Capp) *cappv1alpha1.Capp {
	newCapp := capp.DeepCopy()
	Expect(k8sClient.Create(context.Background(), newCapp)).To(Succeed())
	Eventually(func() string {
		return getCapp(k8sClient, newCapp.Name, newCapp.Namespace).Status.KnativeObjectStatus.ConfigurationStatusFields.LatestReadyRevisionName
	}, testutils.Timeout, testutils.Interval).ShouldNot(Equal(""))

	Eventually(func() string {
		newCapp = getCapp(k8sClient, newCapp.Name, newCapp.Namespace)
		return newCapp.Status.ApplicationLinks.Site
	}, testutils.Timeout, testutils.Interval).ShouldNot(Equal(""), "missing site in capp status")
	return newCapp
}

// createDisabledCapp creates new Disabled Capp instance with a unique name and returns it.
func createDisabledCapp(k8sClient client.Client, capp *cappv1alpha1.Capp) *cappv1alpha1.Capp {
	if capp.Spec.State != testutils.DisabledState {
		panic("Capp must be disabled")
	}

	newCapp := capp.DeepCopy()
	Expect(k8sClient.Create(context.Background(), newCapp)).To(Succeed())
	Eventually(func() string {
		return getCapp(k8sClient, newCapp.Name, newCapp.Namespace).Status.StateStatus.State
	}, testutils.Timeout, testutils.Interval).Should(Equal(testutils.DisabledState))
	return newCapp
}

// getCapp fetches and returns an existing instance of a Capp.
func getCapp(k8sClient client.Client, name string, namespace string) *cappv1alpha1.Capp {
	capp := &cappv1alpha1.Capp{}
	getResource(k8sClient, capp, name, namespace)
	return capp
}

// getCappRevisionNames returns a list of the CappRevision names related to a specific Capp in a namespace.
func getCappRevisionNames(k8sClient client.Client, cappName, namespace, clusterName string) []string {
	revisions := cappv1alpha1.CappRevisionList{}
	labelSelector := client.MatchingLabels{testutils.LabelCappName: cappName}
	listOptions := []client.ListOption{
		labelSelector,
		client.InNamespace(namespace),
	}
	Expect(k8sClient.List(context.TODO(), &revisions, listOptions...)).To(Succeed())

	Expect(k8sClient.List(multicluster.WithMultiClusterContext(context.TODO(), clusterName), &revisions, listOptions...)).To(Succeed())

	names := make([]string, len(revisions.Items))
	for i, revision := range revisions.Items {
		names[i] = revision.Name
	}

	return names
}

// listNamespaces returns a list of namespaces.
func listNamespaces(k8sClient client.Client, labelKey, labelValue string) corev1.NamespaceList {
	namespaces := corev1.NamespaceList{}
	labelSelector := client.MatchingLabels{labelKey: labelValue}
	Expect(k8sClient.List(context.TODO(), &namespaces, labelSelector)).To(Succeed())

	return namespaces
}

// listPlacements returns a list of placemements.
func listPlacements(k8sClient client.Client, labelKey, labelValue string) clusterv1beta1.PlacementList {
	placements := clusterv1beta1.PlacementList{}
	labelSelector := client.MatchingLabels{labelKey: labelValue}
	Expect(k8sClient.List(context.TODO(), &placements, labelSelector)).To(Succeed())

	return placements
}

// getOAuth fetches and returns an existing instance of an OAuth.
func getOAuth(k8sClient client.Client, name string) *configv1.OAuth {
	oauth := &configv1.OAuth{}
	getClusterResource(k8sClient, oauth, name)
	return oauth
}

// updateOAuth updates an existing Oauth instance.
func updateOAuth(k8sClient client.Client, oauth *configv1.OAuth) {
	Eventually(func() error {
		return k8sClient.Update(context.Background(), oauth)
	}, testutils.Timeout, testutils.Interval).Should(Succeed())
}

// getSecret fetches and returns an existing instance of a RCSConfig.
func getRCSConfig(k8sClient client.Client, name, namespace string) *rcsv1alpha1.RCSConfig {
	rcsConfig := &rcsv1alpha1.RCSConfig{}
	getResource(k8sClient, rcsConfig, name, namespace)
	return rcsConfig
}

// updateRCSConfig updates an existing RCSConfig instance.
func updateRCSConfig(k8sClient client.Client, rcsConfig *rcsv1alpha1.RCSConfig) error {
	return k8sClient.Update(context.Background(), rcsConfig)
}

// createPlacement creates a new Placement resource.
func createPlacement(k8sClient client.Client, name string, namespace string, labels map[string]string) {
	placement := mocks.PreparePlacement(name, namespace, labels)
	createResource(k8sClient, &placement)
}

// getSecret fetches and returns an existing instance of a Secret.
func getSecret(k8sClient client.Client, name string, namespace string) *corev1.Secret {
	secret := &corev1.Secret{}
	getResource(k8sClient, secret, name, namespace)
	return secret
}

// getRoleBinding fetches and returns an existing instance of a RoleBinding.
func getRoleBinding(k8sClient client.Client, name string, namespace string) *rbacv1.RoleBinding {
	roleBinding := &rbacv1.RoleBinding{}
	getResource(k8sClient, roleBinding, name, namespace)
	return roleBinding
}

// getServiceAccount fetches and returns an existing instance of a ServiceAccount.
func getServiceAccount(k8sClient client.Client, name string, namespace string) *corev1.ServiceAccount {
	serviceAccount := &corev1.ServiceAccount{}
	getResource(k8sClient, serviceAccount, name, namespace)
	return serviceAccount
}

// getServiceAccountTokenSecret fetches and returns an existing instance of a ServiceAccount's token secret.
func getServiceAccountTokenSecret(k8sClient client.Client, name string, namespace string) *corev1.Secret {
	serviceAccount := getServiceAccount(k8sClient, name, namespace)

	for _, ref := range serviceAccount.Secrets {
		secret := getSecret(k8sClient, ref.Name, namespace)

		if secret.Type == corev1.SecretTypeDockercfg {
			for _, ownerRef := range secret.OwnerReferences {
				tokenSecret := getSecret(k8sClient, ownerRef.Name, namespace)

				if tokenSecret.Type == corev1.SecretTypeServiceAccountToken {
					return tokenSecret
				}
			}
		}
	}

	return nil
}

// createTestNamespace creates a test Namespace object.
func createTestNamespace(k8sClient client.Client, name string) {
	namespace := mocks.PrepareNamespace(name, map[string]string{e2eLabelKey: e2eLabelValue})
	createResource(k8sClient, &namespace)
}

// createTestCapp creates a test Capp object.
func createTestCapp(k8sClient client.Client, name, namespace, site string, labels, annotations map[string]string) *cappv1alpha1.Capp {
	capp := mocks.PrepareCapp(name, namespace, clusterDomain, site, labels, annotations)
	return createCapp(k8sClient, &capp)
}

// createTestCapp creates a test Capp object.
func createTestDisabledCapp(k8sClient client.Client, name, namespace, site string, labels, annotations map[string]string) {
	capp := mocks.PrepareCappWithState(name, namespace, testutils.DisabledState, site, labels, annotations)
	createDisabledCapp(k8sClient, &capp)
}

// createTestCapp creates a test Capp object with hostname.
func createTestCappWithHostname(k8sClient client.Client, name, namespace, hostname, domain string, labels, annotations map[string]string) {
	capp := mocks.PrepareCappWithHostname(name, namespace, hostname, domain, labels, annotations)
	createCapp(k8sClient, &capp)
}

// createTestConfigMap creates a test configMap object.
func createTestConfigMap(k8sClient client.Client, name, namespace string, data map[string]string) {
	configMap := mocks.PrepareConfigMap(name, namespace, data)
	createResource(k8sClient, &configMap)
}

// createTestSecret creates a test Secret object.
func createTestSecret(k8sClient client.Client, name, namespace string) {
	oneSecret := mocks.PrepareSecret(name, namespace, testutils.SecretDataKey, testutils.SecretDataValueEncoded)
	createResource(k8sClient, &oneSecret)
}

// createTestUser creates a test User object.
func createTestUser(k8sClient client.Client, name, namespace string) {
	oneSecret := mocks.PrepareRoleBinding(name, namespace, testutils.AdminKey)
	createResource(k8sClient, &oneSecret)
}

// CreateTestServiceAccount creates a test ServiceAccount object.
func CreateTestServiceAccount(k8sClient client.Client, name, namespace, tokenSecretName, tokenValue, dockerCfgSecretName string) {
	tokenSecret := mocks.PrepareTokenSecret(tokenSecretName, namespace, tokenValue, name)
	createResource(k8sClient, &tokenSecret)

	serviceAccount := mocks.PrepareServiceAccount(name, namespace, "")
	createResource(k8sClient, serviceAccount)
}
