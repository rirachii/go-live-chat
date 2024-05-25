package chatroom_repository

type MessageLogIndex struct {
	RoomID   int `db:"id"`
	LogIndex int `db:"ARRAY_LENGTH(logs,1)"`
}

type ChatroomMessage struct {
	RoomID      int    `db:"id"`
	SenderID    int    `db:"sender_id"`
	MessageText string `db:"text_message"`
}

type ChatroomMessages struct {
	RoomID      int      `db:"id"`
	MessageLogs []string `db:"logs"`
}
