package model

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"time"

	echo "github.com/labstack/echo/v4"
	websocket "nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)


type Chatroom struct {
	RoomID            RoomID
	RoomName          string
	RoomOwner         UserID
	ChatHistory       []*Message
	// ClientConnections map[*websocket.Conn]*ClientInfo
	// AllClients	      []*ClientInfo
	LiveUsers         map[UserID]*ClientInfo
	JoinQueue         chan *ChatroomClient
	LeaveQueue        chan *ChatroomClient
	BroadcastQueue    chan *Message
}

type ChatroomClient struct {
	WebSocket *websocket.Conn
	UserID    UserID
	RoomID    RoomID
}

type ClientInfo struct {
	Client       *ChatroomClient
	Conn         *websocket.Conn
	Role         string
	RoomID       string
	MessageQueue chan *Message
}

type Message struct {
	RoomID    RoomID
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


func NewChatroom(uid UserID, rid RoomID, roomName string) *Chatroom {
	newRoom := &Chatroom{
		RoomID:            rid,
		RoomName:          roomName,
		RoomOwner:         uid,
		ChatHistory:       []*Message{},
		// ClientConnections: make(map[*websocket.Conn]*ClientInfo),
		// AllClients:        []*ClientInfo{},
		LiveUsers:         make(map[UserID]*ClientInfo),
		JoinQueue:         make(chan *ChatroomClient),
		LeaveQueue:        make(chan *ChatroomClient),
		BroadcastQueue:    make(chan *Message),
	}

	return newRoom
}


//Chatroom Functions - maybe create an interface for them?
//Open(), RenderChatroomPage(), HandleNewMessage(), HandleNewMessage(), HandleChatroomLogs()

//Open chatroom
//Loop until: ...
//If recieved a new client on JoinQueue channel
	//start goroutine for client, in a loop wait for client to write to their ws, then we will read it from the websocket and send it to the BroadcastQueue channel for everybody to consume
	//add client as a Chatter in this room
//If recieved a client on Leave queue channel
	//TODO: remove client  from room
//If chat recieves new message on BroadcastQueue channel from some user
	//TODO: Add message to db 
	//Send new message to every LiveUser in the room
func (room *Chatroom) Open() {
	for {
		select {
		case client := <-room.JoinQueue:
			// TODO add to this chat

			log.Println("new user joined!")

			//wait for msg from client, braodcast new msg to room
			go func(client *ChatroomClient) {
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
						RoomID:    RoomID(messageReceived.RoomID),
						From:    UserID(messageReceived.UserID),
						Content: messageReceived.UserMessage,
					}

					//why do we send out the msg we read? where is this used?
					//wait i think we are waiting for new mesage FROM user so then we have to send it out to everybody
					room.BroadcastQueue <- &newMessage
					log.Print(messageReceived)
				}
			}(client)

			//add user to room LiveUsers 
			newClient := &ClientInfo{
				Client:       client,
				Conn:         client.WebSocket,
				Role:         "chatter",
				RoomID:       string(client.RoomID),
				MessageQueue: make(chan *Message),
			}
			room.LiveUsers[client.UserID] = newClient

		case client := <-room.LeaveQueue:
			echo.New().Logger.Printf("user leaving room! joined LEAVE queue!")

			delete(room.LiveUsers, client.UserID)
			client.WebSocket.CloseNow()

		case newMessage := <-room.BroadcastQueue:
			echo.New().Logger.Printf("new message to broadcast -> %i", newMessage)

			// TODO add to db
			room.logMessage(newMessage)

			for _, user := range room.LiveUsers {
				log.Println("attempting to write to every user's ws")

				userWS := user.Conn
				wsWriter, writeErr := userWS.Writer(
					context.TODO(),
					websocket.MessageText,
				)
				// TODO if websocket closed handle it, remove from connections, etc.
				if writeErr != nil {
					log.Println(`error creating ws writer!`, writeErr.Error())
					continue
				} 

				log.Println("writer opened")

				// TODO better way to do this
				chatroomTemplates := template.Must(template.ParseFiles("templates/pages/chatroom.html"))
				singleMessageTemplate := chatroomTemplates.Lookup("single-message")

				templateData := MessageHTML{
					DivID:       "chat-messages",
					PrependMsg:  false,
					DisplayName: string(newMessage.From),
					TextMessage: newMessage.Content,
				}

				log.Printf("msg created: %v", templateData)

				singleMessageTemplate.Execute(
					wsWriter,
					&templateData,
				)

				wsWriter.Close()
			}
		//TODO handle deletion of chatroom, remove and close everything
		// case <-done:
			//handle deletion of 
		}
	}

}


func (room *Chatroom) RenderChatroomPage(c echo.Context) error {

	// TODO check for unauthorized access

	const chatroomTemplate = "chatroom"
	templateData := map[string]string{
		"RoomName": room.RoomName,
		"RoomID":   string(room.RoomID),
	}

	return c.Render(http.StatusOK, chatroomTemplate, templateData)
}


//used in hub_handler
func (room *Chatroom) HandleNewConnection(c echo.Context) error {

	echo.New().Logger.Printf("New websocket connection received! isWebsocket='%s'", c.IsWebSocket())

	if !c.IsWebSocket() {
		return c.NoContent(http.StatusBadRequest)
	}

	userID := c.Param("userID")
	roomID := c.Param("roomID")

	// echo.New().Logger.Printf(" data received... %i", User{
	// 	UserID: UserID(userID),
	// 	RoomID: RoomID(roomID),
	// })

	clientWS, err := websocket.Accept(c.Response().Writer, c.Request(), nil)
	// TODO check err
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	room.JoinQueue <- &ChatroomClient{
		WebSocket: clientWS,
		UserID:    UserID(userID),
		RoomID:    RoomID(roomID),
	}

	return nil

}


func (room *Chatroom) HandleNewMessage(c echo.Context) error {

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
		RoomID:    RoomID(messageReceived.RoomID),
		From:    UserID(messageReceived.UserID),
		Content: messageReceived.UserMessage,
	}

	room.BroadcastQueue <- &newMessage

	return nil
}

func (room *Chatroom) HandleChatroomLogs(c echo.Context) error {
	echo.New().Logger.Printf("get Chat history request received")
	log.Printf("current chat log: %v", room.getChatLogs())

	chatHistory := room.getChatLogs()

	msgsData := map[string][]MessageHTML{}
	const messagesLoopID = "ChatMessages"

	for _, chatMessage := range chatHistory {

		log.Printf(`msg: "[%s]" by: [%s]`, chatMessage.Content, chatMessage.From)

		singleMsgData := MessageHTML{
			DivID:       "chat-messages",
			PrependMsg:  false,
			DisplayName: string(chatMessage.From),
			TextMessage: chatMessage.Content,
		}

		msgsData[messagesLoopID] = append(msgsData[messagesLoopID], singleMsgData)
	}

	log.Print(msgsData)

	const msgsTemplateID = "many-messages"
	return c.Render(http.StatusOK, msgsTemplateID, msgsData)

}


func HandleCreateChatRoom(c echo.Context) error {
	return nil
}



func (room *Chatroom) getChatLogs() []*Message {
	return room.ChatHistory
}

func (room *Chatroom) logMessage(msg *Message) {
	room.ChatHistory = append(room.ChatHistory, msg)
}

func (room *Chatroom) GetRID() string {
	return string(room.RoomID)
}

func (room *Chatroom) GetName() string {
	return room.RoomName
}