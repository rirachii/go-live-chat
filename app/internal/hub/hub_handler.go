package hub

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"log"

// 	"net/http"

// 	echo "github.com/labstack/echo/v4"

// 	chatroom_model "github.com/rirachii/golivechat/app/internal/chatroom/model"
// 	hub_model "github.com/rirachii/golivechat/app/internal/hub/model"
// 	model "github.com/rirachii/golivechat/app/shared/model"

// 	db "github.com/rirachii/golivechat/db"
// 	// hub_svc "github.com/rirachii/golivechat/internal/hub"

// 	chatroom_template "github.com/rirachii/golivechat/templates/chatroom"
// 	hub_template "github.com/rirachii/golivechat/templates/hub"
// )

// type HubHandler struct {
// 	hub hub_model.HubServer
// }
// func (h *HubHandler) Hub() hub_model.HubServer {
// 	return h.hub
// }

// func StartHubHandler() (*HubHandler, error) {

// 	// hubChatrooms, err := getChatroomsDB()
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	hubServer := NewHubServer(
// 		// TODO db for hub chatrooms
// 		[]chatroom_model.Chatroom{},
// 	)
// 	hubHandler := &HubHandler{
// 		hub: hubServer,
// 	}

// 	return hubHandler, nil
// }


