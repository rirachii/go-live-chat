package chatroom_repository

import (
	pgx "github.com/jackc/pgx/v5"
	model "github.com/rirachii/golivechat/app/shared/model"
	repo_model "github.com/rirachii/golivechat/repository/chatroom/model"
)

func NewChatroomRepository(db *pgx.Conn) ChatroomRepository {
	return &chatroomRepository{db: db}
}

func CreateMessageData(msg model.Message) repo_model.MessageData {

	data := repo_model.MessageData{
		RoomId:   msg.RoomID().Int(),
		SenderId: msg.SenderId().Int(),
		Message:  msg.Content(),
	}
	return data

}
