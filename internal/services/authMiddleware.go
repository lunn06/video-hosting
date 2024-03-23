package services

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lunn06/video-hosting/internal/database"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("refreshToken")

		if err != nil {
			slog.Error("AuthMiddleware() error = %v", err)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": fmt.Sprintf("AuthMiddleware() error = %v", err),
			})
			return
		}

		refreshToken, _ := url.QueryUnescape(cookie.Value)
		token, err := database.GetToken(refreshToken)

		if token.CreationTime.Add(time.Duration(cookie.MaxAge)).Compare(time.Now()) == 1 {
			slog.Error("AuthMiddleware() error = %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "INVALID_REFRESH_SESSION: refresh token out of life",
			})
			return
		}

		user, err := database.GetUserByRefreshToken(token.Token)
		if user == nil {
			slog.Error("AuthMiddleware() error = %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "INVALID_REFRESH_SESSION: no user with this token",
			})
			return
		}

		c.Next()
	}
}
