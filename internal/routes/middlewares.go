package routes

import (
	"log"
	"net/http"
	"strings"

	ginhelpers "github.com/adaggerboy/utasksd/pkg/utils/ginHelpers"
	"github.com/adaggerboy/utasksd/pkg/utils/jwt"
	"github.com/gin-gonic/gin"
)

func extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		authValues := strings.Fields(authHeader)
		if len(authValues) == 2 && (strings.ToLower(authValues[0]) == "bearer" || strings.ToLower(authValues[0]) == "jwt") {
			return authValues[1]
		}
	}

	cookie, err := c.Cookie("access_token")
	if err == nil {
		return cookie
	}

	accessToken := c.Query("access_token")
	if accessToken != "" {
		return accessToken
	}

	return ""
}

func GetAuthorizationMiddleware(origin string, secure bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractToken(c)
		if len(tokenString) == 0 {
			return
		}
		token, valid, err := jwt.VerifyToken(tokenString)
		if err != nil || !valid {
			log.Printf("[ERROR] from %s :: Failed to verify jwt token: %s", c.ClientIP(), err)
			ginhelpers.ResetToken(c, origin, secure)
			return
		}

		c.Set("user_id", token.UserID)
		c.Set("user_name", token.Username)
		c.Set("user_pass", token.EncryptedPassword)
	}

}

func GetCORSMiddleware(origin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization, Accept, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
	}
}
