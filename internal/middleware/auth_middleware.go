package middleware

import (
	"fmt"
	"github.com/dana-team/platform-backend/internal/customerrors"
	multicluster "github.com/oam-dev/cluster-gateway/pkg/apis/cluster/transport"
	"k8s.io/apimachinery/pkg/runtime"
	"os"
	"strings"

	"github.com/dana-team/platform-backend/internal/auth"
	"github.com/dana-team/platform-backend/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	zapctrl "sigs.k8s.io/controller-runtime/pkg/log/zap"
)

const (
	httpAuthorizationHeader     = "Authorization"
	httpBearerTokenPrefix       = "Bearer"
	validBearerTokenPartsLength = 2
	httpBearerTokenPrefixIndex  = 0
	httpBearerTokenIndex        = 1
	WebsocketTokenHeader        = "Sec-Websocket-Protocol"
)

const (
	KubeClientCtxKey    = "kubeClient"
	DynamicClientCtxKey = "dynClient"
	TokenCtxKey         = "token"
	ConfigKey           = "config"
)

const (
	envKubeAPIServer      = "KUBE_API_SERVER"
	envInsecureSkipVerify = "INSECURE_SKIP_VERIFY"
)

// TokenAuthMiddleware validates the Authorization header and sets up Kubernetes client.
func TokenAuthMiddleware(tokenProvider auth.TokenProvider, scheme *runtime.Scheme) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger, err := GetLogger(c)
		if AddErrorToContext(c, err) {
			return
		}

		token, err := validateToken(c)
		if err != nil {
			logger.Error("Failed to obtain OpenShift token", zap.Error(err))
			AddErrorToContext(c, customerrors.NewUnauthorizedError("failed to obtain OpenShift token"))
			c.Abort()
			return
		}

		username, err := tokenProvider.ObtainUsername(token, logger)
		if err != nil {
			logger.Error("Failed to get user info", zap.Error(err))
			AddErrorToContext(c, customerrors.NewInternalServerError("failed to get user info"))
			c.Abort()
			return
		}
		userLogger := logger.With(zap.String("user", username))

		config, err := createKubernetesConfig(token, os.Getenv(envKubeAPIServer))
		if err != nil {
			userLogger.Error("Failed to create Kubernetes client config", zap.Error(err))
			AddErrorToContext(c, customerrors.NewInternalServerError("failed to create Kubernetes client config"))
			c.Abort()
			return
		}

		kubeClient, err := kubernetes.NewForConfig(config)
		if err != nil {
			userLogger.Error("Failed to create Kubernetes client", zap.Error(err))
			AddErrorToContext(c, customerrors.NewInternalServerError("failed to create Kubernetes client"))
			c.Abort()
			return
		}

		dynClient, err := client.New(config, client.Options{Scheme: scheme})
		if err != nil {
			userLogger.Error("Failed to create Kubernetes dynamic client", zap.Error(err))
			AddErrorToContext(c, customerrors.NewInternalServerError("failed to create Kubernetes dynamic client"))
			c.Abort()
			return
		}

		opts := zapctrl.Options{Development: true}
		ctrl.SetLogger(zapctrl.New(zapctrl.UseFlagOptions(&opts)))

		// Update the logger with the username
		c.Set(LoggerCtxKey, userLogger)
		c.Set(KubeClientCtxKey, kubeClient)
		c.Set(DynamicClientCtxKey, dynClient)
		c.Set(TokenCtxKey, token)
		c.Set(ConfigKey, config)
		c.Next()
	}
}

// validateToken validates the format and presence of the Authorization token.
func validateToken(c *gin.Context) (string, error) {
	token := c.GetHeader(httpAuthorizationHeader)
	if token == "" {
		return validateTokenFromWS(c)
	}

	tokenParts := strings.Split(token, " ")

	// Check if the token is in the format "Bearer <token>"
	if len(tokenParts) != validBearerTokenPartsLength || tokenParts[httpBearerTokenPrefixIndex] != httpBearerTokenPrefix {
		return "", fmt.Errorf("invalid authentication token")
	}

	return tokenParts[httpBearerTokenIndex], nil
}

// validateTokenFromWS extracts the WebSocket authorization token from the request headers.
func validateTokenFromWS(c *gin.Context) (string, error) {
	token := c.GetHeader(WebsocketTokenHeader)
	if token == "" {
		return "", fmt.Errorf("authorization token not provided")
	}

	return token, nil
}

// createKubernetesConfig creates a new Kubernetes client config using the provided token.
func createKubernetesConfig(token, kubeApiServer string) (*rest.Config, error) {
	skipTlsVerify, err := utils.GetEnvBool(envInsecureSkipVerify, true)
	if err != nil {
		return nil, err
	}

	config := &rest.Config{
		BearerToken: token,
		Host:        kubeApiServer,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: skipTlsVerify,
		},
	}

	config.Wrap(multicluster.NewClusterGatewayRoundTripper)
	return config, nil
}
