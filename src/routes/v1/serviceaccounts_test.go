package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/dana-team/platform-backend/src/middleware"

	"github.com/dana-team/platform-backend/src/controllers"
	"github.com/dana-team/platform-backend/src/utils/testutils"
	"github.com/dana-team/platform-backend/src/utils/testutils/mocks"

	"github.com/stretchr/testify/assert"
)

const (
	serviceAccountNamespace = testutils.TestNamespace + "-serviceaccount"
)

func TestGetServiceAccountToken(t *testing.T) {
	testNamespaceName := serviceAccountNamespace + "-get-token"

	type args struct {
		serviceAccountName string
		namespace          string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		args args
		want want
	}{
		"ShouldSucceedGettingToken": {
			args: args{
				namespace:          testNamespaceName,
				serviceAccountName: testutils.ServiceAccountName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.TokenKey: testutils.Value,
				},
			},
		},
		"ShouldNotSucceedGettingTokenForNonExistingNamespace": {
			args: args{
				namespace:          testNamespaceName + testutils.NonExistentSuffix,
				serviceAccountName: testutils.ServiceAccountName,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.ErrorKey: fmt.Sprintf("%s, %s",
						fmt.Sprintf(controllers.ErrCouldNotGetServiceAccount, testutils.ServiceAccountName, testNamespaceName+testutils.NonExistentSuffix),
						fmt.Sprintf("serviceaccounts %q not found", testutils.ServiceAccountName),
					),
					testutils.ReasonKey: testutils.ReasonNotFound,
				},
			},
		},
		"ShouldNotSucceedGettingTokenForNonExistingName": {
			args: args{
				namespace:          testNamespaceName,
				serviceAccountName: testutils.ServiceAccountName + testutils.NonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.ErrorKey: fmt.Sprintf("%s, %s",
						fmt.Sprintf(controllers.ErrCouldNotGetServiceAccount, testutils.ServiceAccountName+testutils.NonExistentSuffix, testNamespaceName),
						fmt.Sprintf("serviceaccounts %q not found", testutils.ServiceAccountName+testutils.NonExistentSuffix),
					),
					testutils.ReasonKey: testutils.ReasonNotFound,
				},
			},
		},
	}

	setup()
	mocks.CreateTestNamespace(fakeClient, testNamespaceName)
	mocks.CreateTestServiceAccountWithToken(fakeClient, testNamespaceName, testutils.ServiceAccountName, "token-secret", "value", "docker-cfg")

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			params := url.Values{}

			baseURI := fmt.Sprintf("/v1/namespaces/%s/serviceaccounts/%s/token", test.args.namespace, test.args.serviceAccountName)
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

func TestCreateServiceAccountToken(t *testing.T) {
	testNamespaceName := serviceAccountNamespace + "-create-serviceaccount"

	type args struct {
		serviceAccountName         string
		namespace                  string
		existingServiceAccountName string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		args args
		want want
	}{
		"ShouldSucceedCreatingServiceAccount": {
			args: args{
				namespace:                  testNamespaceName,
				serviceAccountName:         testutils.ServiceAccountName,
				existingServiceAccountName: "",
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.NameKey: testutils.ServiceAccountName,
				},
			},
		},
		"ShouldNotSucceedCreatingServiceAccountWithExistingName": {
			args: args{
				namespace:                  testNamespaceName,
				serviceAccountName:         testutils.ServiceAccountName + "-new",
				existingServiceAccountName: testutils.ServiceAccountName + "-new",
			},
			want: want{
				statusCode: http.StatusConflict,
				response: map[string]interface{}{
					testutils.ErrorKey: fmt.Sprintf("%s, %s",
						fmt.Sprintf(controllers.ErrCouldNotCreateServiceAccount, testutils.ServiceAccountName+"-new", testNamespaceName),
						fmt.Sprintf("serviceaccounts %q already exists", testutils.ServiceAccountName+"-new"),
					),
					testutils.ReasonKey: testutils.ReasonAlreadyExists,
				},
			},
		},
	}

	setup()
	mocks.CreateTestNamespace(fakeClient, testNamespaceName)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			if test.args.existingServiceAccountName != "" {
				mocks.CreateTestServiceAccount(fakeClient, testNamespaceName, test.args.existingServiceAccountName, "")
			}
			params := url.Values{}

			baseURI := fmt.Sprintf("/v1/namespaces/%s/serviceaccounts/%s", test.args.namespace, test.args.serviceAccountName)
			request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s?%s", baseURI, params.Encode()), nil)
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

func TestDeleteServiceAccountToken(t *testing.T) {
	testNamespaceName := serviceAccountNamespace + "-delete-serviceaccount"

	type args struct {
		serviceAccountName string
		namespace          string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		args args
		want want
	}{
		"ShouldSucceedDeletingServiceAccount": {
			args: args{
				namespace:          testNamespaceName,
				serviceAccountName: testutils.ServiceAccountName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.MessageKey: fmt.Sprintf("Deleted serviceAccount successfully %q", testutils.ServiceAccountName),
				},
			},
		},
		"ShouldNotSucceedDeletingNonExistingServiceAccount": {
			args: args{
				namespace:          testNamespaceName,
				serviceAccountName: testutils.ServiceAccountName + testutils.NonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.ErrorKey: fmt.Sprintf("%s, %s",
						fmt.Sprintf(controllers.ErrCouldNotDeleteServiceAccount, testutils.ServiceAccountName+testutils.NonExistentSuffix, testNamespaceName),
						fmt.Sprintf("serviceaccounts %q not found", testutils.ServiceAccountName+testutils.NonExistentSuffix),
					),
					testutils.ReasonKey: testutils.ReasonNotFound,
				},
			},
		},
		"ShouldNotSucceedDeletingServiceAccountInNonExistingNamespace": {
			args: args{
				namespace:          testNamespaceName + testutils.NonExistentSuffix,
				serviceAccountName: testutils.ServiceAccountName,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.ErrorKey: fmt.Sprintf("%s, %s",
						fmt.Sprintf(controllers.ErrCouldNotDeleteServiceAccount, testutils.ServiceAccountName, testNamespaceName+testutils.NonExistentSuffix),
						fmt.Sprintf("serviceaccounts %q not found", testutils.ServiceAccountName),
					),
					testutils.ReasonKey: testutils.ReasonNotFound,
				},
			},
		},
	}

	setup()
	mocks.CreateTestNamespace(fakeClient, testNamespaceName)
	mocks.CreateTestServiceAccount(fakeClient, testNamespaceName, testutils.ServiceAccountName, "")

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			params := url.Values{}

			baseURI := fmt.Sprintf("/v1/namespaces/%s/serviceaccounts/%s", test.args.namespace, test.args.serviceAccountName)
			request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s?%s", baseURI, params.Encode()), nil)
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

func TestGetServiceAccounts(t *testing.T) {
	testNamespaceName := serviceAccountNamespace + "-get"
	type pagination struct {
		limit string
		page  string
	}

	type requestURI struct {
		namespace        string
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
		"ShouldSucceedGettingServiceAccounts": {
			requestURI: requestURI{
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.CountKey:           2,
					testutils.ServiceAccountsKey: []string{fmt.Sprintf("%s-1", testutils.ServiceAccountName), fmt.Sprintf("%s-2", testutils.ServiceAccountName)},
				},
			},
		},
		"ShouldSucceedGettingAllServiceAccountsWithLimitOf2": {
			requestURI: requestURI{
				namespace:        testNamespaceName,
				paginationParams: pagination{limit: "2", page: "1"},
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.CountKey:           2,
					testutils.ServiceAccountsKey: []string{fmt.Sprintf("%s-1", testutils.ServiceAccountName), fmt.Sprintf("%s-2", testutils.ServiceAccountName)},
				},
			},
		},
	}

	setup()
	mocks.CreateTestNamespace(fakeClient, testNamespaceName)
	mocks.CreateTestServiceAccount(fakeClient, testNamespaceName, fmt.Sprintf("%s-1", testutils.ServiceAccountName), "")
	mocks.CreateTestServiceAccount(fakeClient, testNamespaceName, fmt.Sprintf("%s-2", testutils.ServiceAccountName), "")

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			params := url.Values{}
			if test.requestURI.paginationParams.limit != "" {
				params.Add(middleware.LimitCtxKey, test.requestURI.paginationParams.limit)
			}

			if test.requestURI.paginationParams.page != "" {
				params.Add(middleware.PageCtxKey, test.requestURI.paginationParams.page)
			}

			baseURI := fmt.Sprintf("/v1/namespaces/%s/serviceaccounts", test.requestURI.namespace)
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

func TestGetServiceAccount(t *testing.T) {
	testNamespaceName := serviceAccountNamespace + "-get-serviceaccount"

	type args struct {
		serviceAccountName         string
		existingServiceAccountName string
		namespace                  string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		args args
		want want
	}{
		"ShouldSucceedGettingServiceAccount": {
			args: args{
				namespace:                  testNamespaceName,
				serviceAccountName:         testutils.ServiceAccountName,
				existingServiceAccountName: testutils.ServiceAccountName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.NameKey: testutils.ServiceAccountName,
				},
			},
		},
		"ShouldFailGettingNonExistingServiceAccount": {
			args: args{
				namespace:                  testNamespaceName,
				serviceAccountName:         testutils.ServiceAccountName + testutils.NonExistentSuffix,
				existingServiceAccountName: "",
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.ReasonKey: testutils.ReasonNotFound,
					testutils.ErrorKey:  fmt.Sprintf("%s, %s", fmt.Sprintf(controllers.ErrCouldNotGetServiceAccount, testutils.ServiceAccountName+testutils.NonExistentSuffix, testNamespaceName), fmt.Sprintf("serviceaccounts %q not found", testutils.ServiceAccountName+testutils.NonExistentSuffix)),
				},
			},
		},
	}

	setup()
	mocks.CreateTestNamespace(fakeClient, testNamespaceName)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			params := url.Values{}
			if test.args.existingServiceAccountName != "" {
				mocks.CreateTestServiceAccount(fakeClient, test.args.namespace, test.args.existingServiceAccountName, "")
			}
			baseURI := fmt.Sprintf("/v1/namespaces/%s/serviceaccounts/%s", test.args.namespace, test.args.serviceAccountName)
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
