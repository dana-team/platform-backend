package e2e_tests

import (
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

var _ = Describe("Validate CappRevision routes and functionality", func() {
	var namespaceName, oneCappName, secondCappName string
	var namespace corev1.Namespace
	var oneCappRevisionNames, secondCappRevisionNames []string

	BeforeEach(func() {
		namespaceName = generateName(e2eNamespace)
		namespace = mocks.PrepareNamespace(namespaceName, map[string]string{e2eLabelKey: e2eLabelValue})
		createResource(k8sClient, &namespace)

		oneCappName = generateName("a-" + testCappRevisionName)
		oneCapp := mocks.PrepareCapp(oneCappName, namespaceName, clusterDomain, nil, nil)
		createCapp(k8sClient, &oneCapp)
		oneCappRevisionNames = getCappRevisionNames(k8sClient, oneCappName, namespaceName)

		secondCappName = generateName("b-" + testCappRevisionName)
		secondCapp := mocks.PrepareCapp(secondCappName, namespaceName, clusterDomain, nil, nil)
		createCapp(k8sClient, &secondCapp)
		secondCappRevisionNames = getCappRevisionNames(k8sClient, secondCappName, namespaceName)
	})

	Context("Validate get CappRevisions route", func() {
		It("Should get all CappRevisions in a namespace", func() {
			baseURI := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.CapprevisionsKey)
			status, response := performAuthorizedHTTPRequest(httpClient, nil, http.MethodGet, baseURI, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.CapprevisionsKey: append(oneCappRevisionNames, secondCappRevisionNames...),
				testutils.CountKey:         2,
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})

		It("Should get all CappRevisions with a specific labelSelector in a namespace", func() {
			baseURI := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.CapprevisionsKey)
			params := url.Values{}
			params.Add(testutils.LabelSelectorKey, fmt.Sprintf("%s=%s", testutils.LabelCappName, secondCappName))

			status, response := performAuthorizedHTTPRequest(httpClient, nil, http.MethodGet, fmt.Sprintf("%s?%s", baseURI, params.Encode()), "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.CapprevisionsKey: secondCappRevisionNames,
				testutils.CountKey:         1,
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})

		It("Should fail getting CappRevisions with an invalid labelSelector", func() {
			baseURI := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.CapprevisionsKey)
			params := url.Values{}
			params.Add(testutils.LabelSelectorKey, fmt.Sprintf("%s=%s", testutils.LabelCappName, secondCappName+" invalid"))

			status, response := performAuthorizedHTTPRequest(httpClient, nil, http.MethodGet, fmt.Sprintf("%s?%s", baseURI, params.Encode()), "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.DetailsKey: "found 'invalid', expected: ',' or 'end of string'",
				testutils.ErrorKey:   testutils.OperationFailed,
			}

			Expect(status).Should(Equal(http.StatusBadRequest))
			compareResponses(response, expectedResponse)
		})

		It("Should get no CappRevisions with valid labelSelector", func() {
			baseURI := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.CapprevisionsKey)
			params := url.Values{}
			params.Add(testutils.LabelSelectorKey, fmt.Sprintf("%s=%s", testutils.LabelCappName, secondCappName+testutils.NonExistentSuffix))

			status, response := performAuthorizedHTTPRequest(httpClient, nil, http.MethodGet, fmt.Sprintf("%s?%s", baseURI, params.Encode()), "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.CapprevisionsKey: nil,
				testutils.CountKey:         0,
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})
	})

	Context("Validate get CappRevision route", func() {
		It("Should get a specific CappRevision in a namespace", func() {
			baseURI := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.CapprevisionsKey, oneCappRevisionNames[0])
			status, response := performAuthorizedHTTPRequest(httpClient, nil, http.MethodGet, baseURI, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.MetadataKey:    types.Metadata{Name: oneCappRevisionNames[0], Namespace: namespaceName},
				testutils.AnnotationsKey: nil,
				testutils.LabelsKey:      []types.KeyValue{{Key: testutils.LabelCappName, Value: oneCappName}},
				testutils.SpecKey:        mocks.PrepareCappRevisionSpec(),
				testutils.StatusKey:      mocks.PrepareCappRevisionStatus(),
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})

		It("Should handle a not found CappRevision in a namespace", func() {
			baseURI := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.CapprevisionsKey, oneCappName+testutils.NonExistentSuffix)
			status, response := performAuthorizedHTTPRequest(httpClient, nil, http.MethodGet, baseURI, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.DetailsKey: fmt.Sprintf("%s.%s %q not found", testutils.CapprevisionsKey, cappv1alpha1.GroupVersion.Group, oneCappName+testutils.NonExistentSuffix),
				testutils.ErrorKey:   testutils.OperationFailed,
			}

			cappRevision := mocks.PrepareCappRevision(oneCappName+testutils.NonExistentSuffix, namespaceName, nil, nil)
			Expect(doesResourceExist(k8sClient, &cappRevision)).To(BeFalse())
			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})

		It("Should handle a get of a CappRevision in a not found namespace", func() {
			baseURI := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName+testutils.NonExistentSuffix, testutils.CapprevisionsKey, oneCappName)
			status, response := performAuthorizedHTTPRequest(httpClient, nil, http.MethodGet, baseURI, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.DetailsKey: fmt.Sprintf("%s.%s %q not found", testutils.CapprevisionsKey, cappv1alpha1.GroupVersion.Group, oneCappName),
				testutils.ErrorKey:   testutils.OperationFailed,
			}

			cappRevision := mocks.PrepareCappRevision(oneCappName, namespaceName+testutils.NonExistentSuffix, nil, nil)
			Expect(doesResourceExist(k8sClient, &cappRevision)).To(BeFalse())
			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})
	})
})
