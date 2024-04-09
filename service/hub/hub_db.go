package hub_service

import "github.com/rirachii/golivechat/model"



type ChatroomsDTO struct {
	Chatrooms []ChatroomInfoDTO
}


type ChatroomInfoDTO struct {
	RoomID   model.RoomID
	RoomName string
}

type ChatroomDTO struct {
	RoomID   model.RoomID
	RoomName  string `db:"room_name"`
	OwnerID  model.UserID
	IsPublic bool
	IsActive bool
}

type dbChatroom struct {
	RoomID int `db:"id"`
	RoomName  string `db:"room_name"`
	OwnerID int `db:"owner_id"`
	IsPublic bool `db:"is_public"`
	IsActive bool `db:"is_active"`
}

type dbChatroomInfo struct {
	RoomID int `db:"id"`
	RoomName  string `db:"room_name"`

}

