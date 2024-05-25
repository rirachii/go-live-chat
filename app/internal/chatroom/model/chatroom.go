package chatroom_model

import (
	model "github.com/rirachii/golivechat/app/shared/model"
)

type Chatroom interface {
	Info() ChatroomInfo
	Id() model.RoomID
	Name() string
	IsPublic() bool
	Messages() []*ChatroomMessage
	SavedMessages() []*ChatroomMessage
	LastSavedMessage() int

	ActiveSubscribers() map[model.UserID]*Subscriber
	AddSubscriber(subscriber *Subscriber)
	RemoveSubscriber(uid model.UserID)

	EnqueueJoin(user *ChatroomUser)
	EnqueueLeave(subscriber *Subscriber)
	Broadcast(message *ChatroomMessage)

	Open()
	Close()
	// AcceptConnection(c echo.Context)
}

// Chatroom information to create the chatroom
type ChatroomInfo struct {
	RoomID    model.RoomID
	RoomName  string
	RoomOwner model.UserID
	IsPublic  bool
}

type CreateChatroomRequest struct {
	RoomID    model.RoomID
	RoomName  string
	RoomOwner model.UserID
	IsPublic  bool
}
