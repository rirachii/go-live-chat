package chatroom

import (
	"log"

	"github.com/labstack/echo/v4"
	chatroom_model "github.com/rirachii/golivechat/app/internal/chatroom/model"
	"github.com/rirachii/golivechat/app/shared/model"
	"nhooyr.io/websocket"
)


func (room *chatroom) StartChatroom() {

}
func (room *chatroom) CloseChatroom() {
	room.closeChannel <- 1
}

// called by the hub to tell chatroom to accept a new connection
func (room *chatroom) AcceptConnection(c echo.Context, uid model.UserID, rid model.RoomID) error {

	// get username from db, inject into chatroom user
	// username := create


	w := c.Response().Writer
	r := c.Request()

	userConn, connErr := websocket.Accept(w, r, nil)
	if connErr != nil {
		log.Printf("%v", connErr)
		return ErrWebsocketConnectionFailed
	}

	user := chatroom_model.NewChatroomUser(
		userConn, uid, rid, "TODO",
	)
	room.EnqueueJoin(user)

	return nil
	
}

