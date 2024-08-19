package v1

import (
	"encoding/json"
	"errors"
	"github.com/dana-team/platform-backend/src/auth"
	"github.com/dana-team/platform-backend/src/middleware"
	"github.com/dana-team/platform-backend/src/utils/testutils"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockTokenProvider struct {
	username string
	token    string
	err      error
}

const (
	tokenKey        = "token"
	validTokenKey   = "valid_token"
	validUser       = "valid_user"
	invalidUser     = "invalid_user"
	validPassword   = "valid_password"
	invalidPassword = "invalid_password"
)

func (m MockTokenProvider) ObtainToken(username, password string, logger *zap.Logger, ctx *gin.Context) (string, error) {
	return m.token, m.err
}

func (m MockTokenProvider) ObtainUsername(token string, logger *zap.Logger) (string, error) {
	return m.username, m.err
}

// setupLogin sets up a router for the Login routes
func setupLogin(tokenProvider auth.TokenProvider) (*gin.Engine, error) {
	r := gin.New()

	mockLogger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	r.Use(middleware.LoggerMiddleware(mockLogger))
	r.Use(middleware.ErrorHandlingMiddleware())
	r.POST("/login", Login(tokenProvider))

	return r, nil
}

func TestLogin(t *testing.T) {
	type args struct {
		tokenProvider auth.TokenProvider
		username      string
		password      string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		args args
		want want
	}{
		"ShouldSucceedObtainingToken": {
			args: args{
				tokenProvider: MockTokenProvider{token: validTokenKey, err: nil},
				username:      validUser,
				password:      validPassword,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					tokenKey: validTokenKey,
				},
			},
		},
		"ShouldFailWithInvalidPayload": {
			args: args{
				tokenProvider: MockTokenProvider{},
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response: map[string]interface{}{
					testutils.ErrorKey:  "Authorization header not found",
					testutils.ReasonKey: metav1.StatusReasonBadRequest,
				},
			},
		},
		"ShouldFailWithInvalidCredentials": {
			args: args{
				tokenProvider: MockTokenProvider{err: auth.ErrInvalidCredentials},
				username:      invalidUser,
				password:      invalidPassword,
			},
			want: want{
				statusCode: http.StatusUnauthorized,
				response: map[string]interface{}{
					testutils.ErrorKey:  "invalid credentials",
					testutils.ReasonKey: metav1.StatusReasonUnauthorized,
				},
			},
		},
		"ShouldFailWithInternalServerError": {
			args: args{
				tokenProvider: MockTokenProvider{err: errors.New("some internal error")},
				username:      validUser,
				password:      validPassword,
			},
			want: want{
				statusCode: http.StatusInternalServerError,
				response: map[string]interface{}{
					testutils.ErrorKey:  "some internal error",
					testutils.ReasonKey: metav1.StatusReasonInternalError,
				},
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			router, err := setupLogin(test.args.tokenProvider)
			assert.NoError(t, err)

			baseURI := "/login"
			request, err := http.NewRequest(http.MethodPost, baseURI, nil)
			assert.NoError(t, err)
			request.Header.Set(testutils.ContentType, testutils.ApplicationJson)
			if test.args.username != "" {
				request.SetBasicAuth(test.args.username, test.args.password)
			}

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			assert.Equal(t, test.want.statusCode, writer.Code)

			var response map[string]interface{}
			err = json.Unmarshal(writer.Body.Bytes(), &response)
			assert.NoError(t, err)

			wantResponseJSON, err := json.Marshal(test.want.response)
			assert.NoError(t, err)
			var wantResponseNormalized map[string]interface{}
			err = json.Unmarshal(wantResponseJSON, &wantResponseNormalized)
			assert.NoError(t, err)
			assert.Equal(t, wantResponseNormalized, response)
		})
	}
}
