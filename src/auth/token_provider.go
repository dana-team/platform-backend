package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// TokenProvider defines an interface for obtaining a token.
type TokenProvider interface {
	ObtainToken(username, password string, logger *zap.Logger, ctx *gin.Context) (string, error)
	ObtainUsername(token string, logger *zap.Logger) (string, error)
}

// DefaultTokenProvider is a default implementation of TokenProvider.
type DefaultTokenProvider struct{}

func (d DefaultTokenProvider) ObtainToken(username, password string, logger *zap.Logger, ctx *gin.Context) (string, error) {
	return ObtainOpenshiftToken(username, password, logger, ctx)
}

func (d DefaultTokenProvider) ObtainUsername(token string, logger *zap.Logger) (string, error) {
	return ObtainOpenshiftUsername(token, logger)
}
