package e2e_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dana-team/platform-backend/internal/utils/testutils"
	"github.com/dana-team/platform-backend/internal/utils/testutils/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate Namespace routes and functionality", func() {
	var oneNamespaceName, secondNamespaceName string

	firstNamespaceSubstring := "a-"
	secondNamespaceSubstring := "b-"

	BeforeEach(func() {
		oneNamespaceName = generateName(firstNamespaceSubstring + e2eNamespace)
		createTestNamespace(k8sClient, oneNamespaceName)

		secondNamespaceName = generateName(secondNamespaceSubstring + e2eNamespace)
		createTestNamespace(k8sClient, secondNamespaceName)
	})

	Context("Validate get Namespaces route", func() {
		It("Should get namespaces", func() {
			uri := fmt.Sprintf("%s/v1/%s", platformURL, testutils.NamespaceKey)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.CountKey: 2,
				testutils.NamespaceKey: []map[string]interface{}{
					{
						testutils.NameKey: oneNamespaceName,
					},
					{
						testutils.NameKey: secondNamespaceName,
					},
				},
			}

			Expect(status).Should(Equal(http.StatusOK))
			Expect(expectedResponse[testutils.CountKey]).To(BeNumerically("<=", response[testutils.CountKey]))
			for _, ns := range expectedResponse[testutils.NamespaceKey].([]map[string]interface{}) {
				Expect(ns).To(BeElementOf(response[testutils.NamespaceKey]))
			}
		})

		It("Should get all namespaces in a namespace with limit of 50", func() {
			limit := "50"
			page := "1"

			uri := fmt.Sprintf("%s/v1/%s?limit=%s&page=%s", platformURL, testutils.NamespaceKey, limit, page)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.CountKey: 2,
				testutils.NamespaceKey: []map[string]interface{}{
					{
						testutils.NameKey: oneNamespaceName,
					},
					{
						testutils.NameKey: secondNamespaceName,
					},
				},
			}

			Expect(status).Should(Equal(http.StatusOK))
			Expect(expectedResponse[testutils.CountKey]).To(BeNumerically("<=", response[testutils.CountKey]))
			for _, ns := range expectedResponse[testutils.NamespaceKey].([]map[string]interface{}) {
				Expect(ns).To(BeElementOf(response[testutils.NamespaceKey]))
			}
		})

		It("Should get one namespace in a namespace with limit of 1 and page 1", func() {
			limit := "1"
			page := "1"

			uri := fmt.Sprintf("%s/v1/%s?limit=%s&page=%s", platformURL, testutils.NamespaceKey, limit, page)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.CountKey: 1,
				testutils.NamespaceKey: []map[string]interface{}{
					{
						testutils.NameKey: oneNamespaceName,
					},
				},
			}

			Expect(status).Should(Equal(http.StatusOK))
			Expect(expectedResponse[testutils.CountKey]).To(BeNumerically("==", response[testutils.CountKey]))
		})

		It("Should get one namespaces in a namespace with limit of 1 and page 2", func() {
			limit := "1"
			page := "2"

			uri := fmt.Sprintf("%s/v1/%s?limit=%s&page=%s", platformURL, testutils.NamespaceKey, limit, page)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.CountKey: 1,
				testutils.NamespaceKey: []map[string]interface{}{
					{
						testutils.NameKey: secondNamespaceName,
					},
				},
			}

			Expect(status).Should(Equal(http.StatusOK))
			Expect(expectedResponse[testutils.CountKey]).To(BeNumerically("<=", response[testutils.CountKey]))
		})

		It("Should not get namespaces with limit of 1 and page 3", func() {
			limit := "1"
			page := "3"

			uri := fmt.Sprintf("%s/v1/%s?limit=%s&page=%s", platformURL, testutils.NamespaceKey, limit, page)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			Expect(status).Should(Equal(http.StatusOK))
			Expect(testutils.KubeSystem).ToNot(BeElementOf(response[testutils.NamespaceKey]))
		})

		It("Should not get namespaces without a Managed label", func() {
			uri := fmt.Sprintf("%s/v1/%s", platformURL, testutils.NamespaceKey)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			Expect(status).Should(Equal(http.StatusOK))
			Expect(testutils.KubeSystem).ToNot(BeElementOf(response[testutils.NamespaceKey]))
		})
	})

	Context("Validate get Namespace route", func() {
		It("Should get a specific namespace", func() {
			uri := fmt.Sprintf("%s/v1/%s/%s", platformURL, testutils.NamespaceKey, oneNamespaceName)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.NameKey: oneNamespaceName,
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})

		It("Should handle getting a not found namespace", func() {
			uri := fmt.Sprintf("%s/v1/%s/%s", platformURL, testutils.NamespaceKey, oneNamespaceName+testutils.NonExistentSuffix)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf("%s %q not found", testutils.NamespaceKey, oneNamespaceName+testutils.NonExistentSuffix),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			ns := mocks.PrepareNamespace(oneNamespaceName+testutils.NonExistentSuffix, map[string]string{})
			Expect(doesResourceExist(k8sClient, &ns)).To(BeFalse())
			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})
	})

	Context("Validate create Namespace route", func() {
		It("Should create a new namespace", func() {
			newNamespaceName := generateName(e2eNamespace)

			uri := fmt.Sprintf("%s/v1/%s", platformURL, testutils.NamespaceKey)
			requestData := mocks.PrepareNamespaceType(newNamespaceName)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPost, uri, "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.NameKey: newNamespaceName,
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)

			ns := mocks.PrepareNamespace(newNamespaceName, map[string]string{})
			Eventually(func() bool {
				return doesResourceExist(k8sClient, &ns)
			}, testutils.Timeout, testutils.Interval).Should(BeTrue())

			deleteResource(k8sClient, &ns)
		})

		It("Should fail creating a new namespace with invalid body", func() {
			uri := fmt.Sprintf("%s/v1/%s", platformURL, testutils.NamespaceKey)
			requestData := mocks.PrepareNamespaceType("")
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPost, uri, "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  "Key: 'Namespace.Name' Error:Field validation for 'Name' failed on the 'required' tag",
				testutils.ReasonKey: testutils.ReasonBadRequest,
			}

			Expect(status).Should(Equal(http.StatusBadRequest))
			compareResponses(response, expectedResponse)
		})

		It("Should handle already existing Namespace on creation", func() {
			uri := fmt.Sprintf("%s/v1/%s", platformURL, testutils.NamespaceKey)
			requestData := mocks.PrepareNamespaceType(oneNamespaceName)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPost, uri, "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf("%s %q already exists", testutils.NamespaceKey, oneNamespaceName),
				testutils.ReasonKey: testutils.ReasonAlreadyExists,
			}

			Expect(status).Should(Equal(http.StatusConflict))
			compareResponses(response, expectedResponse)
		})
	})

	Context("Validate delete Namespace route", func() {
		It("Should delete namespace", func() {
			uri := fmt.Sprintf("%s/v1/%s/%s", platformURL, testutils.NamespaceKey, oneNamespaceName)
			status, response := performHTTPRequest(httpClient, nil, http.MethodDelete, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.MessageKey: fmt.Sprintf("Deleted namespace successfully %q", oneNamespaceName),
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)

			ns := mocks.PrepareNamespace(oneNamespaceName, map[string]string{})
			Eventually(func() bool {
				return doesResourceExist(k8sClient, &ns)
			}, testutils.Timeout, testutils.Interval).Should(BeFalse())
		})

		It("Should handle deletion of not found namespace", func() {
			uri := fmt.Sprintf("%s/v1/%s/%s", platformURL, testutils.NamespaceKey, oneNamespaceName+testutils.NonExistentSuffix)
			status, response := performHTTPRequest(httpClient, nil, http.MethodDelete, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf("%s %q not found", testutils.NamespaceKey, oneNamespaceName+testutils.NonExistentSuffix),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			ns := mocks.PrepareNamespace(oneNamespaceName+testutils.NonExistentSuffix, map[string]string{})
			Expect(doesResourceExist(k8sClient, &ns)).To(BeFalse())
			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})
	})
})
