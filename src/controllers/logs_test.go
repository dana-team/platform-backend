package controllers

import (
	"context"
	"fmt"
	"github.com/dana-team/platform-backend/src/utils/testutils"
	"github.com/dana-team/platform-backend/src/utils/testutils/mocks"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
	"testing"

	corev1 "k8s.io/api/core/v1"
)

const (
	pod1    = testutils.PodName + "-1"
	pod2    = testutils.PodName + "-2"
	pod3    = testutils.PodName + "-3"
	cappPod = testutils.PodName
)

func TestFetchCappPodName(t *testing.T) {
	type args struct {
		podName string
		pods    *corev1.PodList
	}
	type want struct {
		name  string
		found bool
	}

	mockPodList := &corev1.PodList{
		Items: []corev1.Pod{
			{ObjectMeta: metav1.ObjectMeta{Name: pod1}},
			{ObjectMeta: metav1.ObjectMeta{Name: pod2}},
			{ObjectMeta: metav1.ObjectMeta{Name: pod3}},
		},
	}

	cases := map[string]struct {
		args args
		want want
	}{
		"EmptyPodNameShouldReturnFirstPod": {
			args: args{
				podName: "",
				pods:    mockPodList,
			},
			want: want{
				name:  pod1,
				found: true,
			},
		},
		"ShouldFindExistingPodName": {
			args: args{
				podName: pod2,
				pods:    mockPodList,
			},
			want: want{
				name:  pod2,
				found: true,
			},
		},
		"NonExistingPodNameShouldReturnFalse": {
			args: args{
				podName: testutils.PodName + testutils.NonExistentSuffix,
				pods:    mockPodList,
			},
			want: want{
				name:  testutils.PodName + testutils.NonExistentSuffix,
				found: false,
			},
		},
		"EmptyPodListShouldReturnFalse": {
			args: args{
				podName: pod1,
				pods:    &corev1.PodList{},
			},
			want: want{
				name:  pod1,
				found: false,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			actualName, actualFound := FetchCappPodName(tc.args.podName, tc.args.pods)
			if actualName != tc.want.name {
				t.Errorf("Expected pod name %s, but got %s", tc.want.name, actualName)
			}
			if actualFound != tc.want.found {
				t.Errorf("Expected found %v, but got %v", tc.want.found, actualFound)
			}
		})
	}
}

func TestFetchPodLogs(t *testing.T) {
	type args struct {
		client        kubernetes.Interface
		namespace     string
		podName       string
		containerName string
		previous      bool
	}
	type want struct {
		errContains string
	}

	cases := map[string]struct {
		args args
		want want
	}{
		"ShouldFailGettingLogsOnNonExistingPod": {
			args: args{
				client:        fakeClient,
				namespace:     testutils.TestNamespace,
				podName:       testutils.PodName + testutils.NonExistentSuffix,
				containerName: testutils.TestContainerName,
				previous:      true,
			},
			want: want{
				errContains: "error opening log stream",
			},
		},
	}

	setup()
	mocks.CreateTestPod(fakeClient, testutils.TestNamespace, pod1, "", false)
	mockLogger, _ := zap.NewDevelopment()

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := FetchPodLogs(context.TODO(), tc.args.client, tc.args.namespace, tc.args.podName, tc.args.containerName, tc.args.previous, mockLogger)
			if tc.want.errContains == "" {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			} else {
				if err == nil || !strings.Contains(err.Error(), tc.want.errContains) {
					t.Errorf("Expected error containing '%s', but got: %v", tc.want.errContains, err)
				}
			}
		})
	}
}

func TestFetchCappLogs(t *testing.T) {
	type args struct {
		client        kubernetes.Interface
		namespace     string
		cappName      string
		containerName string
		podName       string
		previous      bool
	}
	type want struct {
		errContains string
	}

	cases := map[string]struct {
		args args
		want want
	}{
		"ShouldReturnNoPodsForNonExistingCapp": {
			args: args{
				client:        fakeClient,
				namespace:     testutils.TestNamespace,
				cappName:      testutils.CappName + testutils.NonExistentSuffix,
				containerName: testutils.TestContainerName,
				podName:       cappPod,
				previous:      true,
			},
			want: want{
				errContains: "no pods found for Capp",
			},
		},
		"ShouldFailOnNonExistingPod": {
			args: args{
				client:        fakeClient,
				namespace:     testutils.TestNamespace,
				cappName:      testutils.CappName,
				containerName: testutils.TestContainerName,
				podName:       testutils.PodName + testutils.NonExistentSuffix,
				previous:      false,
			},
			want: want{
				errContains: fmt.Sprintf("no pods found for Capp %q in namespace %q", testutils.CappName, testutils.TestNamespace),
			},
		},
	}

	setup()
	mocks.CreateTestPod(fakeClient, testutils.TestNamespace, cappPod, testutils.CappName, false)
	mockLogger, _ := zap.NewDevelopment()

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := FetchCappLogs(context.TODO(), tc.args.client, tc.args.namespace, tc.args.cappName, tc.args.containerName, tc.args.podName, tc.args.previous, mockLogger)
			if tc.want.errContains == "" {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			} else {
				if err == nil || !strings.Contains(err.Error(), tc.want.errContains) {
					t.Errorf("Expected error containing '%s', but got: %v", tc.want.errContains, err)
				}
			}
		})
	}
}
