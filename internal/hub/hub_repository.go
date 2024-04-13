package hub_service

import (
	"context"
	"fmt"
	"log"

	pgx "github.com/jackc/pgx/v5"
	"github.com/rirachii/golivechat/model"
	hub_model "github.com/rirachii/golivechat/model/hub"
)

type HubRepository interface {
	CreateRoom(context.Context, hub_model.CreateRoomRequest) (dbChatroom, error)
	GetRoomByID(context.Context, hub_model.GetChatroomRequest) (dbChatroom, error)
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
	room hub_model.CreateRoomRequest,
) (room_data dbChatroom, err error) {

	table_fields := "(room_name, owner_id, is_public, is_active, logs)"
	stmt := `INSERT INTO chatrooms %s 
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id, room_name, is_public, is_active
				`

	var (
		name      = room.RoomName
		owner_id  = room.UserID
		is_public = room.IsPublic
		is_active = room.IsActive
		logs      = "{}"
	)

	query := fmt.Sprintf(stmt, table_fields)

	var res dbChatroom
	row := r.db.QueryRow(ctx, query, name, owner_id, is_public, is_active, logs)
	err = row.Scan(
		&res.RoomID,
		&res.RoomName,
		&res.IsPublic,
		&res.IsActive,
	)
	if err != nil {
		return dbChatroom{}, err
	}

	return res, err
}

func (r *hubRepository) GetRoomByID(ctx context.Context, req hub_model.GetChatroomRequest) (dbChatroom, error) {

	roomID, err := model.RIDToInt(req.RoomID)
	if err != nil {
		return dbChatroom{}, err
	}

	fields := "id, room_name, owner_id, is_public, is_active"
	query := fmt.Sprintf(`SELECT %s FROM chatrooms WHERE id = $1`, fields)

	row := r.db.QueryRow(ctx, query, roomID)

	var res dbChatroom
	scanErr := row.Scan(
		&res.RoomID,
		&res.RoomName,
		&res.OwnerID,
		&res.IsPublic,
		&res.IsActive,
	)
	if scanErr != nil {
		return dbChatroom{}, scanErr
	}

	return res, nil
}

func (r *hubRepository) GetRoomInfoByID(ctx context.Context, data hub_model.GetChatroomRequest) (dbChatroomInfo, error) {

	roomID, err := model.RIDToInt(data.RoomID)
	if err != nil {
		return dbChatroomInfo{}, err
	}

	fields := "id, room_name"
	query := fmt.Sprintf(`SELECT %s FROM chatrooms WHERE id = $1`, fields)

	row := r.db.QueryRow(ctx, query, roomID)

	var res dbChatroomInfo
	scanErr := row.Scan(&res.RoomID, &res.RoomName)
	if scanErr != nil {
		return dbChatroomInfo{}, scanErr
	}

	return res, nil
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
