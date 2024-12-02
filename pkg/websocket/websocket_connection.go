package websocket

import (
	"fmt"

	"github.com/PRYVT/utils/pkg/auth"
	"github.com/PRYVT/utils/pkg/interfaces"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type WebsocketConnectionInterfacer interface {
	WriteJSON(v interface{}) error
	ReadForDisconnect()
}

type WebsocketConnection struct {
	connection      *websocket.Conn
	isConnected     bool
	isAuthenticated bool
	userUuid        uuid.UUID
}

func NewWebsocketConnection(conn *websocket.Conn) interfaces.WebsocketConnecter {
	wC := &WebsocketConnection{connection: conn}
	go wC.ReadForDisconnect()
	return wC
}

func (wC *WebsocketConnection) IsConnected() bool {
	return wC.isConnected
}

func (wC *WebsocketConnection) IsAuthenticated() bool {
	return wC.isAuthenticated
}

func (wC *WebsocketConnection) WriteJSON(v interface{}) error {
	if !wC.isAuthenticated {
		return fmt.Errorf("WebsocketConnection is not connected or authenticated")
	}
	err := wC.connection.WriteJSON(v)
	if err != nil {
		log.Warn().Err(err).Msg("Error while writing WriteJSON")
	}
	return err
}

func (wC *WebsocketConnection) ReadForDisconnect() {
	wC.isConnected = true
	for {
		authRequest := AuthRequest{}
		err := wC.connection.ReadJSON(&authRequest)
		if err != nil {
			log.Debug().Err(err).Msg("Error while reading from websocket connection")
			wC.isAuthenticated = false
			wC.connection.Close()
			wC.isConnected = false
			break
		} else {
			log.Debug().Interface("authReq", authRequest).Msg("Received auth request")
			_, err = auth.VerifyToken(authRequest.Token)
			if err != nil {
				log.Debug().Err(err).Msg("Error while verifying token")
				wC.isAuthenticated = false
				wC.connection.Close()
				wC.isConnected = false
				break
			}
			uuid, err := auth.GetUserUuidFromToken(authRequest.Token)
			if err != nil {
				log.Debug().Err(err).Msg("Error while getting user uuid from token")
				wC.isAuthenticated = false
				wC.connection.Close()
				wC.isConnected = false
				break
			}
			wC.userUuid = uuid
			wC.isAuthenticated = true
		}
	}
}
