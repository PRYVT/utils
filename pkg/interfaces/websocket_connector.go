package interfaces

type WebsocketConnecter interface {
	WriteJSON(v interface{}) error
	ReadForDisconnect()
	IsAuthenticated() bool
	IsConnected() bool
}
