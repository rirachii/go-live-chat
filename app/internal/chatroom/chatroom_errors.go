package chatroom

import "errors"

var (
	ErrWebsocketConnectionFailed = errors.New("error connecting to websocket")
)