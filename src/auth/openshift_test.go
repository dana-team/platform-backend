package auth

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestObtainOpenshiftToken(t *testing.T) {
	os.Setenv("INSECURE_SKIP_VERIFY", "true")
	os.Setenv("KUBE_CLIENT_ID", "clientid")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/auth" {
			w.Header().Set("Location", "https://example.com/callback?code=testcode")
			w.WriteHeader(http.StatusFound)
			return
		}

		if r.URL.Path == "/token" {
			if r.FormValue("code") == "testcode" {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"access_token": "test_access_token", "token_type": "bearer"}`))
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
			return
		}
	}))
	defer ts.Close()

	os.Setenv("KUBE_AUTH_URL", ts.URL+"/auth")
	os.Setenv("KUBE_TOKEN_URL", ts.URL+"/token")

	type args struct {
		username string
		password string
	}
	type want struct {
		expectedToken    string
		expectedError    error
		mockServerStatus int
	}
	cases := map[string]struct {
		args args
		want want
	}{
		"ShouldSuccessObtainingToken": {
			args: args{
				username: "test_user",
				password: "test_password",
			},
			want: want{
				expectedToken: "test_access_token",
				expectedError: nil,
			},
		},
		"ShouldFailWithUnauthorized": {
			args: args{
				username: "invalid_user",
				password: "invalid_password",
			},
			want: want{
				expectedToken:    "",
				expectedError:    ErrInvalidCredentials,
				mockServerStatus: http.StatusUnauthorized,
			},
		},
		"ShouldFailWithFailedRequest": {
			args: args{
				username: "test_user",
				password: "test_password",
			},
			want: want{
				expectedToken:    "",
				expectedError:    errors.New("unexpected status code: 500, body: Internal Server Error"),
				mockServerStatus: http.StatusInternalServerError,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if tc.want.mockServerStatus != 0 {
				ts.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tc.want.mockServerStatus)
					if tc.want.mockServerStatus != http.StatusFound {
						_, _ = w.Write([]byte("Internal Server Error"))
					}
				})
			}

			logger, _ := zap.NewProduction()
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequest(http.MethodPost, ts.URL+"/auth", nil).WithContext(context.Background())

			token, err := ObtainOpenshiftToken(tc.args.username, tc.args.password, logger, c)
			if err != nil && tc.want.expectedError == nil {
				t.Errorf("ObtainOpenshiftToken() unexpected error: %v", err)
				return
			}
			if tc.want.expectedError != nil && (err == nil || err.Error() != tc.want.expectedError.Error()) {
				t.Errorf("ObtainOpenshiftToken() expected error: %v, got: %v", tc.want.expectedError, err)
				return
			}
			if token != tc.want.expectedToken {
				t.Errorf("ObtainOpenshiftToken() expected token: %v, got: %v", tc.want.expectedToken, token)
			}
		})
	}
}

func TestObtainOpenshiftUsername(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/apis/user.openshift.io/v1/users/~" {
			if r.Header.Get("Authorization") != "Bearer valid_token" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"metadata": {"name": "test_user"}}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	// Set environment variables
	os.Setenv("INSECURE_SKIP_VERIFY", "true")
	os.Setenv("KUBE_USERINFO_URL", ts.URL+"/apis/user.openshift.io/v1/users/~")

	type args struct {
		token string
	}
	type want struct {
		expectedUsername string
		expectedError    error
	}
	cases := map[string]struct {
		args args
		want want
	}{
		"ShouldSuccessObtainingUsername": {
			args: args{
				token: "valid_token",
			},
			want: want{
				expectedUsername: "test_user",
				expectedError:    nil,
			},
		},
		"ShouldFailWithUnauthorized": {
			args: args{
				token: "invalid_token",
			},
			want: want{
				expectedUsername: "",
				expectedError:    errors.New("failed to obtain OpenShift username: failed to fetch userinfo, status code: 401"),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			logger, _ := zap.NewProduction()
			username, err := ObtainOpenshiftUsername(tc.args.token, logger)
			if err != nil && tc.want.expectedError == nil {
				t.Errorf("ObtainOpenshiftUsername() unexpected error: %v", err)
				return
			}
			if tc.want.expectedError != nil && (err == nil || err.Error() != tc.want.expectedError.Error()) {
				t.Errorf("ObtainOpenshiftUsername() expected error: %v, got: %v", tc.want.expectedError, err)
				return
			}
			if username != tc.want.expectedUsername {
				t.Errorf("ObtainOpenshiftUsername() expected username: %v, got: %v", tc.want.expectedUsername, username)
			}
		})
	}
}
