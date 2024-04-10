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
