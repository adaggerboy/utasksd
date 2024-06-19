package ginhelpers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ResetToken(c *gin.Context, domain string, secure bool) {
	c.SetCookie("access_token", "", 0, "/", domain, secure, true)
}

func SetToken(c *gin.Context, domain string, secure bool, duration int, token string) {
	c.SetCookie("access_token", token, duration, "/", domain, secure, true)
}

func GetID(ctx *gin.Context) (idPointer *int) {
	stringId := ctx.Param("id")
	id, err := strconv.Atoi(stringId)
	if err != nil {
		Error(ctx, http.StatusBadRequest, nil)
		return nil
	}
	return &id
}

func Error(ctx *gin.Context, code int, message []string) {
	if len(message) == 0 {
		ctx.String(code, http.StatusText(code)+" ")
	} else if len(message) == 1 {
		ctx.String(code, message[0]+" ")
	} else {
		ctx.JSON(code, message)
	}
}

func LogError(clientIP string, message string, err error) {
	log.Printf("[ERROR] from %s :: %s: %s", clientIP, message, err)
}

func WebError(ctx *gin.Context, code int, message []string) {
	if len(message) == 0 {
		ctx.String(code, http.StatusText(code)+" ")
	} else if len(message) == 1 {
		ctx.String(code, message[0]+" ")
	} else {
		ctx.JSON(code, message)
	}
}

func PrintErrors(ctx *gin.Context, rctx *GinWrapper) {
	Error(ctx, rctx.GetStatusCode(), rctx.GetPublicErrors())
	for _, v := range rctx.GetPrivateErrors() {
		LogError(ctx.ClientIP(), "Request error", v)
	}
}

func WebPrintErrors(ctx *gin.Context, rctx *GinWrapper) {
	if rctx.GetStatusCode() == http.StatusUnauthorized {
		ctx.Header("Location", "/web/login")
		ctx.Status(http.StatusSeeOther)
		return
	}
	WebError(ctx, rctx.GetStatusCode(), rctx.GetPublicErrors())
	for _, v := range rctx.GetPrivateErrors() {
		LogError(ctx.ClientIP(), "Request error", v)
	}
}

// Should be limited

func ReadJSON[V any](ctx *gin.Context) (data *V, err error) {
	data = nil
	bytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return
	}
	data = new(V)
	err = json.Unmarshal(bytes, data)
	if err != nil {
		return
	}
	return
}

func createWrapper(ctx *gin.Context, web bool) (nctx *GinWrapper, err error) {
	nctx, err = NewGinWrapperContext(ctx)
	if err != nil {
		LogError(ctx.ClientIP(), "Failed to create context", err)
		return
	}
	nctx.GetDatabase()
	if !nctx.IsActive() {
		if web {
			WebPrintErrors(ctx, nctx)
		} else {
			PrintErrors(ctx, nctx)
		}
		err = fmt.Errorf("can't open database")
	}
	return
}

func closeWrapper(ctx *gin.Context, nctx *GinWrapper, web bool) bool {
	if !nctx.IsActive() {
		if err := nctx.GetDatabase().RollbackClose(); err != nil {
			nctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("failed to rollback and close db connection: %s", err))
		}
		if web {
			WebPrintErrors(ctx, nctx)
		} else {
			PrintErrors(ctx, nctx)
		}
		return true
	} else if err := nctx.GetDatabase().Close(); err != nil {
		nctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("failed to commit and close db connection: %s", err))
		if web {
			WebPrintErrors(ctx, nctx)
		} else {
			PrintErrors(ctx, nctx)
		}
		return true
	} else {
		return false
	}
}

func WrapContextDataPermitted[V any](wrapped func(ctx *gin.Context, nctx *GinWrapper, data V) (code int, obj interface{})) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		data, err := ReadJSON[V](ctx)
		if err != nil || data == nil {
			Error(ctx, http.StatusBadRequest, []string{fmt.Sprintf("Invalid input JSON: %s", err)})
			return
		}
		nctx, err := NewGinWrapperContext(ctx)
		if err != nil {
			LogError(ctx.ClientIP(), "Failed to create context", err)
			return
		}

		code, obj := wrapped(ctx, nctx, *data)

		if !nctx.IsActive() {
			PrintErrors(ctx, nctx)
		} else if obj != nil {
			ctx.JSON(code, obj)
		} else {
			ctx.Status(code)
		}
	}
}

func WrapContextWebPermitted(wrapped func(ctx *gin.Context, nctx *GinWrapper) (code int, cType string, data []byte)) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		nctx, err := NewGinWrapperContext(ctx)
		if err != nil {
			LogError(ctx.ClientIP(), "Failed to create context", err)
			return
		}

		code, cType, data := wrapped(ctx, nctx)

		if !nctx.IsActive() {
			WebPrintErrors(ctx, nctx)
		} else {
			ctx.Data(code, cType, data)
		}
	}
}

func WrapContext(wrapped func(ctx *gin.Context, nctx *GinWrapper) (code int, obj interface{})) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		nctx, err := createWrapper(ctx, false)
		if err != nil {
			return
		}

		code, obj := wrapped(ctx, nctx)

		if closeWrapper(ctx, nctx, false) {
		} else if obj != nil {
			ctx.JSON(code, obj)
		} else {
			ctx.Status(code)
		}
	}
}

func WrapContextData[V any](wrapped func(ctx *gin.Context, nctx *GinWrapper, data V) (code int, obj interface{})) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		data, err := ReadJSON[V](ctx)
		if err != nil || data == nil {
			Error(ctx, http.StatusBadRequest, []string{fmt.Sprintf("Invalid input JSON: %s", err)})
			return
		}
		nctx, err := createWrapper(ctx, false)
		if err != nil {
			return
		}

		code, obj := wrapped(ctx, nctx, *data)

		if closeWrapper(ctx, nctx, false) {
		} else if obj != nil {
			ctx.JSON(code, obj)
		} else {
			ctx.Status(code)
		}
	}
}

func WrapContextWeb(wrapped func(ctx *gin.Context, nctx *GinWrapper) (code int, cType string, data []byte)) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		nctx, err := createWrapper(ctx, true)
		if err != nil {
			return
		}

		code, cType, data := wrapped(ctx, nctx)

		if closeWrapper(ctx, nctx, true) {
		} else {
			ctx.Data(code, cType, data)
		}
	}
}

func WrapContextID(wrapped func(ctx *gin.Context, nctx *GinWrapper, id int) (code int, obj interface{})) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id := GetID(ctx)
		if id == nil {
			Error(ctx, http.StatusBadRequest, []string{"Invalid or missing ID"})
			return
		}
		nctx, err := createWrapper(ctx, false)
		if err != nil {
			return
		}

		code, obj := wrapped(ctx, nctx, *id)

		if closeWrapper(ctx, nctx, false) {
		} else if obj != nil {
			ctx.JSON(code, obj)
		} else {
			ctx.Status(code)
		}
	}
}

func WrapContextDataID[V any](wrapped func(ctx *gin.Context, nctx *GinWrapper, id int, data V) (code int, obj interface{})) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id := GetID(ctx)
		if id == nil {
			Error(ctx, http.StatusBadRequest, []string{"Invalid or missing ID"})
			return
		}
		data, err := ReadJSON[V](ctx)
		if err != nil || data == nil {
			Error(ctx, http.StatusBadRequest, []string{fmt.Sprintf("Invalid input JSON: %s", err)})
			return
		}
		nctx, err := createWrapper(ctx, false)
		if err != nil {
			return
		}

		code, obj := wrapped(ctx, nctx, *id, *data)

		if closeWrapper(ctx, nctx, false) {
		} else if obj != nil {
			ctx.JSON(code, obj)
		} else {
			ctx.Status(code)
		}
	}
}

func WrapContextWebID(wrapped func(ctx *gin.Context, nctx *GinWrapper, id int) (code int, cType string, data []byte)) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id := GetID(ctx)
		if id == nil {
			Error(ctx, http.StatusBadRequest, []string{"Invalid or missing ID"})
			return
		}
		nctx, err := createWrapper(ctx, true)
		if err != nil {
			return
		}

		code, cType, data := wrapped(ctx, nctx, *id)

		if closeWrapper(ctx, nctx, true) {
		} else {
			ctx.Data(code, cType, data)
		}
	}
}
