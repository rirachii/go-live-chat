package hub_repository_model

type ChatroomInfo struct {
	Id   int    `db:"id"`
	Name string `db:"room_name"`
}


type Chatroom struct {
	Id       int    `db:"id"`
	Name     string `db:"room_name"`
	Owner  int    `db:"owner_id"`
	IsPublic bool   `db:"is_public"`
	IsActive bool   `db:"is_active"`
}

