package e2e_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dana-team/platform-backend/src/utils/testutils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/util/retry"
	"net/http"
	"net/url"
	"strings"

	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	"github.com/dana-team/platform-backend/src/controllers"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/dana-team/platform-backend/src/utils/testutils/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func getCappClusterDomain(site string) string {
	var domain string
	parts := strings.Split(clusterDomain, ".")
	if len(parts) > 2 {
		domain = strings.Join(parts[len(parts)-2:], ".")
	}

	return fmt.Sprintf("apps.%s.%s", site, domain)
}

func addPlacementToRCSConfig(newPlacementName, regionName, environmentName string) {
	labels := map[string]string{}
	if regionName != "" && environmentName != "" {
		labels = map[string]string{
			e2eLabelKey:                            e2eLabelValue,
			testutils.PlacementRegionLabelKey:      regionName,
			testutils.PlacementEnvironmentLabelKey: environmentName,
		}
	} else if regionName != "" {
		labels = map[string]string{
			e2eLabelKey:                       e2eLabelValue,
			testutils.PlacementRegionLabelKey: regionName,
		}
	} else if environmentName != "" {
		labels = map[string]string{
			e2eLabelKey:                            e2eLabelValue,
			testutils.PlacementEnvironmentLabelKey: environmentName,
		}
	}

	createPlacement(k8sClient, newPlacementName, placementNS, labels)

	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		rcsConfig := getRCSConfig(k8sClient, rcsConfigName, rcsConfigNamespace)
		rcsConfig.Spec.Placements = append(rcsConfig.Spec.Placements, newPlacementName)
		return updateRCSConfig(k8sClient, rcsConfig)
	})
	Expect(err).Should(Not(HaveOccurred()))
}

var _ = Describe("Validate Capp routes and functionality", func() {
	var namespaceName, oneCappName, secondCappName string
	var oneLabelKey, oneLabelValue, secondLabelKey, secondLabelValue string
	var oneCappSite, secondCappSite string

	BeforeEach(func() {
		namespaceName = generateName(e2eNamespace)
		createTestNamespace(k8sClient, namespaceName)

		oneCappName = generateName("a-" + testCappName)
		oneLabelKey = generateName(e2eLabelKey)
		oneLabelValue = generateName(e2eLabelValue)
		oneCapp := createTestCapp(k8sClient, oneCappName, namespaceName, placementName, map[string]string{oneLabelKey: oneLabelValue}, nil)
		oneCappSite = oneCapp.Status.ApplicationLinks.Site

		secondCappName = generateName("b-" + testCappName)
		secondLabelKey = generateName(e2eLabelKey)
		secondLabelValue = generateName(e2eLabelValue)
		secondCapp := createTestCapp(k8sClient, secondCappName, namespaceName, placementName, map[string]string{secondLabelKey: secondLabelValue}, nil)
		secondCappSite = secondCapp.Status.ApplicationLinks.Site
	})

	Context("Validate get Capps route", func() {
		It("Should get all Capps in a namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.CappsKey)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.CappsKey: []types.CappSummary{
					{Name: oneCappName, Images: []string{CappImageName}, URL: fmt.Sprintf("https://%s-%s.%s", oneCappName, namespaceName, getCappClusterDomain(oneCappSite))},
					{Name: secondCappName, Images: []string{CappImageName}, URL: fmt.Sprintf("https://%s-%s.%s", secondCappName, namespaceName, getCappClusterDomain(secondCappSite))},
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
					{Name: oneCappName, Images: []string{CappImageName}, URL: fmt.Sprintf("https://%s-%s.%s", oneCappName, namespaceName, getCappClusterDomain(oneCappSite))},
					{Name: secondCappName, Images: []string{CappImageName}, URL: fmt.Sprintf("https://%s-%s.%s", secondCappName, namespaceName, getCappClusterDomain(secondCappSite))},
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
					{Name: oneCappName, Images: []string{CappImageName}, URL: fmt.Sprintf("https://%s-%s.%s", oneCappName, namespaceName, getCappClusterDomain(oneCappSite))},
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
					{Name: secondCappName, Images: []string{CappImageName}, URL: fmt.Sprintf("https://%s-%s.%s", secondCappName, namespaceName, getCappClusterDomain(secondCappSite))},
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
					{Name: secondCappName, Images: []string{CappImageName}, URL: fmt.Sprintf("https://%s-%s.%s", secondCappName, namespaceName, getCappClusterDomain(secondCappSite))},
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
				testutils.ErrorKey:  controllers.ErrParsingLabelSelector,
				testutils.ReasonKey: testutils.ReasonBadRequest,
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
				testutils.SpecKey:     mocks.PrepareCappSpec(placementName),
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
				testutils.ErrorKey:  fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, oneCappName+testutils.NonExistentSuffix),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			capp := mocks.PrepareCapp(oneCappName+testutils.NonExistentSuffix, namespaceName, clusterDomain, placementName, nil, nil)
			Expect(doesResourceExist(k8sClient, &capp)).To(BeFalse())
			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})

		It("Should handle a get of a Capp in a not found namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName+testutils.NonExistentSuffix, testutils.CappsKey, oneCappName)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, oneCappName),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			capp := mocks.PrepareCapp(oneCappName, namespaceName+testutils.NonExistentSuffix, clusterDomain, placementName, nil, nil)
			Expect(doesResourceExist(k8sClient, &capp)).To(BeFalse())
			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})
	})

	Context("Validate create Capp route", func() {
		It("Should create Capp in a namespace with a given site", func() {
			newCappName := generateName(testCappName)

			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.CappsKey)
			requestData := mocks.PrepareCreateCappType(newCappName, placementName, []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue}}, nil)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPost, uri, "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.MetadataKey:    types.Metadata{Name: newCappName, Namespace: namespaceName},
				testutils.LabelsKey:      []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue}},
				testutils.AnnotationsKey: []types.KeyValue{{Key: testutils.LastUpdatedCappLabel, Value: e2eUser}},
				testutils.SpecKey:        mocks.PrepareCappSpec(placementName),
				testutils.StatusKey:      cappv1alpha1.CappStatus{},
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)

			capp := mocks.PrepareCapp(oneCappName, namespaceName, clusterDomain, placementName, nil, nil)
			Eventually(func() bool {
				return doesResourceExist(k8sClient, &capp)
			}, testutils.Timeout, testutils.Interval).Should(BeTrue())
		})

		It("Should create Capp in a namespace with the site being determined using region and environment parameters", func() {
			newCappName := generateName(testCappName)

			baseURI := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.CappsKey)
			requestData := mocks.PrepareCreateCappType(newCappName, "", []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue}}, nil)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			newPlacementName := generateName(testPlacementName)
			regionName := generateName(testRegionName)
			environmentName := generateName(testEnvironment)
			addPlacementToRCSConfig(newPlacementName, regionName, environmentName)

			params := url.Values{}
			params.Add(testutils.PlacementRegionKey, regionName)
			params.Add(testutils.PlacementEnvironmentKey, environmentName)

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPost, fmt.Sprintf("%s?%s", baseURI, params.Encode()), "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.MetadataKey:    types.Metadata{Name: newCappName, Namespace: namespaceName},
				testutils.LabelsKey:      []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue}},
				testutils.AnnotationsKey: []types.KeyValue{{Key: testutils.LastUpdatedCappLabel, Value: e2eUser}},
				testutils.SpecKey:        mocks.PrepareCappSpec(newPlacementName),
				testutils.StatusKey:      cappv1alpha1.CappStatus{},
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)

			capp := mocks.PrepareCapp(oneCappName, namespaceName, clusterDomain, newPlacementName, nil, nil)
			Eventually(func() bool {
				return doesResourceExist(k8sClient, &capp)
			}, testutils.Timeout, testutils.Interval).Should(BeTrue())
		})

		It("Should create Capp in a namespace with the site being determined using region parameter", func() {
			newCappName := generateName(testCappName)

			baseURI := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.CappsKey)
			requestData := mocks.PrepareCreateCappType(newCappName, "", []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue}}, nil)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			newPlacementName := generateName(testPlacementName)
			regionName := generateName(testRegionName)
			addPlacementToRCSConfig(newPlacementName, regionName, "")

			params := url.Values{}
			params.Add(testutils.PlacementRegionKey, regionName)

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPost, fmt.Sprintf("%s?%s", baseURI, params.Encode()), "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.MetadataKey:    types.Metadata{Name: newCappName, Namespace: namespaceName},
				testutils.LabelsKey:      []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue}},
				testutils.AnnotationsKey: []types.KeyValue{{Key: testutils.LastUpdatedCappLabel, Value: e2eUser}},
				testutils.SpecKey:        mocks.PrepareCappSpec(newPlacementName),
				testutils.StatusKey:      cappv1alpha1.CappStatus{},
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)

			capp := mocks.PrepareCapp(oneCappName, namespaceName, clusterDomain, newPlacementName, nil, nil)
			Eventually(func() bool {
				return doesResourceExist(k8sClient, &capp)
			}, testutils.Timeout, testutils.Interval).Should(BeTrue())
		})

		It("Should create Capp in a namespace with the site being determined using environment parameter", func() {
			newCappName := generateName(testCappName)

			baseURI := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.CappsKey)
			requestData := mocks.PrepareCreateCappType(newCappName, "", []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue}}, nil)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			newPlacementName := generateName(testPlacementName)
			environmentName := generateName(testEnvironment)
			addPlacementToRCSConfig(newPlacementName, "", environmentName)

			params := url.Values{}
			params.Add(testutils.PlacementEnvironmentKey, environmentName)

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPost, fmt.Sprintf("%s?%s", baseURI, params.Encode()), "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.MetadataKey:    types.Metadata{Name: newCappName, Namespace: namespaceName},
				testutils.LabelsKey:      []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue}},
				testutils.AnnotationsKey: []types.KeyValue{{Key: testutils.LastUpdatedCappLabel, Value: e2eUser}},
				testutils.SpecKey:        mocks.PrepareCappSpec(newPlacementName),
				testutils.StatusKey:      cappv1alpha1.CappStatus{},
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)

			capp := mocks.PrepareCapp(oneCappName, namespaceName, clusterDomain, newPlacementName, nil, nil)
			Eventually(func() bool {
				return doesResourceExist(k8sClient, &capp)
			}, testutils.Timeout, testutils.Interval).Should(BeTrue())
		})

		It("Should create Capp in a namespace with the first matching Placement", func() {
			newCappName := generateName(testCappName)

			baseURI := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.CappsKey)
			requestData := mocks.PrepareCreateCappType(newCappName, "", []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue}}, nil)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			oneNewPlacementName := generateName("a-" + testPlacementName)
			secondNewPlacementName := generateName("b-" + testPlacementName)

			regionName := generateName(testRegionName)
			environmentName := generateName(testEnvironment)
			addPlacementToRCSConfig(oneNewPlacementName, regionName, environmentName)
			addPlacementToRCSConfig(secondNewPlacementName, regionName, environmentName)

			params := url.Values{}
			params.Add(testutils.PlacementRegionKey, regionName)
			params.Add(testutils.PlacementEnvironmentKey, environmentName)

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPost, fmt.Sprintf("%s?%s", baseURI, params.Encode()), "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.MetadataKey:    types.Metadata{Name: newCappName, Namespace: namespaceName},
				testutils.LabelsKey:      []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue}},
				testutils.AnnotationsKey: []types.KeyValue{{Key: testutils.LastUpdatedCappLabel, Value: e2eUser}},
				testutils.SpecKey:        mocks.PrepareCappSpec(oneNewPlacementName),
				testutils.StatusKey:      cappv1alpha1.CappStatus{},
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)

			capp := mocks.PrepareCapp(oneCappName, namespaceName, clusterDomain, oneNewPlacementName, nil, nil)
			Eventually(func() bool {
				return doesResourceExist(k8sClient, &capp)
			}, testutils.Timeout, testutils.Interval).Should(BeTrue())
		})

		It("Should create Capp in a namespace with the given site despite using region and environment parameters", func() {
			newCappName := generateName(testCappName)

			baseURI := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.CappsKey)
			requestData := mocks.PrepareCreateCappType(newCappName, placementName, []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue}}, nil)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			newPlacementName := generateName(testPlacementName)
			regionName := generateName(testRegionName)
			environmentName := generateName(testEnvironment)
			addPlacementToRCSConfig(newPlacementName, regionName, environmentName)

			params := url.Values{}
			params.Add(testutils.PlacementRegionKey, regionName)
			params.Add(testutils.PlacementEnvironmentKey, environmentName)

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPost, fmt.Sprintf("%s?%s", baseURI, params.Encode()), "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.MetadataKey:    types.Metadata{Name: newCappName, Namespace: namespaceName},
				testutils.LabelsKey:      []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue}},
				testutils.AnnotationsKey: []types.KeyValue{{Key: testutils.LastUpdatedCappLabel, Value: e2eUser}},
				testutils.SpecKey:        mocks.PrepareCappSpec(placementName),
				testutils.StatusKey:      cappv1alpha1.CappStatus{},
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)

			capp := mocks.PrepareCapp(oneCappName, namespaceName, clusterDomain, newPlacementName, nil, nil)
			Eventually(func() bool {
				return doesResourceExist(k8sClient, &capp)
			}, testutils.Timeout, testutils.Interval).Should(BeTrue())
		})

		It("Should fail in creation of Capp without a given site and without matching placement", func() {
			newCappName := generateName(testCappName)

			baseURI := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.CappsKey)
			requestData := mocks.PrepareCreateCappType(newCappName, "", []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue}}, nil)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			regionName := generateName(testRegionName)
			environmentName := generateName(testEnvironment)

			params := url.Values{}
			params.Add(testutils.PlacementRegionKey, regionName)
			params.Add(testutils.PlacementEnvironmentKey, environmentName)

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPost, fmt.Sprintf("%s?%s", baseURI, params.Encode()), "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  "No matching Placements found",
				testutils.ReasonKey: testutils.ReasonBadRequest,
			}

			Expect(status).Should(Equal(http.StatusBadRequest))
			compareResponses(response, expectedResponse)
		})

		It("Should fail in creation of Capp without a given site and without query parameters", func() {
			newCappName := generateName(testCappName)

			baseURI := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.CappsKey)
			requestData := mocks.PrepareCreateCappType(newCappName, "", []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue}}, nil)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			params := url.Values{}

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPost, fmt.Sprintf("%s?%s", baseURI, params.Encode()), "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf("%q and/or %q query parameters must be set when Site is unspecified in request body", testutils.PlacementEnvironmentKey, testutils.PlacementRegionKey),
				testutils.ReasonKey: testutils.ReasonBadRequest,
			}

			Expect(status).Should(Equal(http.StatusBadRequest))
			compareResponses(response, expectedResponse)
		})

		It("Should fail in creation of Capp with bad request body", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.CappsKey)
			requestData := mocks.PrepareCreateCappType("", placementName, []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue}}, nil)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPost, uri, "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  "Key: 'CreateCapp.Metadata.Name' Error:Field validation for 'Name' failed on the 'required' tag",
				testutils.ReasonKey: testutils.ReasonBadRequest,
			}

			Expect(status).Should(Equal(http.StatusBadRequest))
			compareResponses(response, expectedResponse)
		})

		It("Should handle already existing Capp on creation", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.CappsKey)
			requestData := mocks.PrepareCreateCappType(oneCappName, placementName, []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue}}, nil)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPost, uri, "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf("%s.%s %q already exists", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, oneCappName),
				testutils.ReasonKey: testutils.ReasonAlreadyExists,
			}

			Expect(status).Should(Equal(http.StatusConflict))
			compareResponses(response, expectedResponse)
		})
	})

	Context("Validate update Capp route", func() {
		It("Should update an existing Capp in a namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.CappsKey, oneCappName)
			requestData := mocks.PrepareCreateCappType(oneCappName, placementName, []types.KeyValue{{Key: oneLabelKey, Value: oneLabelValue + "-updated"}}, nil)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPut, uri, "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.MetadataKey: types.Metadata{Name: oneCappName, Namespace: namespaceName},
				testutils.LabelsKey:   []types.KeyValue{{Key: oneLabelKey, Value: oneLabelValue + "-updated"}},
				testutils.SpecKey:     mocks.PrepareCappSpec(placementName),
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
			requestData := mocks.PrepareCreateCappType(oneCappName+testutils.NonExistentSuffix, placementName, []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue + "-updated"}}, nil)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPut, uri, "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, oneCappName+testutils.NonExistentSuffix),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			capp := mocks.PrepareCapp(oneCappName+testutils.NonExistentSuffix, namespaceName, clusterDomain, placementName, nil, nil)
			Expect(doesResourceExist(k8sClient, &capp)).To(BeFalse())
			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})

		It("Should handle update of not found namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName+testutils.NonExistentSuffix, testutils.CappsKey, oneCappName)
			requestData := mocks.PrepareCreateCappType(oneCappName, placementName, []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue + "-updated"}}, nil)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPut, uri, "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, oneCappName),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			capp := mocks.PrepareCapp(oneCappName, namespaceName+testutils.NonExistentSuffix, clusterDomain, placementName, nil, nil)
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
				testutils.ErrorKey:  fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, oneCappName+testutils.NonExistentSuffix),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			capp := mocks.PrepareCapp(oneCappName+testutils.NonExistentSuffix, namespaceName, clusterDomain, placementName, nil, nil)
			Expect(doesResourceExist(k8sClient, &capp)).To(BeFalse())
			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})

		It("Should update state of an existing Capp to enabled", func() {
			disabledCappName := generateName("disabled-" + testCappName)
			createTestCapp(k8sClient, disabledCappName, namespaceName, placementName, map[string]string{}, nil)

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
				testutils.ErrorKey:  fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, oneCappName),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			capp := mocks.PrepareCapp(oneCappName, namespaceName+testutils.NonExistentSuffix, clusterDomain, placementName, nil, nil)
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
				testutils.ErrorKey:  fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, oneCappName+testutils.NonExistentSuffix),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			capp := mocks.PrepareCapp(oneCappName+testutils.NonExistentSuffix, namespaceName, clusterDomain, placementName, nil, nil)
			Expect(doesResourceExist(k8sClient, &capp)).To(BeFalse())
			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})

		It("Should get state of an existing disabled Capp", func() {
			disabledCappName := generateName("disabled-" + testCappName)
			createTestDisabledCapp(k8sClient, disabledCappName, namespaceName, placementName, map[string]string{}, nil)

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
				return capp.Status.StateStatus.State == testutils.DisabledState
			}, testutils.Timeout, testutils.Interval).Should(BeTrue())
		})

		It("Should handle a get state request for a non existing Capp", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s/state", platformURL, namespaceName+testutils.NonExistentSuffix, testutils.CappsKey, oneCappName)

			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, oneCappName),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			capp := mocks.PrepareCapp(oneCappName, namespaceName+testutils.NonExistentSuffix, clusterDomain, placementName, nil, nil)
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
				testutils.ErrorKey:  fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, oneCappName),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			capp := mocks.PrepareCapp(oneCappName, namespaceName+testutils.NonExistentSuffix, clusterDomain, placementName, nil, nil)
			Expect(doesResourceExist(k8sClient, &capp)).To(BeFalse())
			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})
	})

	Context("Validate get dns of a certain Capp route", func() {
		It("Should get records of an existing Capp", func() {
			hostname := "dns-capp"
			domain := "dana-dev.com"
			cappName := generateName(hostname)
			createTestCappWithHostname(k8sClient, cappName, namespaceName, cappName, domain, nil, nil)
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s/dns", platformURL, namespaceName, testutils.CappsKey, oneCappName)

			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.RecordsKey: []types.DNS{
					{Status: corev1.ConditionTrue, Name: fmt.Sprintf("%s.%s", hostname, domain)},
				},
			}
			Expect(status).Should(Equal(http.StatusOK))
			Expect(response[testutils.RecordsKey], expectedResponse[testutils.RecordsKey])
		})

		It("Should handle get dns of a non existing Capp", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s/dns", platformURL, namespaceName, testutils.CappsKey, oneCappName+testutils.NonExistentSuffix)

			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, oneCappName+testutils.NonExistentSuffix),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			capp := mocks.PrepareCapp(oneCappName+testutils.NonExistentSuffix, namespaceName, clusterDomain, placementName, nil, nil)
			Expect(doesResourceExist(k8sClient, &capp)).To(BeFalse())
			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})

		It("Should handle a get dns request for a non existing Capp", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s/dns", platformURL, namespaceName+testutils.NonExistentSuffix, testutils.CappsKey, oneCappName)

			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, oneCappName),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			capp := mocks.PrepareCapp(oneCappName, namespaceName+testutils.NonExistentSuffix, clusterDomain, placementName, nil, nil)
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

			capp := mocks.PrepareCapp(oneCappName, namespaceName, clusterDomain, placementName, nil, nil)
			Eventually(func() bool {
				return doesResourceExist(k8sClient, &capp)
			}, testutils.Timeout, testutils.Interval).Should(BeFalse())
		})

		It("Should handle deletion of not found Capp", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.CappsKey, oneCappName+testutils.NonExistentSuffix)
			status, response := performHTTPRequest(httpClient, nil, http.MethodDelete, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, oneCappName+testutils.NonExistentSuffix),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			capp := mocks.PrepareCapp(oneCappName+testutils.NonExistentSuffix, namespaceName, clusterDomain, placementName, nil, nil)
			Expect(doesResourceExist(k8sClient, &capp)).To(BeFalse())
			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})

		It("Should handle deletion of Capp in a not found namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName+testutils.NonExistentSuffix, testutils.CappsKey, oneCappName)
			status, response := performHTTPRequest(httpClient, nil, http.MethodDelete, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, oneCappName),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			capp := mocks.PrepareCapp(oneCappName, namespaceName+testutils.NonExistentSuffix, clusterDomain, placementName, nil, nil)
			Expect(doesResourceExist(k8sClient, &capp)).To(BeFalse())
			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})
	})
})
