package templates

import (
	t "github.com/rirachii/golivechat/templates"
)

// "chatroom" template
var ChatroomPage = t.TemplateData{
	TemplateName: "chatroom",
}

type TemplateChatroomPage struct {
	RoomID   string
	RoomName string
}

// "chatroom-connection" template
var ChatroomConnection = t.TemplateData{
	TemplateName: "chatroom-connection",
}

type TemplateChatroomConnection struct {
	RoomID string
}

const WebsocketDivID = "chat-messages"

// "many-messages" template
var ManyMessages = t.TemplateData{
	TemplateName: "many-messages",
}

type TemplateManyMessages struct {
	ChatMessages []TemplateSingleMessage
}

// "single-message" template
var SingleMessage = t.TemplateData{
	TemplateName: "single-message",
}

type TemplateSingleMessage struct {
	DivID       string
	PrependMsg  bool
	DisplayName string
	TextMessage string
}

// helpers
func PrepareMessage(
	divID string,
	prependMsg bool,
	displayName string,
	textMessage string,
) TemplateSingleMessage {

	return TemplateSingleMessage{
		DivID:       divID,
		PrependMsg:  prependMsg,
		DisplayName: displayName,
		TextMessage: textMessage,
	}

}
