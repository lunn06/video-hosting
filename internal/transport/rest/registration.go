package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
// @Failure 400 "error: Failed to read body"
// @Failure 422 "error: Failed create email, because it exceeds the character limit or backwards"
// @Failure 422 "error: Failed create channel_name, because it exceeds the character limit or backwards"
// @Failure 422 "error: Failed create password, because it exceeds the character limit or backwards"
// @Failure 500 "error: Failed to hash password. Please, try again later"
// @Failure 409 "error: email or channel already been use"
// @Router /api/auth/registration [post]
func Registration(c *gin.Context) {
	body := models.RegisterRequest{}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	if len(body.Email) > 255 || body.Email == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Failed create email, because it exceeds the character limit or backwards",
		})
		return
	}
	if len(body.ChannelName) > 255 || body.ChannelName == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Failed create channel_name, because it exceeds the character limit or backwards",
		})
		return
	}
	if len(body.Password) > 72 || body.Password == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Failed create password, because it exceeds the character limit or backwards",
		})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error on the server. Please, try again later",
		})
		return
	}
	user := models.User{
		Email:       body.Email,
		ChannelName: body.ChannelName,
		Password:    string(hash),
	}
	err = database.InsertUser(user)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "email or channel already been use",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Registration was successful",
	})
}
