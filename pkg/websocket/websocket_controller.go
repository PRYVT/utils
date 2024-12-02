package websocket

import (
	"net/http"

	"github.com/PRYVT/utils/pkg/interfaces"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type WSController struct {
	eventHandler interfaces.EventHandler
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewWsController(eventHandler interfaces.EventHandler) *WSController {

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	return &WSController{eventHandler: eventHandler}
}

func (w *WSController) OnRequest(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Warn().Err(err).Msg("Error while upgrading connection")

	} else {
		w.eventHandler.AddWebsocketConnection(NewWebsocketConnection(conn))
	}
}
