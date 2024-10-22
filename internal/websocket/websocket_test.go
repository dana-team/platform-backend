package websocket

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_WebSocket_Register(t *testing.T) {
	tests := []struct {
		name       string
		token      string
		setupToken bool
		wantErr    bool
	}{
		{
			name:       "TokenExists",
			token:      "valid-token",
			setupToken: true,
			wantErr:    false,
		},
		{
			name:       "TokenNotExists",
			token:      "",
			setupToken: false,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/ws", nil)
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			if tt.setupToken {
				c.Set("token", tt.token)
			}

			// Initialize WebSocket handler
			wsHandler := NewWebSocket(nil)

			// Set up a WebSocket server to test the upgrade
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c, _ := gin.CreateTestContext(w)
				c.Request = r

				if tt.setupToken {
					c.Set("token", tt.token)
				}

				_, err := wsHandler.Register(c)
				if (err != nil) != tt.wantErr {
					t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				}
			}))
			defer server.Close()

			// Connect to the WebSocket server
			url := "ws" + server.URL[4:]
			ws, _, err := websocket.DefaultDialer.Dial(url, nil)

			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if ws == nil {
					t.Errorf("Expected a WebSocket connection, got nil")
				} else {
					ws.Close()
				}
			}
		})
	}
}
