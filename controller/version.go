package controller

import (
	"api-versioning-me/db"
	"api-versioning-me/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)


// version/:file_id
func GetVersionsByFileIdHandler(c *gin.Context) {
	q := c.Param("file_id")
	var v models.Version
	if err := models.GetVersionsByFileId(20,q,&v,db.Db).Error; err != nil {
		log.Fatalf("Failed to get versions by file id. %+v", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"versions": v.Versions,
	})
}