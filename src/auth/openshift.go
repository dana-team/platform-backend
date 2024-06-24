package auth

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/dana-team/platform-backend/src/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"net/url"
	"os"
)

const (
	httpXCSRFTokenHeader    = "X-CSRF-Token"
	httpLocationHeader      = "Location"
	httpAuthorizationHeader = "Authorization"
	httpBearerTokenPrefix   = "Bearer"
	stateToken              = "state"
	keyCodeParam            = "code"
)

const (
	envInsecureSkipVerify = "INSECURE_SKIP_VERIFY"
	envKubeClientID       = "KUBE_CLIENT_ID"
	envKubeAuthURL        = "KUBE_AUTH_URL"
	envKubeTokenURL       = "KUBE_TOKEN_URL"
	envKubeUserInfoURL    = "KUBE_USERINFO_URL"
)

// OpenshiftUserInfo represents the structure of the user info response from OpenShift.
type OpenshiftUserInfo struct {
	Metadata struct {
		Name string `json:"name"`
	} `json:"metadata"`
}

// ObtainOpenshiftToken obtains an OAuth token from OpenShift using the provided username and password.
// It returns the access token on success, or an error on failure.
// This code was taken from https://gist.github.com/kadel/c30b3085e2e90a93393b99a2b39f4806, with minor adjustments.
func ObtainOpenshiftToken(username, password string, logger *zap.Logger, ctx *gin.Context) (string, error) {
	skipTlsVerify, err := utils.GetEnvBool(envInsecureSkipVerify, true)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to parse environment variable %q", envInsecureSkipVerify), zap.Error(err))
		return "", err
	}

	httpClient := createHTTPClient(skipTlsVerify)
	ctxWithClient := context.WithValue(ctx.Request.Context(), oauth2.HTTPClient, httpClient)
	codeHTTPClient := createCodeHTTPClient(skipTlsVerify)
	conf := getOAuthConfig()

	authCodeReq, err := createAuthCodeRequest(conf, username, password)
	if err != nil {
		logger.Error("failed to create AuthCode request", zap.Error(err))
		return "", err
	}

	logger.Debug(fmt.Sprintf("tyring to get auth code from: %q", authCodeReq.URL.String()))
	authCodeResp, err := codeHTTPClient.Do(authCodeReq)
	if err != nil {
		logger.Error("failed to obtain auth code", zap.Error(err))
		return "", fmt.Errorf("failed to obtain auth code: %v", err)
	}

	if err := checkAuthCodeResponse(authCodeResp); err != nil {
		logger.Error("failed to check AuthCode response", zap.Error(err))
		return "", err
	}
	logger.Debug(fmt.Sprintf("fetched auth code successfully from: %q with status: %q", authCodeResp.Request.URL.String(), authCodeResp.StatusCode))

	code, err := parseAuthCode(authCodeResp)
	if err != nil {
		logger.Error("failed to parse AuthCode", zap.Error(err))
		return "", err
	}

	return exchangeCodeForToken(ctxWithClient, conf, code)
}

// ObtainOpenshiftUsername fetches the username from the OpenShift userinfo endpoint.
func ObtainOpenshiftUsername(token string, logger *zap.Logger) (string, error) {
	userInfo, err := fetchOpenshiftUserInfo(token, logger)
	if err != nil {
		logger.Error("failed to fetch Openshift user info", zap.Error(err))
		return "", fmt.Errorf("failed to obtain OpenShift username: %v", err)
	}

	return userInfo.Metadata.Name, nil
}

// getOAuthConfig returns an OAuth2 configuration based on environment variables.
func getOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID: os.Getenv(envKubeClientID),
		Endpoint: oauth2.Endpoint{
			AuthURL:  os.Getenv(envKubeAuthURL),
			TokenURL: os.Getenv(envKubeTokenURL),
		},
	}
}

// createHTTPClient creates an HTTP client with TLS configuration.
func createHTTPClient(skipTlsVerify bool) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: skipTlsVerify,
			},
		},
	}
}

// createCodeHTTPClient creates an HTTP client for requesting the authorization code.
// It creates a special client used only to request code, using the
// CheckRedirect field which ensures that it is not following 302 redirects,
// and the Location header can be parsed to get code from it.
func createCodeHTTPClient(skipTlsVerify bool) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: skipTlsVerify,
			},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

// createAuthCodeRequest creates an HTTP request for obtaining the authorization code
// from the authorization URL, using the provided username and password for basic authentication.
func createAuthCodeRequest(conf *oauth2.Config, username, password string) (*http.Request, error) {
	authCodeUrl := conf.AuthCodeURL(stateToken, oauth2.AccessTypeOffline)
	authCodeReq, err := http.NewRequest("GET", authCodeUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth code request: %v", err)
	}
	authCodeReq.Header.Set(httpXCSRFTokenHeader, "1")
	authCodeReq.SetBasicAuth(username, password)
	return authCodeReq, nil
}

// checkAuthCodeResponse checks the HTTP response for obtaining the authorization code.
// It returns an error if the status code indicates a failure or if the credentials are invalid.
func checkAuthCodeResponse(authCodeResp *http.Response) error {
	body, err := io.ReadAll(authCodeResp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	if authCodeResp.StatusCode == http.StatusUnauthorized {
		return ErrInvalidCredentials
	}

	if authCodeResp.StatusCode != http.StatusFound {
		return fmt.Errorf("unexpected status code: %d, body: %s", authCodeResp.StatusCode, string(body))
	}

	return nil
}

// parseAuthCode parses the authorization code from the Location header
// of the HTTP response.
func parseAuthCode(authCodeResp *http.Response) (string, error) {
	urlLocation, err := url.Parse(authCodeResp.Header.Get(httpLocationHeader))
	if err != nil {
		return "", fmt.Errorf("failed to parse %q header: %v", httpLocationHeader, err)
	}

	code := urlLocation.Query().Get(keyCodeParam)
	if code == "" {
		return "", fmt.Errorf("authorization code not found in %q header", httpLocationHeader)
	}

	return code, nil
}

// exchangeCodeForToken exchanges the authorization code for an access token
// using the provided OAuth2 configuration.
func exchangeCodeForToken(ctx context.Context, conf *oauth2.Config, code string) (string, error) {
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		return "", fmt.Errorf("failed to exchange code for token: %v", err)
	}
	return tok.AccessToken, nil
}

// fetchOpenshiftUserInfo retrieves user information from the OpenShift userinfo endpoint.
func fetchOpenshiftUserInfo(token string, logger *zap.Logger) (*OpenshiftUserInfo, error) {
	userInfo := OpenshiftUserInfo{}

	skipTlsVerify, err := utils.GetEnvBool(envInsecureSkipVerify, true)
	if err != nil {
		return nil, err
	}

	httpClient := createHTTPClient(skipTlsVerify)

	req, err := createUserInfoRequest(token, os.Getenv(envKubeUserInfoURL))
	if err != nil {
		return nil, fmt.Errorf("failed to create user info request: %v", err)
	}

	logger.Debug(fmt.Sprintf("trying to fetch user info from: %q", req.URL.String()))
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch userinfo: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch userinfo, status code: %d", resp.StatusCode)
	}

	logger.Debug(fmt.Sprintf("fetched user info successfully from: %q with status: %q", resp.Request.URL.String(), resp.StatusCode))

	if err := decodeUserInfoResponse(resp, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode userinfo response: %v", err)
	}

	return &userInfo, nil
}

// createUserInfoRequest creates an HTTP request for fetching userinfo from the OpenShift userinfo endpoint.
func createUserInfoRequest(token, userInfoUrl string) (*http.Request, error) {
	req, err := http.NewRequest("GET", userInfoUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create user info request: %v", err)
	}
	req.Header.Set(httpAuthorizationHeader, fmt.Sprintf("%s %s", httpBearerTokenPrefix, token))

	return req, nil
}

// decodeUserInfoResponse reads the body of the HTTP response from the userinfo endpoint and decodes the JSON response body into the OpenshiftUserInfo struct.
func decodeUserInfoResponse(resp *http.Response, userInfo *OpenshiftUserInfo) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read userinfo response body: %v", err)
	}

	if err := json.Unmarshal(body, userInfo); err != nil {
		return fmt.Errorf("failed to decode userinfo response: %v", err)
	}

	return nil
}
