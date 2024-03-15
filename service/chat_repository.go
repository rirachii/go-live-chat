package service

// import (
// 	"context"
// 	"fmt"

// 	model "github.com/rirachii/golivechat/model"
// 	chat "github.com/rirachii/golivechat/model/chat"
// )

// type ChatRepository interface {
// 	CreateRoom(ctx context.Context, user *model.User) (*model.User, error)
// 	LogMessage(ctx context.Context, email string) (*model.User, error)
// 	// AddAdmin(ctx context.Context, email string) (*model.User, error)
// }

// type chatRepo struct {
// 	db
// }

// func NewChatRepository(db DBTX) ChatRepository {
// 	return &chatRepo{db: db}
// }
