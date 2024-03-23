package rest

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
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
// @Failure 400 "error: Failed to read body"
// @Failure 422 "error: Email entered incorrectly, because it exceeds the character limit or backwards"
// @Failure 422 "error: Invalid password size"
// @Failure 403 "error: Invalid email or password"
// @Failure 500 "error: Invalid to create token"
// @Failure 400 "error: Invalid to insert token"
// @Router /api/auth/login [post]
func RefreshTokens(c *gin.Context) {
	refreshCookie, err := c.Request.Cookie("refreshToken")
	if err != nil {
		slog.Error(fmt.Sprintf("RefreshToken() error = %v, can't fetch refresh cookies", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "InternalServerError, please try again later",
		})
		return
	}
	refreshToken, _ := url.QueryUnescape(refreshCookie.Value)

	token, err := database.PopToken(refreshToken)
	if err != nil {
		slog.Error(fmt.Sprintf("RefreshToken() error = %v, can't delete or select from db", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "InternalServerError, please try again later",
		})
		return
	}

	if token.CreationTime.Add(time.Duration(refreshCookie.MaxAge)).Compare(time.Now()) == 1 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "INVALID_REFRESH_SESSION: refresh token out of life",
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

	err = database.InsertToken(user.Id, refreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid to insert token",
		})
		return
	}

	jwtCookie := http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		MaxAge:   refreshLife,
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
		"message":      "Authentication was successful",
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}
