package chat_model


type MessageRequest struct {
	RoomID      string `json:"room-id"`
	UserID      string `json:"user-id"`
	UserMessage string `json:"chat-message"`
}