package controllers

import (
	"context"
	"github.com/dana-team/platform-backend/internal/types"
	"github.com/dana-team/platform-backend/internal/utils"
	"github.com/dana-team/platform-backend/internal/utils/testutils"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

const baseNsName = "test-namespace"

func TestCreateNamespace(t *testing.T) {
	existingNSName := baseNsName + "-exists"
	nsToCreate := baseNsName + "-create"
	type requestParams struct {
		namespace string
	}

	type want struct {
		response    types.Namespace
		errorStatus metav1.StatusReason
	}
	cases := map[string]struct {
		requestParams requestParams
		want          want
	}{
		"ShouldSucceedCreatingNamespace": {
			requestParams: requestParams{
				namespace: nsToCreate,
			},
			want: want{
				response: types.Namespace{
					Name: baseNsName + "-create",
				},
				errorStatus: metav1.StatusSuccess,
			},
		},
		"ShouldFailCreatingExistingNamespace": {
			requestParams: requestParams{
				namespace: existingNSName,
			},
			want: want{
				response:    types.Namespace{},
				errorStatus: metav1.StatusReasonAlreadyExists,
			},
		},
	}
	setup()
	namespaceController := NewNamespaceController(fakeClient, context.TODO(), logger)
	createTestNamespace(existingNSName, map[string]string{})

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			response, err := namespaceController.CreateNamespace(test.requestParams.namespace)
			if test.want.errorStatus != metav1.StatusSuccess {
				reason := err.(errors.APIStatus).Status().Reason

				assert.Equal(t, test.want.errorStatus, reason)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.want.response, response)
		})
	}
}

func TestGetNamespace(t *testing.T) {
	nsName := baseNsName + "-get"
	type requestParams struct {
		namespace string
	}

	type want struct {
		response    types.Namespace
		errorStatus metav1.StatusReason
	}
	cases := map[string]struct {
		requestParams requestParams
		want          want
	}{
		"ShouldSucceedFetchNamespace": {
			requestParams: requestParams{
				namespace: nsName,
			},
			want: want{
				response: types.Namespace{
					Name: nsName,
				},
				errorStatus: metav1.StatusSuccess,
			},
		},
		"ShouldFailToFetchNotFound": {
			requestParams: requestParams{
				namespace: nsName + testutils.NonExistentSuffix,
			},
			want: want{
				response:    types.Namespace{},
				errorStatus: metav1.StatusReasonNotFound,
			},
		},
	}
	setup()
	namespaceController := NewNamespaceController(fakeClient, context.TODO(), logger)
	createTestNamespace(nsName, map[string]string{})

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			response, err := namespaceController.GetNamespace(test.requestParams.namespace)
			if test.want.errorStatus != metav1.StatusSuccess {
				reason := err.(errors.APIStatus).Status().Reason

				assert.Equal(t, test.want.errorStatus, reason)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.want.response, response)
		})
	}
}

func TestGetNamespaces(t *testing.T) {
	firstNsName := baseNsName + "-1"
	secondNsName := baseNsName + "-2"

	type want struct {
		response    types.NamespaceList
		errorStatus metav1.StatusReason
	}
	cases := map[string]struct {
		want want
	}{
		"ShouldSucceedFetchingNamespaces": {
			want: want{
				response:    types.NamespaceList{Count: 2, Namespaces: []types.Namespace{{Name: firstNsName}, {Name: secondNsName}}},
				errorStatus: metav1.StatusSuccess,
			},
		},
	}
	setup()
	namespaceController := NewNamespaceController(fakeClient, context.TODO(), logger)
	createTestNamespace(firstNsName, map[string]string{utils.ManagedLabel: utils.ManagedLabelValue})
	createTestNamespace(secondNsName, map[string]string{utils.ManagedLabel: utils.ManagedLabelValue})

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			response, err := namespaceController.GetNamespaces()
			if test.want.errorStatus != metav1.StatusSuccess {
				reason := err.(errors.APIStatus).Status().Reason

				assert.Equal(t, test.want.errorStatus, reason)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.want.response, response)
		})
	}
}

func TestDeleteNamespace(t *testing.T) {
	nsToDelete := baseNsName + "-delete"
	type requestParams struct {
		namespace string
	}

	type want struct {
		errorStatus metav1.StatusReason
	}
	cases := map[string]struct {
		requestParams requestParams
		want          want
	}{
		"ShouldSucceedDeletingNamespace": {
			requestParams: requestParams{
				namespace: nsToDelete,
			},
			want: want{
				errorStatus: metav1.StatusSuccess,
			},
		},
		"ShouldFailDeletingNotFound": {
			requestParams: requestParams{
				namespace: baseNsName + testutils.NonExistentSuffix,
			},
			want: want{
				errorStatus: metav1.StatusReasonNotFound,
			},
		},
	}
	setup()
	namespaceController := NewNamespaceController(fakeClient, context.TODO(), logger)
	createTestNamespace(nsToDelete, map[string]string{})

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			err := namespaceController.DeleteNamespace(test.requestParams.namespace)
			if test.want.errorStatus != metav1.StatusSuccess {
				reason := err.(errors.APIStatus).Status().Reason

				assert.Equal(t, test.want.errorStatus, reason)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
