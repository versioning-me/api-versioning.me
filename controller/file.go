package controller

import (
	"api-versioning-me/db"
	"api-versioning-me/models"
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// POST "/file"
func FileUploadHandler(c *gin.Context) {
	file, fileHeader, _ := c.Request.FormFile("file")

	f, err := models.NewFile(file, fileHeader)
	if err != nil {
		log.Fatalf("Failed to New File obj. : %s", err)
	}
	v := models.NewVersion(f)

	objAttrs, err := models.StoreGCS("uploadfile-versioning-me-dev", f)
	if err != nil {
		log.Printf("Failed to upload file. ERROR: %s", err)
		return
	}
	f.Url = fmt.Sprintf("https://storage.googleapis.com/%s/%s", objAttrs.Bucket, objAttrs.Name)

	if err := TxCreateVersionAndFile(v, f, db.Db); err != nil {
		log.Fatalf("Failed to commit. rollback... %s", err.Error())
	}

	c.Redirect(http.StatusFound, "/files")
	// c.JSON(http.StatusOK, gin.H{
	// 	"msg": fmt.Sprintf("%+v uploaded!", *f),
	// 	"url":      f.Url,
	// 	"fileName": objAttrs.Name,
	// })
}

func TxCreateVersionAndFile(v *models.Version, f *models.UploadedFile, db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// データベース操作をトランザクション内で行います（ここからは'db'ではなく'tx'を使います）
		// エラーを返した場合はロールバックされます
		if err := models.StoreUploadedFile(f, tx).Error;err != nil {
			return err
		}
		var version models.Version
		if db.Where("name = ?", v.Name).Last(&version).RecordNotFound() == true {
			log.Println(version.Name)
			if err := models.StoreVersion(v, tx).Error;err != nil {
				return err
			}
		}
		// nilを返すとコミットされます
		return nil
	})
}


// GET "/files"
func GetUploadedFilesHandler(c *gin.Context) {
	var f models.UploadedFile
	if err := models.GetUploadedFiles(20,&f,db.Db).Error; err != nil {
		log.Fatalf("Failed to get uploaded files. %+v", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"files": f.UploadedFiles,
	})
}

// GET "/files/:version_name"
func GetUploadedFilesByVersionNameHandler(c *gin.Context) {
	q := c.Param("version_name")
	var f models.UploadedFile
	if err := models.GetUploadedFilesByVersionName(20,q,&f,db.Db).Error; err != nil {
		log.Fatalf("Failed to get uploaded files by version name. %+v", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"files": f.UploadedFiles,
	})
}

