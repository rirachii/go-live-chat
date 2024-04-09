package chatroom_service

import "github.com/rirachii/golivechat/model"

type RepoLogMessageRequest struct {
	RoomID   int
	SenderID int
	Message  string
}

// message index starts from 1. First element is the index 1.
type RepoGetMessageRequest struct {
	UserID       int
	RoomID       int
	MessageIndex int
}

type RepoChatMsgLogsRequest struct {
	UserID int
	RoomID int
}

type RepoUserDisplayNameRequest struct {
	UserID int
}

type ChatMessageDTO struct {
	RoomID      model.RoomID
	SenderID    model.UserID
	MessageText string
}

type ChatroomAdminDTO struct {
	AdminID model.UserID
	Role    string
}

type UserDisplayNameDTO struct {
	UserID      model.UserID
	DisplayName string
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

type dbUserDisplayName struct {
	UserID      int    `db:"id"`
	DisplayName string `db:"display_name"`
}
