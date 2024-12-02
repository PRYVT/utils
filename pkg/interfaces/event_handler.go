package interfaces

import "github.com/L4B0MB4/EVTSRC/pkg/models"

type EventHandler interface {
	HandleEvent(event models.Event) error
	AddWebsocketConnection(conn WebsocketConnecter)
}
