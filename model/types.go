package model

import "strconv"

type UserID string
type RoomID string

type UserRequest struct {
	UserID   UserID
	Username string
	RoomID   RoomID
}

func UID(uid string) UserID            { return UserID(uid) }
func UIDToInt(uid UserID) (int, error) { return strconv.Atoi(string(uid)) }
func IntToUID(id int) UserID           { return UID(strconv.Itoa(id)) }

func RID(rid string) RoomID            { return RoomID(rid) }
func RIDToInt(rid RoomID) (int, error) { return strconv.Atoi(string(rid)) }
func IntToRID(id int) RoomID           { return RID(strconv.Itoa(id)) }

type ChatroomInfo struct {
	RoomID    RoomID
	RoomName  string
	RoomOwner UserID
	IsPublic  bool
}

type Message struct {
	RoomID        RoomID // room id
	SenderUID     UserID // user's id
	SenderName string // user name
	Content       string // content of message
}
