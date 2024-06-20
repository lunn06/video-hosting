package rest

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {
	// single file
	file, _ := c.FormFile("file")
	log.Println(file.Filename)

	// Upload the file to specific dst.
	c.SaveUploadedFile(file, "./videos")

	c.JSON(http.StatusOK, gin.H{
		"message": "Uploade successful",
	})
}
