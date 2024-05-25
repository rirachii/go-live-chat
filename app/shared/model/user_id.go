package model

import "strconv"

type UserID struct {
	id    string
	idNum int
}

func (uid UserID) String() string { return uid.id }
func (uid UserID) Int() int       { return uid.idNum }

func NewUserID(uid int) UserID {

	userId := UserID{
		id:    strconv.Itoa(uid),
		idNum: uid,
	}

	return userId
}
