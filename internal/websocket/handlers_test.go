package websocket

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func Test_WebsocketHandlers(t *testing.T) {
	type args struct {
		c          *gin.Context
		conn       *websocket.Conn
		stream     io.ReadCloser
		formatFunc func(string) string
	}
	type want struct {
		response []string
	}
	cases := map[string]struct {
		args args
		want want
	}{
		"ShouldSuccess": {
			args: args{
				formatFunc: func(s string) string {
					return "formatted: " + s
				},
			},
			want: want{
				response: []string{
					"formatted: line1",
					"formatted: line2",
					"formatted: line3",
				},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request, _ = http.NewRequest(http.MethodGet, "/ws", nil)
			tc.args.c = c

			// Setup the WebSocket server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				upgrader := websocket.Upgrader{}
				conn, err := upgrader.Upgrade(w, r, nil)
				if err != nil {
					t.Fatalf("Failed to upgrade connection: %v", err)
				}
				tc.args.conn = conn

				// Simulate the input stream
				pr, pw := io.Pipe()
				tc.args.stream = pr
				go func() {
					defer pw.Close()
					for _, line := range []string{"line1", "line2", "line3"} {
						_, err = pw.Write([]byte(line + "\n"))
						time.Sleep(100 * time.Millisecond) // simulate delay between lines
					}
				}()

				Stream(tc.args.c, tc.args.conn, tc.args.stream, tc.args.formatFunc)
			}))
			defer server.Close()

			// Connect to the WebSocket server
			url := "ws" + strings.TrimPrefix(server.URL, "http")
			ws, _, err := websocket.DefaultDialer.Dial(url, nil)
			if err != nil {
				t.Fatalf("Failed to dial WebSocket server: %v", err)
			}
			defer ws.Close()

			// Read messages from the WebSocket
			var messages []string
			for i := 0; i < len(tc.want.response); i++ {
				_, msg, err := ws.ReadMessage()
				if err != nil {
					t.Fatalf("Failed to read message: %v", err)
				}
				messages = append(messages, string(msg))
			}

			// Verify the messages
			if len(messages) != len(tc.want.response) {
				t.Fatalf("Expected %d messages, got %d", len(tc.want.response), len(messages))
			}

			for i, msg := range messages {
				if msg != tc.want.response[i] {
					t.Errorf("Expected message %d to be %q, got %q", i, tc.want.response[i], msg)
				}
			}
		})
	}
}
