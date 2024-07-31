package controllers

import (
	"context"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/dana-team/platform-backend/src/utils"
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
					Count: 2,
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
	podController := NewPodController(fakeClient, context.TODO(), logger)
	createTestNamespace(namespaceName, utils.AddManagedLabel(map[string]string{}))
	mocks.CreateTestPod(fakeClient, namespaceName, pod1, testutils.CappName, false)
	mocks.CreateTestPod(fakeClient, namespaceName, pod2, testutils.CappName, true)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			response, err := podController.GetPods(test.args.namespace, test.args.cappName)
			if test.want.error != "" {
				assert.ErrorContains(t, err, test.want.error)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.want.response, response)
		})
	}
}
