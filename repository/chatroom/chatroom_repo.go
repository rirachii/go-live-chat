package chatroom_repository

import (
	"context"

	pgx "github.com/jackc/pgx/v5"
	repo_model "github.com/rirachii/golivechat/repository/chatroom/model"
)

type ChatroomRepository interface {
	LogMessage(ctx context.Context, msg repo_model.MessageData) (repo_model.MessageLogIndex, error)
	LogMessageReturn(ctx context.Context, msg repo_model.MessageData) (repo_model.ChatroomMessage, error)
	GetChatroomMessages(ctx context.Context, id int) (repo_model.ChatroomMessages, error)
}

type chatroomRepository struct {
	db *pgx.Conn
}

func (repo *chatroomRepository) LogMessage(
	ctx context.Context, d repo_model.MessageData,
) (repo_model.MessageLogIndex, error) {

	const (
		table   = "chatrooms"
		msgLogs = "logs"
	)

	var (
		chatroom_id = d.RoomId
		sender_id   = d.SenderId
		msg         = d.Message
	)

	query := createLogMessageSQLStatement(table, msgLogs)
	dbRes := repo.db.QueryRow(ctx, query, sender_id, msg, chatroom_id)

	var res repo_model.MessageLogIndex

	scanErr := dbRes.Scan(&res.LogIndex)
	if scanErr != nil {
		return res, scanErr
	}

	return res, nil
}

// logs message, and then does another query to return its data.
func (repo *chatroomRepository) LogMessageReturn(
	ctx context.Context, d repo_model.MessageData,
) (repo_model.ChatroomMessage, error) {

	var (
		tblName     = "chatrooms"
		msgLogs     = "logs"
		chatroom_id = d.RoomId
		sender_id   = d.SenderId
		msg         = d.Message
	)

	// log message
	logQuery := createLogMessageSQLStatement(tblName, msgLogs)
	logRes := repo.db.QueryRow(ctx, logQuery, sender_id, msg, chatroom_id)

	var msgIndex int
	indexScanErr := logRes.Scan(&msgIndex)
	if indexScanErr != nil {
		return repo_model.ChatroomMessage{}, indexScanErr
	}

	// get message
	getQuery := createGetMessageSQLStatement(tblName, msgLogs)
	getRes := repo.db.QueryRow(ctx, getQuery, msgIndex)

	var res repo_model.ChatroomMessage
	msgScanErr := getRes.Scan(&res.SenderID, &res.MessageText)
	if msgScanErr != nil {
		return repo_model.ChatroomMessage{}, msgScanErr
	}

	return res, nil
}

func (repo *chatroomRepository) GetChatroomMessages(
	ctx context.Context, chatroomID int,
) (repo_model.ChatroomMessages, error) {

	query := `SELECT id, logs
				FROM chatrooms
				WHERE id = $1`

	dbRes := repo.db.QueryRow(ctx, query, chatroomID)

	var dbRoomID int
	var dbLogs []string
	scanErr := dbRes.Scan(&dbRoomID, &dbLogs)
	if scanErr != nil {
		// log.Print(dbRes)
		return repo_model.ChatroomMessages{}, scanErr
	}

	chatroomLogs := repo_model.ChatroomMessages{
		RoomID:      dbRoomID,
		MessageLogs: dbLogs,
	}

	return chatroomLogs, nil
}
