package controller

import (
	"api-versioning-me/db"
	"api-versioning-me/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)


// GET "/file"
func GetFilesHandler(c *gin.Context) {
	var f models.File
	if err := models.GetFiles(20,&f,db.Db).Error; err != nil {
		log.Fatalf("Failed to get version by name. %+v", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"files": f.Files,
	})
}

