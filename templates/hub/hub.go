package templates

import (
	model "github.com/rirachii/golivechat/model"
	t "github.com/rirachii/golivechat/templates"
)


// "hub" template

var HubPage = t.TemplateData{
	TemplateName: "hub",
}

// "hub-chatrooms" template

var HubChatrooms = t.TemplateData{
	TemplateName: "hub-chatrooms",
}

type TemplateHubChatrooms struct {
	Rooms []ChatroomTemplateData
}

type ChatroomTemplateData struct {
	Chatroom hubChatroomData
}

type hubChatroomData struct {
	RoomID   model.RoomID
	RoomName string
}

func PrepareChatroomData(roomID model.RoomID, roomName string) ChatroomTemplateData {

	chatroomData := hubChatroomData{
		RoomID:   roomID,
		RoomName: roomName,
	}

	renderChatroom := ChatroomTemplateData{
		Chatroom: chatroomData,
	}

	return renderChatroom
}
