package chatroom_repository

import "fmt"

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
