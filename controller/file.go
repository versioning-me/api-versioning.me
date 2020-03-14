package controller

import (
	"api-versioning-me/db"
	"api-versioning-me/models"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)


func GetFilesHandler(c *gin.Context) {
	var f models.File
	if err := models.GetFiles(20,&f,db.Db).Error; err != nil {
		log.Fatalf("Failed to get version by name. %+v", err)
	}
	fmt.Println(f.Files)
	c.JSON(http.StatusOK, gin.H{
		"files": f.Files,
	})
}

// POST "/new"
func NewVersionAndNewFileHandler(c *gin.Context) {
	file, fileHeader, _ := c.Request.FormFile("file")

	v,err := models.NewVersion(file, fileHeader)
	if err != nil {
		log.Fatalf("Failed to New File obj. : %s", err)
	}

	// デフォルトでtitleに現在のfilenameを入れておく
	// Todo
	// POST時にtitleとdetailを入れる？一旦空欄stringで
	f := models.NewFile()
	f.Title = fileHeader.Filename
	if c.PostForm("title") != "" {
		f.Title = c.PostForm("title")
	}
	f.Detail = c.PostForm("detail")

	objAttrs, err := models.StoreGCS("uploadfile-versioning-me-dev", f, v)
	if err != nil {
		log.Printf("Failed to upload file. ERROR: %s", err)
		return
	}
	v.Url = fmt.Sprintf("https://storage.googleapis.com/%s/%s", objAttrs.Bucket, objAttrs.Name)

	if err := TxInsertVersionAndInsertFile(v, f, db.Db); err != nil {
		log.Fatalf("Failed to commit. rollback... %s", err.Error())
	}

	//c.Redirect(http.StatusFound, "/file")
	// c.JSON(http.StatusOK, gin.H{
	// 	"msg": fmt.Sprintf("%+v uploaded!", *f),
	// 	"url":      f.Url,
	// 	"fileName": objAttrs.Name,
	// })
	GetFilesHandler(c)
}