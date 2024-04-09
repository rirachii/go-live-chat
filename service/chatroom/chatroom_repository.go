package chatroom_service

import (
	"context"
	"fmt"
	"log"

	pgx "github.com/jackc/pgx/v5"
)

type ChatroomRepository interface {
	LogMessage(context.Context, RepoLogMessageRequest) (dbMessageLogIndex, error)
	LogMessageAndReturn(context.Context, RepoLogMessageRequest) (dbChatMessage, error)
	GetChatroomMessages(context.Context, RepoChatMsgLogsRequest) (dbChatroomLogs, error)
	GetUserDisplayName(context.Context, RepoUserDisplayNameRequest) (dbUserDisplayName, error)
	// GetAdmins
	// AddAdmins
}

type chatroomRepository struct {
	db *pgx.Conn
}

func NewChatroomRepository(db *pgx.Conn) ChatroomRepository {
	return &chatroomRepository{db: db}
}

func (repo *chatroomRepository) LogMessage(ctx context.Context, req RepoLogMessageRequest) (dbMessageLogIndex, error) {

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
func (repo *chatroomRepository) LogMessageAndReturn(ctx context.Context, req RepoLogMessageRequest) (dbChatMessage, error) {

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

func (repo *chatroomRepository) GetChatroomMessages(ctx context.Context, req RepoChatMsgLogsRequest) (dbChatroomLogs, error) {

	query := `SELECT id, logs
				FROM chatrooms
				WHERE id = $1`

	dbRes := repo.db.QueryRow(ctx, query, req.RoomID)

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

func (r *chatroomRepository) GetUserDisplayName(ctx context.Context, req RepoUserDisplayNameRequest) (dbUserDisplayName, error) {
	query := `SELECT id, display_name FROM user_display_names WHERE id = $1`

	var id int
	var displayName string

	err := r.db.QueryRow(ctx, query, id).Scan(&id, &displayName)
	if err != nil {
		return dbUserDisplayName{}, err
	}

	res := dbUserDisplayName{
		UserID:      id,
		DisplayName: displayName,
	}

	return res, err
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
