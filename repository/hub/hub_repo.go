package hub_repository

import (
	"context"
	"fmt"

	pgx "github.com/jackc/pgx/v5"
	repo_model "github.com/rirachii/golivechat/repository/hub/model"
)

type HubRepository interface {
	CreateRoom(ctx context.Context, r repo_model.RoomData) (repo_model.Chatroom, error)
	GetRoomByID(ctx context.Context, id int) (repo_model.Chatroom, error)
	GetRoomsPublic(ctx context.Context) ([]repo_model.ChatroomInfo, error)
}

type hubRepository struct {
	db *pgx.Conn
}

func (r *hubRepository) CreateRoom(
	ctx context.Context,
	room repo_model.RoomData,
) (room_data repo_model.Chatroom, err error) {

	table_fields := "(room_name, owner_id, is_public, is_active, logs)"
	stmt := `INSERT INTO chatrooms %s 
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id, room_name, is_public, is_active
				`

	var (
		name      = room.Name
		owner_id  = room.Owner
		is_public = room.IsPublic
		is_active = room.IsActive
		logs      = "{}"
	)

	query := fmt.Sprintf(stmt, table_fields)

	var res repo_model.Chatroom
	row := r.db.QueryRow(ctx, query, name, owner_id, is_public, is_active, logs)
	err = row.Scan(
		&res.Id,
		&res.Name,
		&res.IsPublic,
		&res.IsActive,
	)
	if err != nil {
		return repo_model.Chatroom{}, err
	}

	return res, err
}

func (r *hubRepository) GetRoomByID(ctx context.Context, id int) (repo_model.Chatroom, error) {

	rid := id

	fields := "id, room_name, owner_id, is_public, is_active"
	query := fmt.Sprintf(`SELECT %s FROM chatrooms WHERE id = $1`, fields)

	row := r.db.QueryRow(ctx, query, rid)

	var res repo_model.Chatroom
	scanErr := row.Scan(
		&res.Id,
		&res.Name,
		&res.Owner,
		&res.IsPublic,
		&res.IsActive,
	)
	if scanErr != nil {
		return repo_model.Chatroom{}, scanErr
	}

	return res, nil
}

func (r *hubRepository) GetRoomInfoByID(ctx context.Context, id int) (repo_model.ChatroomInfo, error) {


	fields := "id, room_name"
	query := fmt.Sprintf(`SELECT %s FROM chatrooms WHERE id = $1`, fields)

	rid := id
	row := r.db.QueryRow(ctx, query, rid)

	var res repo_model.ChatroomInfo
	scanErr := row.Scan(&res.Id, &res.Name)
	if scanErr != nil {
		return repo_model.ChatroomInfo{}, scanErr
	}

	return res, nil
}

func (r *hubRepository) GetRoomsPublic(ctx context.Context) ([]repo_model.ChatroomInfo, error) {

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
		// log.Println("error with query")

		return nil, query_err
	}

	chatroomsData := make([]repo_model.ChatroomInfo, 0)

	var (
		Id   int
		Name string
	)

	// log.Print(rows.Next())
	// log.Print(rows.Values())

	_, err := pgx.ForEachRow(
		rows, []any{&Id, &Name},
		func() error {

			rowData := repo_model.ChatroomInfo{
				Id:   Id,
				Name: Name,
			}

			chatroomsData = append(chatroomsData, rowData)

			return nil
		})

	if err != nil {
		// log.Println("error with rows", err)

		return nil, err
	}

	return chatroomsData, nil
}
