package handler

import (
	"log"
	"net/http"

	"github.com/zx06/saut/internal/pkg/assets_connector"
	"github.com/zx06/saut/internal/pkg/client"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func TerminalHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"message": err,
		})
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Printf("close websocket error: %s\n", err)
		}
	}(conn)
	//todo: parse params and create connector
	ac := assets_connector.SSHConnector{}
	client.WSHandler(conn, &ac)
}
