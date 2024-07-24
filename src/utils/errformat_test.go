package utils

import (
	"errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"testing"
)

func TestFormatK8sError(t *testing.T) {
	type args struct {
		err          error
		errorMessage string
	}
	type want struct {
		expectedCode    int
		expectedMessage string
	}
	cases := map[string]struct {
		args args
		want want
	}{
		"ShouldReturnK8sErrorWithSpecificCode": {
			args: args{
				err: &k8serrors.StatusError{
					ErrStatus: metav1.Status{
						Status: metav1.StatusFailure,
						Code:   http.StatusNotFound,
						Reason: "NotFound",
					},
				},
				errorMessage: "resource not found",
			},
			want: want{
				expectedCode:    http.StatusNotFound,
				expectedMessage: "resource not found",
			},
		},
		"ShouldReturnDefaultErrorForNonK8sError": {
			args: args{
				err:          errors.New("some other error"),
				errorMessage: "a non-K8s error occurred",
			},
			want: want{
				expectedCode:    http.StatusInternalServerError,
				expectedMessage: "a non-K8s error occurred",
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			err := FormatK8sError(tc.args.err, tc.args.errorMessage)

			var statusErr *k8serrors.StatusError
			if !errors.As(err, &statusErr) {
				t.Fatalf("expected error to be of type StatusError, got %T", err)
			}

			if int(statusErr.ErrStatus.Code) != tc.want.expectedCode {
				t.Errorf("FormatK8sError() code = %d, want %d", statusErr.ErrStatus.Code, tc.want.expectedCode)
			}

			if statusErr.ErrStatus.Message != tc.want.expectedMessage {
				t.Errorf("FormatK8sError() message = %s, want %s", statusErr.ErrStatus.Message, tc.want.expectedMessage)
			}
		})
	}
}
