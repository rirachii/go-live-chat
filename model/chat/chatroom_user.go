package chat_model

import (
	model "github.com/rirachii/golivechat/model"
	user_model "github.com/rirachii/golivechat/model/user"

	websocket "nhooyr.io/websocket"
)

type ChatroomUser struct {
	client       *ChatroomClient
	userInfo     user_model.UserInfo
	rid          model.RoomID
	role         string
	messageQueue chan model.Message
}

type ChatroomClient struct {
	Conn        *websocket.Conn
	UserRequest model.UserRequest
}

func NewChatroomUser(
	client *ChatroomClient,
	user user_model.UserInfo,
	roomID model.RoomID,
	userRole string,
) *ChatroomUser {

	chatroomUser := &ChatroomUser{
		client:       client,
		userInfo:     user,
		rid:          roomID,
		role:         userRole,
		messageQueue: make(chan model.Message),
	}

	return chatroomUser
}

// ChatroomUser
func (user ChatroomUser) UserInfo() user_model.UserInfo { return user.userInfo }
func (user ChatroomUser) ID() model.UserID      { return user.UserInfo().ID }
func (user ChatroomUser) Username() string          { return user.UserInfo().Username }

func (user ChatroomUser) RoomID() model.RoomID                  { return user.rid }
func (user ChatroomUser) Role() string                          { return user.role }
func (user *ChatroomUser) PendingMessages() chan model.Message { return user.messageQueue }
func (user *ChatroomUser) Client() *ChatroomClient              { return user.client }

// ChatroomClient
func (client *ChatroomClient) Websocket() *websocket.Conn { return client.Conn }
func (client ChatroomClient) UserID() model.UserID        { return client.UserRequest.UserID }
func (client ChatroomClient) Username() string       { return client.UserRequest.Username }
func (client ChatroomClient) RoomID() model.RoomID        { return client.UserRequest.RoomID }
