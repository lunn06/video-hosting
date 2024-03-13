package rest

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lunn06/video-hosting/internal/config"
	"github.com/lunn06/video-hosting/internal/database"
	"github.com/lunn06/video-hosting/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// @BasePath /

// Authorization godoc
// @Summary авторизирует пользователя
// @Schemes application/json
// @Description accepts json sent by the user as input and authorize it
// @Tags authorization
// @Accept json
// @Produce json
// @Param input body models.LoginRequest true "account info"
// @Success 200 "message: Authorization was successful"
// @Failure 400
// @Router /login [post]
func Authorization(c *gin.Context) {
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
	user := &models.User{}
	user, err := database.GetUser(body.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}
	var jwtSecretKey = []byte(config.CFG.JWTSecretKey)

	payload := jwt.MapClaims{
		"sub": user.Email,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString(jwtSecretKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid to create token",
		})
		return
	}
	err = database.InsertToken(user.Id, t)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid to insert token",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"access_token": t,
	})
}
