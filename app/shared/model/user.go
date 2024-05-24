package model

type UserInfo struct {
	Id   UserID
	Name string
}

type UserRequest struct {
	UserId UserID
	RoomId RoomID
}

func CreateUserInfo(id UserID, name string) UserInfo {
	return UserInfo{
		Id:   id,
		Name: name,
	}
}
func CreateUserRequest(uid UserID, rid RoomID) UserRequest {
	return UserRequest{
		UserId: uid,
		RoomId: rid,
	}
}
