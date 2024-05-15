package chatroom_model

import (
	model "github.com/rirachii/golivechat/app/internal/shared/model"
	websocket "nhooyr.io/websocket"
)

type Subscriber struct {
	Conn        *websocket.Conn
	Request model.UserRequest
}

// ChatroomClient
func (subscriber *Subscriber) Websocket() *websocket.Conn { return subscriber.Conn }
func (subscriber Subscriber) ID() model.UserID        { return subscriber.Request.UserID }
func (subscriber Subscriber) Name() string            { return subscriber.Request.Username }
func (subscriber Subscriber) RoomID() model.RoomID        { return subscriber.Request.RoomID }
