package chat_model

import (
	model "github.com/rirachii/golivechat/model"
	websocket "nhooyr.io/websocket"
)

type ChatroomUser struct {
	client       *ChatroomClient
	uid          model.UserID
	rid          model.RoomID
	role         string
	messageQueue chan *model.Message
}

type ChatroomClient struct {
	Conn        *websocket.Conn
	UserRequest model.UserRequest
}

func NewChatroomUser(
	client *ChatroomClient,
	userID model.UserID,
	roomID model.RoomID,
	userRole string,
) *ChatroomUser {

	chatroomUser := &ChatroomUser{
		client:       client,
		uid:          userID,
		rid:          roomID,
		role:         userRole,
		messageQueue: make(chan *model.Message),
	}

	return chatroomUser
}

func (user ChatroomUser) ID() model.UserID {
	return user.uid
}

func (user ChatroomUser) RoomID() model.RoomID {
	return user.rid
}

func (user ChatroomUser) Role() string {
	return user.role
}

func (user *ChatroomUser) PendingMessages() chan *model.Message {
	return user.messageQueue
}

func (user *ChatroomUser) Client() *ChatroomClient {
	return user.client

}

func (client *ChatroomClient) Websocket() *websocket.Conn {
	return client.Conn
}

func (client ChatroomClient) UserID() model.UserID {
	return client.UserRequest.UserID
}

func (client ChatroomClient) RoomID() model.RoomID {
	return client.UserRequest.RoomID
}
