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

func (m MessageRequest) SenderId() model.UserID { return m.UserID }
func (m MessageRequest) Content() string        { return m.ChatMessage }
func (m MessageRequest) Room() string           { return m.RoomID }
