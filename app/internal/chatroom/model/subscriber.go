package chatroom_model

import (
	model "github.com/rirachii/golivechat/app/shared/model"
)

// Users become Subscribers when they join the chatroom
type Subscriber struct {
	user         *ChatroomUser
	userInfo     model.UserInfo
	rid          model.RoomID
	role         string
	messageQueue chan model.Message
}

func NewSubscriber(
	user *ChatroomUser,
	role string,
) *Subscriber {
	subscriber := Subscriber{
		user:         user,
		userInfo:     user.Info(),
		rid:          user.roomId,
		role:         role,
		messageQueue: make(chan model.Message),
	}

	return &subscriber
}
func (subscriber *Subscriber) User() *ChatroomUser                 { return subscriber.user }
func (subscriber Subscriber) UserInfo() model.UserInfo             { return subscriber.userInfo }
func (subscriber Subscriber) Id() model.UserID                     { return subscriber.userInfo.Id }
func (subscriber Subscriber) Name() string                         { return subscriber.userInfo.Name }
func (subscriber Subscriber) RoomID() model.RoomID                 { return subscriber.rid }
func (subscriber Subscriber) Role() string                         { return subscriber.role }

func (subscriber *Subscriber) PendingMessages() chan model.Message { return subscriber.messageQueue }
