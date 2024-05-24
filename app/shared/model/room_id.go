package model

import "strconv"

type RoomID struct {
	id    string
	idNum int
}

func (uid RoomID) ID() string { return uid.id }
func (uid RoomID) IntID() int { return uid.idNum }

func NewRoomID(uid int) RoomID {

	roomId := RoomID{
		id:    strconv.Itoa(uid),
		idNum: uid,
	}

	return roomId
}
