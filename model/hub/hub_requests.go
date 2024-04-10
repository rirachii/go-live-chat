package hub_model

import (
	model "github.com/rirachii/golivechat/model"
)

type CreateRoomRequest struct {
	RoomName string `json:"room-name"`
	UserID   model.UserID
	IsPublic bool
	IsActive bool
}


type GetPublicChatroomsRequest struct {
	IsPublic bool
	IsActive bool
}


type GetChatroomsRequest struct {
	UserID   model.UserID
	IsPublic bool
	IsActive bool
}


type GetChatroomRequest struct {
	UserID model.UserID
	RoomID model.RoomID
}

type JoinRoomRequest struct {
	// TODO user ID instaed of user display name
	UserID model.UserID
	RoomID string `json:"room-id"`
}

type LeaveRoomRequest struct {
	// TODO user ID instaed of user display name
	UserID model.UserID
	RoomID string `json:"room-id"`
}

// TODO add user token
type RoomRequest struct {
	UserID string
	RoomID string `param:"roomID" json:"room-id"`
}
