package hub_model


import (
	model "github.com/rirachii/golivechat/app/shared/model"
	chatroom_model "github.com/rirachii/golivechat/app/internal/chatroom/model"
)

type HubServer interface {

	Start()
	Close()

	RegisterRoom(chatroom_model.Chatroom)
	UnregisterRoom(chatroom_model.Chatroom)

	PublicChatrooms() map[model.RoomID]chatroom_model.Chatroom
	Chatroom(model.RoomID) chatroom_model.Chatroom

}



