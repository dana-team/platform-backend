package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	"github.com/dana-team/platform-backend/src/controllers"
	"github.com/dana-team/platform-backend/src/middleware"
	"github.com/dana-team/platform-backend/src/utils/testutils"
	"github.com/dana-team/platform-backend/src/utils/testutils/mocks"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/dana-team/platform-backend/src/types"
	"github.com/stretchr/testify/assert"
)

func TestGetCapps(t *testing.T) {
	testNamespaceName := testutils.CappNamespace + "-get"

	type selector struct {
		keys   []string
		values []string
	}

	type pagination struct {
		limit string
		page  string
	}

	type requestURI struct {
		namespace        string
		labelSelector    selector
		paginationParams pagination
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		requestURI requestURI
		want       want
	}{
		"ShouldSucceedGettingCapps": {
			requestURI: requestURI{
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.CountKey: 4,
					testutils.CappsKey: []types.CappSummary{
						{Name: testutils.CappName + "-1", URL: fmt.Sprintf("https://%s-%s.%s", testutils.CappName+"-1", testNamespaceName, testutils.Domain), Images: []string{testutils.CappImage}},
						{Name: testutils.CappName + "-2", URL: fmt.Sprintf("https://%s-%s.%s", testutils.CappName+"-2", testNamespaceName, testutils.Domain), Images: []string{testutils.CappImage}},
						{Name: testutils.CappName + "-3", URL: fmt.Sprintf("https://%s.%s", testutils.Hostname, testutils.Domain), Images: []string{testutils.CappImage}},
						{Name: testutils.CappName + "-4", URL: fmt.Sprintf("https://%s.%s", testutils.Hostname, testutils.Domain), Images: []string{testutils.CappImage}},
					},
				},
			},
		},
		"ShouldSucceedGettingAllCappsWithLimitOf4": {
			requestURI: requestURI{
				namespace:        testNamespaceName,
				paginationParams: pagination{limit: "4", page: "1"},
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.CountKey: 4,
					testutils.CappsKey: []types.CappSummary{
						{Name: testutils.CappName + "-1", URL: fmt.Sprintf("https://%s-%s.%s", testutils.CappName+"-1", testNamespaceName, testutils.Domain), Images: []string{testutils.CappImage}},
						{Name: testutils.CappName + "-2", URL: fmt.Sprintf("https://%s-%s.%s", testutils.CappName+"-2", testNamespaceName, testutils.Domain), Images: []string{testutils.CappImage}},
						{Name: testutils.CappName + "-3", URL: fmt.Sprintf("https://%s.%s", testutils.Hostname, testutils.Domain), Images: []string{testutils.CappImage}},
						{Name: testutils.CappName + "-4", URL: fmt.Sprintf("https://%s.%s", testutils.Hostname, testutils.Domain), Images: []string{testutils.CappImage}},
					},
				},
			},
		},
		"ShouldSucceedGettingCappsWithLabelSelector": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				labelSelector: selector{
					keys:   []string{testutils.LabelKey + "-1"},
					values: []string{testutils.LabelValue + "-1"},
				},
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.CountKey: 1,
					testutils.CappsKey: []types.CappSummary{
						{Name: testutils.CappName + "-1", URL: fmt.Sprintf("https://%s-%s.%s", testutils.CappName+"-1", testNamespaceName, testutils.Domain), Images: []string{testutils.CappImage}},
					},
				},
			},
		},
		"ShouldFailGettingCappsWithInvalidLabelSelector": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				labelSelector: selector{
					keys:   []string{testutils.LabelKey + "-1"},
					values: []string{testutils.LabelValue + " 1"},
				},
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response: map[string]interface{}{
					testutils.ErrorKey:  controllers.ErrParsingLabelSelector,
					testutils.ReasonKey: metav1.StatusReasonBadRequest,
				},
			},
		},
		"ShouldSucceedGettingNoCappsWithLabelSelector": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				labelSelector: selector{
					keys:   []string{testutils.LabelKey + testutils.NonExistentSuffix},
					values: []string{testutils.LabelValue + testutils.NonExistentSuffix},
				},
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.CountKey: 0,
					testutils.CappsKey: nil,
				},
			},
		},
	}

	setup()
	mocks.CreateTestNamespace(fakeClient, testNamespaceName)
	mocks.CreateTestCapp(dynClient, testutils.CappName+"-1", testNamespaceName, testutils.Domain, testutils.SiteName,
		map[string]string{testutils.LabelKey + "-1": testutils.LabelValue + "-1"}, nil)
	mocks.CreateTestCapp(dynClient, testutils.CappName+"-2", testNamespaceName, testutils.Domain, testutils.SiteName,
		map[string]string{testutils.LabelKey + "-2": testutils.LabelValue + "-2"}, nil)
	mocks.CreateTestCappWithHostname(dynClient, testutils.CappName+"-3", testNamespaceName, testutils.Hostname, testutils.Domain,
		map[string]string{testutils.LabelKey + "-3": testutils.LabelValue + "-3"}, nil)
	mocks.CreateTestCappWithHostname(dynClient, testutils.CappName+"-4", testNamespaceName, testutils.Hostname, testutils.Domain,
		map[string]string{testutils.LabelKey + "-4": testutils.LabelValue + "-4"}, nil)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			params := url.Values{}
			for i, key := range test.requestURI.labelSelector.keys {
				params.Add(testutils.LabelSelectorKey, fmt.Sprintf("%s=%s", key, test.requestURI.labelSelector.values[i]))
			}

			baseURI := fmt.Sprintf("/v1/namespaces/%s/capps", test.requestURI.namespace)
			if test.requestURI.paginationParams.limit != "" {
				params.Add(middleware.LimitCtxKey, test.requestURI.paginationParams.limit)
			}

			if test.requestURI.paginationParams.page != "" {
				params.Add(middleware.PageCtxKey, test.requestURI.paginationParams.page)
			}

			request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?%s", baseURI, params.Encode()), nil)
			assert.NoError(t, err)
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			assert.Equal(t, test.want.statusCode, writer.Code)

			var response map[string]interface{}
			err = json.Unmarshal(writer.Body.Bytes(), &response)
			assert.NoError(t, err)

			wantResponseJSON, err := json.Marshal(test.want.response)
			assert.NoError(t, err)
			var wantResponseNormalized map[string]interface{}
			err = json.Unmarshal(wantResponseJSON, &wantResponseNormalized)
			assert.NoError(t, err)
			assert.Equal(t, wantResponseNormalized, response)
		})
	}
}

func TestGetCapp(t *testing.T) {
	testNamespaceName := testutils.CappNamespace + "-get-one"

	type requestURI struct {
		name      string
		namespace string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		requestURI requestURI
		want       want
	}{
		"ShouldSucceedGettingCapp": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				name:      testutils.CappName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.MetadataKey:    types.Metadata{Name: testutils.CappName, Namespace: testNamespaceName},
					testutils.LabelsKey:      []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue}},
					testutils.AnnotationsKey: nil,
					testutils.SpecKey:        mocks.PrepareCappSpec(testutils.SiteName),
					testutils.StatusKey:      mocks.PrepareCappStatus(testutils.CappName, testNamespaceName, testutils.Domain),
				},
			},
		},
		"ShouldHandleNotFoundCapp": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				name:      testutils.CappName + testutils.NonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.ReasonKey: metav1.StatusReasonNotFound,
					testutils.ErrorKey: fmt.Sprintf("%v, %v",
						fmt.Sprintf(controllers.ErrCouldNotGetCapp, testutils.CappName+testutils.NonExistentSuffix, testNamespaceName),
						fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, testutils.CappName+testutils.NonExistentSuffix),
					),
				},
			},
		},
		"ShouldHandleNotFoundNamespace": {
			requestURI: requestURI{
				namespace: testNamespaceName + testutils.NonExistentSuffix,
				name:      testutils.CappName,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.ReasonKey: metav1.StatusReasonNotFound,
					testutils.ErrorKey: fmt.Sprintf("%v, %v",
						fmt.Sprintf(controllers.ErrCouldNotGetCapp, testutils.CappName, testNamespaceName+testutils.NonExistentSuffix),
						fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, testutils.CappName),
					),
				},
			},
		},
	}

	setup()
	mocks.CreateTestNamespace(fakeClient, testNamespaceName)
	mocks.CreateTestCapp(dynClient, testutils.CappName, testNamespaceName, testutils.Domain, testutils.SiteName, map[string]string{testutils.LabelKey: testutils.LabelValue}, nil)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			baseURI := fmt.Sprintf("/v1/namespaces/%s/capps/%s", test.requestURI.namespace, test.requestURI.name)
			request, err := http.NewRequest(http.MethodGet, baseURI, nil)
			assert.NoError(t, err)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			assert.Equal(t, test.want.statusCode, writer.Code)

			var response map[string]interface{}
			err = json.Unmarshal(writer.Body.Bytes(), &response)
			assert.NoError(t, err)

			wantResponseJSON, err := json.Marshal(test.want.response)
			assert.NoError(t, err)
			var wantResponseNormalized map[string]interface{}
			err = json.Unmarshal(wantResponseJSON, &wantResponseNormalized)
			assert.NoError(t, err)
			assert.Equal(t, wantResponseNormalized, response)
		})
	}
}

func TestGetCappState(t *testing.T) {
	testNamespaceName := testutils.CappNamespace + "-get-state"

	type requestURI struct {
		name      string
		namespace string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		requestURI requestURI
		want       want
	}{
		"ShouldSucceedGettingEnabledCappState": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				name:      fmt.Sprintf("%s-%s", testutils.CappName, testutils.EnabledState),
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.LastReadyRevision:   fmt.Sprintf("%s-%s-%s", testutils.CappName, testutils.EnabledState, "00001"),
					testutils.LastCreatedRevision: fmt.Sprintf("%s-%s-%s", testutils.CappName, testutils.EnabledState, "00001"),
					testutils.StateKey:            testutils.EnabledState,
				},
			},
		},
		"ShouldSucceedGettingDisabledCappState": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				name:      fmt.Sprintf("%s-%s", testutils.CappName, testutils.DisabledState),
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.LastReadyRevision:   testutils.NoRevision,
					testutils.LastCreatedRevision: testutils.NoRevision,
					testutils.StateKey:            testutils.DisabledState,
				},
			},
		},
		"ShouldHandleNotFoundCapp": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				name:      testutils.CappName + testutils.NonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.ErrorKey: fmt.Sprintf("%v, %v",
						fmt.Sprintf(controllers.ErrCouldNotGetCapp, testutils.CappName+testutils.NonExistentSuffix, testNamespaceName),
						fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, testutils.CappName+testutils.NonExistentSuffix),
					),
					testutils.ReasonKey: metav1.StatusReasonNotFound,
				},
			},
		},
		"ShouldHandleNotFoundNamespace": {
			requestURI: requestURI{
				namespace: testNamespaceName + testutils.NonExistentSuffix,
				name:      testutils.CappName,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.ErrorKey: fmt.Sprintf("%v, %v",
						fmt.Sprintf(controllers.ErrCouldNotGetCapp, testutils.CappName, testNamespaceName+testutils.NonExistentSuffix),
						fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, testutils.CappName)),
					testutils.ReasonKey: metav1.StatusReasonNotFound,
				},
			},
		},
	}
	setup()
	mocks.CreateTestNamespace(fakeClient, testNamespaceName)
	mocks.CreateTestCappWithState(dynClient, fmt.Sprintf("%s-%s", testutils.CappName, testutils.EnabledState),
		testNamespaceName, testutils.EnabledState, testutils.SiteName, map[string]string{}, map[string]string{})
	mocks.CreateTestCappWithState(dynClient, fmt.Sprintf("%s-%s", testutils.CappName, testutils.DisabledState),
		testNamespaceName, testutils.DisabledState, testutils.SiteName, map[string]string{}, map[string]string{})

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			baseURI := fmt.Sprintf("/v1/namespaces/%s/capps/%s/state", test.requestURI.namespace, test.requestURI.name)
			request, err := http.NewRequest(http.MethodGet, baseURI, nil)
			assert.NoError(t, err)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			assert.Equal(t, test.want.statusCode, writer.Code)

			var response map[string]interface{}
			err = json.Unmarshal(writer.Body.Bytes(), &response)
			assert.NoError(t, err)

			wantResponseJSON, err := json.Marshal(test.want.response)
			assert.NoError(t, err)
			var wantResponseNormalized map[string]interface{}
			err = json.Unmarshal(wantResponseJSON, &wantResponseNormalized)
			assert.NoError(t, err)
			assert.Equal(t, wantResponseNormalized, response)
		})
	}
}

func TestGetCappDNS(t *testing.T) {
	testNamespaceName := testutils.CappNamespace + "-get-dns"

	type requestURI struct {
		name      string
		namespace string
	}

	type dnsParams struct {
		readyStatus   corev1.ConditionStatus
		syncedStatus  corev1.ConditionStatus
		isConditioned bool
		hostname      string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		requestURI requestURI
		want       want
		records    []dnsParams
		cappName   string
	}{
		"ShouldSucceedGettingRecords": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				name:      fmt.Sprintf("%s-%s", testutils.CappName, testutils.Available),
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.RecordsKey: []types.DNS{
						{Status: corev1.ConditionFalse, Name: fmt.Sprintf("%s.%s", testutils.Hostname+"-1", testutils.DefaultZone)},
						{Status: corev1.ConditionTrue, Name: fmt.Sprintf("%s.%s", testutils.Hostname+"-2", testutils.DefaultZone)},
						{Status: corev1.ConditionUnknown, Name: fmt.Sprintf("%s.%s", testutils.Hostname+"-3", testutils.DefaultZone)}},
				},
			},
			records: []dnsParams{
				{readyStatus: corev1.ConditionFalse, syncedStatus: corev1.ConditionTrue, isConditioned: true,
					hostname: fmt.Sprintf("%s.%s", testutils.Hostname+"-1", testutils.DefaultZone)},
				{readyStatus: corev1.ConditionTrue, syncedStatus: corev1.ConditionTrue, isConditioned: true,
					hostname: fmt.Sprintf("%s.%s", testutils.Hostname+"-2", testutils.DefaultZone)},
				{readyStatus: corev1.ConditionUnknown, syncedStatus: corev1.ConditionTrue, isConditioned: true,
					hostname: fmt.Sprintf("%s.%s", testutils.Hostname+"-3", testutils.DefaultZone)},
			},
			cappName: fmt.Sprintf("%s-%s", testutils.CappName, testutils.Available),
		},
		"ShouldSucceedGettingUnknownDNS": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				name:      fmt.Sprintf("%s-%s", testutils.CappName, testutils.Unknown),
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.RecordsKey: []types.DNS{
						{Status: corev1.ConditionUnknown, Name: fmt.Sprintf("%s.%s", testutils.Hostname+"-1", testutils.DefaultZone)}},
				},
			},
			records: []dnsParams{
				{readyStatus: corev1.ConditionUnknown, syncedStatus: corev1.ConditionFalse, isConditioned: true,
					hostname: fmt.Sprintf("%s.%s", testutils.Hostname+"-1", testutils.DefaultZone)},
			},

			cappName: fmt.Sprintf("%s-%s", testutils.CappName, testutils.Unknown),
		},
		"ShouldSucceedUnavailableDNS": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				name:      fmt.Sprintf("%s-%s", testutils.CappName, testutils.Unavailable),
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.RecordsKey: []types.DNS{
						{Status: corev1.ConditionFalse, Name: fmt.Sprintf("%s.%s", testutils.Hostname+"-1", testutils.DefaultZone)}},
				},
			},
			records: []dnsParams{
				{readyStatus: corev1.ConditionFalse, syncedStatus: corev1.ConditionFalse, isConditioned: true,
					hostname: fmt.Sprintf("%s.%s", testutils.Hostname+"-1", testutils.DefaultZone)},
			},

			cappName: fmt.Sprintf("%s-%s", testutils.CappName, testutils.Unavailable),
		},
		"ShouldSucceedAvailableDNS": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				name:      fmt.Sprintf("%s-%s", testutils.CappName, testutils.Available+"1"),
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.RecordsKey: []types.DNS{
						{Status: corev1.ConditionTrue, Name: fmt.Sprintf("%s.%s", testutils.Hostname+"-1", testutils.DefaultZone)}},
				},
			},
			records: []dnsParams{
				{readyStatus: corev1.ConditionTrue, syncedStatus: corev1.ConditionTrue, isConditioned: true, hostname: fmt.Sprintf("%s.%s", testutils.Hostname+"-1", testutils.DefaultZone)},
			},

			cappName: fmt.Sprintf("%s-%s", testutils.CappName, testutils.Available+"1"),
		},
		"ShouldHandleNotFoundCapp": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				name:      testutils.CappName + testutils.NonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.ErrorKey:  fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, testutils.CappName+testutils.NonExistentSuffix),
					testutils.ReasonKey: metav1.StatusReasonNotFound,
				},
			},
		},
		"ShouldHandleNotFoundNamespace": {
			requestURI: requestURI{
				namespace: testNamespaceName + testutils.NonExistentSuffix,
				name:      testutils.CappName,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.ErrorKey:  fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, testutils.CappName),
					testutils.ReasonKey: metav1.StatusReasonNotFound,
				},
			},
		},
	}

	setup()
	mocks.CreateTestNamespace(fakeClient, testNamespaceName)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {

			if test.cappName != "" {
				mocks.CreateTestCapp(dynClient, test.cappName, testNamespaceName, testutils.Domain, testutils.SiteName, map[string]string{}, map[string]string{})
			}

			for i, dns := range test.records {

				if !dns.isConditioned {
					mocks.CreateTestCNAMERecordWithoutConditions(dynClient, test.cappName+strconv.Itoa(i), test.cappName, testNamespaceName, dns.hostname)
				} else {
					mocks.CreateTestCNAMERecord(dynClient, test.cappName+strconv.Itoa(i), test.cappName, testNamespaceName, dns.hostname, dns.readyStatus, dns.syncedStatus)
				}
			}

			baseURI := fmt.Sprintf("/v1/namespaces/%s/capps/%s/dns", test.requestURI.namespace, test.requestURI.name)
			request, err := http.NewRequest(http.MethodGet, baseURI, nil)
			assert.NoError(t, err)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			assert.Equal(t, test.want.statusCode, writer.Code)

			var response map[string]interface{}
			err = json.Unmarshal(writer.Body.Bytes(), &response)
			assert.NoError(t, err)

			wantResponseJSON, err := json.Marshal(test.want.response)
			assert.NoError(t, err)
			var wantResponseNormalized map[string]interface{}
			err = json.Unmarshal(wantResponseJSON, &wantResponseNormalized)
			assert.NoError(t, err)
			assert.Equal(t, wantResponseNormalized, response)
		})
	}
}

func TestCreateCapp(t *testing.T) {
	testNamespaceName := testutils.CappNamespace + "-create"

	type queryParams struct {
		keys   []string
		values []string
	}

	type requestURI struct {
		query     queryParams
		namespace string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		requestURI  requestURI
		want        want
		requestData interface{}
	}{
		"ShouldSucceedCreatingCappWithSite": {
			requestURI: requestURI{
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.MetadataKey:    types.Metadata{Name: testutils.CappName, Namespace: testNamespaceName},
					testutils.LabelsKey:      []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue}},
					testutils.AnnotationsKey: nil,
					testutils.SpecKey:        mocks.PrepareCappSpec(testutils.SiteName),
					testutils.StatusKey:      cappv1alpha1.CappStatus{},
				},
			},
			requestData: mocks.PrepareCreateCappType(testutils.CappName, testutils.SiteName, []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue}}, nil),
		},
		"ShouldSucceedCreatingCappWithPlacementRegionAndEnvironment": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				query:     queryParams{keys: []string{testutils.PlacementEnvironmentKey, testutils.PlacementRegionKey}, values: []string{testutils.EnvironmentName + "-1", testutils.RegionName + "-1"}},
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.MetadataKey:    types.Metadata{Name: testutils.CappName + "-2", Namespace: testNamespaceName},
					testutils.LabelsKey:      []types.KeyValue{{Key: testutils.LabelKey + "-2", Value: testutils.LabelValue + "-2"}},
					testutils.AnnotationsKey: nil,
					testutils.SpecKey:        mocks.PrepareCappSpec(testutils.PlacementName + "-1"),
					testutils.StatusKey:      cappv1alpha1.CappStatus{},
				},
			},
			requestData: mocks.PrepareCreateCappType(testutils.CappName+"-2", "", []types.KeyValue{{Key: testutils.LabelKey + "-2", Value: testutils.LabelValue + "-2"}}, nil),
		},
		"ShouldSucceedCreatingCappWithOnlyPlacementRegion": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				query:     queryParams{keys: []string{testutils.PlacementRegionKey}, values: []string{testutils.RegionName + "-2"}},
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.MetadataKey:    types.Metadata{Name: testutils.CappName + "-3", Namespace: testNamespaceName},
					testutils.LabelsKey:      []types.KeyValue{{Key: testutils.LabelKey + "-3", Value: testutils.LabelValue + "-3"}},
					testutils.AnnotationsKey: nil,
					testutils.SpecKey:        mocks.PrepareCappSpec(testutils.PlacementName + "-2"),
					testutils.StatusKey:      cappv1alpha1.CappStatus{},
				},
			},
			requestData: mocks.PrepareCreateCappType(testutils.CappName+"-3", "", []types.KeyValue{{Key: testutils.LabelKey + "-3", Value: testutils.LabelValue + "-3"}}, nil),
		},
		"ShouldSucceedCreatingCappWithOnlyPlacementEnvironment": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				query:     queryParams{keys: []string{testutils.PlacementEnvironmentKey}, values: []string{testutils.EnvironmentName + "-3"}},
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.MetadataKey:    types.Metadata{Name: testutils.CappName + "-4", Namespace: testNamespaceName},
					testutils.LabelsKey:      []types.KeyValue{{Key: testutils.LabelKey + "-4", Value: testutils.LabelValue + "-4"}},
					testutils.AnnotationsKey: nil,
					testutils.SpecKey:        mocks.PrepareCappSpec(testutils.PlacementName + "-3"),
					testutils.StatusKey:      cappv1alpha1.CappStatus{},
				},
			},
			requestData: mocks.PrepareCreateCappType(testutils.CappName+"-4", "", []types.KeyValue{{Key: testutils.LabelKey + "-4", Value: testutils.LabelValue + "-4"}}, nil),
		},
		"ShouldSucceedCreatingCappWithFirstMatchingPlacement": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				query:     queryParams{keys: []string{testutils.PlacementEnvironmentKey, testutils.PlacementEnvironmentKey}, values: []string{testutils.EnvironmentName + "-4", testutils.RegionName + "-4"}},
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.MetadataKey:    types.Metadata{Name: testutils.CappName + "-5", Namespace: testNamespaceName},
					testutils.LabelsKey:      []types.KeyValue{{Key: testutils.LabelKey + "-5", Value: testutils.LabelValue + "-5"}},
					testutils.AnnotationsKey: nil,
					testutils.SpecKey:        mocks.PrepareCappSpec(testutils.PlacementName + "-4"),
					testutils.StatusKey:      cappv1alpha1.CappStatus{},
				},
			},
			requestData: mocks.PrepareCreateCappType(testutils.CappName+"-5", "", []types.KeyValue{{Key: testutils.LabelKey + "-5", Value: testutils.LabelValue + "-5"}}, nil),
		},
		"ShouldSucceedCreatingCappWithSiteEvenIfPlacementQueryIsSet": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				query:     queryParams{keys: []string{testutils.PlacementEnvironmentKey, testutils.PlacementEnvironmentKey}, values: []string{testutils.EnvironmentName + "-5", testutils.RegionName + "-5"}},
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.MetadataKey:    types.Metadata{Name: testutils.CappName + "-6", Namespace: testNamespaceName},
					testutils.LabelsKey:      []types.KeyValue{{Key: testutils.LabelKey + "-6", Value: testutils.LabelValue + "-6"}},
					testutils.AnnotationsKey: nil,
					testutils.SpecKey:        mocks.PrepareCappSpec(testutils.SiteName),
					testutils.StatusKey:      cappv1alpha1.CappStatus{},
				},
			},
			requestData: mocks.PrepareCreateCappType(testutils.CappName+"-6", testutils.SiteName, []types.KeyValue{{Key: testutils.LabelKey + "-6", Value: testutils.LabelValue + "-6"}}, nil),
		},
		"ShouldFailCreatingCappWithoutSiteAndWithoutMatchingPlacement": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				query:     queryParams{keys: []string{testutils.PlacementEnvironmentKey, testutils.PlacementEnvironmentKey}, values: []string{testutils.EnvironmentName + testutils.NonExistentSuffix, testutils.RegionName + testutils.NonExistentSuffix}},
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response: map[string]interface{}{
					testutils.ErrorKey:  "No matching Placements found",
					testutils.ReasonKey: metav1.StatusReasonBadRequest,
				},
			},
			requestData: mocks.PrepareCreateCappType(testutils.CappName+"-7", "", []types.KeyValue{{Key: testutils.LabelKey + "-7", Value: testutils.LabelValue + "-7"}}, nil),
		},
		"ShouldFailCreatingCappWithoutSetAndWithoutQuery": {
			requestURI: requestURI{
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response: map[string]interface{}{
					testutils.ErrorKey:  fmt.Sprintf("%q and/or %q query parameters must be set when Site is unspecified in request body", testutils.PlacementEnvironmentKey, testutils.PlacementRegionKey),
					testutils.ReasonKey: metav1.StatusReasonBadRequest,
				},
			},
			requestData: mocks.PrepareCreateCappType(testutils.CappName+"-8", "", []types.KeyValue{{Key: testutils.LabelKey + "-8", Value: testutils.LabelValue + "-8"}}, nil),
		},
		"ShouldFailWithBadRequestBody": {
			requestURI: requestURI{
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response: map[string]interface{}{
					testutils.ErrorKey:  "Key: 'CreateCapp.Metadata.Name' Error:Field validation for 'Name' failed on the 'required' tag",
					testutils.ReasonKey: metav1.StatusReasonBadRequest,
				},
			},
			requestData: mocks.PrepareCreateCappType("", testutils.SiteName, []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue}}, nil),
		},
		"ShouldHandleAlreadyExists": {
			requestURI: requestURI{
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusConflict,
				response: map[string]interface{}{
					testutils.ErrorKey: fmt.Sprintf("%v, %v",
						fmt.Sprintf(controllers.ErrCouldNotCreateCapp, testutils.CappName+"-1", testNamespaceName),
						fmt.Sprintf("%s.%s %q already exists", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, testutils.CappName+"-1")),
					testutils.ReasonKey: metav1.StatusReasonAlreadyExists,
				},
			},
			requestData: mocks.PrepareCreateCappType(testutils.CappName+"-1", testutils.SiteName, []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue}}, nil),
		},
	}

	setup()
	mocks.CreateTestNamespace(fakeClient, testNamespaceName)
	mocks.CreateTestPlacement(dynClient, testutils.PlacementName+"-1", testNamespaceName, map[string]string{testutils.PlacementRegionLabelKey: testutils.RegionName + "-1", testutils.PlacementEnvironmentLabelKey: testutils.EnvironmentName + "-1"})
	mocks.CreateTestPlacement(dynClient, testutils.PlacementName+"-2", testNamespaceName, map[string]string{testutils.PlacementRegionLabelKey: testutils.RegionName + "-2"})
	mocks.CreateTestPlacement(dynClient, testutils.PlacementName+"-3", testNamespaceName, map[string]string{testutils.PlacementEnvironmentLabelKey: testutils.EnvironmentName + "-3"})
	mocks.CreateTestPlacement(dynClient, testutils.PlacementName+"-4", testNamespaceName, map[string]string{testutils.PlacementRegionLabelKey: testutils.RegionName + "-4", testutils.PlacementEnvironmentLabelKey: testutils.EnvironmentName + "-4"})
	mocks.CreateTestPlacement(dynClient, testutils.PlacementName+"-5", testNamespaceName, map[string]string{testutils.PlacementRegionLabelKey: testutils.RegionName + "-4", testutils.PlacementEnvironmentLabelKey: testutils.EnvironmentName + "-4"})
	mocks.CreateTestCapp(dynClient, testutils.CappName+"-1", testNamespaceName, testutils.Domain, testutils.SiteName, map[string]string{testutils.LabelKey + "-1": testutils.LabelValue + "-1"}, nil)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			params := url.Values{}
			for i, key := range test.requestURI.query.keys {
				params.Add(key, test.requestURI.query.values[i])
			}

			payload, err := json.Marshal(test.requestData)
			assert.NoError(t, err)

			baseURI := fmt.Sprintf("/v1/namespaces/%s/capps", test.requestURI.namespace)
			request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s?%s", baseURI, params.Encode()), bytes.NewBuffer(payload))
			assert.NoError(t, err)
			request.Header.Set(testutils.ContentType, testutils.ApplicationJson)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			assert.Equal(t, test.want.statusCode, writer.Code)

			var response map[string]interface{}
			err = json.Unmarshal(writer.Body.Bytes(), &response)
			assert.NoError(t, err)

			wantResponseJSON, err := json.Marshal(test.want.response)
			assert.NoError(t, err)
			var wantResponseNormalized map[string]interface{}
			err = json.Unmarshal(wantResponseJSON, &wantResponseNormalized)
			assert.NoError(t, err)
			assert.Equal(t, wantResponseNormalized, response)
		})
	}
}

func TestUpdateCapp(t *testing.T) {
	testNamespaceName := testutils.CappNamespace + "-update"

	type requestURI struct {
		name      string
		namespace string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		requestURI  requestURI
		want        want
		requestData interface{}
	}{
		"ShouldSucceedUpdatingCapp": {
			requestURI: requestURI{
				name:      testutils.CappName,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.MetadataKey:    types.Metadata{Name: testutils.CappName, Namespace: testNamespaceName},
					testutils.LabelsKey:      []types.KeyValue{{Key: testutils.LabelKey + "-updated", Value: testutils.LabelValue + "-updated"}},
					testutils.AnnotationsKey: nil,
					testutils.SpecKey:        mocks.PrepareCappSpec(testutils.SiteName),
					testutils.StatusKey:      mocks.PrepareCappStatus(testutils.CappName, testNamespaceName, testutils.Domain),
				},
			},
			requestData: mocks.PrepareUpdateCappType(testutils.SiteName, []types.KeyValue{{Key: testutils.LabelKey + "-updated", Value: testutils.LabelValue + "-updated"}}, nil),
		},
		"ShouldHandleNotFoundCapp": {
			requestURI: requestURI{
				name:      testutils.CappName + testutils.NonExistentSuffix,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.ErrorKey: fmt.Sprintf("%v, %v",
						fmt.Sprintf(controllers.ErrCouldNotGetCapp, testutils.CappName+testutils.NonExistentSuffix, testNamespaceName),
						fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, testutils.CappName+testutils.NonExistentSuffix)),
					testutils.ReasonKey: metav1.StatusReasonNotFound,
				},
			},
			requestData: mocks.PrepareUpdateCappType(testutils.SiteName, []types.KeyValue{{Key: testutils.LabelKey + "-updated", Value: testutils.LabelValue + "-updated"}}, nil),
		},
		"ShouldHandleNotFoundNamespace": {
			requestURI: requestURI{
				name:      testutils.CappName,
				namespace: testNamespaceName + testutils.NonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.ErrorKey: fmt.Sprintf("%v, %v",
						fmt.Sprintf(controllers.ErrCouldNotGetCapp, testutils.CappName, testNamespaceName+testutils.NonExistentSuffix),
						fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, testutils.CappName)),
					testutils.ReasonKey: metav1.StatusReasonNotFound,
				},
			},
			requestData: mocks.PrepareUpdateCappType(testutils.SiteName, []types.KeyValue{{Key: testutils.LabelKey + "-updated", Value: testutils.LabelValue + "-updated"}}, nil),
		},
	}

	setup()
	mocks.CreateTestCapp(dynClient, testutils.CappName, testNamespaceName, testutils.Domain, testutils.SiteName, map[string]string{testutils.LabelKey: testutils.LabelValue}, nil)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			payload, err := json.Marshal(test.requestData)
			assert.NoError(t, err)

			baseURI := fmt.Sprintf("/v1/namespaces/%s/capps/%s", test.requestURI.namespace, test.requestURI.name)
			request, err := http.NewRequest(http.MethodPut, baseURI, bytes.NewBuffer(payload))
			assert.NoError(t, err)
			request.Header.Set(testutils.ContentType, testutils.ApplicationJson)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			assert.Equal(t, test.want.statusCode, writer.Code)

			var response map[string]interface{}
			err = json.Unmarshal(writer.Body.Bytes(), &response)
			assert.NoError(t, err)

			wantResponseJSON, err := json.Marshal(test.want.response)
			assert.NoError(t, err)
			var wantResponseNormalized map[string]interface{}
			err = json.Unmarshal(wantResponseJSON, &wantResponseNormalized)
			assert.NoError(t, err)
			assert.Equal(t, wantResponseNormalized, response)
		})
	}
}

func TestEditCappState(t *testing.T) {
	testNamespaceName := testutils.CappNamespace + "-update"

	type requestURI struct {
		name      string
		namespace string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		requestURI  requestURI
		want        want
		requestData interface{}
	}{
		"ShouldSucceedEditingState": {
			requestURI: requestURI{
				name:      testutils.CappName,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.NameKey:  testutils.CappName,
					testutils.StateKey: testutils.DisabledState,
				},
			},
			requestData: types.CappState{State: testutils.DisabledState},
		},
		"ShouldHandleNotFoundCapp": {
			requestURI: requestURI{
				name:      testutils.CappName + testutils.NonExistentSuffix,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.ErrorKey: fmt.Sprintf("%v, %v",
						fmt.Sprintf(controllers.ErrCouldNotGetCapp, testutils.CappName+testutils.NonExistentSuffix, testNamespaceName),
						fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, testutils.CappName+testutils.NonExistentSuffix)),
					testutils.ReasonKey: metav1.StatusReasonNotFound,
				},
			},
			requestData: types.CappState{State: testutils.DisabledState},
		},
		"ShouldHandleNotFoundNamespace": {
			requestURI: requestURI{
				name:      testutils.CappName,
				namespace: testNamespaceName + testutils.NonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.ErrorKey: fmt.Sprintf("%v, %v",
						fmt.Sprintf(controllers.ErrCouldNotGetCapp, testutils.CappName, testNamespaceName+testutils.NonExistentSuffix),
						fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, testutils.CappName)),
					testutils.ReasonKey: metav1.StatusReasonNotFound,
				},
			},
			requestData: types.CappState{State: testutils.DisabledState},
		},
		"ShouldHandleStateNotAllowed": {
			requestURI: requestURI{
				name:      testutils.CappName,
				namespace: testNamespaceName + testutils.NonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response: map[string]interface{}{
					testutils.ErrorKey:  "Key: 'CappState.State' Error:Field validation for 'State' failed on the 'oneof' tag",
					testutils.ReasonKey: metav1.StatusReasonBadRequest,
				},
			},
			requestData: types.CappState{State: "blabla"},
		},
	}

	setup()
	mocks.CreateTestCapp(dynClient, testutils.CappName, testNamespaceName, testutils.Domain, testutils.SiteName, map[string]string{testutils.LabelKey: testutils.LabelValue}, nil)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			payload, err := json.Marshal(test.requestData)
			assert.NoError(t, err)

			baseURI := fmt.Sprintf("/v1/namespaces/%s/capps/%s/state", test.requestURI.namespace, test.requestURI.name)
			request, err := http.NewRequest(http.MethodPut, baseURI, bytes.NewBuffer(payload))
			assert.NoError(t, err)
			request.Header.Set(testutils.ContentType, testutils.ApplicationJson)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			assert.Equal(t, test.want.statusCode, writer.Code)

			var response map[string]interface{}
			err = json.Unmarshal(writer.Body.Bytes(), &response)
			assert.NoError(t, err)

			wantResponseJSON, err := json.Marshal(test.want.response)
			assert.NoError(t, err)
			var wantResponseNormalized map[string]interface{}
			err = json.Unmarshal(wantResponseJSON, &wantResponseNormalized)
			assert.NoError(t, err)
			assert.Equal(t, wantResponseNormalized, response)
		})
	}
}

func TestDeleteCapp(t *testing.T) {
	testNamespaceName := testutils.CappNamespace + "-delete"

	type requestURI struct {
		name      string
		namespace string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		requestParams requestURI
		want          want
	}{
		"ShouldSucceedDeletingCapp": {
			requestParams: requestURI{
				name:      testutils.CappName,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.MessageKey: fmt.Sprintf("Deleted capp %q in namespace %q successfully", testutils.CappName, testNamespaceName),
				},
			},
		},
		"ShouldHandleNotFoundCapp": {
			requestParams: requestURI{
				name:      testutils.CappName + testutils.NonExistentSuffix,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.ErrorKey: fmt.Sprintf("%v, %v",
						fmt.Sprintf(controllers.ErrCouldNotDeleteCapp, testutils.CappName+testutils.NonExistentSuffix, testNamespaceName),
						fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, testutils.CappName+testutils.NonExistentSuffix)),
					testutils.ReasonKey: metav1.StatusReasonNotFound,
				},
			},
		},
		"ShouldHandleNotFoundNamespace": {
			requestParams: requestURI{
				name:      testutils.CappName,
				namespace: testNamespaceName + testutils.NonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.ErrorKey: fmt.Sprintf("%v, %v",
						fmt.Sprintf(controllers.ErrCouldNotDeleteCapp, testutils.CappName, testNamespaceName+testutils.NonExistentSuffix),
						fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, testutils.CappName)),
					testutils.ReasonKey: metav1.StatusReasonNotFound,
				},
			},
		},
	}

	setup()
	mocks.CreateTestCapp(dynClient, testutils.CappName, testNamespaceName, testutils.Domain, testutils.SiteName, map[string]string{testutils.LabelKey: testutils.LabelValue}, nil)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			baseURI := fmt.Sprintf("/v1/namespaces/%s/capps/%s", test.requestParams.namespace, test.requestParams.name)
			request, err := http.NewRequest(http.MethodDelete, baseURI, nil)
			assert.NoError(t, err)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			assert.Equal(t, test.want.statusCode, writer.Code)

			var response map[string]interface{}
			err = json.Unmarshal(writer.Body.Bytes(), &response)
			assert.NoError(t, err)

			wantResponseJSON, err := json.Marshal(test.want.response)
			assert.NoError(t, err)
			var wantResponseNormalized map[string]interface{}
			err = json.Unmarshal(wantResponseJSON, &wantResponseNormalized)
			assert.NoError(t, err)
			assert.Equal(t, wantResponseNormalized, response)
		})
	}
}
