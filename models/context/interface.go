package context

import (
	"context"

	"github.com/adaggerboy/utasksd/pkg/database"
)

type IContext interface {
	context.Context

	GetUserID() (userID int, exists bool)
	GetUserName() (username string, exists bool)
	AddPublicError(code int, message string)
	AddPrivateError(code int, err error)

	IsActive() bool

	SetStatusCode(code int)

	GetDatabase() database.IDatabase
}

type Session struct {
	Username          string `json:"username"`
	UserID            int    `json:"id"`
	EncryptedPassword string `json:"encr"`
}
