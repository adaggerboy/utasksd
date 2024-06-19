package ginhelpers

import (
	"context"
	"fmt"
	"sync"
	"time"

	ctxModel "github.com/adaggerboy/utasksd/models/context"
	"github.com/adaggerboy/utasksd/pkg/database"
	"github.com/gin-gonic/gin"
)

type GinWrapper struct {
	ctxModel.IContext
	goContext context.Context

	internalMtx      sync.RWMutex
	publicErrors     []string
	privateErrors    []error
	actualStatusCode int

	verifiedUserID int
	authorized     bool

	dbUsername, dbPass string

	db *database.IDatabase

	active bool
}

func NewGinWrapperContext(ctx *gin.Context) (out *GinWrapper, err error) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		userID = 0
	}
	newUserID, isValidType := userID.(int)
	if !isValidType {
		return nil, fmt.Errorf("user_id variable has incorrect type")
	}

	name, exists := ctx.Get("user_name")
	if !exists {
		name = ""
	}
	newDbName, isValidType := name.(string)
	if !isValidType {
		return nil, fmt.Errorf("user_name variable has incorrect type")
	}

	pass, exists := ctx.Get("user_pass")
	if !exists {
		pass = ""
	}
	newDbPass, isValidType := pass.(string)
	if !isValidType {
		return nil, fmt.Errorf("user_pass variable has incorrect type")
	}

	return &GinWrapper{
		goContext:        ctx,
		publicErrors:     []string{},
		privateErrors:    []error{},
		actualStatusCode: 200,
		verifiedUserID:   newUserID,
		authorized:       exists,
		active:           true,
		dbUsername:       newDbName,
		dbPass:           newDbPass,
		db:               nil,
	}, nil
}

//Wrapper interface

func (ctx *GinWrapper) GetUserID() (userID int, exists bool) {
	ctx.internalMtx.RLock()
	defer ctx.internalMtx.RUnlock()
	return ctx.verifiedUserID, ctx.authorized
}

func (ctx *GinWrapper) AddPublicError(code int, message string) {
	ctx.internalMtx.Lock()
	defer ctx.internalMtx.Unlock()
	ctx.publicErrors = append(ctx.publicErrors, message)
	ctx.actualStatusCode = code
	ctx.active = false
}
func (ctx *GinWrapper) AddPrivateError(code int, err error) {
	ctx.internalMtx.Lock()
	defer ctx.internalMtx.Unlock()
	ctx.privateErrors = append(ctx.privateErrors, err)
	ctx.actualStatusCode = code
	ctx.active = false
}

func (ctx *GinWrapper) SetStatusCode(code int) {
	ctx.internalMtx.Lock()
	defer ctx.internalMtx.Unlock()
	ctx.actualStatusCode = code
	if code >= 400 {
		ctx.active = false
	}
}

func (ctx *GinWrapper) IsActive() bool {
	ctx.internalMtx.RLock()
	defer ctx.internalMtx.RUnlock()
	return ctx.active
}

//Go context

func (ctx *GinWrapper) Deadline() (deadline time.Time, ok bool) {
	return ctx.goContext.Deadline()
}

func (ctx *GinWrapper) Done() <-chan struct{} {
	return ctx.goContext.Done()
}

func (ctx *GinWrapper) Err() error {
	return ctx.goContext.Err()
}

func (ctx *GinWrapper) Value(key any) any {
	return ctx.goContext.Value(key)
}

//Gin accessible

func (ctx *GinWrapper) GetPrivateErrors() []error {
	ctx.internalMtx.RLock()
	defer ctx.internalMtx.RUnlock()
	return ctx.privateErrors
}

func (ctx *GinWrapper) GetPublicErrors() []string {
	ctx.internalMtx.RLock()
	defer ctx.internalMtx.RUnlock()
	return ctx.publicErrors
}

func (ctx *GinWrapper) GetStatusCode() (code int) {
	ctx.internalMtx.RLock()
	defer ctx.internalMtx.RUnlock()
	return ctx.actualStatusCode
}

func (ctx *GinWrapper) GetDatabase() database.IDatabase {
	if ctx.db != nil {
		return *ctx.db
	} else {

		if len(ctx.dbUsername) > 0 {
			db, err := database.ConnectDatabase(ctx.goContext, ctx.dbUsername, ctx.dbPass)
			if err != nil {
				ctx.AddPrivateError(500, err)
				ctx.active = false
			}
			ctx.db = &db
			return db
		} else {
			ctx.AddPublicError(401, "Need authorization")
			ctx.active = false
			return nil
		}

	}
}

func (ctx *GinWrapper) GetUserName() (username string, exists bool) {
	return ctx.dbUsername, len(ctx.dbUsername) > 0
}
