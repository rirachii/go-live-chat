package model

// Text Messages only for now
type Message struct {
	roomId RoomID
	userId UserID
	// change to payload.
	content  string // content of message
	metadata MessageMetadata
}

// TODO
type MessageMetadata struct {
}

func CreateMessage(request UserRequest, textMessage string, metadata MessageMetadata) Message {
	message := Message{
		roomId:   request.RoomId,
		userId:   request.UserId,
		content:  textMessage,
		metadata: metadata,
	}
	return message
}

func (m Message) RoomID() RoomID            { return m.roomId }
func (m Message) SenderId() UserID          { return m.userId }
// func (m Message) SenderInfo() UserInfo      { return m.userInfo }
// func (m Message) SenderName() string        { return m.userInfo.Name }
func (m Message) Content() string           { return m.content }
func (m Message) Metadata() MessageMetadata { return m.metadata }
