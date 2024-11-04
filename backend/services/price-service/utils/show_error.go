package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func ShowError(statusCode int64, message string, context *gin.Context) {
	context.JSON(int(statusCode), gin.H{"message": message})
}

func ShowErrorSocket(ws *websocket.Conn, message string) {
	var errorMsg = fmt.Sprintf(`
	"message": %s
	`, message)
	ws.WriteMessage(websocket.TextMessage, []byte(errorMsg))
}
