package hub_service

import (
	"context"
	"fmt"
	"log"

	pgx "github.com/jackc/pgx/v5"
	hub_model "github.com/rirachii/golivechat/model/hub"
)

type HubRepository interface {
	CreateRoom(context.Context, *hub_model.CreateRoomRequest) (dbChatroomInfo, error)
	GetRoomByID(context.Context, *hub_model.GetRoomRequest) (dbChatroomInfo, error)
	GetRoomsPublic(context.Context) ([]dbChatroomInfo, error)
}

type hubRepository struct {
	db *pgx.Conn
}

func NewHubRepository(db *pgx.Conn) HubRepository {
	return &hubRepository{db: db}
}

func (r *hubRepository) CreateRoom(
	ctx context.Context,
	room *hub_model.CreateRoomRequest,
) (room_data dbChatroomInfo, err error) {


	const (
		cmd   = "INSERT INTO %s %s VALUES %s RETURNING id"
		table = "chatrooms"
		table_fields = "(room_name,	owner_id, is_public, is_active, logs)"
		values = "($1, $2, $3, $4, $5)"
	)

	var (
		name      = room.RoomName
		owner_id  = room.UserID
		is_public = room.IsPublic
		is_active = room.IsActive
		logs      = "{}"
		query     = fmt.Sprintf(cmd, table, table_fields, values)
	)

	var createdRoomID int
	row := r.db.QueryRow(ctx, query, name, owner_id, is_public, is_active, logs)
	err = row.Scan(&createdRoomID); if err != nil {
		return dbChatroomInfo{}, err
	}

	room_data = dbChatroomInfo{
		RoomID: createdRoomID,
		RoomName: name,
	}

	return room_data, err
}

func (r *hubRepository) GetRoomByID(ctx context.Context, data *hub_model.GetRoomRequest) (dbChatroomInfo, error) {
	return dbChatroomInfo{}, nil
}

func (r *hubRepository) GetRoomsPublic(ctx context.Context) ([]dbChatroomInfo, error) {

	const (
		cmd         = "SELECT %s FROM %s WHERE %s AND %s;"
		data_fields = "id, room_name"
		table       = "chatrooms"
		cond1       = "is_public = TRUE"
		cond2       = "is_active = TRUE"
	)

	var (
		query = fmt.Sprintf(cmd, data_fields, table, cond1, cond2)
	)

	rows, query_err := r.db.Query(ctx, query)

	if query_err != nil {
		log.Println("error with query")
		
		return nil, query_err
	}

	chatroomsData := make([]dbChatroomInfo, 0)

	var (
		roomID   int
		roomName string
	)

	// log.Print(rows.Next())
	// log.Print(rows.Values())


	_, err := pgx.ForEachRow(
		rows, []any{&roomID, &roomName},
		func() error {


			rowData := dbChatroomInfo{
				RoomID:   roomID,
				RoomName: roomName,
			}

			chatroomsData = append(chatroomsData, rowData)

			return nil
		})

	if err != nil {
		log.Println("error with rows", err)

		return nil, err
	}


	return chatroomsData, nil
}
