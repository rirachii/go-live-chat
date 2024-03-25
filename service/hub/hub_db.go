package hub_service


type ChatroomsDTO struct {
	Chatrooms []ChatroomDTO
}
type ChatroomDTO struct {
	RoomID   int
	RoomName string
}

type ChatroomStatusDTO struct {
	RoomID   int
	OwnerID  int
	IsPublic bool
	IsActive bool
	// Admins int[]
}

type dbChatroomInfo struct {
	RoomID int `db:"id"`
	RoomName  string `db:"room_name"`
}