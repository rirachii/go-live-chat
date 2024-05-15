package chatroom_model

import (
	echo "github.com/labstack/echo/v4"
	model "github.com/rirachii/golivechat/app/internal/shared/model"
)

type Chatroom interface {
	Info() ChatroomInfo
	ID() model.RoomID
	Name() string
	IsPublic() bool
	LastSavedLogIndex() int

	ActiveUsers() map[model.UserID]*ChatroomUser
	AddUser(user *ChatroomUser)
	RemoveUser(uid model.UserID)

	EnqueueJoin(client *ChatroomClient)
	EnqueueLeave(user *ChatroomUser)
	Broadcast(message model.Message)

	Open()
	Close()
	AcceptConnection(c echo.Context)
}

type ChatroomInfo struct {
	RoomID    model.RoomID
	RoomName  string
	RoomOwner model.UserID
	IsPublic  bool
}

type Message struct {
	RoomID   model.RoomID
	SenderID model.UserID // user's id
	Sender   string       // user name
	Content  string       // content of message
}
