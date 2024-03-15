package hub_model

import (
	model "github.com/rirachii/golivechat/model"
	chat "github.com/rirachii/golivechat/model/chat"

	echo "github.com/labstack/echo/v4"
)

type ChatroomsHub struct {
	chatrooms       map[model.RoomID]*chat.Chatroom
	userChatrooms   map[model.UserID]UserSetOfChatrooms // userid -> their chatrooms (room ids)
	registerQueue   chan *model.UserRequest
	unregisterQueue chan *model.UserRequest
}

// goes into db
type UserSetOfChatrooms struct {
	ChatroomsSet map[model.RoomID]bool
}

type SetOfChatrooms struct {
	Chatrooms map[model.RoomID]bool
}

func NewChatroomsHub() *ChatroomsHub {
	hub := &ChatroomsHub{
		chatrooms:       make(map[model.RoomID]*chat.Chatroom),
		userChatrooms:   make(map[model.UserID]UserSetOfChatrooms),
		registerQueue:   make(chan *model.UserRequest),
		unregisterQueue: make(chan *model.UserRequest),
	}

	return hub
}

// Handle adding/removing users to rooms
func (hub *ChatroomsHub) Run() {
	for {
		select {
		case userReq := <-hub.registerQueue:

			// register them to the approrpiate chat
			go hub.registerUser(userReq)

		case userReq := <-hub.unregisterQueue:

			// unregister them
			go hub.unregisterUser(userReq)

		}
	}
}

func (hub *ChatroomsHub) Register(userRequest *model.UserRequest) {

	hub.registerQueue <- userRequest

}

func (hub *ChatroomsHub) Unregister(userRequest *model.UserRequest) {
	hub.unregisterQueue <- userRequest

}

func (hub *ChatroomsHub) AddandOpenRoom(newChatRoom *chat.Chatroom) error {

	roomID := newChatRoom.ID()
	hub.chatrooms[roomID] = newChatRoom
	go newChatRoom.Open()

	// TODO check errors
	return nil
}

func (hub ChatroomsHub) Chatroom(roomID model.RoomID) *chat.Chatroom {
	chatroom, ok := hub.chatrooms[roomID]
	if !ok {
		return nil
	}
	return chatroom
}

func (hub ChatroomsHub) Chatrooms() map[model.RoomID]*chat.Chatroom {
	return hub.chatrooms
}

func (rooms *UserSetOfChatrooms) RegisterRoom(roomID model.RoomID) {
	// TODO be aware does not check if the room is already registered to them
	rooms.ChatroomsSet[roomID] = true
}
func (rooms *UserSetOfChatrooms) UnregisterRoom(roomID model.RoomID) {
	// TODO does not check for errors
	delete(rooms.ChatroomsSet, roomID)
}

func (hub *ChatroomsHub) registerUser(userReq *model.UserRequest) {

	var (
		userID = userReq.UserID
		roomID = userReq.RoomID
	)

	// TODO use DB
	userChatrooms, ok := hub.userChatrooms[userID]
	if !ok {
		hub.userChatrooms[userID] = UserSetOfChatrooms{
			ChatroomsSet: make(map[model.RoomID]bool),
		}
		userChatrooms = hub.userChatrooms[userID]
		// log.Printf("user [%s] not found", clientID)
	}

	userChatrooms.RegisterRoom(roomID)

	echo.New().Logger.Debugf("Registering %s to room %s", userID, roomID)
	echo.New().Logger.Debugf("User [%s] rooms: %i", userID, userChatrooms)
}

func (hub *ChatroomsHub) unregisterUser(userReq *model.UserRequest) {

	var (
		userID = userReq.UserID
		roomID = userReq.RoomID
	)

	// userChatrooms, ok := hub.UserChatrooms[clientID]
	// if !ok {
	// 	log.Printf("user [%s] not found", clientID)
	// }

	// userChatrooms.UnregisterRoom(roomID)

	// chatRoom := hub.ChatRooms[roomID]

	// chatRoom.LeaveQueue <- client

	echo.New().Logger.Debugf("Unregistering %s from room %s", userID, roomID)

}
