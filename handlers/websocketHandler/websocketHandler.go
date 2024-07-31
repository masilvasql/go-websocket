package websocketHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/masilvasql/go-websocket/pkg/samwebsocket"
	"net/http"
)

func Handle(c *gin.Context) {
	sessionID := c.Param("sessionId")

	samupgrader := samwebsocket.NewUpgrader(1024, 1024, true)

	conn, err := samupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	samwebsocket.AddClient(sessionID, conn)

	go samwebsocket.HandleMessages(sessionID, conn)
}
