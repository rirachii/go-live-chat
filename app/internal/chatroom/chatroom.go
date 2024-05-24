package chatroom

import (
	"context"
	"os/user"

	chatroom_model "github.com/rirachii/golivechat/app/internal/chatroom/model"
	model "github.com/rirachii/golivechat/app/shared/model"
	"nhooyr.io/websocket/wsjson"
)

// TODO. Use them
type ChatroomStatus int

const (
	Open ChatroomStatus = iota
	Closed
	Loading
)

// The fields for a chatroom.
//
// Chat messages are saved locally first, then entered in a queue to save in DB.
//
// Users are [ChatroomUser] when they join and become registered as a [Subscriber] when successful
type chatroom struct {
	status            ChatroomStatus
	info              chatroom_model.ChatroomInfo
	activeSubscribers map[model.UserID]*chatroom_model.Subscriber
	messagesLog       []*chatroom_model.ChatroomMessage

	// TODO change to private fields and make function to add to queues
	joiningQueue   chan *chatroom_model.ChatroomUser
	leavingQueue   chan *chatroom_model.Subscriber
	broadcastQueue chan *chatroom_model.ChatroomMessage

	// holds the index of the last saved chat message
	lastSavedMessage   int
	messageSavingQueue chan *chatroom_model.ChatroomMessage
	// holds the index of the last saved chat message

	closeChannel chan int
}

// Creates a Chatroom Interface.
// TODO: change input to a request that takes JSON stuff.
func CreateChatroom(roomInfo chatroom_model.ChatroomInfo) chatroom_model.Chatroom {
	newChatroom := chatroom{
		status:            Closed,
		info:              roomInfo,
		messagesLog:       make([]*chatroom_model.ChatroomMessage, 0),
		activeSubscribers: make(map[model.UserID]*chatroom_model.Subscriber),

		joiningQueue:   make(chan *chatroom_model.ChatroomUser),
		leavingQueue:   make(chan *chatroom_model.Subscriber),
		broadcastQueue: make(chan *chatroom_model.ChatroomMessage),

		lastSavedMessage:   -1,
		messageSavingQueue: make(chan *chatroom_model.ChatroomMessage),

		closeChannel: make(chan int),
	}

	return &newChatroom
}

func (room chatroom) Info() chatroom_model.ChatroomInfo { return room.info }
func (room chatroom) ID() model.RoomID                  { return room.info.RoomID }
func (room chatroom) Name() string                      { return room.info.RoomName }
func (room chatroom) IsPublic() bool                    { return room.info.IsPublic }

func (room chatroom) Messages() []*chatroom_model.ChatroomMessage { return room.messagesLog }
func (room chatroom) SavedMessages() []*chatroom_model.ChatroomMessage {
	return room.messagesLog[:room.lastSavedMessage+1]
}
func (room chatroom) UnsavedMessages() []*chatroom_model.ChatroomMessage {
	return room.messagesLog[room.lastSavedMessage+1:]
}

func (room chatroom) LastSavedMessage() int { return room.lastSavedMessage }

func (room *chatroom) ActiveSubscribers() map[model.UserID]*chatroom_model.Subscriber {
	return room.activeSubscribers
}

func (room *chatroom) AddSubscriber(subscriber *chatroom_model.Subscriber) {
	room.activeSubscribers[subscriber.ID()] = subscriber
}
func (room *chatroom) RemoveSubscriber(uid model.UserID) {
	delete(room.activeSubscribers, uid)
}

func (room *chatroom) EnqueueJoin(user *chatroom_model.ChatroomUser) {
	room.joiningQueue <- user
}
func (room *chatroom) EnqueueLeave(subscriber *chatroom_model.Subscriber) {
	room.leavingQueue <- subscriber
}
func (room *chatroom) Broadcast(message *chatroom_model.ChatroomMessage) {
	room.broadcastQueue <- message
}

func (room *chatroom) Open()  { room.StartChatroom() }
func (room *chatroom) Close() { room.CloseChatroom() }

func (room *chatroom) LogMessage(m *chatroom_model.ChatroomMessage) {
	room.messagesLog = append(room.messagesLog, m)
}
func (room *chatroom) SaveMessage(m *chatroom_model.ChatroomMessage) { room.messageSavingQueue <- m }

func (room *chatroom) ListenToSubscriber(s *chatroom_model.Subscriber) {
	
	ws := s.User().Websocket()
	userId := s.Id()

	ctx := context.Background()

	for {
		request := chatroom_model.NewMessageRequest(userId)
		readErr := wsjson.Read(ctx, ws, &request)
		if readErr != nil {
			// remove subscriber from room
			room.EnqueueLeave(s)
			return
		}

		// message is empty
		if len(request.ChatMessage) == 0 {
			continue
		}
		
		userInfo := model.CreateUserInfo(request.UserID, req)
		newMessage := model.CreateMessage(
			model.CreateUserRequest(request.))
		)

		
		chatroomMessage := chatroom_model.NewChatroomMessage(newMessage)
		
	}

}