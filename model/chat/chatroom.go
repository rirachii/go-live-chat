package chat_model

import (
	echo "github.com/labstack/echo/v4"
	model "github.com/rirachii/golivechat/model"
)

type Chatroom interface {
	Info() model.ChatroomInfo
	ID() model.RoomID
	Name() string
	IsPublic() bool
	ChatLogs() []model.Message
	IndexLastSavedLog() int

	ActiveUsers() map[model.UserID]*ChatroomUser
	AddUser(user *ChatroomUser)
	RemoveUser(uid model.UserID)
	EnqueueJoin(client *ChatroomClient)
	EnqueueLeave(client *ChatroomUser)
	EnqueueMessageBroadcast(message model.Message)

	Open()
	Close()
	AcceptConnection(c echo.Context, userReq model.UserRequest) error

	ListenToUserWS(user *ChatroomUser)
	SendMessageToUser(user *ChatroomUser, msg model.Message)
	LogMessage(msg model.Message)
}
