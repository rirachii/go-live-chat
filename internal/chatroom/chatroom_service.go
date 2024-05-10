package chatroom_service

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	user_repo "github.com/rirachii/golivechat/internal/user"
	model "github.com/rirachii/golivechat/model"
	chat_model "github.com/rirachii/golivechat/model/chat"
	user_model "github.com/rirachii/golivechat/model/user"
)

type ChatroomService interface {
	SaveMessage(ctx context.Context, req chat_model.SaveUserMessageRequest) (chat_model.ChatMessageDTO, error)
	SaveChatroomMessages(ctx context.Context, req chat_model.SaveChatLogsRequest) error
	GetChatroomMessages(ctx context.Context, req chat_model.GetChatLogsRequest) ([]chat_model.ChatMessageDTO, error)
	GetUsernameByID(ctx context.Context, req user_model.GetUsernameByIDRequest) (username string, err error)
}

type chatroomService struct {
	chatroomRepository ChatroomRepository
	userRepository     user_repo.UserRepository
	timeout            time.Duration
}

func (svc *chatroomService) ChatroomRepo() ChatroomRepository   { return svc.chatroomRepository }
func (svc *chatroomService) UserRepo() user_repo.UserRepository { return svc.userRepository }
func (svc *chatroomService) TimeoutDuration() time.Duration     { return svc.timeout }

func NewChatroomService(chatRepo ChatroomRepository, userRepo user_repo.UserRepository) ChatroomService {
	return &chatroomService{
		chatroomRepository: chatRepo,
		userRepository:     userRepo,
		timeout:            time.Duration(2) * time.Second,
	}
}

func (svc *chatroomService) GetChatroomMessages(
	ctx context.Context, req chat_model.GetChatLogsRequest,
) ([]chat_model.ChatMessageDTO, error) {

	roomID, ridErr := model.RIDToInt(req.RoomID)

	if ridErr != nil {
		return []chat_model.ChatMessageDTO{}, errors.New("RoomID is not a number")
	}

	dbRes, dbErr := svc.ChatroomRepo().GetChatroomMessages(ctx, roomID)
	if dbErr != nil {

		return []chat_model.ChatMessageDTO{}, dbErr
	}
	log.Printf("getting chatroom messages from db for chatroom[%s]: %d messages found",
		req.RoomID, len(dbRes.MsgLogs),
	)

	dbRoomID := model.IntToRID(dbRes.RoomID)
	chatroomMessages := []chat_model.ChatMessageDTO{}
	for _, m := range dbRes.MsgLogs {

		// log.Printf("converting: %+v, of type %T", m, m)
		msgSender, msgText := convertLogMsgToDTO(m)

		msg := chat_model.ChatMessageDTO{
			RoomID:      dbRoomID,
			SenderID:    msgSender,
			MessageText: msgText,
		}
		chatroomMessages = append(chatroomMessages, msg)
	}

	// log.Print(chatroomMessages)
	return chatroomMessages, nil
}

func (svc *chatroomService) SaveMessage(
	ctx context.Context, req chat_model.SaveUserMessageRequest,
) (chat_model.ChatMessageDTO, error) {

	userID, uidErr := model.UIDToInt(req.UserID)
	roomID, ridErr := model.RIDToInt(req.RoomID)

	if uidErr != nil || ridErr != nil {
		return chat_model.ChatMessageDTO{}, errors.New("RoomID or UserID is not a number")
	}

	userMsg := req.UserMessage

	dbReq := RepoLogMessage{
		RoomID:   roomID,
		SenderID: userID,
		Message:  userMsg,
	}

	dbRes, dbErr := svc.ChatroomRepo().LogMessageAndReturn(ctx, dbReq)
	if dbErr != nil {
		return chat_model.ChatMessageDTO{}, dbErr
	}

	data := chat_model.ChatMessageDTO{
		RoomID:      model.IntToRID(dbReq.RoomID),
		SenderID:    model.IntToUID(dbRes.SenderID),
		MessageText: dbRes.MessageText,
	}

	return data, nil
}

func (svc *chatroomService) SaveChatroomMessages(
	ctx context.Context, req chat_model.SaveChatLogsRequest,
) error {

	roomID, _ := model.RIDToInt(req.RoomID)
	chatLogs := req.ChatLogs

	for _, log := range chatLogs {

		senderID, _ := model.UIDToInt(log.UserID)

		logMsgRequest := RepoLogMessage{
			RoomID:   roomID,
			SenderID: senderID,
			Message:  log.UserMessage,
		}

		svc.ChatroomRepo().LogMessage(ctx, logMsgRequest)

	}

	return nil
}

func (svc *chatroomService) GetUsernameByID(
	ctx context.Context,
	req user_model.GetUsernameByIDRequest,
) (string, error) {

	uid := req.UserID

	uidNum, err := model.UIDToInt(uid)
	if err != nil {
		return "", errors.New("uid must be a valid number")
	}

	dbUsername, err := svc.UserRepo().GetUsernameByID(ctx, uidNum)
	if err != nil || dbUsername.Username == "" {
		return "", errors.New("could not find username with uid: " + string(uid))
	}

	username := dbUsername.Username

	return username, nil
}

// msg in expected format: "({id},{text})"
func convertLogMsgToDTO(msgData string) (model.UserID, string) {

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

	return model.UID(senderID), textMsg
}
