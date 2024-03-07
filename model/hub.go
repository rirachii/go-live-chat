package model

import (
	echo "github.com/labstack/echo/v4"
)

type ChatroomsHub struct {
	ChatRooms       map[RoomID]*Chatroom
	UserChatrooms   map[UserID]UserSetOfChatrooms // userid -> their chatrooms (room ids)
	RegisterQueue   chan *UserRequest
	UnregisterQueue chan *UserRequest
}

type UserRequest struct {
	UserID UserID
	RoomID RoomID
}
type UserID string
type RoomID string

type ChatroomData struct {
	RoomID   RoomID
	RoomName string
}

type UserSetOfChatrooms struct {
	ChatroomsSet map[RoomID]bool
}


type SetOfChatrooms struct {
	Chatrooms map[RoomID]bool
}

// Handle adding/removing users to rooms
func (hub *ChatroomsHub) Run() {
	for {
		select {
		case userReq := <-hub.RegisterQueue:

			// register them to the approrpiate chat
			go hub.registerUser(userReq)

		case userReq := <-hub.UnregisterQueue:

			// unregister them
			go hub.unRegisterUser(userReq)

		}
	}
}

func (hub *ChatroomsHub) AddandOpenRoom(newChatRoom *Chatroom) error {

	roomID := newChatRoom.GetRID()
	hub.ChatRooms[roomID] = newChatRoom
	go newChatRoom.Open()

	// TODO check errors
	return nil
}

func (hub ChatroomsHub) GetChatroom(roomID RoomID) *Chatroom {
	getChatroom, ok := hub.ChatRooms[roomID]
	if !ok {
		return nil
	}
	return getChatroom
}

func (rooms *UserSetOfChatrooms) RegisterRoom(roomID RoomID) {
	// TODO be aware does not check if the room is already registered to them
	rooms.ChatroomsSet[roomID] = true
}
func (rooms *UserSetOfChatrooms) UnregisterRoom(roomID RoomID) {
	// TODO does not check for errors
	delete(rooms.ChatroomsSet, roomID)
}

func (hub *ChatroomsHub) registerUser(userReq *UserRequest) {

	var (
		userID = userReq.UserID
		roomID = userReq.RoomID
	)

	// TODO use DB
	userChatrooms, ok := hub.UserChatrooms[userID]
	if !ok {
		hub.UserChatrooms[userID] = UserSetOfChatrooms{
			ChatroomsSet: make(map[RoomID]bool),
		}
		userChatrooms = hub.UserChatrooms[userID]
		// log.Printf("user [%s] not found", clientID)
	}

	userChatrooms.RegisterRoom(roomID)

	echo.New().Logger.Debugf("Registering %s to room %s", userID, roomID)
	echo.New().Logger.Debugf("User [%s] rooms: %i", userID, userChatrooms)
}


func (hub *ChatroomsHub) unRegisterUser(userReq *UserRequest) {

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