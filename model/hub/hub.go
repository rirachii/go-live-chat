package hub_model

import (
	model "github.com/rirachii/golivechat/model"
	chatroom_model "github.com/rirachii/golivechat/model/chat"

	echo "github.com/labstack/echo/v4"
)

type ChatroomsHub struct {
	PublicActiveChatrooms  map[model.RoomID]chatroom_model.Chatroom
	PrivateActiveChatrooms map[model.RoomID]chatroom_model.Chatroom
	registerQueue          chan chatroom_model.Chatroom
	unregisterQueue        chan chatroom_model.Chatroom
}

func CreateChatroomsHub(chatrooms []chatroom_model.Chatroom) *ChatroomsHub {
	hub := &ChatroomsHub{
		PublicActiveChatrooms:  make(map[model.RoomID]chatroom_model.Chatroom),
		PrivateActiveChatrooms: make(map[model.RoomID]chatroom_model.Chatroom),
		registerQueue:          make(chan chatroom_model.Chatroom),
		unregisterQueue:        make(chan chatroom_model.Chatroom),
	}

	// fill in with data from db
	hub.populate(chatrooms)

	return hub
}

// Handle adding/removing users to rooms
func (hub *ChatroomsHub) Run() {
	for {
		select {
		case room := <-hub.registerQueue:

			// register chatroom
			go hub.registerRoom(room)

		case room := <-hub.unregisterQueue:

			// unregister chatroom
			go hub.unregisterRoom(room)

		}
	}
}

func (hub *ChatroomsHub) RegisterRoom(room chatroom_model.Chatroom) {

	hub.registerQueue <- room

}

func (hub *ChatroomsHub) UnregisterRoom(room chatroom_model.Chatroom) {
	hub.unregisterQueue <- room

}

// returns nil when no chatroom is found.
func (hub ChatroomsHub) Chatroom(roomID model.RoomID) chatroom_model.Chatroom {

	pubChatroom, pubFound := hub.PublicActiveChatrooms[roomID]

	if pubFound {
		return pubChatroom
	}

	privChatroom, privFound := hub.PrivateActiveChatrooms[roomID]
	if privFound {
		return privChatroom
	}

	return nil

}

func (hub *ChatroomsHub) PublicChatrooms() map[model.RoomID]chatroom_model.Chatroom {
	return hub.PublicActiveChatrooms
}

func (hub *ChatroomsHub) PrivateChatrooms() map[model.RoomID]chatroom_model.Chatroom {
	return hub.PrivateActiveChatrooms
}

func (hub *ChatroomsHub) registerRoom(chatroom chatroom_model.Chatroom) {

	chatroomID := chatroom.ID()
	chatroomIsPublic := chatroom.IsPublic()

	if chatroomIsPublic {
		hub.PublicActiveChatrooms[chatroomID] = chatroom
	} else {
		hub.PrivateActiveChatrooms[chatroomID] = chatroom
	}

}

func (hub *ChatroomsHub) unregisterRoom(chatroom chatroom_model.Chatroom) {

	echo.New().Logger.Printf("Unregistering room %s: %s", chatroom.ID(), chatroom.Name())

	if chatroom.IsPublic() {
		delete(hub.PublicActiveChatrooms, chatroom.ID())
	} else {
		delete(hub.PrivateActiveChatrooms, chatroom.ID())
	}

	// go cchatroom.Close()

}

func (hub *ChatroomsHub) populate(chatrooms []chatroom_model.Chatroom) {
	
	
	for _, room := range chatrooms {
		hub.addChatroom(room)
	}


}

func (hub *ChatroomsHub) addChatroom(chatroom chatroom_model.Chatroom) {

	chatroomID := chatroom.ID()

	if chatroom.IsPublic() {
		hub.PublicActiveChatrooms[chatroomID] = chatroom
	} else {
		hub.PrivateActiveChatrooms[chatroomID] = chatroom
	}

}
