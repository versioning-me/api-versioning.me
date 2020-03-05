package controller

import (
	"api-versioning-me/db"
	"api-versioning-me/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)


// POST "/new"
func NewVersionAndNewFileHandler(c *gin.Context) {
	file, fileHeader, _ := c.Request.FormFile("file")

	v,err := models.NewVersion(file, fileHeader)
	if err != nil {
		log.Fatalf("Failed to New File obj. : %s", err)
	}

	// var title, detail string
	// Todo
	// POST時にtitleとdetailを入れる？一旦空欄stringで
	title := c.PostForm("title")
	detail := c.PostForm("detail")
	f := models.NewFileWithTitle(title, detail)

	objAttrs, err := models.StoreGCS("uploadfile-versioning-me-dev", f, v)
	if err != nil {
		log.Printf("Failed to upload file. ERROR: %s", err)
		return
	}
	v.Url = fmt.Sprintf("https://storage.googleapis.com/%s/%s", objAttrs.Bucket, objAttrs.Name)

	if err := TxInsertVersionAndInsertFile(v, f, db.Db); err != nil {
		log.Fatalf("Failed to commit. rollback... %s", err.Error())
	}

	c.Redirect(http.StatusFound, "/files")
	// c.JSON(http.StatusOK, gin.H{
	// 	"msg": fmt.Sprintf("%+v uploaded!", *f),
	// 	"url":      f.Url,
	// 	"fileName": objAttrs.Name,
	// })
}

// file insert  version insert
func TxInsertVersionAndInsertFile(v *models.Version, f *models.File, db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// データベース操作をトランザクション内で行います（ここからは'db'ではなく'tx'を使います）
		// エラーを返した場合はロールバックされます
		// fileIDを取得
		if err := models.StoreFile(f, tx).Error;err != nil {
			return err
		}

		// insertしたfile.idをそのままversion.file_idにする。
		v.FileId = f.Id

		// version.versionの最新を+1する。
		models.GenerateVersionId(v,tx)

		if err := models.StoreNewVersion(v,tx).Error;err != nil {
			return err
		}
		// nilを返すとコミットされます
		return nil
	})
}

// PUT "/version/:file_id"
func UpdateFileAndNewVersionHandler(c *gin.Context) {
	file, fileHeader, _ := c.Request.FormFile("file")

	v,err := models.NewVersion(file, fileHeader)
	if err != nil {
		log.Fatalf("Failed to New File obj. : %s", err)
	}

	f := models.NewFile()
	models.GetFileById(c.Param("file_id"), f, db.Db)

	objAttrs, err := models.StoreGCS("uploadfile-versioning-me-dev", f, v)
	if err != nil {
		log.Printf("Failed to upload file. ERROR: %s", err)
		return
	}
	v.Url = fmt.Sprintf("https://storage.googleapis.com/%s/%s", objAttrs.Bucket, objAttrs.Name)

	if err := TxUpdateFileAndInsertVersion(v, f, db.Db); err != nil {
		log.Fatalf("Failed to commit. rollback... %s", err.Error())
	}

	c.Redirect(http.StatusFound, "/files")
	// c.JSON(http.StatusOK, gin.H{
	// 	"msg": fmt.Sprintf("%+v uploaded!", *f),
	// 	"url":      f.Url,
	// 	"fileName": objAttrs.Name,
	// })
}


// file update version insert
func TxUpdateFileAndInsertVersion(v *models.Version, f *models.File, db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// データベース操作をトランザクション内で行います（ここからは'db'ではなく'tx'を使います）
		// エラーを返した場合はロールバックされます
		// fileIDを取得
		if err := models.UpdateFile(f, tx).Error;err != nil {
			return err
		}

		// insertしたfile.idをそのままversion.file_idにする。
		v.FileId = f.Id

		// version.versionの最新を+1する。
		models.UpdateVersionNum(v,tx)

		log.Printf("file: %+v\nversion: %+v", f, v)

		if err := models.StoreNewVersion(v,tx).Error;err != nil {
			return err
		}
		// nilを返すとコミットされます
		return nil
	})
}
