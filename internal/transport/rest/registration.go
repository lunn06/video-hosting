package rest

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lunn06/video-hosting/internal/database"
	"github.com/lunn06/video-hosting/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func Registration(c *gin.Context) {
	body := models.RegisterRequest{}
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
	if len(body.ChannelName) > 255 && len(body.ChannelName) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Filed create channel_name, because it exceeds the character limit or backwards",
		})
	}
	if len(body.Password) > 255 && len(body.Password) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Filed create password, because it exceeds the character limit or backwards",
		})
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
	}

	id := uuid.New().String()
	user := models.User{
		Id:               id,
		Email:            body.Email,
		ChannelName:      body.ChannelName,
		Password:         string(hash),
		RegistrationTime: time.Now(),
	}
	tx := database.DB.MustBegin()
	_, err = tx.Exec("INSERT INTO users (id, email, channel_name, password, registration_time) VALUES ($1, $2, $3, $4, $5)", user.Id, user.Email, user.ChannelName, user.Password, user.RegistrationTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"fail": "This email is already in use",
		})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"message": "Congratulations",
	})
}
