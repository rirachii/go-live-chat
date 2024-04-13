package chatroom_service

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/rirachii/golivechat/model"
	chat_model "github.com/rirachii/golivechat/model/chat"
)

type ChatroomService interface {
	LogMessage(ctx context.Context, req chat_model.SaveUserMessageRequest) (chat_model.ChatMessageDTO, error)
	LogChatroomMessages(ctx context.Context, req chat_model.SaveChatLogsRequest) error
	GetChatroomMessages(ctx context.Context, req chat_model.GetChatLogsRequest) ([]chat_model.ChatMessageDTO, error)
}

type chatroomService struct {
	chatroomRepository ChatroomRepository
	timeout            time.Duration
}

func (svc *chatroomService) Repo() ChatroomRepository       { return svc.chatroomRepository }
func (svc *chatroomService) TimeoutDuration() time.Duration { return svc.timeout }

func NewChatroomService(chatRepo ChatroomRepository) ChatroomService {
	return &chatroomService{
		chatroomRepository: chatRepo,
		timeout:            time.Duration(2) * time.Second,
	}
}

func (svc *chatroomService) LogMessage(
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

	dbRes, dbErr := svc.Repo().LogMessageAndReturn(ctx, dbReq)
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

func (svc *chatroomService) GetChatroomMessages(
	ctx context.Context, req chat_model.GetChatLogsRequest,
) ([]chat_model.ChatMessageDTO, error) {

	roomID, ridErr := model.RIDToInt(req.RoomID)

	if ridErr != nil {
		return []chat_model.ChatMessageDTO{}, errors.New("RoomID is not a number")
	}

	dbRes, dbErr := svc.Repo().GetChatroomMessages(ctx, roomID)
	if dbErr != nil {

		return []chat_model.ChatMessageDTO{}, dbErr
	}
	log.Printf("getting chatroom messages from db: %+v", dbRes)

	dbRoomID := model.IntToRID(dbRes.RoomID)
	chatroomMessages := []chat_model.ChatMessageDTO{}
	for _, m := range dbRes.MsgLogs {

		log.Printf("converting: %+v, of type %T", m, m)
		msgSender, msgText := convertLogMsgToDTO(m)

		msg := chat_model.ChatMessageDTO{
			RoomID:      dbRoomID,
			SenderID:    msgSender,
			MessageText: msgText,
		}
		chatroomMessages = append(chatroomMessages, msg)
	}

	log.Print(chatroomMessages)
	return chatroomMessages, nil
}

func (svc *chatroomService) LogChatroomMessages(
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

		svc.Repo().LogMessage(ctx, logMsgRequest)

	}

	return nil
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
