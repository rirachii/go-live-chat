package hub_model

import (
	model "github.com/rirachii/golivechat/model"
)

type CreateRoomRequest struct {
	UserID   model.UserID
	RoomName string `json:"room-name"`
	IsPublic bool
	IsActive bool
}


type GetChatroomsRequest struct {
	UserID   model.UserID
	IsPublic bool
	IsActive bool
}


type GetRoomRequest struct {
	UserID model.UserID
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
