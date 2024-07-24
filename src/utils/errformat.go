package utils

import (
	"errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

const (
	statusUnknown = "Unknown"
)

// FormatK8sError formats a Kubernetes error to include relevant details.
func FormatK8sError(err error, errorMessage string) error {
	var statusErr *k8serrors.StatusError
	if errors.As(err, &statusErr) {
		return &k8serrors.StatusError{
			ErrStatus: metav1.Status{
				Status:  metav1.StatusFailure,
				Code:    statusErr.ErrStatus.Code,
				Reason:  statusErr.ErrStatus.Reason,
				Details: statusErr.ErrStatus.Details,
				Message: errorMessage,
			},
		}
	}

	return &k8serrors.StatusError{
		ErrStatus: metav1.Status{
			Status:  metav1.StatusFailure,
			Code:    http.StatusInternalServerError,
			Reason:  statusUnknown,
			Message: errorMessage,
		},
	}
}
