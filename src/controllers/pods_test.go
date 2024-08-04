package controllers

import (
	"github.com/dana-team/platform-backend/src/types"
	"github.com/dana-team/platform-backend/src/utils"
	"github.com/dana-team/platform-backend/src/utils/pagination"
	"github.com/dana-team/platform-backend/src/utils/testutils"
	"github.com/dana-team/platform-backend/src/utils/testutils/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetPods(t *testing.T) {
	namespaceName := testutils.TestNamespace + "-getPods"
	type args struct {
		namespace string
		cappName  string
		limit     int
		page      int
	}
	type want struct {
		response types.GetPodsResponse
		error    string
	}
	cases := map[string]struct {
		args args
		want want
	}{
		"ShouldSucceedGettingAllPods": {
			args: args{
				namespace: namespaceName,
				cappName:  testutils.CappName,
			},
			want: want{
				response: types.GetPodsResponse{
					ListMetadata: types.ListMetadata{Count: 2},
					Pods: []types.Pod{
						{PodName: pod1},
						{PodName: pod2},
					},
				},
			},
		},
		"ShouldNotFindPodsInNonExistingNamespaces": {
			args: args{
				cappName:  testutils.CappName + testutils.NonExistentSuffix,
				namespace: namespaceName + testutils.NonExistentSuffix,
			},
			want: want{
				response: types.GetPodsResponse{},
			},
		},
	}

	setup()
	createTestNamespace(namespaceName, utils.AddManagedLabel(map[string]string{}))
	mocks.CreateTestPod(fakeClient, namespaceName, pod1, testutils.CappName, false)
	mocks.CreateTestPod(fakeClient, namespaceName, pod2, testutils.CappName, true)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			c := mocks.GinContext()
			mocks.SetPaginationValues(c, test.args.limit, test.args.page)
			podController := NewPodController(fakeClient, c, logger)

			limit, page, _ := pagination.ExtractPaginationParamsFromCtx(c)
			response, err := podController.GetPods(test.args.namespace, test.args.cappName, limit, page)
			if test.want.error != "" {
				assert.ErrorContains(t, err, test.want.error)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.want.response, response)
		})
	}
}
