package controller

import (
	"api-versioning-me/db"
	"api-versioning-me/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// GET "/versions"
func GetVersionsHandler(c *gin.Context) {
	var v models.Version
	if err := models.GetVersions(20,&v,db.Db).Error; err != nil {
		log.Fatalf("Failed to get version by name. %+v", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"versions": v.Versions,
	})
}
