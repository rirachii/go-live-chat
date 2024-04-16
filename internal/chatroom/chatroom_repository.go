package chatroom_service

import (
	"context"
	"fmt"
	"log"

	pgx "github.com/jackc/pgx/v5"
)

type ChatroomRepository interface {
	LogMessage(ctx context.Context, msgData RepoLogMessage) (dbMessageLogIndex, error)
	LogMessageAndReturn(ctx context.Context, msgData RepoLogMessage) (dbChatMessage, error)
	GetChatroomMessages(ctx context.Context, id int) (dbChatroomLogs, error)
	// GetAdmins
	// AddAdmins
}

type chatroomRepository struct {
	db *pgx.Conn
}

func NewChatroomRepository(db *pgx.Conn) ChatroomRepository {
	return &chatroomRepository{db: db}
}

func (repo *chatroomRepository) LogMessage(ctx context.Context, req RepoLogMessage) (dbMessageLogIndex, error) {

	const (
		table   = "chatrooms"
		msgLogs = "logs"
	)

	var (
		chatroom_id = req.RoomID
		sender_id   = req.SenderID
		msg         = req.Message
	)

	query := createLogMessageSQLStatement(table, msgLogs)
	dbRes := repo.db.QueryRow(ctx, query, sender_id, msg, chatroom_id)

	var res dbMessageLogIndex

	scanErr := dbRes.Scan(&res.LogIndex)
	if scanErr != nil {
		return res, scanErr
	}

	return res, nil
}

// logs message, and then does another query to return its data.
func (repo *chatroomRepository) LogMessageAndReturn(ctx context.Context, req RepoLogMessage) (dbChatMessage, error) {

	var (
		tblName     = "chatrooms"
		msgLogs     = "logs"
		chatroom_id = req.RoomID
		sender_id   = req.SenderID
		msg         = req.Message
	)

	// log message
	logQuery := createLogMessageSQLStatement(tblName, msgLogs)
	logRes := repo.db.QueryRow(ctx, logQuery, sender_id, msg, chatroom_id)

	var msgIndex int
	indexScanErr := logRes.Scan(&msgIndex)
	if indexScanErr != nil {
		return dbChatMessage{}, indexScanErr
	}

	// get message
	getQuery := createGetMessageSQLStatement(tblName, msgLogs)
	getRes := repo.db.QueryRow(ctx, getQuery, msgIndex)
	var res dbChatMessage
	msgScanErr := getRes.Scan(&res.SenderID, &res.MessageText)
	if msgScanErr != nil {
		return dbChatMessage{}, msgScanErr
	}

	return res, nil
}

func (repo *chatroomRepository) GetChatroomMessages(ctx context.Context, chatroomID int) (dbChatroomLogs, error) {

	query := `SELECT id, logs
				FROM chatrooms
				WHERE id = $1`

	dbRes := repo.db.QueryRow(ctx, query, chatroomID)

	var dbRoomID int
	var dbLogs []string
	scanErr := dbRes.Scan(&dbRoomID, &dbLogs)
	if scanErr != nil {
		log.Print(dbRes)
		return dbChatroomLogs{}, scanErr
	}

	chatroomLogs := dbChatroomLogs{
		RoomID:  dbRoomID,
		MsgLogs: dbLogs,
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
