package rest

import "github.com/gin-gonic/gin"

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept string
// @Produce stirng
// @Success 200 {string} pong
// @Router /ping [get]
func Ping(c *gin.Context) {
	c.String(200, "pong")
}
