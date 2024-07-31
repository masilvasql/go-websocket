package main

import (
	"github.com/gin-gonic/gin"
	"github.com/masilvasql/go-websocket/handlers/websocketHandler"
)

func main() {
	router := gin.Default()
	router.GET("/ws/:sessionId", websocketHandler.Handle)

	router.GET("/", func(c *gin.Context) {
		c.File("./public/index.html")
	})

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
