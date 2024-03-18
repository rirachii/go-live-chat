package hub_service

import (
	"context"
	"log"
	"time"

	hub_model "github.com/rirachii/golivechat/model/hub"
)

type HubService interface {
	CreateRoom(ctx context.Context, req *hub_model.CreateRoomRequest) (ChatroomDTO, error)
	GetRoom(ctx context.Context, req *hub_model.GetRoomRequest) (ChatroomDTO, error)
	GetRoomsPublic(ctx context.Context, req *hub_model.GetChatroomsRequest) (ChatroomsDTO, error)
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
	req *hub_model.CreateRoomRequest,
) (res ChatroomDTO, err error) {

	chatroom, err := svc.Repo().CreateRoom(ctx, req); if err != nil {
		return ChatroomDTO{}, err
	}

	res = ChatroomDTO(chatroom)

	return res, nil
}

func (svc *hubService) GetRoom(
	ctx context.Context,
	req *hub_model.GetRoomRequest,
) (res ChatroomDTO, err error) {

	return ChatroomDTO{}, nil

}

func (svc *hubService) GetRoomsPublic(
	ctx context.Context,
	req *hub_model.GetChatroomsRequest,
) (res ChatroomsDTO, err error) {

	chatrooms, err := svc.Repo().GetRoomsPublic(ctx)
	if err != nil {
		log.Println("error with repo")
		return ChatroomsDTO{}, err
	}

	for _, room := range chatrooms {
		roomDTO := ChatroomDTO(room)

		res.Chatrooms = append(res.Chatrooms, roomDTO)

	}

	return res, nil
}
