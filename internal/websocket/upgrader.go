package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"os"
	"regexp"
)

// DefaultUpgrader creates a new WebSocket upgrader with default buffer sizes and an origin check based on regex from an environment variable.
func DefaultUpgrader() *websocket.Upgrader {
	allowedOriginRegex := os.Getenv("ALLOWED_ORIGIN_REGEX")
	if allowedOriginRegex == "" {
		allowedOriginRegex = ".*" // Allow all if the environment variable is not set
	}
	originRegex, err := regexp.Compile(allowedOriginRegex)
	if err != nil {
		panic("Invalid origin regex")
	}

	return &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			check := originRegex.MatchString(origin)

			return check
		},
	}
}
