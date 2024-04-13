package chat_model

import "github.com/rirachii/golivechat/model"

type SaveUserMessageRequest struct {
	UserID      model.UserID
	RoomID      model.RoomID
	UserMessage string
}

type SaveChatLogsRequest struct {
	RoomID   model.RoomID
	ChatLogs []SaveUserMessageRequest
}

type GetMessageRequest struct {
	UserID model.UserID
	RoomID model.RoomID
}

type GetChatLogsRequest struct {
	RoomID model.RoomID
}


type ChatMessageDTO struct {
	RoomID      model.RoomID
	SenderID    model.UserID
	MessageText string
}

type ChatroomAdminDTO struct {
	AdminID model.UserID
	Role    string
}

type UserDisplayNameDTO struct {
	UserID      model.UserID
	DisplayName string
}
