package e2e_tests

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dana-team/platform-backend/src/controllers"
	"github.com/dana-team/platform-backend/src/utils/testutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate ServiceAccount token routes and functionality", func() {
	var namespaceName, serviceAccountName string

	BeforeEach(func() {
		namespaceName = generateName(e2eNamespace)
		createTestNamespace(k8sClient, namespaceName)

		serviceAccountName = generateName(testServiceAccountName)
		tokenRequestSecret := fmt.Sprintf("%s-token-request", serviceAccountName)
		createTestSecret(k8sClient, tokenRequestSecret, namespaceName)
		CreateTestServiceAccount(k8sClient, serviceAccountName, namespaceName, "token-secret", "value", "docker-cfg")
		Eventually(func() bool {
			tokenSecret := getServiceAccountTokenSecret(k8sClient, serviceAccountName, namespaceName)
			_, ok := tokenSecret.Data[testutils.TokenKey]
			return ok
		}, testutils.Timeout, testutils.Interval).Should(BeTrue())
	})

	Context("Validate create token route", func() {
		It("Should create a token for a ServiceAccount", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s/token", platformURL, namespaceName, testutils.ServiceAccountParam, serviceAccountName)
			status, response := performHTTPRequest(httpClient, nil, http.MethodPost, uri, "", "", userToken)

			Expect(status).Should(Equal(http.StatusOK))
			Expect(response[testutils.TokenKey]).ShouldNot(BeEmpty())
		})
		It("Should create a multiple tokens for a ServiceAccount", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s/token", platformURL, namespaceName, testutils.ServiceAccountParam, serviceAccountName)
			firstStatus, firstResponse := performHTTPRequest(httpClient, nil, http.MethodPost, uri, "", "", userToken)

			Expect(firstStatus).Should(Equal(http.StatusOK))
			Expect(firstResponse[testutils.TokenKey]).ShouldNot(BeEmpty())
			firstToken := firstResponse[testutils.TokenKey].(string)

			// Sleep for 1 second to ensure that the requests are seperated
			time.Sleep(1 * time.Second)

			secondStatus, secondResponse := performHTTPRequest(httpClient, nil, http.MethodPost, uri, "", "", userToken)
			Expect(secondStatus).Should(Equal(http.StatusOK))
			Expect(secondResponse[testutils.TokenKey]).ShouldNot(BeEmpty())
			secondToken := secondResponse[testutils.TokenKey].(string)

			Expect(firstToken).ShouldNot(Equal(secondToken))
		})

		It("Should create a token with a different expiration time depending on query params", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s/token?%s=%s", platformURL, namespaceName, testutils.ServiceAccountParam, serviceAccountName, testutils.ExpirationSecondsParam, "7200")
			status, response := performHTTPRequest(httpClient, nil, http.MethodPost, uri, "", "", userToken)

			Expect(status).Should(Equal(http.StatusOK))
			Expect(response[testutils.TokenKey]).ShouldNot(BeEmpty())
			laterExpiration, err := time.Parse(time.RFC3339, response[testutils.ExpiresKey].(string))
			Expect(err).ShouldNot(HaveOccurred())

			uri = fmt.Sprintf("%s/v1/namespaces/%s/%s/%s/token?%s=%s", platformURL, namespaceName, testutils.ServiceAccountParam, serviceAccountName, testutils.ExpirationSecondsParam, "3600")
			status, response = performHTTPRequest(httpClient, nil, http.MethodPost, uri, "", "", userToken)
			Expect(status).Should(Equal(http.StatusOK))
			Expect(response[testutils.TokenKey]).ShouldNot(BeEmpty())
			earlierExpiration, err := time.Parse(time.RFC3339, response[testutils.ExpiresKey].(string))
			Expect(err).ShouldNot(HaveOccurred())
			isActualAfterExpected := laterExpiration.After(earlierExpiration)
			Expect(isActualAfterExpected).Should(BeTrue())
		})

		It("Should handle an invalid string in expirationSeconds param", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s/token?%s=%s", platformURL, namespaceName, testutils.ServiceAccountParam, serviceAccountName, testutils.ExpirationSecondsParam, testutils.TestName)
			status, _ := performHTTPRequest(httpClient, nil, http.MethodPost, uri, "", "", userToken)

			Expect(status).Should(Equal(http.StatusBadRequest))
		})
		It("Should handle a not found ServiceAccount in a namespace", func() {
			tokenRequestSecret := fmt.Sprintf("%s-token-request", serviceAccountName+testutils.NonExistentSuffix)
			createTestSecret(k8sClient, tokenRequestSecret, namespaceName)
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s/token", platformURL, namespaceName, testutils.ServiceAccountParam, serviceAccountName+testutils.NonExistentSuffix)
			status, response := performHTTPRequest(httpClient, nil, http.MethodPost, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf(controllers.ErrCouldNotCreateTokenRequest, serviceAccountName+testutils.NonExistentSuffix, namespaceName),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})
	})

	Context("Validate revoke token route", func() {
		It("Should revoke a token for a ServiceAccount", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s/token", platformURL, namespaceName, testutils.ServiceAccountParam, serviceAccountName)
			status, response := performHTTPRequest(httpClient, nil, http.MethodDelete, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.MessageKey: fmt.Sprintf("Revoked tokens for ServiceAccount %q", serviceAccountName),
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})

		It("Should handle a not found ServiceAccount in a namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s/token", platformURL, namespaceName, testutils.ServiceAccountParam, serviceAccountName+testutils.NonExistentSuffix)
			status, response := performHTTPRequest(httpClient, nil, http.MethodDelete, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf(controllers.ErrCouldNotGetServiceAccount, serviceAccountName+testutils.NonExistentSuffix, namespaceName),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})

		It("Should handle a not found token request secret in a namespace", func() {
			newServiceAccount := generateName(testServiceAccountName)
			CreateTestServiceAccount(k8sClient, newServiceAccount, namespaceName, "token-secret", "value", "docker-cfg")
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s/token", platformURL, namespaceName, testutils.ServiceAccountParam, newServiceAccount)
			status, response := performHTTPRequest(httpClient, nil, http.MethodDelete, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf(controllers.ErrCouldNotGetTokenRequestSecret, newServiceAccount, namespaceName),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})
	})
})
