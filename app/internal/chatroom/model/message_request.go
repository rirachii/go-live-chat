package chatroom_model

import (
	"github.com/rirachii/golivechat/app/shared/model"
)

type MessageRequest struct {
    UserID      model.UserID
    ChatMessage string `json:"chat-message"`
    RoomID      string `json:"room-id"`
}
func NewMessageRequest(uid model.UserID) MessageRequest {
	msgRequest := MessageRequest{
		UserID: uid,
	}
	return msgRequest
}
