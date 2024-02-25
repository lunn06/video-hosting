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
		c.JSON(http.StatusBadRequest, gin.H{
			// неправильно набран email
			"error": "Email entered incorrectly, because it exceeds the character limit or backwards",
		})
		return
	}
	if len(body.Password) > 72 || body.Password == "" {
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
	c.SetCookie(
		// cookie name
		"Authorization",
		// cookie value
		t,
		// cookie lifetime
		3600*24*30,
		// path on the server to which the cookies are applied, path applies to the current one
		"",
		// domain for which cookies are valid are applied to the current one
		"",
		// cookies are sent via http
		false,
		// accessible only to server requests and not accessible to JS reading and modification
		true)
	c.JSON(http.StatusOK, gin.H{
		"message": "Authorization was successful",
	})
}
