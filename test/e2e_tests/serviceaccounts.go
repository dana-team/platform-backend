package e2e_tests

import (
	"fmt"
	"github.com/dana-team/platform-backend/src/controllers"
	"github.com/dana-team/platform-backend/src/utils/testutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = Describe("Validate ServiceAccounts routes and functionality", func() {
	var namespaceName, serviceAccountName string

	BeforeEach(func() {
		namespaceName = generateName(e2eNamespace)
		createTestNamespace(k8sClient, namespaceName)

		serviceAccountName = generateName(testServiceAccountName)
		CreateTestServiceAccount(k8sClient, serviceAccountName, namespaceName, "token-secret", "value", "docker-cfg")
	})

	Context("Validate get ServiceAccount token route", func() {
		It("Should get token from service account", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s/token", platformURL, namespaceName, testutils.ServiceAccountsKey, serviceAccountName)
			fmt.Println(uri)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			tokenSecret := getServiceAccountTokenSecret(k8sClient, serviceAccountName, namespaceName)
			expectedResponse := map[string]interface{}{
				testutils.TokenKey: string(tokenSecret.Data[testutils.TokenKey]),
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})

		It("Should handle a not found ServiceAccount in a namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s/token", platformURL, namespaceName, testutils.ServiceAccountsKey, serviceAccountName+testutils.NonExistentSuffix)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf(controllers.ErrCouldNotGetServiceAccount, serviceAccountName+testutils.NonExistentSuffix, namespaceName),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})

		It("Should handle a get of a ServiceAccount in a not found namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s/token", platformURL, namespaceName+testutils.NonExistentSuffix, testutils.ServiceAccountsKey, serviceAccountName)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf(controllers.ErrCouldNotGetServiceAccount, serviceAccountName, namespaceName+testutils.NonExistentSuffix),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})
	})
})
