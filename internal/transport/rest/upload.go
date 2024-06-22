package rest

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /

// Upload godoc
// @Summary upload a FILE
// @Schemes application/json
// @Description accepts file sent by the user as input and upload it
// @Tags uploading
// @Accept json
// @Produce json
// @Success 200 "message: Uploade was successful"
// @Router /api/auth/upload [post]
func Upload(c *gin.Context) {
	// single file
	file, _ := c.FormFile("file")
	log.Println(file.Filename)

	c.Header("Access-Control-Allow-Origin", "*")

	log.Println(file.Size, file)

	// Upload the file to specific dst.
	c.SaveUploadedFile(file, "./videos/"+file.Filename)

	c.JSON(http.StatusOK, gin.H{
		"message": "Uploade successful",
	})
}
