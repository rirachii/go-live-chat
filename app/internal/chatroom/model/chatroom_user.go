package chatroom_model

import (
	model "github.com/rirachii/golivechat/app/shared/model"
	websocket "nhooyr.io/websocket"
)

// ChatroomUser.
// Users become subscribers when they join a chatroom
type ChatroomUser struct {
	wsConn *websocket.Conn
	userId model.UserID
	roomId model.RoomID
}

func (user *ChatroomUser) Websocket() *websocket.Conn { return user.wsConn }
func (user ChatroomUser) ID() model.UserID            { return user.userId }

// func (user ChatroomUser) Name() string                { return user.userInfo.Name }
func (user ChatroomUser) RoomID() model.RoomID { return user.roomId }

func NewChatroomUser(ws *websocket.Conn, uid model.UserID, rid model.RoomID) *ChatroomUser {

	chatroomUser := ChatroomUser{
		wsConn: ws,
		userId: uid,
		roomId: rid,
	}

	return &chatroomUser

}
