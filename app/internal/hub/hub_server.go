package hub

import (
	chatroom_model "github.com/rirachii/golivechat/app/internal/chatroom/model"
	hub_model "github.com/rirachii/golivechat/app/internal/hub/model"
	model "github.com/rirachii/golivechat/app/shared/model"
)

type hubServer struct {
	publicChatrooms  map[model.RoomID]chatroom_model.Chatroom
	privateChatrooms map[model.RoomID]chatroom_model.Chatroom

	registerQueue    chan chatroom_model.Chatroom
	unregisterQueue  chan chatroom_model.Chatroom

	closeChannel chan int
}

func NewHubServer(chatrooms []chatroom_model.Chatroom) hub_model.HubServer {

	hub := &hubServer{
		publicChatrooms:  make(map[model.RoomID]chatroom_model.Chatroom),
		privateChatrooms: make(map[model.RoomID]chatroom_model.Chatroom),
		registerQueue:    make(chan chatroom_model.Chatroom),
		unregisterQueue:  make(chan chatroom_model.Chatroom),
	}

	return hub
}

func (hub *hubServer) Start() {
	for {
		select {
		case room := <-hub.registerQueue:

			// register chatroom
			go hub.register(room)

		case room := <-hub.unregisterQueue:

			// unregister chatroom
			go hub.unregister(room)

		}
	}
}

func (hub *hubServer) Close() {
	hub.closeChannel <- 1
}

func (hub *hubServer) RegisterRoom(room chatroom_model.Chatroom) {

	hub.registerQueue <- room

}

func (hub *hubServer) UnregisterRoom(room chatroom_model.Chatroom) {
	hub.unregisterQueue <- room

}

// returns nil when no chatroom is found.
func (hub hubServer) Chatroom(roomID model.RoomID) chatroom_model.Chatroom {

	pubChatroom, pubFound := hub.publicChatrooms[roomID]

	if pubFound {
		return pubChatroom
	}

	privChatroom, privFound := hub.privateChatrooms[roomID]
	if privFound {
		return privChatroom
	}

	return nil

}


func (hub *hubServer) PublicChatrooms() map[model.RoomID]chatroom_model.Chatroom {
	return hub.publicChatrooms
}

func (hub *hubServer) PrivateChatrooms() map[model.RoomID]chatroom_model.Chatroom {
	return hub.privateChatrooms
}


func (hub *hubServer) register(chatroom chatroom_model.Chatroom) {

	chatroomID := chatroom.Id()
	chatroomIsPublic := chatroom.IsPublic()

	if chatroomIsPublic {
		hub.publicChatrooms[chatroomID] = chatroom
	} else {
		hub.privateChatrooms[chatroomID] = chatroom
	}

}

func (hub *hubServer) unregister(chatroom chatroom_model.Chatroom) {

	// echo.New().Logger.Printf("Unregistering room %s: %s", chatroom.ID(), chatroom.Name())

	if chatroom.IsPublic() {
		delete(hub.publicChatrooms, chatroom.Id())
	} else {
		delete(hub.privateChatrooms, chatroom.Id())
	}

	// go cchatroom.Close()

}


func (hub *hubServer) populate(chatrooms []chatroom_model.Chatroom) {
	
	
	for _, room := range chatrooms {
		hub.addChatroom(room)
	}

}

func (hub *hubServer) addChatroom(chatroom chatroom_model.Chatroom) {

	chatroomID := chatroom.Id()

	if chatroom.IsPublic() {
		hub.publicChatrooms[chatroomID] = chatroom
	} else {
		hub.privateChatrooms[chatroomID] = chatroom
	}

}