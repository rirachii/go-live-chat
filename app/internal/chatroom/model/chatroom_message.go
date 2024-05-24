package chatroom_model

import (
	model "github.com/rirachii/golivechat/app/shared/model"
)

type ChatroomMessage struct {
	id      int
	message model.Message
	isSaved bool
}

func NewChatroomMessage(message model.Message) *ChatroomMessage {
	chatroomMessage := ChatroomMessage{
		message: message,
		isSaved: false,
	}
	return &chatroomMessage
}

func (chatroom_message ChatroomMessage) Id() int                { return chatroom_message.id }
func (chatroom_message ChatroomMessage) Message() model.Message { return chatroom_message.message }
func (chatroom_message ChatroomMessage) IsSaved() bool          { return chatroom_message.isSaved }

// marks message as saved.
func (chatroom_message *ChatroomMessage) Saved() { chatroom_message.isSaved = true }
