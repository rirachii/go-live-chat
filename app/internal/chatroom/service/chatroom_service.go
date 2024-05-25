package chatroom_service

import (
	"context"
	// "errors"
	"log"
	"strconv"
	"strings"
	"time"

	model "github.com/rirachii/golivechat/app/shared/model"
	chatroom_repository "github.com/rirachii/golivechat/repository/chatroom"
)



type ChatroomService interface {
	SaveMessage(ctx context.Context, req SaveMessageRequest) (chatroom_repository.ChatroomMessage, error)
	SaveChatroomMessages(ctx context.Context, req SaveChatroomMessagesRequest) error
	GetChatroomMessages(ctx context.Context, req GetChatroomLogsRequest) ([]chatroom_repository.ChatroomMessage, error)
	// GetUsernameByID(ctx context.Context, req user_model.GetUsernameByIDRequest) (username string, err error)
}

type chatroomService struct {
	chatroomRepository chatroom_repository.ChatroomRepository
	// userRepository     user_repo.UserRepository
	timeout            time.Duration
}

func (svc *chatroomService) ChatroomRepo() chatroom_repository.ChatroomRepository   { return svc.chatroomRepository }
// func (svc *chatroomService) UserRepo() user_repo.UserRepository { return svc.userRepository }
func (svc *chatroomService) TimeoutDuration() time.Duration     { return svc.timeout }

func NewChatroomService(chatRepo chatroom_repository.ChatroomRepository, 
	// userRepo user_repo.UserRepository,
) ChatroomService {
	return &chatroomService{
		chatroomRepository: chatRepo,
		// userRepository:     userRepo,
		timeout:            time.Duration(2) * time.Second,
	}
}


func (svc *chatroomService) GetChatroomMessages(
	ctx context.Context, req GetChatroomLogsRequest,
) ([]chatroom_repository.ChatroomMessage, error) {

	roomId := req.RoomID.Int()

	dbRes, dbErr := svc.ChatroomRepo().GetChatroomMessages(ctx, roomId)
	if dbErr != nil {

		return []chatroom_repository.ChatroomMessage{}, dbErr
	}
	log.Printf("getting chatroom messages from db for chatroom[%s]: %d messages found",
		req.RoomID.String(), len(dbRes.MessageLogs),
	)

	dbRoomID := model.NewRoomID(dbRes.RoomID)
	chatroomMessages := []chatroom_repository.ChatroomMessage{}
	for _, m := range dbRes.MessageLogs {

		// log.Printf("converting: %+v, of type %T", m, m)
		// converts data to readable data
		senderId, messageContent := parseMessage(m)

		msg := chatroom_repository.ChatroomMessage{
			RoomID: dbRoomID.Int(),
			SenderID: senderId.Int(),
			MessageText: messageContent,
		}

		chatroomMessages = append(chatroomMessages, msg)
	}

	// log.Print(chatroomMessages)
	return chatroomMessages, nil
}

func (svc *chatroomService) SaveMessage(
	ctx context.Context, req SaveMessageRequest,
) (chatroom_repository.ChatroomMessage, error) {

	uid := req.UserID
	rid := req.RoomID
	msgContent := req.MessageContent

	msgRequest := model.CreateUserRequest(uid, rid)
	
	message := model.CreateMessage(msgRequest, msgContent, model.MessageMetadata{})

	dbReq := chatroom_repository.CreateMessageData(message)

	dbRes, dbErr := svc.ChatroomRepo().LogMessageReturn(ctx, dbReq)
	if dbErr != nil {
		return chatroom_repository.ChatroomMessage{}, dbErr
	}

	return dbRes, nil
}

func (svc *chatroomService) SaveChatroomMessages(
	ctx context.Context, req SaveChatroomMessagesRequest,
) error {

	roomID := req.RoomID
	chatLogs := req.ChatLogs

	for _, log := range chatLogs {

		senderID := log.UserID
		
		msgRequest := model.CreateUserRequest(senderID, roomID)
		msg := model.CreateMessage(msgRequest, log.MessageContent, model.MessageMetadata{})
		msgData := chatroom_repository.CreateMessageData(msg)

		svc.ChatroomRepo().LogMessage(ctx, msgData)

	}

	return nil
}

// func (svc *chatroomService) GetUsernameByID(
// 	ctx context.Context,
// 	req user_model.GetUsernameByIDRequest,
// ) (string, error) {

// 	uid := req.UserID

// 	uidNum, err := model.UIDToInt(uid)
// 	if err != nil {
// 		return "", errors.New("uid must be a valid number")
// 	}

// 	dbUsername, err := svc.UserRepo().GetUsernameByID(ctx, uidNum)
// 	if err != nil || dbUsername.Username == "" {
// 		return "", errors.New("could not find username with uid: " + string(uid))
// 	}

// 	username := dbUsername.Username

// 	return username, nil
// }


func createChatroomService() (ChatroomService, error){


	return &chatroomService{}, nil
}


// msg in expected format: "({id},{text})"
func parseMessage(msgData string) (model.UserID, string) {

	// remove "(" and ")" from ends
	trimData := strings.TrimPrefix(msgData, `(`)
	trimData = strings.TrimSuffix(trimData, `)`)

	// split by first comma ","
	splitData := strings.SplitN(trimData, ",", 2)

	senderID := splitData[0]

	textMsg := splitData[1]

	// double quotes are added onto sentences saved in db. words do not.
	if len(strings.Split(textMsg, " ")) > 1 {
		textMsg = strings.TrimPrefix(textMsg, `"`)
		textMsg = strings.TrimSuffix(textMsg, `"`)
	}

	id, _ := strconv.Atoi(senderID)

	return model.NewUserID(id), textMsg
}