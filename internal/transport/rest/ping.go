package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// Ping godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce html
// @Success 200 html pong
// @Router /ping [get]
func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
