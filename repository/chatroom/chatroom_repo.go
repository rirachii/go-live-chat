package chatroom_repository

import (
	"context"
	"fmt"

	pgx "github.com/jackc/pgx/v5"
)


type ChatroomRepository interface {
	LogMessage(ctx context.Context, msg messageData) (MessageLogIndex, error)
	LogMessageReturn(ctx context.Context, msg messageData) (ChatroomMessage, error)
	GetChatroomMessages(ctx context.Context, id int) (ChatroomMessages, error)
}

type chatroomRepository struct {
	db *pgx.Conn
}


func (repo *chatroomRepository) LogMessage(ctx context.Context, d messageData) (MessageLogIndex, error) {

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

	var res MessageLogIndex

	scanErr := dbRes.Scan(&res.LogIndex)
	if scanErr != nil {
		return res, scanErr
	}

	return res, nil
}

// logs message, and then does another query to return its data.
func (repo *chatroomRepository) LogMessageReturn(ctx context.Context, d messageData) (ChatroomMessage, error) {

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
		return ChatroomMessage{}, indexScanErr
	}

	// get message
	getQuery := createGetMessageSQLStatement(tblName, msgLogs)
	getRes := repo.db.QueryRow(ctx, getQuery, msgIndex)

	var res ChatroomMessage
	msgScanErr := getRes.Scan(&res.SenderID, &res.MessageText)
	if msgScanErr != nil {
		return ChatroomMessage{}, msgScanErr
	}

	return res, nil
}

func (repo *chatroomRepository) GetChatroomMessages(ctx context.Context, chatroomID int) (ChatroomMessages, error) {

	query := `SELECT id, logs
				FROM chatrooms
				WHERE id = $1`

	dbRes := repo.db.QueryRow(ctx, query, chatroomID)

	var dbRoomID int
	var dbLogs []string
	scanErr := dbRes.Scan(&dbRoomID, &dbLogs)
	if scanErr != nil {
		// log.Print(dbRes)
		return ChatroomMessages{}, scanErr
	}

	chatroomLogs := ChatroomMessages{
		RoomID:  dbRoomID,
		MessageLogs: dbLogs,
	}

	return chatroomLogs, nil
}

/*
returns a statement where:
$1 = sender_id,
$2 = message,
$3 = chatroom_id
*/
func createLogMessageSQLStatement(tblName string, logsColName string) string {

	msgData := "ROW($1, $2)::CHAT_MESSAGE"
	stmt := `UPDATE %[1]s
				SET %[2]s = ARRAY_APPEND(%[2]s, %[3]s) 
				WHERE id = $3
				RETURNING id, ARRAY_LENGTH(%[2]s, 1);`

	query := fmt.Sprintf(stmt, tblName, logsColName, msgData)

	return query
}

/*
returns a statement where:
$1 = log_index,
$2 = chatroom_id,
*/
func createGetMessageSQLStatement(tblName string, colName string) string {

	stmt := `SELECT (%[1]v[$1]).sender_id, (%[1]v)[$1].msg
			FROM %[2]v 
			WHERE id = $2;`

	query := fmt.Sprintf(stmt, colName, tblName)

	return query
}
