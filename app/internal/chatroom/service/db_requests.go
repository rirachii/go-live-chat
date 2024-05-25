package chatroom_service

import model "github.com/rirachii/golivechat/app/shared/model"

type SaveMessageRequest struct {
	UserID         model.UserID
	RoomID         model.RoomID
	MessageContent string
}

type SaveChatroomMessagesRequest struct {
	RoomID   model.RoomID
	ChatLogs []SaveMessageRequest
}

type GetChatroomLogsRequest struct {
	RoomID model.RoomID
}

type GetUsernameByIdRequest struct {
	UserID model.UserID
}
