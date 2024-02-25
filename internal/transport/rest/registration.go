package rest

import (
	"net/http"
	"time"

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
			"error": "Failed to hash password. Please, try again later",
		})
		return
	}
	tx := database.DB.MustBegin()
	var lastUserID uint32
	err = tx.Get(&lastUserID, "SELECT id FROM users ORDER BY id DESC LIMIT 1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error on the server. Please, try again later",
		})
		return
	}
	newUserID := lastUserID + 1
	usr := models.User{
		Id:               newUserID,
		Email:            body.Email,
		ChannelName:      body.ChannelName,
		Password:         string(hash),
		RegistrationTime: time.Now(),
	}
	defer tx.Commit()
	result, err := tx.Exec("INSERT INTO users (id, email, channel_name, password, registration_time) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING", usr.Id, usr.Email, usr.ChannelName, usr.Password, usr.RegistrationTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error on the server. Please, try again later",
		})
		return
	}
	if rowsAffected, err := result.RowsAffected(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error on the server. Please, try again later",
		})
		return
	} else {
		if rowsAffected == 0 {
			c.JSON(http.StatusConflict, gin.H{
				"error": "This email or channel has already been use",
			})
			return
		}
	}
	result, err = tx.Exec("INSERT INTO users_roles (user_id, role_id) VALUES ($1, $2)", usr.Id, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error on the server. Please, try again later",
		})
		tx.MustExec("DELETE FROM users WHERE id=$1", usr.Id)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Registration was successful",
	})
}
