package chatroom_model

import (
	model "github.com/rirachii/golivechat/app/shared/model"
	websocket "nhooyr.io/websocket"
)

// ChatroomUser.
// Users become subscribers when they join a chatroom
type ChatroomUser struct {
	wsConn *websocket.Conn
	id     model.UserID
	roomId model.RoomID
	name   string
}

func (user *ChatroomUser) Websocket() *websocket.Conn { return user.wsConn }
func (user ChatroomUser) Id() model.UserID            { return user.id }
func (user ChatroomUser) Name() string                { return user.name }
func (user ChatroomUser) Info() model.UserInfo        { return model.CreateUserInfo(user.id, user.name) }

// func (user ChatroomUser) Name() string                { return user.userInfo.Name }
func (user ChatroomUser) RoomID() model.RoomID { return user.roomId }

func NewChatroomUser(ws *websocket.Conn, uid model.UserID, rid model.RoomID, username string) *ChatroomUser {

	chatroomUser := ChatroomUser{
		wsConn: ws,
		id: uid,
		roomId: rid,
		name: username,
	}

	return &chatroomUser

}
