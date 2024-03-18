package chat_model

type MessageRequest struct {
	JWT         string `json:"user-id"`
	UserMessage string `json:"chat-message"`
	RoomID      string `json:"room-id"`
}

// for web socket connection
type ConnectionRequest struct {
	JWT    string
	RoomID string `param:"roomID"`
}
