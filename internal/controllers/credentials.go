package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/adaggerboy/utasksd/config"
	"github.com/adaggerboy/utasksd/models/context"
	"github.com/adaggerboy/utasksd/models/generic"
	"github.com/adaggerboy/utasksd/models/reqresp"
	"github.com/adaggerboy/utasksd/pkg/database"
	"github.com/adaggerboy/utasksd/pkg/utils/jwt"
)

func RequireAuth(ctx context.IContext) (userID int, valid bool) {
	userID, ok := ctx.GetUserID()
	if !ok {
		ctx.AddPublicError(http.StatusUnauthorized, "Unauthorized")
	}
	return userID, ok

}

func RequireStrictUser(ctx context.IContext, userID int) bool {
	uid, ok := ctx.GetUserID()
	if !ok {
		ctx.AddPublicError(http.StatusUnauthorized, "Unauthorized")
		return false
	} else if uid != userID {
		ctx.AddPublicError(http.StatusForbidden, "Access denied")
		return false
	} else {
		return true
	}
}

//Credentials API controllers

func UpdateCredentials(ctx context.IContext, userID int, req reqresp.UpdatePasswordRequest) {
	user, err := ctx.GetDatabase().ReadUser(ctx, userID)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("read database UpdateCredentials(): %s", err))
		return
	} else if user == nil {
		ctx.AddPublicError(http.StatusNotFound, "User not found")
		return
	}
	err = ctx.GetDatabase().ChangeDBUserPassword(ctx, *user.Username, req.Secret)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("write database UpdateCredentials(): %s", err))
		return
	}
}

func SetUserPermissions(ctx context.IContext, userID int, user generic.User) {
	if user.IsActive == nil || user.IsAdmin == nil || user.IsDirector == nil {
		ctx.AddPublicError(http.StatusBadRequest, "All permissions should be provided")
	}
	user2, err := ctx.GetDatabase().ReadUser(ctx, userID)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("read database UpdateCredentials(): %s", err))
		return
	} else if user2 == nil {
		ctx.AddPublicError(http.StatusNotFound, "User not found")
		return
	}

	err = ctx.GetDatabase().SetUserPermissions(ctx, userID, *user2.Username, *user.IsActive, *user.IsAdmin, *user.IsDirector)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("write database UpdateCredentials(): %s", err))
		return
	}
}

func UpdateMyCredentials(ctx context.IContext, req reqresp.UpdatePasswordRequest) {
	username, ok := ctx.GetUserName()
	if !ok {
		ctx.AddPublicError(http.StatusUnauthorized, "Unauthorized")
	}
	ctx.GetDatabase().ChangeDBUserPassword(ctx, username, req.Secret)
}

func Login(ctx context.IContext, username, clientPassword string) (token string) {
	db, err := database.ConnectDatabase(ctx, username, clientPassword)
	if err != nil {
		ctx.AddPublicError(http.StatusForbidden, "Incorrect password")
		return
	}
	userID, isActive, err := db.GetUserByUsername(ctx, username)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("read database Login(): %s", err))
		return
	} else if userID == nil {
		ctx.AddPublicError(http.StatusNotFound, "User not found")
		return
	} else if !isActive {
		ctx.AddPublicError(http.StatusGone, "User is not active")
		return
	}
	db.Close()

	token, err = jwt.SignToken(context.Session{
		UserID:            *userID,
		Username:          username,
		EncryptedPassword: clientPassword,
	}, username, time.Now().Add(time.Duration(config.GlobalConfig.HTTPServer.TokenDuration)*time.Second))
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("signing jwt token Login(): %s", err))
		return
	}
	return token
}

func CreateUser(ctx context.IContext, db database.IDatabase, user generic.User) (userID int) {
	if user.AvatarPath == nil {
		user.AvatarPath = new(string)
		*user.AvatarPath = "null"
	}
	userID = 0
	id, err := db.CreateUser(ctx, user)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("write database CreateUser(): %s", err))
		return
	} else if id == nil {
		ctx.AddPublicError(http.StatusConflict, "User can't be created, possibly already exists")
		return
	}
	return *id
}

func CreateUserWithPassword(ctx context.IContext, user reqresp.CreateUserWithPasswordRequest) (userID int) {
	db, err := database.GetDatabase(ctx)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("get database CreateUserWithPassword(): %s", err))
		return
	}
	err = db.CreateDBUser(ctx, *user.User.Username, user.Secret)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("write database CreateUserWithPassword(): %s", err))
		ctx.AddPublicError(http.StatusConflict, "User exists")
		db.RollbackClose()
		return
	}
	userID = CreateUser(ctx, db, user.User)
	if !ctx.IsActive() {
		db.DeleteDBUser(ctx, *user.User.Username)
		db.RollbackClose()
		return
	}
	db.Close()
	return
}
