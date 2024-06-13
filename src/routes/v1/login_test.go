package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/dana-team/platform-backend/src/auth"
	"github.com/dana-team/platform-backend/src/middleware"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockTokenProvider struct {
	Username string
	Token    string
	Err      error
}

func (m MockTokenProvider) ObtainToken(username, password string, logger *zap.Logger, ctx *gin.Context) (string, error) {
	return m.Token, m.Err
}

func (m MockTokenProvider) ObtainUsername(token string, logger *zap.Logger) (string, error) {
	return m.Username, m.Err
}
func Test_Login(t *testing.T) {
	type args struct {
		tokenProvider auth.TokenProvider
	}
	type want struct {
		statusCode int
		response   map[string]string
	}
	cases := map[string]struct {
		args    args
		request interface{}
		want    want
	}{
		"ShouldSuccessObtainingToken": {
			args: args{
				tokenProvider: MockTokenProvider{Token: "valid_token", Err: nil},
			},
			request: map[string]string{
				"username": "valid_user",
				"password": "valid_password",
			},
			want: want{
				statusCode: http.StatusOK,
				response:   map[string]string{"token": "valid_token"},
			},
		},
		"ShouldFailWithInvalidPayload": {
			args: args{
				tokenProvider: MockTokenProvider{},
			},
			request: `invalid payload`,
			want: want{
				statusCode: http.StatusBadRequest,
				response:   map[string]string{"error": "Invalid request payload"},
			},
		},
		"ShouldFailWithInvalidCredentials": {
			args: args{
				tokenProvider: MockTokenProvider{Err: auth.ErrInvalidCredentials},
			},
			request: map[string]string{
				"username": "invalid_user",
				"password": "invalid_password",
			},
			want: want{
				statusCode: http.StatusUnauthorized,
				response:   map[string]string{"error": "Invalid credentials"},
			},
		},
		"ShouldFailWithInternalServerError": {
			args: args{
				tokenProvider: MockTokenProvider{Err: errors.New("some internal error")},
			},
			request: map[string]string{
				"username": "user",
				"password": "password",
			},
			want: want{
				statusCode: http.StatusInternalServerError,
				response:   map[string]string{"error": "Internal server error"},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			r := gin.New()

			mockLogger, _ := zap.NewDevelopment()
			r.Use(middleware.LoggerMiddleware(mockLogger))
			r.POST("/login", Login(tc.args.tokenProvider))

			payload, _ := json.Marshal(tc.request)

			req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tc.want.statusCode, w.Code)

			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tc.want.response, response)
		})
	}
}
