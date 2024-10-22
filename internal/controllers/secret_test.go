package controllers

import (
	"fmt"
	"github.com/dana-team/platform-backend/internal/customerrors"
	"github.com/dana-team/platform-backend/internal/types"
	"github.com/dana-team/platform-backend/internal/utils"
	"github.com/dana-team/platform-backend/internal/utils/pagination"
	"github.com/dana-team/platform-backend/internal/utils/testutils"
	"github.com/dana-team/platform-backend/internal/utils/testutils/mocks"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
	"testing"
)

func TestGetSecret(t *testing.T) {
	namespaceName := testutils.SecretNamespace + "-get"
	type requestParams struct {
		name      string
		namespace string
	}

	type want struct {
		response    types.GetSecretResponse
		errorStatus metav1.StatusReason
	}

	cases := map[string]struct {
		requestParams requestParams
		want          want
	}{
		"ShouldSucceedGettingSecret": {
			requestParams: requestParams{
				namespace: namespaceName,
				name:      testutils.SecretName + "-1",
			},
			want: want{
				response: types.GetSecretResponse{
					Type:       string(corev1.SecretTypeOpaque),
					SecretName: testutils.SecretName + "-1",
					Data: []types.KeyValue{
						{
							Key:   testutils.LabelKey + "-1",
							Value: testutils.SecretDataValue,
						},
					},
				},
				errorStatus: metav1.StatusSuccess,
			},
		},
		"ShouldHandleNotFoundSecret": {
			requestParams: requestParams{
				namespace: namespaceName,
				name:      testutils.SecretName + testutils.NonExistentSuffix,
			},
			want: want{
				response:    types.GetSecretResponse{},
				errorStatus: metav1.StatusReasonNotFound,
			},
		},
		"ShouldHandleNotFoundNamespace": {
			requestParams: requestParams{
				namespace: namespaceName + testutils.NonExistentSuffix,
				name:      testutils.SecretName + testutils.NonExistentSuffix,
			},
			want: want{
				response:    types.GetSecretResponse{},
				errorStatus: metav1.StatusReasonNotFound,
			},
		},
	}
	setup()
	secretController := NewSecretController(fakeClient, mocks.GinContext(), logger)
	createTestNamespace(namespaceName, utils.AddManagedLabel(map[string]string{}))
	createTestSecret(testutils.SecretName+"-1", namespaceName, utils.AddManagedLabel(map[string]string{}))
	createTestSecret(testutils.SecretName+"-2", namespaceName, utils.AddManagedLabel(map[string]string{}))
	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			response, err := secretController.GetSecret(test.requestParams.namespace, test.requestParams.name)
			if test.want.errorStatus != metav1.StatusSuccess {
				reason := err.(customerrors.ErrorWithStatusCode).StatusReason()

				assert.Equal(t, test.want.errorStatus, reason)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.want.response, response)
		})

	}
}
func TestGetSecrets(t *testing.T) {
	namespaceName := testutils.SecretNamespace + "-getmany"
	type requestParams struct {
		namespace string
		limit     int
		page      int
	}
	type want struct {
		response types.GetSecretsResponse
	}
	cases := map[string]struct {
		requestParams requestParams
		want          want
	}{
		"ShouldSucceedGettingAllSecrets": {
			requestParams: requestParams{
				namespace: namespaceName,
			},
			want: want{
				response: types.GetSecretsResponse{
					ListMetadata: types.ListMetadata{Count: 2},
					Secrets: []types.Secret{
						{NamespaceName: namespaceName, SecretName: testutils.SecretName + "-1", Type: string(corev1.SecretTypeOpaque)},
						{NamespaceName: namespaceName, SecretName: testutils.SecretName + "-2", Type: string(corev1.SecretTypeOpaque)},
					},
				},
			},
		},
		"ShouldNotFindSecretsInNonExistingNamespaces": {
			requestParams: requestParams{
				namespace: namespaceName + testutils.NonExistentSuffix,
			},
			want: want{
				response: types.GetSecretsResponse{},
			},
		},
	}
	setup()
	createTestNamespace(namespaceName, utils.AddManagedLabel(map[string]string{}))
	createTestSecret(testutils.SecretName+"-1", namespaceName, utils.AddManagedLabel(map[string]string{}))
	createTestSecret(testutils.SecretName+"-2", namespaceName, utils.AddManagedLabel(map[string]string{}))
	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			c := mocks.GinContext()
			mocks.SetPaginationValues(c, test.requestParams.limit, test.requestParams.page)
			secretController := NewSecretController(fakeClient, c, logger)

			limit, page, _ := pagination.ExtractPaginationParamsFromCtx(c)
			response, err := secretController.GetSecrets(test.requestParams.namespace, limit, page)
			assert.NoError(t, err)
			assert.Equal(t, test.want.response, response)
		})
	}
}
func TestCreateSecret(t *testing.T) {
	namespaceName := testutils.SecretNamespace + "-create"
	type requestParams struct {
		request   types.CreateSecretRequest
		namespace string
	}
	type want struct {
		response    types.CreateSecretResponse
		errorStatus metav1.StatusReason
	}
	cases := map[string]struct {
		requestParams requestParams
		want          want
	}{
		"ShouldSucceedCreatingSecret": {
			requestParams: requestParams{
				namespace: namespaceName,
				request: types.CreateSecretRequest{
					SecretName: testutils.SecretName,
					Type:       strings.ToLower(string(corev1.SecretTypeOpaque)),
					Data: []types.KeyValue{
						{
							Key:   testutils.SecretDataKey,
							Value: testutils.SecretDataValue},
					},
				},
			},
			want: want{
				response: types.CreateSecretResponse{
					Type:          string(corev1.SecretTypeOpaque),
					SecretName:    testutils.SecretName,
					NamespaceName: namespaceName,
				},
				errorStatus: metav1.StatusSuccess,
			},
		},
		"ShouldFailCreatingSecretWithInvalidType": {
			requestParams: requestParams{
				namespace: namespaceName,
				request: types.CreateSecretRequest{
					SecretName: testutils.SecretName,
					Type:       testutils.InvalidSecretType,
					Data: []types.KeyValue{
						{
							Key:   testutils.SecretDataKey,
							Value: testutils.SecretDataValue,
						},
					},
				},
			},
			want: want{
				response:    types.CreateSecretResponse{},
				errorStatus: metav1.StatusReasonBadRequest,
			},
		},
		"ShouldFailCreatingSecretThatAlreadyExists": {
			requestParams: requestParams{
				namespace: namespaceName,
				request: types.CreateSecretRequest{
					SecretName: testutils.SecretName + "-1",
					Type:       strings.ToLower(string(corev1.SecretTypeOpaque)),
					Data: []types.KeyValue{
						{
							Key:   testutils.SecretDataKey,
							Value: testutils.SecretDataValue,
						},
					},
				},
			},
			want: want{
				response:    types.CreateSecretResponse{},
				errorStatus: metav1.StatusReasonAlreadyExists,
			},
		},
	}
	setup()
	secretController := NewSecretController(fakeClient, mocks.GinContext(), logger)
	createTestNamespace(namespaceName, utils.AddManagedLabel(map[string]string{}))
	createTestSecret(testutils.SecretName+"-1", namespaceName, utils.AddManagedLabel(map[string]string{}))
	createTestSecret(testutils.SecretName+"-2", namespaceName, utils.AddManagedLabel(map[string]string{}))
	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			response, err := secretController.CreateSecret(test.requestParams.namespace, test.requestParams.request)
			if test.want.errorStatus != metav1.StatusSuccess {
				reason := err.(customerrors.ErrorWithStatusCode).StatusReason()

				assert.Equal(t, test.want.errorStatus, reason)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.want.response, response)
		})
	}
}

func TestUpdateSecret(t *testing.T) {
	namespaceName := testutils.SecretNamespace + "-update"
	type requestParams struct {
		request   types.UpdateSecretRequest
		name      string
		namespace string
	}
	type want struct {
		response    types.UpdateSecretResponse
		errorStatus metav1.StatusReason
	}
	cases := map[string]struct {
		requestParams requestParams
		want          want
	}{
		"ShouldSucceedUpdatingSecret": {
			requestParams: requestParams{
				namespace: namespaceName,
				name:      testutils.SecretName + "-1",
				request: mocks.PrepareSecretRequestType(
					[]types.KeyValue{
						{Key: testutils.SecretDataKey, Value: testutils.SecretDataNewValue},
					},
				),
			},
			want: want{
				response: types.UpdateSecretResponse{
					Type:          string(corev1.SecretTypeOpaque),
					SecretName:    testutils.SecretName + "-1",
					NamespaceName: namespaceName,
					Data: []types.KeyValue{
						{Key: testutils.SecretDataKey, Value: testutils.SecretDataNewValue},
					},
				},
				errorStatus: metav1.StatusSuccess,
			},
		},
		"ShouldFailUpdatingNonExistingSecret": {
			requestParams: requestParams{
				namespace: namespaceName,
				name:      testutils.SecretName + testutils.NonExistentSuffix,
				request: mocks.PrepareSecretRequestType(
					[]types.KeyValue{
						{Key: testutils.SecretDataKey, Value: testutils.SecretDataNewValue},
					},
				),
			},
			want: want{
				response:    types.UpdateSecretResponse{},
				errorStatus: metav1.StatusReasonNotFound,
			},
		},
	}
	setup()
	secretController := NewSecretController(fakeClient, mocks.GinContext(), logger)
	createTestNamespace(namespaceName, utils.AddManagedLabel(map[string]string{}))
	createTestSecret(testutils.SecretName+"-1", namespaceName, utils.AddManagedLabel(map[string]string{}))
	createTestSecret(testutils.SecretName+"-2", namespaceName, utils.AddManagedLabel(map[string]string{}))
	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			response, err := secretController.UpdateSecret(test.requestParams.namespace, test.requestParams.name, test.requestParams.request)
			if test.want.errorStatus != metav1.StatusSuccess {
				reason := err.(customerrors.ErrorWithStatusCode).StatusReason()

				assert.Equal(t, test.want.errorStatus, reason)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.want.response, response)
		})
	}
}
func TestDeleteSecret(t *testing.T) {
	namespaceName := testutils.SecretNamespace + "-delete"
	type requestParams struct {
		name      string
		namespace string
	}
	type want struct {
		response    types.DeleteSecretResponse
		errorStatus metav1.StatusReason
	}
	cases := map[string]struct {
		requestParams requestParams
		want          want
	}{
		"ShouldSucceedDeletingSecret": {
			requestParams: requestParams{
				namespace: namespaceName,
				name:      testutils.SecretName + "-1",
			},
			want: want{
				response: types.DeleteSecretResponse{
					Message: fmt.Sprintf("Deleted secret %q in namespace %q successfully", testutils.SecretName+"-1", namespaceName),
				},
				errorStatus: metav1.StatusSuccess,
			},
		},
		"ShouldFailDeletingNonExistingSecret": {
			requestParams: requestParams{
				namespace: namespaceName,
				name:      testutils.SecretName + testutils.NonExistentSuffix,
			},
			want: want{
				response:    types.DeleteSecretResponse{},
				errorStatus: metav1.StatusReasonNotFound,
			},
		},
	}
	setup()
	secretController := NewSecretController(fakeClient, mocks.GinContext(), logger)
	createTestNamespace(namespaceName, utils.AddManagedLabel(map[string]string{}))
	createTestSecret(testutils.SecretName+"-1", namespaceName, utils.AddManagedLabel(map[string]string{}))
	createTestSecret(testutils.SecretName+"-2", namespaceName, utils.AddManagedLabel(map[string]string{}))
	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			response, err := secretController.DeleteSecret(test.requestParams.namespace, test.requestParams.name)
			if test.want.errorStatus != metav1.StatusSuccess {
				reason := err.(customerrors.ErrorWithStatusCode).StatusReason()

				assert.Equal(t, test.want.errorStatus, reason)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.want.response, response)
		})
	}
}
