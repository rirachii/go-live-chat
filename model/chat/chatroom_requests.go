package chat_model

import "github.com/rirachii/golivechat/model"

type SendMessageRequest struct {
	UserID      model.UserID
	MessageText string `json:"chat-message"`
	RoomID      string `json:"room-id"`
}

// for web socket connection
type ConnectionRequest struct {
	UserID model.UserID
	RoomID string `param:"roomID"`
}
