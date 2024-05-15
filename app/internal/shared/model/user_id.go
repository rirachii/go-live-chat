package model

import "strconv"

type UserID struct {
	id    string
	idNum int
}

func (uid UserID) ID() string { return uid.id }
func (uid UserID) IntID() int { return uid.idNum }

func NewUserID(uid int) UserID {

	userId := UserID{
		id:    strconv.Itoa(uid),
		idNum: uid,
	}

	return userId
}
