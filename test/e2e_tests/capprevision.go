package e2e_tests

import (
	"fmt"
	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	"github.com/dana-team/platform-backend/src/controllers"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/dana-team/platform-backend/src/utils/testutils"
	"github.com/dana-team/platform-backend/src/utils/testutils/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = Describe("Validate CappRevision routes and functionality", func() {
	var namespaceName, oneCappName, secondCappName string
	var oneCappRevisionNames, secondCappRevisionNames []string
	var site string

	BeforeEach(func() {
		namespaceName = generateName(e2eNamespace)
		createTestNamespace(k8sClient, namespaceName)

		oneCappName = generateName("a-" + testCappRevisionName)
		oneCapp := createTestCapp(k8sClient, oneCappName, namespaceName, placementName, nil, nil)
		oneCappRevisionNames = getCappRevisionNames(k8sClient, oneCappName, namespaceName, oneCapp.Status.ApplicationLinks.Site)

		secondCappName = generateName("b-" + testCappRevisionName)
		secondCapp := createTestCapp(k8sClient, secondCappName, namespaceName, placementName, nil, nil)
		secondCappRevisionNames = getCappRevisionNames(k8sClient, secondCappName, namespaceName, secondCapp.Status.ApplicationLinks.Site)

		site = oneCapp.Status.ApplicationLinks.Site
	})

	Context("Validate get CappRevisions from Capp route", func() {
		It("Should get all CappRevisions of a specific capp in a namespace with limit of 50", func() {
			limit := "50"
			page := "1"

			uri := fmt.Sprintf("%s/v1/namespaces/%s/capps/%s/%s?limit=%s&page=%s", platformURL, namespaceName, oneCappName, testutils.CapprevisionsKey, limit, page)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.CapprevisionsKey: oneCappRevisionNames,
				testutils.CountKey:         len(oneCappRevisionNames),
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})

		It("Should get one CappRevision of a specific capp in a namespace with limit of 1 and page 1", func() {
			limit := "1"
			page := "1"

			uri := fmt.Sprintf("%s/v1/namespaces/%s/capps/%s/%s?limit=%s&page=%s", platformURL, namespaceName, oneCappName, testutils.CapprevisionsKey, limit, page)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.CapprevisionsKey: oneCappRevisionNames,
				testutils.CountKey:         1,
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})

		It("Should not get CappRevisions in a namespace with limit of 1 and page 2", func() {
			limit := "1"
			page := "2"

			uri := fmt.Sprintf("%s/v1/namespaces/%s/capps/%s/%s?limit=%s&page=%s", platformURL, namespaceName, oneCappName, testutils.CapprevisionsKey, limit, page)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.CapprevisionsKey: nil,
				testutils.CountKey:         0,
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})

		It("Should not get CappRevisions of a specific capp with limit of 1 and page 10000", func() {
			limit := "1"
			page := "10000"

			uri := fmt.Sprintf("%s/v1/namespaces/%s/capps/%s/%s?limit=%s&page=%s", platformURL, namespaceName, oneCappName, testutils.CapprevisionsKey, limit, page)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.CapprevisionsKey: nil,
				testutils.CountKey:         0,
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})

		It("Should get all CappRevisions of a specific capp in a namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/capps/%s/%s", platformURL, namespaceName, secondCappName, testutils.CapprevisionsKey)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.CapprevisionsKey: secondCappRevisionNames,
				testutils.CountKey:         len(secondCappRevisionNames),
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})

		It("Should fail getting CappRevisions with an invalid capp name", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/capps/%s/%s", platformURL, namespaceName, oneCappName+testutils.NonExistentSuffix, testutils.CapprevisionsKey)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, oneCappName+testutils.NonExistentSuffix),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})
	})

	Context("Validate get CappRevision from Capp route", func() {
		It("Should get a specific CappRevision in a namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/capps/%s/%s/%s", platformURL, namespaceName, oneCappName, testutils.CapprevisionsKey, oneCappRevisionNames[0])
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			annotations := []types.KeyValue{
				{Key: testutils.HasPlacementLabel, Value: site},
				{Key: testutils.LastUpdatedCappLabel, Value: fmt.Sprintf("%s:%s:%s:%s", testutils.System, testutils.ServiceAccount, testutils.RcsDeployerSystem, testutils.RcsOcmDeployerControllerManager)}}

			expectedResponse := map[string]interface{}{
				testutils.MetadataKey:    types.Metadata{Name: oneCappRevisionNames[0], Namespace: namespaceName},
				testutils.AnnotationsKey: annotations,
				testutils.LabelsKey:      []types.KeyValue{{Key: testutils.LabelCappName, Value: oneCappName}},
				testutils.SpecKey: mocks.PrepareCappRevisionSpec(
					placementName,
					mocks.ConvertKeyValueSliceToMap([]types.KeyValue{{Key: testutils.ManagedByLabel, Value: testutils.Rcs}}),
					mocks.ConvertKeyValueSliceToMap(annotations)),
				testutils.StatusKey: mocks.PrepareCappRevisionStatus(),
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})

		It("Should handle a not found CappRevision in a namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/capps/%s/%s/%s", platformURL, namespaceName, oneCappName, testutils.CapprevisionsKey, oneCappName+testutils.NonExistentSuffix)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.ErrorKey: fmt.Sprintf("%s, %s",
					fmt.Sprintf(controllers.ErrCouldNotGetCappRevision, oneCappName+testutils.NonExistentSuffix, namespaceName),
					fmt.Sprintf("%s.%s %q not found", testutils.CapprevisionsKey, cappv1alpha1.GroupVersion.Group, oneCappName+testutils.NonExistentSuffix),
				),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			cappRevision := mocks.PrepareCappRevision(site, oneCappName+testutils.NonExistentSuffix, namespaceName, nil, nil)
			Expect(doesResourceExist(k8sClient, &cappRevision)).To(BeFalse())
			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})

		It("Should handle a get of a CappRevision in a not found namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/capps/%s/%s/%s", platformURL, namespaceName+testutils.NonExistentSuffix, oneCappName, testutils.CapprevisionsKey, oneCappName)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, oneCappName),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			cappRevision := mocks.PrepareCappRevision(site, oneCappName, namespaceName+testutils.NonExistentSuffix, nil, nil)
			Expect(doesResourceExist(k8sClient, &cappRevision)).To(BeFalse())
			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})
	})

	Context("Validate get CappRevisions from cluster route", func() {
		It("Should get all CappRevisions in a cluster", func() {
			uri := fmt.Sprintf("%s/v1/clusters/%s/namespaces/%s/%s", platformURL, site, namespaceName, testutils.CapprevisionsKey)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			allCapps := append(oneCappRevisionNames, secondCappRevisionNames...)
			expectedResponse := map[string]interface{}{
				testutils.CapprevisionsKey: allCapps,
				testutils.CountKey:         len(allCapps),
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})

		It("Should fail getting CappRevisions with an invalid cluster", func() {
			uri := fmt.Sprintf("%s/v1/clusters/%s/namespaces/%s/%s", platformURL, site+testutils.NonExistentSuffix, namespaceName, testutils.CapprevisionsKey)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf("Could not list capp revisions, no such cluster %s", site+testutils.NonExistentSuffix),
				testutils.ReasonKey: "",
			}

			Expect(status).Should(Equal(http.StatusInternalServerError))
			compareResponses(response, expectedResponse)
		})

		It("Should get all CappRevisions in a namespace with limit of 50", func() {
			limit := "50"
			page := "1"

			uri := fmt.Sprintf("%s/v1/clusters/%s/namespaces/%s/%s?limit=%s&page=%s", platformURL, site, namespaceName, testutils.CapprevisionsKey, limit, page)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			allCapps := append(oneCappRevisionNames, secondCappRevisionNames...)
			expectedResponse := map[string]interface{}{
				testutils.CapprevisionsKey: allCapps,
				testutils.CountKey:         len(allCapps),
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})

		It("Should get one CappRevision in a namespace with limit of 1 and page 1", func() {
			limit := "1"
			page := "1"

			uri := fmt.Sprintf("%s/v1/clusters/%s/namespaces/%s/%s?limit=%s&page=%s", platformURL, site, namespaceName, testutils.CapprevisionsKey, limit, page)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.CapprevisionsKey: oneCappRevisionNames,
				testutils.CountKey:         1,
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})

		It("Should not get CappRevisions with limit of 1 and page 10000", func() {
			limit := "1"
			page := "10000"

			uri := fmt.Sprintf("%s/v1/clusters/%s/namespaces/%s/%s?limit=%s&page=%s", platformURL, site, namespaceName, testutils.CapprevisionsKey, limit, page)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.CapprevisionsKey: nil,
				testutils.CountKey:         0,
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})
	})

	Context("Validate get CappRevision from cluster route", func() {
		It("Should get a specific CappRevision in a namespace", func() {
			uri := fmt.Sprintf("%s/v1/clusters/%s/namespaces/%s/%s/%s", platformURL, site, namespaceName, testutils.CapprevisionsKey, oneCappRevisionNames[0])
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			annotations := []types.KeyValue{
				{Key: testutils.HasPlacementLabel, Value: site},
				{Key: testutils.LastUpdatedCappLabel, Value: fmt.Sprintf("%s:%s:%s:%s", testutils.System, testutils.ServiceAccount, testutils.RcsDeployerSystem, testutils.RcsOcmDeployerControllerManager)}}

			expectedResponse := map[string]interface{}{
				testutils.MetadataKey:    types.Metadata{Name: oneCappRevisionNames[0], Namespace: namespaceName},
				testutils.AnnotationsKey: annotations,
				testutils.LabelsKey:      []types.KeyValue{{Key: testutils.LabelCappName, Value: oneCappName}},
				testutils.SpecKey: mocks.PrepareCappRevisionSpec(
					placementName,
					mocks.ConvertKeyValueSliceToMap([]types.KeyValue{{Key: testutils.ManagedByLabel, Value: testutils.Rcs}}),
					mocks.ConvertKeyValueSliceToMap(annotations)),
				testutils.StatusKey: mocks.PrepareCappRevisionStatus(),
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})

		It("Should handle a not found CappRevision in a namespace", func() {
			uri := fmt.Sprintf("%s/v1/clusters/%s/namespaces/%s/%s/%s", platformURL, site, namespaceName, testutils.CapprevisionsKey, oneCappName+testutils.NonExistentSuffix)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.ErrorKey: fmt.Sprintf("%s, %s",
					fmt.Sprintf(controllers.ErrCouldNotGetCappRevision, oneCappName+testutils.NonExistentSuffix, namespaceName),
					fmt.Sprintf("%s.%s %q not found", testutils.CapprevisionsKey, cappv1alpha1.GroupVersion.Group, oneCappName+testutils.NonExistentSuffix),
				),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			cappRevision := mocks.PrepareCappRevision(site, oneCappName+testutils.NonExistentSuffix, namespaceName, nil, nil)
			Expect(doesResourceExist(k8sClient, &cappRevision)).To(BeFalse())
			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})

		It("Should handle a get of a CappRevision in a not found namespace", func() {
			uri := fmt.Sprintf("%s/v1/clusters/%s/namespaces/%s/%s/%s", platformURL, site, namespaceName+testutils.NonExistentSuffix, testutils.CapprevisionsKey, oneCappName)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf("%s.%s %q not found", testutils.CapprevisionsKey, cappv1alpha1.GroupVersion.Group, oneCappName),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			cappRevision := mocks.PrepareCappRevision(site, oneCappName, namespaceName+testutils.NonExistentSuffix, nil, nil)
			Expect(doesResourceExist(k8sClient, &cappRevision)).To(BeFalse())
			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})
	})
})
