package rest

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lunn06/video-hosting/internal/database"
)

// @BasePath /auth/api/

// RefreshTokens godoc
// @Summary refresh user's tokens
// @Schemes application/json
// @Description accept json and refresh user refresh and access tokens
// @Tags authorization
// @Accept json
// @Produce json
// @Param input body models.LoginRequest true "account info"
// @Success 200 "message: RefreshTokens was successful"
// @Failure 401 "error: Invalid to get refresh token from cookie"
// @Failure 500 "error: Invalid to pop token"
// @Failure 500 "error: Invalid to insert token"
// @Failure 500 "error: Invalid to create token"
// @Router /api/auth/refresh [post]
func RefreshTokens(c *gin.Context) {
	refreshUUID, err := c.Cookie("refreshToken")
	if err != nil {
		slog.Error(fmt.Sprintf("RefreshToken() error = %v, can't fetch refresh cookies", err))
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid to get refresh token from cookie",
		})
		return
	}

	token, err := database.PopToken(refreshUUID)
	if err != nil {
		slog.Error(fmt.Sprintf("RefreshToken() error = %v, can't delete or select from db", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "InternalServerError, please try again later",
		})
		return
	}

	if token.CreationTime.Add(time.Duration(RefreshLife)).Compare(time.Now()) < 1 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "INVALID_REFRESH_SESSION: refresh token out of life",
		})
		return
	}

	user, err := database.GetUserByRefreshToken(token.Token)

	accessToken, refreshToken, err := newTokens(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid to create token",
		})
		return
	}

	newRefreshUUID, err := database.InsertToken(user.Id, refreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid to insert token",
		})
		return
	}

	jwtCookie := http.Cookie{
		Name:     "refreshToken",
		Value:    newRefreshUUID,
		MaxAge:   RefreshLife,
		Path:     "/api/auth",
		HttpOnly: true,
	}

	c.SetCookie(
		jwtCookie.Name,
		jwtCookie.Value,
		jwtCookie.MaxAge,
		jwtCookie.Path,
		jwtCookie.Domain,
		jwtCookie.Secure,
		jwtCookie.HttpOnly,
	)

	c.JSON(http.StatusOK, gin.H{
		"message":      "RefreshToken was successful",
		"accessToken":  accessToken,
		"refreshToken": newRefreshUUID,
	})
}
