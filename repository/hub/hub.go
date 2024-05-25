package hub_repository

import (
	pgx "github.com/jackc/pgx/v5"
	model "github.com/rirachii/golivechat/app/shared/model"
	repo_model "github.com/rirachii/golivechat/repository/hub/model"
)



func NewHubRepository(db *pgx.Conn) HubRepository {
	return &hubRepository{db: db}
}

func CreateRoomData(name string, owner model.UserID, isPublic bool) repo_model.RoomData {

	room := repo_model.RoomData{
		Name: name,
		Owner: owner,
		IsPublic: isPublic,
		IsActive: false,
	}

	return room
}