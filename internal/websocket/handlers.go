package websocket

import (
	"bufio"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Stream reads lines from the provided stream, formats each line using the given formatFunc,
// and sends them as WebSocket messages to the client connection until the context is done.
func Stream(c *gin.Context, conn *websocket.Conn, stream io.ReadCloser, formatFunc func(string) string) {
	defer stream.Close()
	defer conn.Close()

	reader := bufio.NewScanner(stream)
	var line string
	for {
		select {
		case <-c.Done():
			return
		default:
			for reader.Scan() {
				line = reader.Text()
				message := formatFunc(line)
				if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
					SendErrorMessage(conn, "Error writing message to WebSocket")
					return
				}
			}
		}
	}
}

// SendErrorMessage sends an error message to the WebSocket client.
func SendErrorMessage(conn *websocket.Conn, errorMsg string) {
	message := "error: " + errorMsg
	_ = conn.WriteMessage(websocket.TextMessage, []byte(message))
}
