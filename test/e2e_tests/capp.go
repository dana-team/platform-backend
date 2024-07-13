package e2e_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/dana-team/platform-backend/src/utils/testutils"
	"github.com/dana-team/platform-backend/src/utils/testutils/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"net/http"
	"net/url"
)

var _ = Describe("Validate Capp routes and functionality", func() {
	var namespaceName, oneCappName, secondCappName string
	var namespace corev1.Namespace
	var oneLabelKey, oneLabelValue, secondLabelKey, secondLabelValue string

	BeforeEach(func() {
		namespaceName = generateName(e2eNamespace)
		namespace = mocks.PrepareNamespace(namespaceName, map[string]string{e2eLabelKey: e2eLabelValue})
		createResource(k8sClient, &namespace)

		oneCappName = generateName("a-" + testCappName)
		oneLabelKey = generateName(e2eLabelKey)
		oneLabelValue = generateName(e2eLabelValue)
		oneCappLabels := map[string]string{
			oneLabelKey: oneLabelValue,
		}
		oneCapp := mocks.PrepareCapp(oneCappName, namespaceName, clusterDomain, oneCappLabels, nil)
		createCapp(k8sClient, &oneCapp)

		secondCappName = generateName("b-" + testCappName)
		secondLabelKey = generateName(e2eLabelKey)
		secondLabelValue = generateName(e2eLabelValue)
		secondCappLabels := map[string]string{
			secondLabelKey: secondLabelValue,
		}
		secondCapp := mocks.PrepareCapp(secondCappName, namespaceName, clusterDomain, secondCappLabels, nil)
		createCapp(k8sClient, &secondCapp)
	})

	Context("Validate get Capps route", func() {
		It("Should get all Capps in a namespace", func() {
			baseURI := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.CappsKey)
			status, response := performAuthorizedHTTPRequest(httpClient, nil, http.MethodGet, baseURI, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.CappsKey: []types.CappSummary{
					{Name: oneCappName, Images: []string{CappImageName}, URL: fmt.Sprintf("https://%s-%s.%s", oneCappName, namespaceName, clusterDomain)},
					{Name: secondCappName, Images: []string{CappImageName}, URL: fmt.Sprintf("https://%s-%s.%s", secondCappName, namespaceName, clusterDomain)},
				},
				testutils.CountKey: 2,
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})

		It("Should get all Capps with a specific labelSelector in a namespace", func() {
			baseURI := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.CappsKey)
			params := url.Values{}
			params.Add(testutils.LabelSelectorKey, fmt.Sprintf("%s=%s", secondLabelKey, secondLabelValue))

			status, response := performAuthorizedHTTPRequest(httpClient, nil, http.MethodGet, fmt.Sprintf("%s?%s", baseURI, params.Encode()), "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.CappsKey: []types.CappSummary{
					{Name: secondCappName, Images: []string{CappImageName}, URL: fmt.Sprintf("https://%s-%s.%s", secondCappName, namespaceName, clusterDomain)},
				},
				testutils.CountKey: 1,
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})

		It("Should fail getting Capps with an invalid labelSelector", func() {
			baseURI := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.CappsKey)
			params := url.Values{}
			params.Add(testutils.LabelSelectorKey, fmt.Sprintf("%s=%s", secondLabelKey, secondLabelValue+" invalid"))

			status, response := performAuthorizedHTTPRequest(httpClient, nil, http.MethodGet, fmt.Sprintf("%s?%s", baseURI, params.Encode()), "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.DetailsKey: "found 'invalid', expected: ',' or 'end of string'",
				testutils.ErrorKey:   testutils.OperationFailed,
			}

			Expect(status).Should(Equal(http.StatusBadRequest))
			compareResponses(response, expectedResponse)
		})

		It("Should get no Capps with valid labelSelector", func() {
			baseURI := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.CappsKey)
			params := url.Values{}
			params.Add(testutils.LabelSelectorKey, fmt.Sprintf("%s=%s", secondLabelKey, secondLabelValue+testutils.NonExistentSuffix))

			status, response := performAuthorizedHTTPRequest(httpClient, nil, http.MethodGet, fmt.Sprintf("%s?%s", baseURI, params.Encode()), "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.CappsKey: nil,
				testutils.CountKey: 0,
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})
	})

	Context("Validate get Capp route", func() {
		It("Should get a specific Capp in a namespace", func() {
			baseURI := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.CappsKey, oneCappName)
			status, response := performAuthorizedHTTPRequest(httpClient, nil, http.MethodGet, baseURI, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.MetadataKey: types.Metadata{Name: oneCappName, Namespace: namespaceName},
				testutils.LabelsKey:   []types.KeyValue{{Key: oneLabelKey, Value: oneLabelValue}},
				testutils.SpecKey:     mocks.PrepareCappSpec(),
			}

			Expect(status).Should(Equal(http.StatusOK))
			Expect(response[testutils.MetadataKey], expectedResponse[testutils.MetadataKey])
			Expect(response[testutils.LabelsKey], expectedResponse[testutils.LabelsKey])
			Expect(response[testutils.SpecKey], expectedResponse[testutils.SpecKey])
		})

		It("Should handle a not found Capp in a namespace", func() {
			baseURI := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.CappsKey, oneCappName+testutils.NonExistentSuffix)
			status, response := performAuthorizedHTTPRequest(httpClient, nil, http.MethodGet, baseURI, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.DetailsKey: fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, oneCappName+testutils.NonExistentSuffix),
				testutils.ErrorKey:   testutils.OperationFailed,
			}

			capp := mocks.PrepareCapp(oneCappName+testutils.NonExistentSuffix, namespaceName, clusterDomain, nil, nil)
			Expect(doesResourceExist(k8sClient, &capp)).To(BeFalse())
			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})

		It("Should handle a get of a Capp in a not found namespace", func() {
			baseURI := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName+testutils.NonExistentSuffix, testutils.CappsKey, oneCappName)
			status, response := performAuthorizedHTTPRequest(httpClient, nil, http.MethodGet, baseURI, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.DetailsKey: fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, oneCappName),
				testutils.ErrorKey:   testutils.OperationFailed,
			}

			capp := mocks.PrepareCapp(oneCappName, namespaceName+testutils.NonExistentSuffix, clusterDomain, nil, nil)
			Expect(doesResourceExist(k8sClient, &capp)).To(BeFalse())
			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})
	})

	Context("Validate create Capp route", func() {
		It("Should create Capp in a namespace", func() {
			newCappName := generateName(testCappName)

			baseURI := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.CappsKey)
			requestData := mocks.PrepareCreateCappType(newCappName, []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue}}, nil)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performAuthorizedHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPost, baseURI, "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.MetadataKey:    types.Metadata{Name: newCappName, Namespace: namespaceName},
				testutils.LabelsKey:      []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue}},
				testutils.AnnotationsKey: nil,
				testutils.SpecKey:        mocks.PrepareCappSpec(),
				testutils.StatusKey:      cappv1alpha1.CappStatus{},
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)

			capp := mocks.PrepareCapp(oneCappName, namespaceName, clusterDomain, nil, nil)
			Eventually(func() bool {
				return doesResourceExist(k8sClient, &capp)
			}, testutils.Timeout, testutils.Interval).Should(BeTrue())
		})

		It("Should fail in creation of Capp with bad request body", func() {
			baseURI := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.CappsKey)
			requestData := mocks.PrepareCreateCappType("", []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue}}, nil)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performAuthorizedHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPost, baseURI, "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.DetailsKey: "Key: 'CreateCapp.Metadata.Name' Error:Field validation for 'Name' failed on the 'required' tag",
				testutils.ErrorKey:   testutils.InvalidRequest,
			}

			Expect(status).Should(Equal(http.StatusBadRequest))
			compareResponses(response, expectedResponse)
		})

		It("Should handle already existing Capp on creation", func() {
			baseURI := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.CappsKey)
			requestData := mocks.PrepareCreateCappType(oneCappName, []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue}}, nil)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performAuthorizedHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPost, baseURI, "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.DetailsKey: fmt.Sprintf("%s.%s %q already exists", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, oneCappName),
				testutils.ErrorKey:   testutils.OperationFailed,
			}

			Expect(status).Should(Equal(http.StatusConflict))
			compareResponses(response, expectedResponse)
		})
	})

	Context("Validate update Capp route", func() {
		It("Should update an existing Capp in a namespace", func() {
			baseURI := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.CappsKey, oneCappName)
			requestData := mocks.PrepareCreateCappType(oneCappName, []types.KeyValue{{Key: oneLabelKey, Value: oneLabelValue + "-updated"}}, nil)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performAuthorizedHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPut, baseURI, "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.MetadataKey: types.Metadata{Name: oneCappName, Namespace: namespaceName},
				testutils.LabelsKey:   []types.KeyValue{{Key: oneLabelKey, Value: oneLabelValue + "-updated"}},
				testutils.SpecKey:     mocks.PrepareCappSpec(),
			}

			Expect(status).Should(Equal(http.StatusOK))
			Expect(response[testutils.MetadataKey], expectedResponse[testutils.MetadataKey])
			Expect(response[testutils.LabelsKey], expectedResponse[testutils.LabelsKey])
			Expect(response[testutils.SpecKey], expectedResponse[testutils.SpecKey])

			Eventually(func() bool {
				capp := getCapp(k8sClient, oneCappName, namespaceName)
				return capp.Labels[oneLabelKey] == oneLabelValue+"-updated"
			}, testutils.Timeout, testutils.Interval).Should(BeTrue())
		})

		It("Should handle update of not found Capp", func() {
			baseURI := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.CappsKey, oneCappName+testutils.NonExistentSuffix)
			requestData := mocks.PrepareCreateCappType(oneCappName+testutils.NonExistentSuffix, []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue + "-updated"}}, nil)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performAuthorizedHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPut, baseURI, "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.DetailsKey: fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, oneCappName+testutils.NonExistentSuffix),
				testutils.ErrorKey:   testutils.OperationFailed,
			}

			capp := mocks.PrepareCapp(oneCappName+testutils.NonExistentSuffix, namespaceName, clusterDomain, nil, nil)
			Expect(doesResourceExist(k8sClient, &capp)).To(BeFalse())
			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})

		It("Should handle update of not found namespace", func() {
			baseURI := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName+testutils.NonExistentSuffix, testutils.CappsKey, oneCappName)
			requestData := mocks.PrepareCreateCappType(oneCappName, []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue + "-updated"}}, nil)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performAuthorizedHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPut, baseURI, "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.DetailsKey: fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, oneCappName),
				testutils.ErrorKey:   testutils.OperationFailed,
			}

			capp := mocks.PrepareCapp(oneCappName, namespaceName+testutils.NonExistentSuffix, clusterDomain, nil, nil)
			Expect(doesResourceExist(k8sClient, &capp)).To(BeFalse())
			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})
	})

	Context("Validate delete Capp route", func() {
		It("Should delete Capp from namespace", func() {
			baseURI := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.CappsKey, oneCappName)
			status, response := performAuthorizedHTTPRequest(httpClient, nil, http.MethodDelete, baseURI, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.MessageKey: fmt.Sprintf("Deleted capp %q in namespace %q successfully", oneCappName, namespaceName),
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)

			capp := mocks.PrepareCapp(oneCappName, namespaceName, clusterDomain, nil, nil)
			Eventually(func() bool {
				return doesResourceExist(k8sClient, &capp)
			}, testutils.Timeout, testutils.Interval).Should(BeFalse())
		})

		It("Should handle deletion of not found Capp", func() {
			baseURI := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.CappsKey, oneCappName+testutils.NonExistentSuffix)
			status, response := performAuthorizedHTTPRequest(httpClient, nil, http.MethodDelete, baseURI, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.DetailsKey: fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, oneCappName+testutils.NonExistentSuffix),
				testutils.ErrorKey:   testutils.OperationFailed,
			}

			capp := mocks.PrepareCapp(oneCappName+testutils.NonExistentSuffix, namespaceName, clusterDomain, nil, nil)
			Expect(doesResourceExist(k8sClient, &capp)).To(BeFalse())
			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})

		It("Should handle deletion of Capp in a not found namespace", func() {
			baseURI := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName+testutils.NonExistentSuffix, testutils.CappsKey, oneCappName)
			status, response := performAuthorizedHTTPRequest(httpClient, nil, http.MethodDelete, baseURI, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.DetailsKey: fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, oneCappName),
				testutils.ErrorKey:   testutils.OperationFailed,
			}

			capp := mocks.PrepareCapp(oneCappName, namespaceName+testutils.NonExistentSuffix, clusterDomain, nil, nil)
			Expect(doesResourceExist(k8sClient, &capp)).To(BeFalse())
			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})
	})
})
