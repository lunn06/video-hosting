package rest

import (
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lunn06/video-hosting/internal/config"
	"github.com/lunn06/video-hosting/internal/database"
	"github.com/lunn06/video-hosting/internal/models"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessLife  = 60 * 30
	refreshLife = 3600 * 24
)

// @BasePath /auth/api/

// Authentication godoc
// @Summary authenticates the user
// @Schemes application/json
// @Description accepts json sent by the user as input and authorize it
// @Tags authorization
// @Accept json
// @Produce json
// @Param input body models.LoginRequest true "account info"
// @Success 200 "message: Authentication was successful"
// @Failure 400 "error: Failed to read body"
// @Failure 422 "error: Email entered incorrectly, because it exceeds the character limit or backwards"
// @Failure 422 "error: Invalid password size"
// @Failure 403 "error: Invalid email or password"
// @Failure 500 "error: Invalid to create token"
// @Failure 400 "error: Invalid to insert token"
// @Router /api/auth/login [post]
func Authentication(c *gin.Context) {
	body := models.LoginRequest{}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	if len(body.Email) > 255 || body.Email == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Email entered incorrectly, because it exceeds the character limit or backwards",
		})
		return
	}
	if len(body.Password) > 72 || body.Password == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Invalid password size",
		})
		return
	}
	user, err := database.GetUser(body.Email)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Invalid email or password",
		})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	userRefreshToken, err := database.GetTokenByUser(*user)
	if err == nil {
		err = database.UpdateTokenTime(*userRefreshToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Server error! Pleas try again later",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Authentication was successful",
		})
		return
	}

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

func newTokens(user models.User) (string, string, error) {
	var jwtSecretKey = []byte(config.CFG.JWTSecretKey)

	accessPayload := jwt.MapClaims{
		"email": user.Email,
		"exp":   accessLife,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessPayload)
	signedAccessToken, err := accessToken.SignedString(jwtSecretKey)
	if err != nil {
		return "", "", err
	}

	refreshPayload := jwt.MapClaims{
		"sub": rand.Int(),
		"exp": refreshLife,
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshPayload)
	signedRefreshToken, err := refreshToken.SignedString(jwtSecretKey)
	if err != nil {
		return "", "", err
	}

	return signedAccessToken, signedRefreshToken, nil
}
