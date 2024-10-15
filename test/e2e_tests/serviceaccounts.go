package e2e_tests

import (
	"fmt"
	"net/http"

	"github.com/dana-team/platform-backend/src/controllers"
	"github.com/dana-team/platform-backend/src/utils/testutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate ServiceAccounts routes and functionality", func() {
	var namespaceName, serviceAccountName, secondServiceAccountName string

	BeforeEach(func() {
		namespaceName = generateName(e2eNamespace)
		createTestNamespace(k8sClient, namespaceName)

		serviceAccountName = "a-" + generateName(testServiceAccountName)
		secondServiceAccountName = "b-" + generateName(testServiceAccountName)
		CreateTestServiceAccount(k8sClient, serviceAccountName, namespaceName, "token-secret", "value", "docker-cfg")
		CreateTestServiceAccount(k8sClient, secondServiceAccountName, namespaceName, "token-secret", "value", "")
		Eventually(func() bool {
			tokenSecret := getServiceAccountTokenSecret(k8sClient, serviceAccountName, namespaceName)
			_, ok := tokenSecret.Data[testutils.TokenKey]
			return ok
		}, testutils.Timeout, testutils.Interval).Should(BeTrue())
	})

	Context("Validate get ServiceAccount token route", func() {
		It("Should get token from service account", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s/token", platformURL, namespaceName, testutils.ServiceAccountParam, serviceAccountName)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)
			tokenSecret := getServiceAccountTokenSecret(k8sClient, serviceAccountName, namespaceName)
			expectedResponse := map[string]interface{}{
				testutils.TokenKey: string(tokenSecret.Data[testutils.TokenKey]),
			}
			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})

		It("Should handle a not found ServiceAccount in a namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s/token", platformURL, namespaceName, testutils.ServiceAccountParam, serviceAccountName+testutils.NonExistentSuffix)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf(controllers.ErrCouldNotGetServiceAccount, serviceAccountName+testutils.NonExistentSuffix, namespaceName),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})

		It("Should handle a get of a ServiceAccount token in a not found namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s/token", platformURL, namespaceName+testutils.NonExistentSuffix, testutils.ServiceAccountParam, serviceAccountName)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf(controllers.ErrCouldNotGetServiceAccount, serviceAccountName, namespaceName+testutils.NonExistentSuffix),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})

		It("Should handle a get of a ServiceAccount", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.ServiceAccountParam, serviceAccountName)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.NameKey: serviceAccountName,
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})

		It("Should handle a get of a non-existing ServiceAccount", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.ServiceAccountParam, serviceAccountName+testutils.NonExistentSuffix)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf(controllers.ErrCouldNotGetServiceAccount, serviceAccountName+testutils.NonExistentSuffix, namespaceName),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})

		It("Should handle the creation of a ServiceAccount", func() {
			newServiceAccountName := generateName(testServiceAccountName)
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.ServiceAccountParam, newServiceAccountName)
			status, response := performHTTPRequest(httpClient, nil, http.MethodPost, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.NameKey: newServiceAccountName,
			}
			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})

		It("Should handle the deletion of a ServiceAccount", func() {
			newServiceAccount := generateName(testServiceAccountName)
			CreateTestServiceAccount(k8sClient, newServiceAccount, namespaceName, "token-secret", "value", "")
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.ServiceAccountParam, newServiceAccount)
			status, response := performHTTPRequest(httpClient, nil, http.MethodDelete, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.MessageKey: fmt.Sprintf("Deleted serviceAccount successfully %q", newServiceAccount),
			}
			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})

		It("Should handle the deletion of a non-existing ServiceAccount", func() {
			nonExistingServiceAccount := generateName(testServiceAccountName)
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.ServiceAccountParam, nonExistingServiceAccount)
			status, response := performHTTPRequest(httpClient, nil, http.MethodDelete, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.ReasonKey: testutils.ReasonNotFound,
				testutils.ErrorKey:  fmt.Sprintf(controllers.ErrCouldNotDeleteServiceAccount, nonExistingServiceAccount, namespaceName),
			}
			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})
	})
	Context("Validate get ServiceAccounts route", func() {
		It("Should get only serviceAccounts with a managed label in a namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.ServiceAccountParam)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.CountKey: 2,
				testutils.ServiceAccountsKey: []string{
					serviceAccountName,
					secondServiceAccountName,
				},
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})

		It("Should get all secretAccounts in a namespace with limit of 50", func() {
			limit := "50"
			page := "1"

			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s?limit=%s&page=%s", platformURL, namespaceName, testutils.ServiceAccountParam, limit, page)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.CountKey: 2,
				testutils.ServiceAccountsKey: []string{
					serviceAccountName,
					secondServiceAccountName,
				},
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})

		It("Should get one ServiceAccount in a namespace with limit of 1 and page 1", func() {
			limit := "1"
			page := "1"

			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s?limit=%s&page=%s", platformURL, namespaceName, testutils.ServiceAccountParam, limit, page)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.CountKey: 1,
				testutils.ServiceAccountsKey: []string{
					serviceAccountName,
				},
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})

		It("Should get one ServiceAccount in a namespace with limit of 1 and page 2", func() {
			limit := "1"
			page := "2"

			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s?limit=%s&page=%s", platformURL, namespaceName, testutils.ServiceAccountParam, limit, page)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.CountKey: 1,
				testutils.ServiceAccountsKey: []string{
					secondServiceAccountName,
				},
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})

		It("Should not get ServiceAccount with limit of 1 and page 3", func() {
			limit := "1"
			page := "3"

			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s?limit=%s&page=%s", platformURL, namespaceName, testutils.ServiceAccountParam, limit, page)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.CountKey:           0,
				testutils.ServiceAccountsKey: nil,
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})
	})
})
