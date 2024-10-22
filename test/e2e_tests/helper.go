package e2e_tests

import (
	"encoding/json"
	"fmt"
	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	"github.com/dana-team/platform-backend/internal/utils/testutils"
	rcsv1alpha1 "github.com/dana-team/rcs-ocm-deployer/api/v1alpha1"
	"github.com/google/go-cmp/cmp/cmpopts"
	. "github.com/onsi/gomega"
	configv1 "github.com/openshift/api/config/v1"
	userv1 "github.com/openshift/api/user/v1"
	"io"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"net/http"
	clusterv1beta1 "open-cluster-management.io/api/cluster/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	e2eUser                = "e2e-user"
	e2ePassword            = "e2e-password"
	e2eNamespace           = "e2e-namespace"
	e2eLabelKey            = "e2e-label-key"
	e2eLabelValue          = "e2e-label-value"
	testCappName           = "e2e-capp"
	testCappRevisionName   = "e2e-capp-revision"
	testConfigMapName      = "e2e-configmap"
	testServiceAccountName = "e2e-serviceaccount"
	testSecretName         = "e2e-secret"
	testUserName           = "e2e-user"
	testEnvironment        = "e2e-environment"
	testRegionName         = "e2e-region"
	testPlacementName      = "e2e-placement"

	rcsConfigName      = "rcs-config"
	rcsConfigNamespace = "rcs-deployer-system"
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
	clusterDomain string
	platformURL   string
	userToken     string
	placementName string
	placementNS   string
)

// newScheme initializes a new scheme by adding the necessary schemes to it.
func newScheme() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(cappv1alpha1.AddToScheme(scheme))
	utilruntime.Must(configv1.AddToScheme(scheme))
	utilruntime.Must(userv1.AddToScheme(scheme))
	utilruntime.Must(clusterv1beta1.AddToScheme(scheme))
	utilruntime.Must(rcsv1alpha1.AddToScheme(scheme))
}

// compareResponses compares two responses and asserts that they are equal.
func compareResponses(response, expectedResponse map[string]interface{}) {
	expectedResponseJSON, err := json.Marshal(expectedResponse)
	Expect(err).ShouldNot(HaveOccurred())

	var expectedResponseNormalized map[string]interface{}
	err = json.Unmarshal(expectedResponseJSON, &expectedResponseNormalized)
	Expect(err).ShouldNot(HaveOccurred())

	compareError(expectedResponseNormalized, response)
	less := func(a, b string) bool { return a < b }
	Expect(response).Should(BeComparableTo(expectedResponseNormalized, cmpopts.SortSlices(less)))
}

// compareError compares two errors and asserts that the response contains the expected response.
func compareError(expectedResponse, response map[string]interface{}) {
	expectedError, expectedResponseHasError := expectedResponse[testutils.ErrorKey]
	responseError, responseHasError := response[testutils.ErrorKey]

	if !expectedResponseHasError {
		Expect(responseHasError).Should(Equal(false), fmt.Sprintf("unexpected error: %s", responseError))
		return
	}

	Expect(responseError).To(ContainSubstring(expectedError.(string)), "error %q does not contain expected error %q", responseError, expectedError)
	delete(response, "error")
	delete(expectedResponse, "error")
}

// prepareAuthorizedHTTPRequest prepares an HTTP request.
// It creates a new HTTP request, sets the content type, and adds authorization headers as needed.
func prepareAuthorizedHTTPRequest(body io.Reader, httpMethod, baseURI, username, password, userToken string) *http.Request {
	request, err := http.NewRequest(httpMethod, baseURI, body)
	Expect(err).NotTo(HaveOccurred())

	request.Header.Set(contentType, applicationJson)
	addAuthorization(request, username, password, userToken)

	return request
}

// addAuthorization adds the appropriate authorization headers to the HTTP request.
// It adds Basic Auth headers if the username is provided, or Bearer token headers if the user token is provided.
func addAuthorization(request *http.Request, username, password, userToken string) {
	if username != "" {
		request.SetBasicAuth(username, password)
	} else if userToken != "" {
		request.Header.Set(httpAuthorizationHeader, fmt.Sprintf("%s %s", httpBearerToken, userToken))
	}
}

// performHTTPRequest makes an HTTP request and returns a response.
func performHTTPRequest(httpClient http.Client, body io.Reader, httpMethod, baseURI, username, password, userToken string) (int, map[string]interface{}) {
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
