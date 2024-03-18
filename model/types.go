package model

type UserID string
type RoomID string

type UserRequest struct {
	UserID UserID
	RoomID RoomID
}

func RID(rid string) RoomID {
	return RoomID(rid)
}

func UID(uid string) UserID {
	return UserID(uid)
}


type ChatroomInfo struct {
	RoomID    RoomID
	RoomName  string
	RoomOwner UserID
}

type Message struct {
	RoomID  RoomID
	From    UserID
	Content string
}