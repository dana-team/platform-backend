package e2e_tests

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/dana-team/platform-backend/src/utils/testutils"
	"github.com/dana-team/platform-backend/src/utils/testutils/mocks"
	configv1 "github.com/openshift/api/config/v1"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	zapctrl "sigs.k8s.io/controller-runtime/pkg/log/zap"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

const (
	platformBackendCappName      = "platform-backend"
	platformBackendCappNamespace = platformBackendCappName + "-system"

	htpassEncoded = e2eUser + ":$apr1$D3cCCeru$cuVKr.hn1cbKrhg5NSaT20"
	oauthName     = "cluster"

	htpasswdSecretNamespace = "openshift-config"
	htpasswdProviderName    = "e2e-test-htpasswd"
	htpasswdSecretName      = "e2e-tests-htpasswd-secret"
	htpasswdKey             = "htpasswd"
	htpasswdType            = "HTPasswd"
	clusterAdminKey         = "cluster-admin"
	clusterIngressName      = "cluster"

	tokenKey = "token"
	loginKey = "login"
)

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)

	SetDefaultEventuallyTimeout(testutils.DefaultEventually)
	RunSpecs(t, "Platform Backend Suite")
}

var _ = SynchronizedBeforeSuite(func() {
	newScheme()
	initClients()
	cleanup()
	createTestUserIdentity()
	getURLFromCapp()
	getTokenFromLogin()
	getClusterIngressDomain()
}, func() {
	newScheme()
	initClients()
	getURLFromCapp()
	getTokenFromLogin()
	getClusterIngressDomain()
})

var _ = SynchronizedAfterSuite(func() {}, func() {
	cleanup()
})

// cleanup deletes all the resources which were created for the e2e testing.
func cleanup() {
	cleanUpTestNamespaces()
	removeHTPasswdProviderFromOAuth()
	removeUserIdentity()
	removeClusterRoleBindingFromUser()
}

// createTestUserIdentity adds a new HTPasswd provider to the OAuth object in
// the cluster and adds a ClusterRoleBinding to the created user.
func createTestUserIdentity() {
	createHTPasswdProviderSecret()
	addHTPasswdProviderToOAuth()
	addClusterRoleBindingToUser()
}

// initClient initializes a k8s client.
func initClients() {
	initHTTPClient()
	initKubeClient()
}

// initKubeClient initializes a Kubernetes client.
func initKubeClient() {
	opts := zapctrl.Options{Development: true}
	ctrl.SetLogger(zapctrl.New(zapctrl.UseFlagOptions(&opts)))

	cfg, err := config.GetConfig()
	Expect(err).NotTo(HaveOccurred())

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())
}

// initKubeClient initializes an HTTP client.
func initHTTPClient() {
	httpClient = http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

// addClusterRoleBindingToUser adds a ClusterRoleBinding to a user.
func addClusterRoleBindingToUser() {
	clusterRole := mocks.PrepareClusterRoleBinding(e2eUser, clusterAdminKey)
	createResource(k8sClient, &clusterRole)
}

// removeClusterRoleBindingFromUser removes a ClusterRoleBinding from a user.
func removeClusterRoleBindingFromUser() {
	clusterRole := mocks.PrepareClusterRoleBinding(e2eUser, clusterAdminKey)
	deleteResource(k8sClient, &clusterRole)
}

// removeUserIdentity removes User and Identity objects from the cluster.
func removeUserIdentity() {
	user := mocks.PrepareUser(e2eUser)
	deleteResource(k8sClient, &user)

	identity := mocks.PrepareIdentity(htpasswdProviderName, e2eUser)
	deleteResource(k8sClient, &identity)
}

// removeHTPasswdProviderFromOAuth removes any newly-added Identity Providers from the OAuth object.
func removeHTPasswdProviderFromOAuth() {
	oauth := getOAuth(k8sClient, oauthName)
	if len(oauth.Spec.IdentityProviders) > 0 {
		oauth.Spec.IdentityProviders = []configv1.IdentityProvider{oauth.Spec.IdentityProviders[0]}
	}
	updateOAuth(k8sClient, oauth)
}

// addHTPasswdProviderToOAuth adds a new HTPasswd provider to the OAuth object.
func addHTPasswdProviderToOAuth() {
	oauth := getOAuth(k8sClient, oauthName)
	htpasswdProvider := mocks.PrepareHTPasswdProvider(htpasswdSecretName)
	identityProvider := mocks.PrepareHTPasswdIdentityProvider(htpasswdProviderName, htpasswdType, htpasswdProvider)

	oauth.Spec.IdentityProviders = append(oauth.Spec.IdentityProviders, identityProvider)
	updateOAuth(k8sClient, oauth)
}

// createHTPasswdProviderSecret creates an HTPasswd secret.
func createHTPasswdProviderSecret() {
	secret := mocks.PrepareSecret(htpasswdSecretName, htpasswdSecretNamespace, htpasswdKey, htpassEncoded)
	createResource(k8sClient, &secret)
}

// getURLFromCapp returns the URL of the Capp the backend platform is exposed at.
func getURLFromCapp() {
	platformCapp := getCapp(k8sClient, platformBackendCappName, platformBackendCappNamespace)
	if platformCapp.Status.KnativeObjectStatus.RouteStatusFields.URL != nil {
		platformURL = platformCapp.Status.KnativeObjectStatus.RouteStatusFields.URL.String()
	}
}

// getTokenFromLogin performs a login request and gets the token back in response.
// It uses Eventually since it may take a while after OAuth object getting updated util
// the user can authenticate to the cluster; this is due to OAuth pods getting restarted.
func getTokenFromLogin() {
	baseURI := fmt.Sprintf("%s/v1/%s", platformURL, loginKey)
	userToken = ""

	Eventually(func() bool {
		_, response := performHTTPRequest(httpClient, nil, http.MethodPost, baseURI, e2eUser, e2ePassword, "")
		token, ok := response[tokenKey]

		if !ok {
			return false
		}

		userToken, ok = token.(string)
		return ok

	}, testutils.Timeout, testutils.Interval).Should(Equal(true))
}

// getClusterIngressDomain returns the ingress domain of an OpenShift cluster
func getClusterIngressDomain() {
	ingress := &configv1.Ingress{}
	getClusterResource(k8sClient, ingress, clusterIngressName)
	clusterDomain = ingress.Spec.Domain
}

// cleanUpTestNamespaces() deletes test namespaces.
func cleanUpTestNamespaces() {
	namespaces := listNamespaces(k8sClient, e2eLabelKey, e2eLabelValue)
	for _, namespace := range namespaces.Items {
		Expect(k8sClient.Delete(context.Background(), &namespace)).To(Succeed())
	}

	Eventually(func() bool {
		namespaces = listNamespaces(k8sClient, e2eLabelKey, e2eLabelValue)
		return len(namespaces.Items) == 0
	}, testutils.Timeout, testutils.Interval).Should(Equal(true))
}
