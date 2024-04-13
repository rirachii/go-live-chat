package handler

import (
	context "context"
	errors "errors"
	fmt "fmt"
	template "html/template"
	log "log"

	echo "github.com/labstack/echo/v4"
	db "github.com/rirachii/golivechat/db"
	chat_svc "github.com/rirachii/golivechat/internal/chatroom"
	model "github.com/rirachii/golivechat/model"
	chat_model "github.com/rirachii/golivechat/model/chat"
	user_model "github.com/rirachii/golivechat/model/user"

	websocket "nhooyr.io/websocket"
	wsjson "nhooyr.io/websocket/wsjson"

	chatroom_template "github.com/rirachii/golivechat/templates/chatroom"
)

func NewChatroom(roomInfo model.ChatroomInfo) chat_model.Chatroom {

	newRoom := chatroom{
		info:             roomInfo,
		chatLogs:         []model.Message{},
		lastSavedChatLog: -1,
		activeUsers:      make(map[model.UserID]*chat_model.ChatroomUser),
		joinQueue:        make(chan *chat_model.ChatroomClient),
		leaveQueue:       make(chan *chat_model.ChatroomUser),
		broadcastQueue:   make(chan model.Message),
	}

	return &newRoom
}

type ChatroomStatus int

const (
	Open ChatroomStatus = iota
	Closed
	Loading
)

type chatroom struct {
	// status   ChatroomStatus
	info model.ChatroomInfo

	chatLogs []model.Message
	// holds the index of the last saved chat message
	lastSavedChatLog int
	activeUsers      map[model.UserID]*chat_model.ChatroomUser

	// TODO change to private fields and make function to add to queues
	joinQueue      chan *chat_model.ChatroomClient
	leaveQueue     chan *chat_model.ChatroomUser
	broadcastQueue chan model.Message
}

func (room chatroom) Info() model.ChatroomInfo  { return room.info }
func (room chatroom) ID() model.RoomID          { return room.info.RoomID }
func (room chatroom) Name() string              { return room.info.RoomName }
func (room chatroom) ChatLogs() []model.Message { return room.chatLogs }
func (room chatroom) IndexLastSavedLog() int    { return room.lastSavedChatLog }
func (room *chatroom) UpdateLastSavedIndex()    { room.lastSavedChatLog = len(room.chatLogs) - 1 }

func (room chatroom) ActiveUsers() map[model.UserID]*chat_model.ChatroomUser { return room.activeUsers }
func (room chatroom) IsPublic() bool                                         { return room.info.IsPublic }

func (room *chatroom) AddUser(user *chat_model.ChatroomUser) { room.activeUsers[user.ID()] = user }
func (room *chatroom) RemoveUser(uid model.UserID)           { delete(room.activeUsers, uid) }

func (room *chatroom) EnqueueJoin(client *chat_model.ChatroomClient) { room.joinQueue <- client }
func (room *chatroom) EnqueueLeave(user *chat_model.ChatroomUser)    { room.leaveQueue <- user }
func (room *chatroom) EnqueueMessageBroadcast(message model.Message) { room.broadcastQueue <- message }

func (room *chatroom) Open() {

	// populate self if needed

	room.populateChatLogsFromDB()

	// TODO close functionality
	for {
		select {
		case client := <-room.joinQueue:
			// TODO add to this chat

			var (
				uid      = client.UserID()
				username = client.Username()
				rid      = client.RoomID()
			)

			userInfo := user_model.UserInfo{
				ID: uid, Username: username,
			}

			user := chat_model.NewChatroomUser(client, userInfo, rid, "chatter")
			go room.ListenToUserWS(user)

			room.AddUser(user)

		case user := <-room.leaveQueue:
			echo.New().Logger.Printf("user leaving room! joined LEAVE queue!")

			user.Client().Websocket().CloseNow()
			room.RemoveUser(user.ID())

		case message := <-room.broadcastQueue:

			echo.New().Logger.Printf("new message to broadcast -> %i", message)

			room.LogMessage(message)
			go room.SaveMessagesToDB()

			for _, user := range room.ActiveUsers() {
				go room.SendMessageToUser(user, message)
			}

			//TODO handle deletion of chatroom, remove and close everything
			// case <-done:
			//handle deletion of
		}
	}

}

func (room *chatroom) Close() {

	// close room. save msgs to db.

}

func (room *chatroom) SaveMessagesToDB() error {

	lastSavedIndex := room.IndexLastSavedLog()

	firstUnsavedIndex := lastSavedIndex + 1

	msgsToSave := room.ChatLogs()[firstUnsavedIndex:]

	chatroomSvc, err := createChatroomService()
	if err != nil {
		return err
	}

	var chatroomMessages chat_model.SaveChatLogsRequest
	chatroomMessages.RoomID = room.ID()

	for _, msg := range msgsToSave {

		saveMsgRequest := chat_model.SaveUserMessageRequest{

			UserID:      msg.SenderUID,
			RoomID:      msg.RoomID,
			UserMessage: msg.Content,
		}

		chatroomMessages.ChatLogs = append(chatroomMessages.ChatLogs, saveMsgRequest)

	}

	ctx := context.TODO()
	err = chatroomSvc.LogChatroomMessages(ctx, chatroomMessages)
	if err != nil {
		return err
	}

	return nil

}

func (room *chatroom) LogMessage(msg model.Message) {
	room.chatLogs = append(room.chatLogs, msg)
}

// used in hub_handler
func (room *chatroom) AcceptConnection(c echo.Context, userReq model.UserRequest) error {

	echo.New().Logger.Printf("New websocket connection received! isWebsocket='%s'", c.IsWebSocket())

	var (
		w = c.Response().Writer
		r = c.Request()
	)

	clientWebsocket, connErr := websocket.Accept(w, r, nil)

	if connErr != nil {
		errorText := fmt.Sprintf("error accepting connection: %v", connErr)
		return errors.New(errorText)
	}

	client := &chat_model.ChatroomClient{
		Conn:        clientWebsocket,
		UserRequest: userReq,
	}

	room.EnqueueJoin(client)

	return nil

}

// run with go routine
// not 100% sure handling leaving is working
func (room *chatroom) ListenToUserWS(user *chat_model.ChatroomUser) {

	ws := user.Client().Websocket()
	userID := user.ID()

	for {
		messageToSend := chat_model.SendMessageRequest{
			UserID: userID,
		}
		readErr := wsjson.Read(context.TODO(), ws, &messageToSend)

		// TODO handle err
		if readErr != nil {
			room.EnqueueLeave(user)
			return
		}

		newMessage := model.Message{
			RoomID:         model.RID(messageToSend.RoomID),
			SenderUsername: user.Username(),
			SenderUID:      user.ID(),
			Content:        messageToSend.MessageText,
		}

		room.EnqueueMessageBroadcast(newMessage)
		// log.Print(messageToSend)
	}
}

func (room *chatroom) SendMessageToUser(user *chat_model.ChatroomUser, msg model.Message) {
	// TODO tell room when an error occurs
	// log.Println("attempting to write to every user's ws")

	userWS := user.Client().Websocket()
	wsWriter, writeErr := userWS.Writer(
		context.TODO(),
		websocket.MessageText,
	)

	// TODO if websocket closed handle it, remove from connections, etc.
	if writeErr != nil {
		log.Println(`error creating ws writer!`, writeErr.Error())
		return
	}

	// TODO better way to do this
	chatroomTemplates := template.Must(template.ParseFiles("templates/chatroom/chatroom.html"))
	singleMessageTemplate := chatroomTemplates.Lookup("single-message")

	templateData := chatroom_template.PrepareMessage(
		chatroom_template.WebsocketDivID,
		false,
		msg.SenderUsername,
		msg.Content,
	)

	log.Printf("msg created: %v", templateData)

	singleMessageTemplate.Execute(
		wsWriter,
		&templateData,
	)

	wsWriter.Close()
}

func (room *chatroom) populateChatLogsFromDB() error {

	failedErr := errors.New("failed to get get chat logs from db")

	chatroom_service, err := createChatroomService()
	if err != nil {
		return errors.Join(failedErr, err)
	}

	ctx := context.Background()
	req := chat_model.GetChatLogsRequest{
		RoomID: room.ID(),
	}

	messages, dbErr := chatroom_service.GetChatroomMessages(ctx, req)
	if dbErr != nil {
		return errors.Join(failedErr, dbErr)
	}

	for _, m := range messages {

		msg := model.Message{
			RoomID:    m.RoomID,
			SenderUID: m.SenderID,
			Content:   m.MessageText,
		}

		room.LogMessage(msg)
	}

	room.UpdateLastSavedIndex()

	return nil
}

func createChatroomService() (chat_svc.ChatroomService, error) {

	dbConn, err := db.ConnectDatabase()
	if err != nil {
		return nil, err
	}

	chatroomRepo := chat_svc.NewChatroomRepository(dbConn.DB())
	chatroomSvc := chat_svc.NewChatroomService(chatroomRepo)

	return chatroomSvc, nil

}
