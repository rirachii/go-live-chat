package chatroom_repository

import (
	model "github.com/rirachii/golivechat/app/shared/model"
	pgx "github.com/jackc/pgx/v5"
)


func NewChatroomRepository(db *pgx.Conn) ChatroomRepository {
	return &chatroomRepository{db: db}
}


func CreateMessageData(msg model.Message) messageData {

	data := messageData{
		RoomId: msg.RoomID().Int(),
		SenderId: msg.SenderId().Int(),
		Message: msg.Content(),

	}
	return data

}