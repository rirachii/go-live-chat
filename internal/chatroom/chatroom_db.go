package chatroom_service


type RepoLogMessage struct {
	RoomID   int
	SenderID int
	Message  string
}

// message index starts from 1. First element is the index 1.
type RepoGetMessage struct {
	UserID       int
	RoomID       int
	MessageIndex int
}

type RepoChatMsgLogs struct {
	UserID int
	RoomID int
}


type dbChatMessage struct {
	RoomID      int    `db:"id"`
	SenderID    int    `db:"sender_id"`
	MessageText string `db:"text_message"`
}

type dbChatroomLogs struct {
	RoomID  int      `db:"id"`
	MsgLogs []string `db:"logs"`
}

type dbMessageLogIndex struct {
	RoomID   int `db:"id"`
	LogIndex int `db:"ARRAY_LENGTH(logs,1)"`
}
