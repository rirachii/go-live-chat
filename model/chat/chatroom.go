package chat_model

import (
	context "context"
	errors "errors"
	fmt "fmt"
	template "html/template"
	log "log"

	echo "github.com/labstack/echo/v4"
	model "github.com/rirachii/golivechat/model"

	websocket "nhooyr.io/websocket"
	wsjson "nhooyr.io/websocket/wsjson"

	chat_template "github.com/rirachii/golivechat/templates/chat"
)

type Chatroom struct {
	info     ChatroomInfo
	chatLogs []*Message
	// ClientConnections map[*websocket.Conn]*ClientInfo
	// AllClients	      []*ClientInfo
	activeUsers map[model.UserID]*ChatroomUser

	// TODO change to private fields and make function to add to queues
	joinQueue      chan *ChatroomClient
	leaveQueue     chan *ChatroomClient
	broadcastQueue chan *Message
}

type ChatroomInfo struct {
	RoomID    model.RoomID
	RoomName  string
	RoomOwner model.UserID
}

type Message struct {
	RoomID  model.RoomID
	From    model.UserID
	Content string
}

func NewChatroom(userReq model.UserRequest, roomName string) *Chatroom {

	var (
		rid = userReq.RoomID
		uid = userReq.UserID
	)

	chatroomInfo := ChatroomInfo{
		RoomID:    rid,
		RoomName:  roomName,
		RoomOwner: uid,
	}

	newRoom := &Chatroom{
		info:     chatroomInfo,
		chatLogs: []*Message{},
		// ClientConnections: make(map[*websocket.Conn]*ClientInfo),
		// AllClients:        []*ClientInfo{},
		activeUsers:    make(map[model.UserID]*ChatroomUser),
		joinQueue:      make(chan *ChatroomClient),
		leaveQueue:     make(chan *ChatroomClient),
		broadcastQueue: make(chan *Message),
	}

	return newRoom
}

func (room Chatroom) Info() ChatroomInfo   { return room.info }
func (room Chatroom) ID() model.RoomID     { return room.info.RoomID }
func (room Chatroom) Name() string         { return room.info.RoomName }
func (room Chatroom) ChatLogs() []*Message { return room.chatLogs }

func (room *Chatroom) ActiveUsers() map[model.UserID]*ChatroomUser  { return room.activeUsers }
func (room *Chatroom) AddUser(uid model.UserID, user *ChatroomUser) { room.activeUsers[uid] = user }
func (room *Chatroom) RemoveUser(uid model.UserID)                  { delete(room.activeUsers, uid) }

func (room *Chatroom) EnqueueJoin(client *ChatroomClient)       { room.joinQueue <- client }
func (room *Chatroom) EnqueueLeave(client *ChatroomClient)      { room.joinQueue <- client }
func (room *Chatroom) EnqueueMessageBroadcast(message *Message) { room.broadcastQueue <- message }

//Chatroom Functions - maybe create an interface for them?
//Open(), RenderChatroomPage(), HandleNewMessage(), HandleNewMessage(), HandleChatroomLogs()

// Open chatroom
// Loop until: ...
// If recieved a new client on JoinQueue channel
// start goroutine for client, in a loop wait for client to write to their ws, then we will read it from the websocket and send it to the BroadcastQueue channel for everybody to consume
// add client as a Chatter in this room
// If recieved a client on Leave queue channel
// TODO: remove client  from room
// If chat recieves new message on BroadcastQueue channel from some user
// TODO: Add message to db
// Send new message to every LiveUser in the room
func (room *Chatroom) Open() {

	// TODO close functionality
	for {
		select {
		case client := <-room.joinQueue:
			// TODO add to this chat

			log.Println("new user joined!")

			// will wait for msg from client, braodcast new msg to room

			go room.clientListenWS(client)

			//add user to room LiveUsers

			var (
				uid = client.UserID()
				rid = client.RoomID()
			)

			user := NewChatroomUser(client, uid, rid, "chatter")
			room.AddUser(user.ID(), user)

		case client := <-room.leaveQueue:
			echo.New().Logger.Printf("user leaving room! joined LEAVE queue!")

			client.Websocket().CloseNow()
			room.RemoveUser(client.UserID())

		case message := <-room.broadcastQueue:

			echo.New().Logger.Printf("new message to broadcast -> %i", message)

			// TODO add to db
			room.logMessage(message)

			for _, user := range room.ActiveUsers() {
				go room.sendMessageToUser(user, message)
			}

			//TODO handle deletion of chatroom, remove and close everything
			// case <-done:
			//handle deletion of
		}
	}

}
func (room *Chatroom) logMessage(msg *Message) {
	room.chatLogs = append(room.chatLogs, msg)
}

// used in hub_handler
func (room *Chatroom) AcceptConnection(c echo.Context, userReq model.UserRequest) error {

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

	client := &ChatroomClient{
		Conn:        clientWebsocket,
		UserRequest: userReq,
	}

	room.EnqueueJoin(client)

	return nil

}

// run with go routine
// not 100% sure handling leaving is working
func (room *Chatroom) clientListenWS(client *ChatroomClient) {

	ws := client.Websocket()

	for {
		var messageReceived MessageRequest
		readErr := wsjson.Read(context.TODO(), ws, &messageReceived)
		// TODO handle err
		if readErr != nil {
			// log.Panicln(readErr.Error())
			room.EnqueueLeave(client)
			return
		}
		newMessage := &Message{
			RoomID:  model.RID(messageReceived.RoomID),
			From:    model.UID(messageReceived.UserID),
			Content: messageReceived.UserMessage,
		}

		//why do we send out the msg we read? where is this used?
		//wait i think we are waiting for new mesage FROM user so then we have to send it out to everybody
		// ^^ yea once we read the message, we broadcast it
		room.EnqueueMessageBroadcast(newMessage)
		log.Print(messageReceived)
	}
}

func (room *Chatroom) sendMessageToUser(user *ChatroomUser, msg *Message) {
	// TODO tell room when an error occurs
	log.Println("attempting to write to every user's ws")

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
	chatroomTemplates := template.Must(template.ParseFiles("templates/pages/chatroom.html"))
	singleMessageTemplate := chatroomTemplates.Lookup("single-message")

	templateData := chat_template.PrepareMessage(
		chat_template.WebsocketDivID,
		false,
		string(msg.From),
		msg.Content,
	)

	log.Printf("msg created: %v", templateData)

	singleMessageTemplate.Execute(
		wsWriter,
		&templateData,
	)

	wsWriter.Close()
}
