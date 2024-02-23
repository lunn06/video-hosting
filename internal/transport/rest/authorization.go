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

func Authorization(c *gin.Context) {
	body := models.LoginRequest{}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	if len(body.Email) > 255 && len(body.Email) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Filed create email, because it exceeds the character limit or backwards",
		})
	}
	if len(body.Password) > 255 && len(body.Password) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid password size",
		})
		return
	}
	user := models.User{}
	if err := database.DB.Get(&user, "SELECT * FROM users WHERE email=$1", body.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}
	var jwtSecretKey = []byte(config.CFG.JWTSecretKey)
	payload := jwt.MapClaims{
		"sub": user.Email,
		"exp": time.Now().Add(time.Hour * 30).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString(jwtSecretKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid to create token",
		})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization",
		t,
		3600*24*30,
		"",
		"",
		false,
		true)
	c.JSON(http.StatusOK, gin.H{
		"message": "Authorization was successful",
	})
}
