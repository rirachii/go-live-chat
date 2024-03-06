package model

import (
	echo "github.com/labstack/echo/v4"
)

type UserID string
type RoomID string

type UserRoom struct {
	UserID UserID
	RoomID RoomID
}

type ChatroomsHub struct {
	ChatRooms       map[RoomID]*Chatroom  		// want to store chatroom data to db and query them by user location
	UserChatrooms   map[UserID]SetOfChatrooms 	// userid -> their chatrooms (room ids) ; also wanna store this info
	RegisterQueue   chan *UserRoom
	UnregisterQueue chan *UserRoom
}

type SetOfChatrooms struct {
	Chatrooms map[RoomID]bool
}


//Handle adding/removing users to rooms 
func (hub *ChatroomsHub) Run() {
	for {
		select {
		//when new user want to register to a room 
		case client := <-hub.RegisterQueue:
			// add them to the approrpiate chat
			clientID, roomID := client.UserID, client.RoomID

			userChatrooms, ok := hub.UserChatrooms[clientID]
			if !ok {
				hub.UserChatrooms[clientID] = SetOfChatrooms{
					Chatrooms: make(map[RoomID]bool),
				}
				userChatrooms = hub.UserChatrooms[clientID]
				// log.Printf("user [%s] not found", clientID)
			}

			userChatrooms.RegisterRoom(roomID)

			// chatRoom := hub.ChatRooms[roomID]
			// // chatRoom.JoinQueue <- client

			echo.New().Logger.Printf("Registering %s to room %s", clientID, roomID)
			echo.New().Logger.Printf("User [%s] rooms: %i", clientID, userChatrooms)

		case client := <-hub.UnregisterQueue:
			clientID, roomID := client.UserID, client.RoomID

			// userChatrooms, ok := hub.UserChatrooms[clientID]
			// if !ok {
			// 	log.Printf("user [%s] not found", clientID)
			// }

			// userChatrooms.UnregisterRoom(roomID)

			// chatRoom := hub.ChatRooms[roomID]

			// chatRoom.LeaveQueue <- client

			echo.New().Logger.Printf("Unregistering %s from room %s", clientID, roomID)

		}
	}
}

func (hub *ChatroomsHub) AddandOpenRoom(newChatRoom *Chatroom) error {
	roomID := newChatRoom.RoomID
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



func (rooms *SetOfChatrooms) RegisterRoom(roomID RoomID) {
	// TODO be aware does not check if the room is already registered to them
	rooms.Chatrooms[roomID] = true
}
func (rooms *SetOfChatrooms) UnregisterRoom(roomID RoomID) {
	delete(rooms.Chatrooms, roomID)
}
