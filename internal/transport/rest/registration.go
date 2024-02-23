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

// @BasePath /

// Registration godoc
// @Summary registers a user
// @Schemes application/json
// @Description accepts json sent by the user as input and registers it
// @Tags registration
// @Accept json
// @Produce json
// @Param input body models.RegisterRequest true "account info"
// @Success 200 "message: Registration was successful"
// @Failure 400
// @Router /registration [post]
func Registration(c *gin.Context) {
	body := models.RegisterRequest{}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	if len(body.Email) > 255 || body.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Filed create email, because it exceeds the character limit or backwards",
		})
		return
	}
	if len(body.ChannelName) > 255 || body.ChannelName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Filed create channel_name, because it exceeds the character limit or backwards",
		})
		return
	}
	if len(body.Password) > 255 || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Filed create password, because it exceeds the character limit or backwards",
		})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
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
		"message": "Registration was successful",
	})
}
