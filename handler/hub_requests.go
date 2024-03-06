package handler

type CreateRoomRequest struct {
	// TODO: user id or accesstoken instead of user display name
	UserID   string `json:"display-name"`
	RoomName string `json:"room-name"`
}

type JoinRoomRequest struct {
	// TODO user ID instaed of user display name
	UserID string `json:"display-name"`
	RoomID string `json:"room-id"`
}

type LeaveRoomRequest struct {
	// TODO user ID instaed of user display name
	UserID string `json:"display-name"`
	RoomID string `json:"room-id"`
}

type ConnectionRequest struct {
	UserID string `query:"userID"`
	RoomID string `param:"roomID"`
}

// TODO add user token
type RoomRequest struct {
	UserID string `query:"userID" json:"display-name"`
	RoomID string `param:"roomID" json:"room-id"`
}
