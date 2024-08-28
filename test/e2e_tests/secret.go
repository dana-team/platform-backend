package e2e_tests

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dana-team/platform-backend/src/types"
	"github.com/dana-team/platform-backend/src/utils/testutils"
	"github.com/dana-team/platform-backend/src/utils/testutils/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
)

var _ = Describe("Validate Secret routes and functionality", func() {
	var namespaceName, oneSecretName, secondSecretName string

	BeforeEach(func() {
		namespaceName = generateName(e2eNamespace)
		createTestNamespace(k8sClient, namespaceName)

		oneSecretName = generateName("a-" + testSecretName)
		createTestSecret(k8sClient, oneSecretName, namespaceName)

		secondSecretName = generateName("b-" + testSecretName)
		createTestSecret(k8sClient, secondSecretName, namespaceName)
	})

	Context("Validate get Secrets route", func() {
		It("Should get only secrets with a managed label in a namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.SecretsKey)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.CountKey: 2,
				testutils.SecretsKey: []map[string]interface{}{
					{
						testutils.SecretNameKey:    oneSecretName,
						testutils.NamespaceNameKey: namespaceName,
						testutils.TypeKey:          string(corev1.SecretTypeOpaque),
					},
					{
						testutils.SecretNameKey:    secondSecretName,
						testutils.NamespaceNameKey: namespaceName,
						testutils.TypeKey:          string(corev1.SecretTypeOpaque),
					},
				},
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})

		It("Should get all secrets in a namespace with limit of 50", func() {
			limit := "50"
			page := "1"

			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s?limit=%s&page=%s", platformURL, namespaceName, testutils.SecretsKey, limit, page)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.CountKey: 2,
				testutils.SecretsKey: []map[string]interface{}{
					{
						testutils.SecretNameKey:    oneSecretName,
						testutils.NamespaceNameKey: namespaceName,
						testutils.TypeKey:          string(corev1.SecretTypeOpaque),
					},
					{
						testutils.SecretNameKey:    secondSecretName,
						testutils.NamespaceNameKey: namespaceName,
						testutils.TypeKey:          string(corev1.SecretTypeOpaque),
					},
				},
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})

		It("Should get one secret in a namespace with limit of 1 and page 1", func() {
			limit := "1"
			page := "1"

			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s?limit=%s&page=%s", platformURL, namespaceName, testutils.SecretsKey, limit, page)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.CountKey: 1,
				testutils.SecretsKey: []map[string]interface{}{
					{
						testutils.SecretNameKey:    oneSecretName,
						testutils.NamespaceNameKey: namespaceName,
						testutils.TypeKey:          string(corev1.SecretTypeOpaque),
					},
				},
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})

		It("Should get one secret in a namespace with limit of 1 and page 2", func() {
			limit := "1"
			page := "2"

			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s?limit=%s&page=%s", platformURL, namespaceName, testutils.SecretsKey, limit, page)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.CountKey: 1,
				testutils.SecretsKey: []map[string]interface{}{
					{
						testutils.SecretNameKey:    secondSecretName,
						testutils.NamespaceNameKey: namespaceName,
						testutils.TypeKey:          string(corev1.SecretTypeOpaque),
					},
				},
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})

		It("Should not get secrets with limit of 1 and page 3", func() {
			limit := "1"
			page := "3"

			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s?limit=%s&page=%s", platformURL, namespaceName, testutils.SecretsKey, limit, page)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.CountKey:   0,
				testutils.SecretsKey: nil,
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})
	})

	Context("Validate get Secret route", func() {
		It("Should get a specific Secret in a namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.SecretsKey, oneSecretName)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.SecretNameKey: oneSecretName,
				testutils.TypeKey:       string(corev1.SecretTypeOpaque),
				testutils.DataKey:       []types.KeyValue{{Key: testutils.SecretDataKey, Value: testutils.SecretDataValue}},
			}

			Expect(status).Should(Equal(http.StatusOK))
			Expect(response[testutils.SecretNameKey]).To(Equal(expectedResponse[testutils.SecretNameKey]))
			Expect(response[testutils.TypeKey]).To(Equal(expectedResponse[testutils.TypeKey]))

			var responseKeyVal []types.KeyValue
			for _, data := range response[testutils.DataKey].([]interface{}) {
				pair := data.(map[string]interface{})
				responseKeyVal = append(responseKeyVal, types.KeyValue{Key: pair[testutils.KeyField].(string), Value: pair[testutils.ValueField].(string)})
			}

			Expect(responseKeyVal).To(Equal(expectedResponse[testutils.DataKey]))
		})

		It("Should handle a not found Secret in a namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.SecretsKey, oneSecretName+testutils.NonExistentSuffix)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf("%s %q not found", testutils.SecretsKey, oneSecretName+testutils.NonExistentSuffix),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			secret := mocks.PrepareSecret(oneSecretName+testutils.NonExistentSuffix, namespaceName, "", "")
			Expect(doesResourceExist(k8sClient, &secret)).To(BeFalse())
			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})

		It("Should handle a get of a ConfigMap in a not found namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName+testutils.NonExistentSuffix, testutils.SecretsKey, oneSecretName)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf("%s %q not found", testutils.SecretsKey, oneSecretName),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			secret := mocks.PrepareSecret(oneSecretName, namespaceName+testutils.NonExistentSuffix, "", "")
			Expect(doesResourceExist(k8sClient, &secret)).To(BeFalse())
			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})
	})

	Context("Validate create Secret route", func() {
		It("Should create a new secret in a namespace", func() {
			newSecretName := generateName(testSecretName)

			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.SecretsKey)
			requestData := mocks.PrepareCreateSecretRequestType(newSecretName, strings.ToLower(string(corev1.SecretTypeOpaque)), "", "", []types.KeyValue{{Key: testutils.SecretDataKey, Value: testutils.SecretDataValue}})
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPost, uri, "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.SecretNameKey:    newSecretName,
				testutils.TypeKey:          corev1.SecretTypeOpaque,
				testutils.NamespaceNameKey: namespaceName,
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)

			secret := mocks.PrepareSecret(newSecretName, namespaceName, "", "")
			Eventually(func() bool {
				return doesResourceExist(k8sClient, &secret)
			}, testutils.Timeout, testutils.Interval).Should(BeTrue())

			Eventually(func() bool {
				secret := getSecret(k8sClient, newSecretName, namespaceName)
				return secret.Type == corev1.SecretTypeOpaque
			}, testutils.Timeout, testutils.Interval).Should(BeTrue())
		})

		It("Should fail creating secret with invalid body", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.SecretsKey)
			requestData := mocks.PrepareCreateSecretRequestType("", strings.ToLower(string(corev1.SecretTypeOpaque)), "", "", []types.KeyValue{{Key: testutils.SecretDataKey, Value: testutils.SecretDataValue}})
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPost, uri, "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  "Key: 'CreateSecretRequest.SecretName' Error:Field validation for 'SecretName' failed on the 'required' tag",
				testutils.ReasonKey: testutils.ReasonBadRequest,
			}

			Expect(status).Should(Equal(http.StatusBadRequest))
			compareResponses(response, expectedResponse)
		})

		It("Should handle already existing secret", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.SecretsKey)
			requestData := mocks.PrepareCreateSecretRequestType(oneSecretName, strings.ToLower(string(corev1.SecretTypeOpaque)), "", "", []types.KeyValue{{Key: testutils.SecretDataKey, Value: testutils.SecretDataValue}})
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPost, uri, "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf("%s %q already exists", testutils.SecretsKey, oneSecretName),
				testutils.ReasonKey: testutils.ReasonAlreadyExists,
			}

			Expect(status).Should(Equal(http.StatusConflict))
			compareResponses(response, expectedResponse)
		})
	})

	Context("Validate update Secret route", func() {
		It("Should update an existing Secret in a namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.SecretsKey, oneSecretName)
			requestData := mocks.PrepareSecretRequestType([]types.KeyValue{{Key: testutils.SecretDataKey, Value: testutils.SecretDataValue + "-updated"}})
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPut, uri, "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.SecretNameKey:    oneSecretName,
				testutils.TypeKey:          corev1.SecretTypeOpaque,
				testutils.NamespaceNameKey: namespaceName,
				testutils.DataKey:          []types.KeyValue{{Key: testutils.SecretDataKey, Value: testutils.SecretDataValue + "-updated"}},
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)

			updatedValue := base64.StdEncoding.EncodeToString([]byte(testutils.SecretDataValue + "-updated"))
			Eventually(func() bool {
				secret := getSecret(k8sClient, oneSecretName, namespaceName)
				return string(secret.Data[testutils.SecretDataKey]) == updatedValue
			}, testutils.Timeout, testutils.Interval).Should(BeTrue())
		})

		It("Should handle update of a not found Secret in a namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.SecretsKey, oneSecretName+testutils.NonExistentSuffix)
			requestData := mocks.PrepareSecretRequestType([]types.KeyValue{{Key: testutils.SecretDataKey, Value: testutils.SecretDataValue + "-updated"}})
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPut, uri, "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf("%s %q not found", testutils.SecretsKey, oneSecretName+testutils.NonExistentSuffix),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})

		It("Should handle update of a Secret in a not found namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName+testutils.NonExistentSuffix, testutils.SecretsKey, oneSecretName)
			requestData := mocks.PrepareSecretRequestType([]types.KeyValue{{Key: testutils.SecretDataKey, Value: testutils.SecretDataValue + "-updated"}})
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPut, uri, "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf("%s %q not found", testutils.SecretsKey, oneSecretName),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})

		It("Should fail to update a Secret with invalid body", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.SecretsKey, oneSecretName)
			requestData := mocks.PrepareSecretRequestType(nil)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPut, uri, "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  "Key: 'UpdateSecretRequest.Data' Error:Field validation for 'Data' failed on the 'required' tag",
				testutils.ReasonKey: testutils.ReasonBadRequest,
			}

			Expect(status).Should(Equal(http.StatusBadRequest))
			compareResponses(response, expectedResponse)
		})
	})

	Context("Validate delete Secret route", func() {
		It("Should delete a Secret in a namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.SecretsKey, oneSecretName)
			status, response := performHTTPRequest(httpClient, nil, http.MethodDelete, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.MessageKey: fmt.Sprintf("Deleted secret %q in namespace %q successfully", oneSecretName, namespaceName),
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)

			secret := mocks.PrepareSecret(oneSecretName, namespaceName, "", "")
			Eventually(func() bool {
				return doesResourceExist(k8sClient, &secret)
			}, testutils.Timeout, testutils.Interval).Should(BeFalse())
		})

		It("Should handle deletion of a not found Secret in a namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.SecretsKey, oneSecretName+testutils.NonExistentSuffix)
			status, response := performHTTPRequest(httpClient, nil, http.MethodDelete, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf("%s %q not found", testutils.SecretsKey, oneSecretName+testutils.NonExistentSuffix),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})

		It("Should handle deletion of a Secret in a not found namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName+testutils.NonExistentSuffix, testutils.SecretsKey, oneSecretName)
			status, response := performHTTPRequest(httpClient, nil, http.MethodDelete, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf("%s %q not found", testutils.SecretsKey, oneSecretName),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})
	})
})
