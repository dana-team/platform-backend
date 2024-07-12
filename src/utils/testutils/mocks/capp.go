package mocks

import (
	"context"
	"fmt"

	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/dana-team/platform-backend/src/utils/testutils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	knativeapis "knative.dev/pkg/apis"
	knativev1 "knative.dev/serving/pkg/apis/serving/v1"
	knativev1beta1 "knative.dev/serving/pkg/apis/serving/v1beta1"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
)

// PrepareCapp returns a mock Capp object.
func PrepareCapp(name, namespace string, labels, annotations map[string]string) cappv1alpha1.Capp {
	return cappv1alpha1.Capp{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Annotations: annotations,
			Labels:      labels,
		},
		Spec:   PrepareCappSpec(),
		Status: PrepareCappStatus(name, namespace),
	}
}

// PrepareCappWithHostname returns a mock Capp object with Hostname set in the spec.
func PrepareCappWithHostname(name, namespace string, labels, annotations map[string]string) cappv1alpha1.Capp {
	return cappv1alpha1.Capp{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Annotations: annotations,
			Labels:      labels,
		},
		Spec:   PrepareCappSpecWithHostname(),
		Status: PrepareCappStatusWithHostname(name, namespace),
	}
}

// PrepareCappSpec returns a mock Capp spec.
func PrepareCappSpec() cappv1alpha1.CappSpec {
	return cappv1alpha1.CappSpec{
		ConfigurationSpec: knativev1.ConfigurationSpec{
			Template: knativev1.RevisionTemplateSpec{
				Spec: knativev1.RevisionSpec{
					PodSpec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Image: testutils.CappImage,
								Name:  testutils.ContainerName,
							},
						},
					},
				},
			},
		},
	}
}

// PrepareCappSpecWithHostname returns a mock Capp spec with Hostname set.
func PrepareCappSpecWithHostname() cappv1alpha1.CappSpec {
	return cappv1alpha1.CappSpec{
		ConfigurationSpec: knativev1.ConfigurationSpec{
			Template: knativev1.RevisionTemplateSpec{
				Spec: knativev1.RevisionSpec{
					PodSpec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Image: testutils.CappImage,
								Name:  testutils.ContainerName,
							},
						},
					},
				},
			},
		},
		RouteSpec: cappv1alpha1.RouteSpec{
			Hostname: testutils.Hostname + "." + testutils.Domain,
		},
	}
}

// PrepareCappStatus returns a mock Capp status.
func PrepareCappStatus(name, namespace string) cappv1alpha1.CappStatus {
	return cappv1alpha1.CappStatus{
		KnativeObjectStatus: knativev1.ServiceStatus{
			RouteStatusFields: knativev1.RouteStatusFields{
				URL: knativeapis.HTTPS(fmt.Sprintf("%s-%s.%s", name, namespace, testutils.Domain)),
			},
		},
	}
}

// PrepareCappStatusWithHostname returns a mock Capp status with Hostname set.
func PrepareCappStatusWithHostname(name, namespace string) cappv1alpha1.CappStatus {
	return cappv1alpha1.CappStatus{
		KnativeObjectStatus: knativev1.ServiceStatus{
			RouteStatusFields: knativev1.RouteStatusFields{
				URL: knativeapis.HTTPS(fmt.Sprintf("%s-%s.%s", name, namespace, testutils.Domain)),
			},
		},
		RouteStatus: cappv1alpha1.RouteStatus{
			DomainMappingObjectStatus: knativev1beta1.DomainMappingStatus{
				URL: knativeapis.HTTPS(fmt.Sprintf("%s.%s", testutils.Hostname, testutils.Domain)),
			},
		},
	}
}

func PrepareUpdateCappType(labels, annotations []types.KeyValue) types.UpdateCapp {
	return types.UpdateCapp{
		Annotations: annotations,
		Labels:      labels,
		Spec:        PrepareCappSpec(),
	}
}

// PrepareCreateCappType returns a CreateCapp type object.
func PrepareCreateCappType(name string, labels, annotations []types.KeyValue) types.CreateCapp {
	return types.CreateCapp{
		Metadata: types.CreateMetadata{
			Name: name,
		},
		Annotations: annotations,
		Labels:      labels,
		Spec:        PrepareCappSpec(),
	}
}

func PrepareCappMetadata(name, namespace string) types.Metadata {
	return types.Metadata{
		Name:      name,
		Namespace: namespace,
	}
}

func PrepareCappSummary(name string, namespace string) types.CappSummary {
	return types.CappSummary{
		Name:   name,
		Images: []string{testutils.CappImage},
		URL:    fmt.Sprintf("https://%s-%s.%s", name, namespace, testutils.Domain),
	}
}

// CreateTestCapp creates a test Capp object.
func CreateTestCapp(name, namespace string, labels, annotations map[string]string, dynClient runtimeClient.WithWatch) {
	cappRevision := PrepareCapp(name, namespace, labels, annotations)
	err := dynClient.Create(context.TODO(), &cappRevision)
	if err != nil {
		panic(err)
	}
}

// CreateTestCappWithHostname creates a test Capp object with a hostname.
func CreateTestCappWithHostname(name, namespace string, labels, annotations map[string]string, dynClient runtimeClient.WithWatch) {
	capp := PrepareCappWithHostname(name, namespace, labels, annotations)
	err := dynClient.Create(context.TODO(), &capp)
	if err != nil {
		panic(err)
	}
}
