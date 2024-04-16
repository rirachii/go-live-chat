package hub_service

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/rirachii/golivechat/model"
	hub_model "github.com/rirachii/golivechat/model/hub"
)

type HubService interface {
	CreateRoom(context.Context, hub_model.CreateRoomRequest) (ChatroomDTO, error)
	// GetRoom(context.Context, hub_model.GetChatroomRequest) (ChatroomDTO, error)
	GetRoomInfo(context.Context, hub_model.GetChatroomRequest) (ChatroomInfoDTO, error)
	GetRoomsPublic(context.Context, hub_model.GetPublicChatroomsRequest) ([]ChatroomDTO, error)
}

type hubService struct {
	hubRepository HubRepository
	timeout       time.Duration
}

func (svc *hubService) Repo() HubRepository            { return svc.hubRepository }
func (svc *hubService) TimeoutDuration() time.Duration { return svc.timeout }

func NewHubService(repository HubRepository) HubService {
	return &hubService{
		hubRepository: repository,
		timeout:       time.Duration(2) * time.Second,
	}
}

func (svc *hubService) CreateRoom(
	ctx context.Context,
	req hub_model.CreateRoomRequest,
) (res ChatroomDTO, err error) {

	chatroom, err := svc.Repo().CreateRoom(ctx, req); if err != nil {
		return ChatroomDTO{}, err
	}

	res = ChatroomDTO{
		RoomID: model.RID(strconv.Itoa(chatroom.RoomID)),
		RoomName: chatroom.RoomName,
		
	}

	return res, nil
}

func (svc *hubService) GetRoomInfo(
	ctx context.Context,
	req hub_model.GetChatroomRequest,
) (ChatroomInfoDTO, error) {

	chatroom, err := svc.Repo().GetRoomByID(ctx, req)
	if err != nil {
		return ChatroomInfoDTO{}, err
	}

	res := ChatroomInfoDTO{
		RoomID: model.IntToRID(chatroom.RoomID),
		RoomName: chatroom.RoomName,

	}

	return res, nil

}

func (svc *hubService) GetRoomsPublic(
	ctx context.Context,
	req hub_model.GetPublicChatroomsRequest,
) ([]ChatroomDTO, error) {

	dbChatrooms, err := svc.Repo().GetRoomsPublic(ctx)
	if err != nil {
		log.Println("error with repo")
		return []ChatroomDTO{}, err
	}


	chatrooms := []ChatroomDTO{}
	for _, room := range dbChatrooms {

		roomDTO := ChatroomDTO{
			RoomID: model.IntToRID(room.RoomID),
			RoomName: room.RoomName,

		}

		chatrooms = append(chatrooms, roomDTO)

	}

	return chatrooms, nil
}
