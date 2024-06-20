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

	c.Header("Access-Control-Allow-Origin", "*")

	log.Println(file.Size, file)

	// Upload the file to specific dst.
	c.SaveUploadedFile(file, "./videos/" + file.Filename)

	c.JSON(http.StatusOK, gin.H{
		"message": "Uploade successful",
	})
}
