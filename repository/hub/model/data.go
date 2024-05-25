package hub_repository_model

import "github.com/rirachii/golivechat/app/shared/model"


type RoomData struct {
    Name string
    Owner   model.UserID
    IsPublic bool
    IsActive bool
}
