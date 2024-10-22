package mocks

import (
	"github.com/dana-team/platform-backend/internal/utils/testutils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PreparePod simulates creating a pod and adding some log lines.
func PreparePod(namespace, podName, cappName string, isMultipleContainers bool) *corev1.Pod {
	labels := map[string]string{}
	containers := []corev1.Container{
		{
			Name:  testutils.TestContainerName,
			Image: testutils.Image,
		},
	}

	if cappName != "" {
		labels = map[string]string{
			testutils.ParentCappLabel: cappName,
		}

		if isMultipleContainers {
			containers = append(containers,
				corev1.Container{
					Name:  testutils.CappName,
					Image: testutils.Image,
				},
			)
		} else {
			containers = []corev1.Container{
				{
					Name:  testutils.CappName,
					Image: testutils.Image,
				},
			}
		}
	}

	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: containers,
		},
	}
}
