package routes

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/adaggerboy/utasksd/pkg/database"
	ginhelpers "github.com/adaggerboy/utasksd/pkg/utils/ginHelpers"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var dataDir = "./data"

const httpRelativePath = "/files/"

func generateUniqueFileName(originalName string) string {
	extension := filepath.Ext(originalName)
	uniqueName := uuid.New().String() + extension
	return uniqueName
}

func DeployAttachmentsRoutes(r gin.IRouter, dataDirSet string) {
	dataDir = dataDirSet
	r.GET(":filename", func(c *gin.Context) {
		filename := c.Param("filename")
		filePath := filepath.Join(dataDir, filename)
		_, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			c.String(http.StatusNotFound, "Not found")
			return
		}
		c.File(filePath)
	})

	r.POST("upload", func(ctx *gin.Context) {

		_, ok := ctx.Get("user_name")
		if !ok {
			ctx.String(http.StatusUnauthorized, "Unauthorized")
			return
		}

		file, header, err := ctx.Request.FormFile("file")
		if err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}

		uniqueName := generateUniqueFileName(header.Filename)
		out, err := os.Create(filepath.Join(dataDir, uniqueName))
		if err != nil {
			ginhelpers.LogError(ctx.ClientIP(), "Creating attachment in filesystem", err)
			ctx.String(http.StatusInternalServerError, "Internal server error")
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			ginhelpers.LogError(ctx.ClientIP(), "Writing attachment to filesystem", err)
			ctx.String(http.StatusInternalServerError, "Internal server error")
			return
		}

		db, err := database.GetDatabase(ctx)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Internal server error")
			return
		}
		if err = db.RegisterAttachment(ctx, httpRelativePath+uniqueName); err != nil {
			db.RollbackClose()
			ginhelpers.LogError(ctx.ClientIP(), "Registering attachment in database", err)
			os.Remove(filepath.Join(dataDir, uniqueName))
			ctx.String(http.StatusInternalServerError, "Internal server error")
			return
		}
		if err = db.Close(); err != nil {
			ctx.String(http.StatusInternalServerError, "Internal server error")
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "path": httpRelativePath + uniqueName})
	})

}
