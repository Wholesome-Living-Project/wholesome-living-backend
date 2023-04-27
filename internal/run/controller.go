package run

import "cmd/http/main.go/internal/user"

type Controller struct {
	storage     *Storage
	userStorage *user.Storage
}

func NewController(storage *Storage, userStorage *user.Storage) *Controller {
	return &Controller{
		storage:     storage,
		userStorage: userStorage,
	}
}

type addUserDataRequest struct {
	UserID        string `json:"userId" bson:"userId"`
	Notifications bool   `json:"notifications" bson:"notifications"`
}
