package e2e_tests

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	. "github.com/onsi/gomega"
	configv1 "github.com/openshift/api/config/v1"
	userv1 "github.com/openshift/api/user/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	e2eUser              = "e2e-user"
	e2ePassword          = "e2e-password"
	e2eNamespace         = "e2e-namespace"
	e2eLabelKey          = "e2e-label-key"
	e2eLabelValue        = "e2e-label-value"
	testCappName         = "e2e-capp"
	testCappRevisionName = "e2e-capp-revision"
	testConfigMapName    = "e2e-configmap"
	testSecretName       = "e2e-secret"
	testUserName         = "e2e-user"
)

const (
	httpAuthorizationHeader = "Authorization"
	httpBearerToken         = "Bearer"
	contentType             = "Content-Type"
	applicationJson         = "application/json"
)

var (
	scheme        = runtime.NewScheme()
	k8sClient     client.Client
	httpClient    http.Client
	logger        logr.Logger
	clusterDomain string
	platformURL   string
	userToken     string
)

func newScheme() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(cappv1alpha1.AddToScheme(scheme))
	utilruntime.Must(configv1.AddToScheme(scheme))
	utilruntime.Must(userv1.AddToScheme(scheme))
}

// createHTTPClient creates an HTTP client with TLS configuration.
func createHTTPClient() http.Client {
	return http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

// compareResponses compares two responses and asserts that they are equal.
func compareResponses(response, expectedResponse map[string]interface{}) {
	expectedResponseJSON, err := json.Marshal(expectedResponse)
	Expect(err).ShouldNot(HaveOccurred())

	var expectedResponseNormalized map[string]interface{}
	err = json.Unmarshal(expectedResponseJSON, &expectedResponseNormalized)
	Expect(err).ShouldNot(HaveOccurred())

	Expect(response).Should(Equal(expectedResponseNormalized))
}

// prepareAuthorizedHTTPRequest returns a request object based on certain parameters.
func prepareAuthorizedHTTPRequest(body io.Reader, httpMethod, baseURI, username, password, userToken string) *http.Request {
	request, err := http.NewRequest(httpMethod, baseURI, body)
	Expect(err).NotTo(HaveOccurred())
	request.Header.Set(contentType, applicationJson)

	if username != "" {
		request.SetBasicAuth(username, password)
	} else if userToken != "" {
		request.Header.Set(httpAuthorizationHeader, fmt.Sprintf("%s %s", httpBearerToken, userToken))
	}

	return request
}

// performAuthorizedHTTPRequest makes an HTTP request and returns a response.
func performAuthorizedHTTPRequest(httpClient http.Client, body io.Reader, httpMethod, baseURI, username, password, userToken string) (int, map[string]interface{}) {
	request := prepareAuthorizedHTTPRequest(body, httpMethod, baseURI, username, password, userToken)

	response, err := httpClient.Do(request)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(response).NotTo(BeNil())

	responseBody, err := io.ReadAll(response.Body)
	Expect(err).ShouldNot(HaveOccurred())

	var jsonResponse map[string]interface{}
	err = json.Unmarshal(responseBody, &jsonResponse)
	Expect(err).NotTo(HaveOccurred())

	return response.StatusCode, jsonResponse
}
