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
	"net/http"
	"net/url"
)

var _ = Describe("Validate Capp routes and functionality", func() {
	var namespaceName, oneCappName, secondCappName string
	var oneLabelKey, oneLabelValue, secondLabelKey, secondLabelValue string

	BeforeEach(func() {
		namespaceName = generateName(e2eNamespace)
		createTestNamespace(k8sClient, namespaceName)

		oneCappName = generateName("a-" + testCappName)
		oneLabelKey = generateName(e2eLabelKey)
		oneLabelValue = generateName(e2eLabelValue)
		createTestCapp(k8sClient, oneCappName, namespaceName, map[string]string{oneLabelKey: oneLabelValue}, nil)

		secondCappName = generateName("b-" + testCappName)
		secondLabelKey = generateName(e2eLabelKey)
		secondLabelValue = generateName(e2eLabelValue)
		createTestCapp(k8sClient, secondCappName, namespaceName, map[string]string{secondLabelKey: secondLabelValue}, nil)
	})

	Context("Validate get Capps route", func() {
		It("Should get all Capps in a namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.CappsKey)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

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

		It("Should get all Capps in a namespace with limit of 50", func() {
			limit := "50"
			page := "1"

			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s?limit=%s&page=%s", platformURL, namespaceName, testutils.CappsKey, limit, page)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

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

		It("Should get one Capp in a namespace with limit of 1 and page 1", func() {
			limit := "1"
			page := "1"

			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s?limit=%s&page=%s", platformURL, namespaceName, testutils.CappsKey, limit, page)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.CappsKey: []types.CappSummary{
					{Name: oneCappName, Images: []string{CappImageName}, URL: fmt.Sprintf("https://%s-%s.%s", oneCappName, namespaceName, clusterDomain)},
				},
				testutils.CountKey: 1,
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})

		It("Should get one Capp in a namespace with limit of 1 and page 2", func() {
			limit := "1"
			page := "2"

			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s?limit=%s&page=%s", platformURL, namespaceName, testutils.CappsKey, limit, page)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.CappsKey: []types.CappSummary{
					{Name: secondCappName, Images: []string{CappImageName}, URL: fmt.Sprintf("https://%s-%s.%s", secondCappName, namespaceName, clusterDomain)},
				},
				testutils.CountKey: 1,
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})

		It("Should not get Capps with limit of 1 and page 3", func() {
			limit := "1"
			page := "3"

			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s?limit=%s&page=%s", platformURL, namespaceName, testutils.CappsKey, limit, page)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.CappsKey: nil,
				testutils.CountKey: 0,
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})

		It("Should get all Capps with a specific labelSelector in a namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.CappsKey)
			params := url.Values{}
			params.Add(testutils.LabelSelectorKey, fmt.Sprintf("%s=%s", secondLabelKey, secondLabelValue))

			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, fmt.Sprintf("%s?%s", uri, params.Encode()), "", "", userToken)

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
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.CappsKey)
			params := url.Values{}
			params.Add(testutils.LabelSelectorKey, fmt.Sprintf("%s=%s", testutils.LabelCappName, testutils.InvalidLabelSelector))

			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, fmt.Sprintf("%s?%s", uri, params.Encode()), "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.DetailsKey: fmt.Sprintf("unable to parse requirement: values[0][%s]: Invalid value: %q: "+
					"a valid label must be an empty string or consist of alphanumeric characters, "+
					"'-', '_' or '.', and must start and end with an alphanumeric character "+
					"(e.g. 'MyValue',  or 'my_value',  or '12345', regex used for validation is '(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?')",
					testutils.LabelCappName, testutils.InvalidLabelSelector),
				testutils.ErrorKey: testutils.OperationFailed,
			}

			Expect(status).Should(Equal(http.StatusBadRequest))
			compareResponses(response, expectedResponse)
		})

		It("Should get no Capps with valid labelSelector", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.CappsKey)
			params := url.Values{}
			params.Add(testutils.LabelSelectorKey, fmt.Sprintf("%s=%s", secondLabelKey, secondLabelValue+testutils.NonExistentSuffix))

			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, fmt.Sprintf("%s?%s", uri, params.Encode()), "", "", userToken)

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
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.CappsKey, oneCappName)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

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
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.CappsKey, oneCappName+testutils.NonExistentSuffix)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

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
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName+testutils.NonExistentSuffix, testutils.CappsKey, oneCappName)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

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

			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.CappsKey)
			requestData := mocks.PrepareCreateCappType(newCappName, []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue}}, nil)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPost, uri, "", "", userToken)
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
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.CappsKey)
			requestData := mocks.PrepareCreateCappType("", []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue}}, nil)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPost, uri, "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.DetailsKey: "Key: 'CreateCapp.Metadata.Name' Error:Field validation for 'Name' failed on the 'required' tag",
				testutils.ErrorKey:   testutils.InvalidRequest,
			}

			Expect(status).Should(Equal(http.StatusBadRequest))
			compareResponses(response, expectedResponse)
		})

		It("Should handle already existing Capp on creation", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.CappsKey)
			requestData := mocks.PrepareCreateCappType(oneCappName, []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue}}, nil)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPost, uri, "", "", userToken)
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
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.CappsKey, oneCappName)
			requestData := mocks.PrepareCreateCappType(oneCappName, []types.KeyValue{{Key: oneLabelKey, Value: oneLabelValue + "-updated"}}, nil)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPut, uri, "", "", userToken)
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
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.CappsKey, oneCappName+testutils.NonExistentSuffix)
			requestData := mocks.PrepareCreateCappType(oneCappName+testutils.NonExistentSuffix, []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue + "-updated"}}, nil)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPut, uri, "", "", userToken)
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
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName+testutils.NonExistentSuffix, testutils.CappsKey, oneCappName)
			requestData := mocks.PrepareCreateCappType(oneCappName, []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue + "-updated"}}, nil)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPut, uri, "", "", userToken)
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

	Context("Validate update state of a certain Capp route", func() {
		It("Should update state of an existing Capp to disabled", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s/state", platformURL, namespaceName, testutils.CappsKey, oneCappName)
			requestData := mocks.PrepareUpdateCappStateType(testutils.DisabledState)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPut, uri, "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.NameKey:  oneCappName,
				testutils.StateKey: testutils.DisabledState,
			}
			Expect(status).Should(Equal(http.StatusOK))
			Expect(response[testutils.NameKey], expectedResponse[testutils.NameKey])
			Expect(response[testutils.StateKey], expectedResponse[testutils.StateKey])

			Eventually(func() bool {
				capp := getCapp(k8sClient, oneCappName, namespaceName)
				return capp.Status.StateStatus.State == testutils.DisabledState
			}, testutils.Timeout, testutils.Interval).Should(BeTrue())
		})

		It("Should handle update state request for non existing Capp", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s/state", platformURL, namespaceName, testutils.CappsKey, oneCappName+testutils.NonExistentSuffix)
			requestData := mocks.PrepareUpdateCappStateType(testutils.DisabledState)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPut, uri, "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.DetailsKey: fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, oneCappName+testutils.NonExistentSuffix),
				testutils.ErrorKey:   testutils.OperationFailed,
			}

			capp := mocks.PrepareCapp(oneCappName+testutils.NonExistentSuffix, namespaceName, clusterDomain, nil, nil)
			Expect(doesResourceExist(k8sClient, &capp)).To(BeFalse())
			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})

		It("Should update state of an existing Capp to enabled", func() {
			disabledCappName := generateName("disabled-" + testCappName)
			createTestCapp(k8sClient, disabledCappName, namespaceName, map[string]string{}, nil)

			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s/state", platformURL, namespaceName, testutils.CappsKey, disabledCappName)
			requestData := mocks.PrepareUpdateCappStateType(testutils.EnabledState)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPut, uri, "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.NameKey:  disabledCappName,
				testutils.StateKey: testutils.EnabledState,
			}
			Expect(status).Should(Equal(http.StatusOK))
			Expect(response[testutils.NameKey], expectedResponse[testutils.NameKey])
			Expect(response[testutils.StateKey], expectedResponse[testutils.StateKey])

			Eventually(func() bool {
				capp := getCapp(k8sClient, disabledCappName, namespaceName)
				return capp.Status.StateStatus.State == testutils.EnabledState
			}, testutils.Timeout, testutils.Interval).Should(BeTrue())
		})

		It("Should handle update state request for a non existing namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s/state", platformURL, namespaceName+testutils.NonExistentSuffix, testutils.CappsKey, oneCappName)
			requestData := mocks.PrepareUpdateCappStateType(testutils.EnabledState)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPut, uri, "", "", userToken)

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

	Context("Validate get state of a certain Capp route", func() {
		It("Should get state of an existing enabled Capp", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s/state", platformURL, namespaceName, testutils.CappsKey, oneCappName)

			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.LastReadyRevision:   oneCappName + "-00001",
				testutils.LastCreatedRevision: oneCappName + "-00001",
				testutils.StateKey:            testutils.EnabledState,
			}
			Expect(status).Should(Equal(http.StatusOK))
			Expect(response[testutils.LastReadyRevision], expectedResponse[testutils.LastReadyRevision])
			Expect(response[testutils.LastCreatedRevision], expectedResponse[testutils.LastCreatedRevision])
			Expect(response[testutils.StateKey], expectedResponse[testutils.StateKey])

			Eventually(func() bool {
				capp := getCapp(k8sClient, oneCappName, namespaceName)
				return capp.Status.StateStatus.State == testutils.EnabledState
			}, testutils.Timeout, testutils.Interval).Should(BeTrue())
		})

		It("Should handle get state of a non existing Capp", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s/state", platformURL, namespaceName, testutils.CappsKey, oneCappName+testutils.NonExistentSuffix)

			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.DetailsKey: fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, oneCappName+testutils.NonExistentSuffix),
				testutils.ErrorKey:   testutils.OperationFailed,
			}

			capp := mocks.PrepareCapp(oneCappName+testutils.NonExistentSuffix, namespaceName, clusterDomain, nil, nil)
			Expect(doesResourceExist(k8sClient, &capp)).To(BeFalse())
			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})

		It("Should get state of an existing disabled Capp", func() {
			disabledCappName := generateName("disabled-" + testCappName)
			createTestCapp(k8sClient, disabledCappName, namespaceName, map[string]string{}, nil)

			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s/state", platformURL, namespaceName, testutils.CappsKey, disabledCappName)

			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.LastReadyRevision:   testutils.NoRevision,
				testutils.LastCreatedRevision: testutils.NoRevision,
				testutils.StateKey:            testutils.DisabledState,
			}

			Expect(status).Should(Equal(http.StatusOK))
			Expect(response[testutils.LastReadyRevision], expectedResponse[testutils.LastReadyRevision])
			Expect(response[testutils.LastCreatedRevision], expectedResponse[testutils.LastCreatedRevision])
			Expect(response[testutils.StateKey], expectedResponse[testutils.StateKey])

			Eventually(func() bool {
				capp := getCapp(k8sClient, disabledCappName, namespaceName)
				return capp.Status.StateStatus.State == testutils.EnabledState
			}, testutils.Timeout, testutils.Interval).Should(BeTrue())
		})

		It("Should handle a get state request for a non existing Capp", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s/state", platformURL, namespaceName+testutils.NonExistentSuffix, testutils.CappsKey, oneCappName)

			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.DetailsKey: fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, oneCappName),
				testutils.ErrorKey:   testutils.OperationFailed,
			}

			capp := mocks.PrepareCapp(oneCappName, namespaceName+testutils.NonExistentSuffix, clusterDomain, nil, nil)
			Expect(doesResourceExist(k8sClient, &capp)).To(BeFalse())
			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})

		It("Should handle a get state request for a non existing namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s/state", platformURL, namespaceName+testutils.NonExistentSuffix, testutils.CappsKey, oneCappName)
			requestData := mocks.PrepareUpdateCappStateType(testutils.EnabledState)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPut, uri, "", "", userToken)
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
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.CappsKey, oneCappName)
			status, response := performHTTPRequest(httpClient, nil, http.MethodDelete, uri, "", "", userToken)

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
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.CappsKey, oneCappName+testutils.NonExistentSuffix)
			status, response := performHTTPRequest(httpClient, nil, http.MethodDelete, uri, "", "", userToken)

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
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName+testutils.NonExistentSuffix, testutils.CappsKey, oneCappName)
			status, response := performHTTPRequest(httpClient, nil, http.MethodDelete, uri, "", "", userToken)

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
