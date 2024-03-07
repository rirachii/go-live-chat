package model

import (
	"context"
	"errors"
	"fmt"

	// "errors"
	"html/template"
	"log"
	"net/http"
	"time"

	echo "github.com/labstack/echo/v4"
	websocket "nhooyr.io/websocket"
	wsjson "nhooyr.io/websocket/wsjson"
)

type Chatroom struct {
	roomID    RoomID
	roomName  string
	roomOwner UserID
	chatLogs  []*Message
	// ClientConnections map[*websocket.Conn]*ClientInfo
	// AllClients	      []*ClientInfo
	LiveUsers map[UserID]*ClientInfo

	// TODO change to private fields and make function to add to queues
	JoinQueue      chan *ChatroomClient
	LeaveQueue     chan *ChatroomClient
	BroadcastQueue chan *Message
}

type ChatroomClient struct {
	WebSocket *websocket.Conn
	UserData  UserRequest
}

type ClientInfo struct {
	Client       *ChatroomClient
	Conn         *websocket.Conn
	Role         string
	RoomID       string
	MessageQueue chan *Message
}

type Message struct {
	RoomID  RoomID
	From    UserID
	Content string
}

type MessageRequest struct {
	RoomID      string `json:"room-id"`
	UserID      string `json:"user-id"`
	UserMessage string `json:"chat-message"`
}

type MessageHTML struct {
	DivID       string
	PrependMsg  bool
	DisplayName string
	TextMessage string
}

func NewChatroom(userReq UserRequest, roomName string) *Chatroom {

	var (
		rid = userReq.RoomID
		uid = userReq.UserID
	)

	newRoom := &Chatroom{
		roomID:    rid,
		roomOwner: uid,
		roomName:  roomName,
		chatLogs:  []*Message{},
		// ClientConnections: make(map[*websocket.Conn]*ClientInfo),
		// AllClients:        []*ClientInfo{},
		LiveUsers:      make(map[UserID]*ClientInfo),
		JoinQueue:      make(chan *ChatroomClient),
		LeaveQueue:     make(chan *ChatroomClient),
		BroadcastQueue: make(chan *Message),
	}

	return newRoom
}

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
	for {
		select {
		case client := <-room.JoinQueue:
			// TODO add to this chat

			log.Println("new user joined!")

			// will wait for msg from client, braodcast new msg to room

			go room.clientListenWS(client)

			//add user to room LiveUsers
			newClient := &ClientInfo{
				Client:       client,
				Conn:         client.WebSocket,
				Role:         "chatter",
				RoomID:       string(client.UserData.RoomID),
				MessageQueue: make(chan *Message),
			}

			room.LiveUsers[client.UserData.UserID] = newClient

		case client := <-room.LeaveQueue:
			echo.New().Logger.Printf("user leaving room! joined LEAVE queue!")

			delete(room.LiveUsers, client.UserData.UserID)
			client.WebSocket.CloseNow()

		case newMessage := <-room.BroadcastQueue:

			// TODO extract into own function
			echo.New().Logger.Printf("new message to broadcast -> %i", newMessage)

			// TODO add to db
			room.logMessage(newMessage)

			for _, user := range room.LiveUsers {
				go room.broadcastToUser(user, newMessage)
			}

			//TODO handle deletion of chatroom, remove and close everything
			// case <-done:
			//handle deletion of
		}
	}

}

func (room *Chatroom) GetChatroomData() ChatroomData {

	// TODO check for unauthorized access, maybe add err as return?

	chatroomData := ChatroomData{
		RoomID:   room.roomID,
		RoomName: room.roomName,
	}

	return chatroomData

}

// used in hub_handler
func (room *Chatroom) AcceptConnection(c echo.Context, userReq UserRequest) error {

	echo.New().Logger.Printf("New websocket connection received! isWebsocket='%s'", c.IsWebSocket())

	var (
		w = c.Response().Writer
		r = c.Request()
	)

	clientWS, connErr := websocket.Accept(w, r, nil)

	if connErr != nil {
		errorText := fmt.Sprintf("error accepting connection: %v", connErr)
		return errors.New(errorText)
	}

	client := &ChatroomClient{
		WebSocket: clientWS,
		UserData:  userReq,
	}

	// echo.New().Logger.Printf("chatroom client: ", client)

	room.JoinQueue <- client

	return nil

}

// TODO remove, i dont think this does anything lol
func (room *Chatroom) ReceiveNewMessage(c echo.Context) error {

	// TODO
	// get msg, log it, send to broadcast channel

	if !c.IsWebSocket() {
		return c.NoContent(http.StatusBadRequest)
	}
	clientWS, ws_err := websocket.Accept(c.Response().Writer, c.Request(), nil)
	// TODO handle ws_err
	_ = ws_err

	r := c.Request()
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	var messageReceived *MessageRequest
	readErr := wsjson.Read(ctx, clientWS, messageReceived)
	// TODO handle error
	_ = readErr

	roomID := c.Param("roomID")
	messageReceived.RoomID = roomID

	newMessage := Message{
		RoomID:  RoomID(messageReceived.RoomID),
		From:    UserID(messageReceived.UserID),
		Content: messageReceived.UserMessage,
	}

	room.BroadcastQueue <- &newMessage

	return nil
}

func (room *Chatroom) GetChatroomHistory(c echo.Context) map[string][]MessageHTML {

	echo.New().Logger.Debugf("Current chat history: %v", room.getChatLogs())

	chatHistory := room.getChatLogs()

	chatHistoryData := map[string][]MessageHTML{}
	const messagesLoopID = "ChatMessages"

	for _, chatMessage := range chatHistory {

		log.Printf(`msg: "[%s]" by: [%s]`, chatMessage.Content, chatMessage.From)

		singleMsgData := MessageHTML{
			DivID:       "chat-messages",
			PrependMsg:  false,
			DisplayName: string(chatMessage.From),
			TextMessage: chatMessage.Content,
		}

		chatHistoryData[messagesLoopID] = append(chatHistoryData[messagesLoopID], singleMsgData)
	}

	return chatHistoryData

}

func (room *Chatroom) GetRID() RoomID {
	return room.roomID
}

func (room *Chatroom) GetName() string {
	return room.roomName
}

func (room *Chatroom) getChatLogs() []*Message {
	return room.chatLogs
}

func (room *Chatroom) logMessage(msg *Message) {
	room.chatLogs = append(room.chatLogs, msg)
}

// run with go routine
// not 100% sure handling leaving is working
func (room *Chatroom) clientListenWS(client *ChatroomClient) {

	ws := client.WebSocket

	for {
		var messageReceived MessageRequest
		readErr := wsjson.Read(context.TODO(), ws, &messageReceived)
		// TODO handle err
		if readErr != nil {
			// log.Panicln(readErr.Error())
			room.LeaveQueue <- client
			return
		}
		newMessage := Message{
			RoomID:  RoomID(messageReceived.RoomID),
			From:    UserID(messageReceived.UserID),
			Content: messageReceived.UserMessage,
		}

		//why do we send out the msg we read? where is this used?
		//wait i think we are waiting for new mesage FROM user so then we have to send it out to everybody
		// ^^ yea once we read the message, we broadcast it
		room.BroadcastQueue <- &newMessage
		log.Print(messageReceived)
	}
}

func (room *Chatroom) broadcastToUser(user *ClientInfo, msg *Message) {
	// TODO tell room when an error occurs
	log.Println("attempting to write to every user's ws")

	userWS := user.Conn
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

	templateData := MessageHTML{
		DivID:       "chat-messages",
		PrependMsg:  false,
		DisplayName: string(msg.From),
		TextMessage: msg.Content,
	}

	log.Printf("msg created: %v", templateData)

	singleMessageTemplate.Execute(
		wsWriter,
		&templateData,
	)

	wsWriter.Close()
}