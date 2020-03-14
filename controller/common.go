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

	// デフォルトでtitleに現在のfilenameを入れておく
	// Todo
	// POST時にtitleとdetailを入れる？一旦空欄stringで
	f := models.NewFile()
	f.Title = fileHeader.Filename
	if c.PostForm("title") != "" {
		f.Title = c.PostForm("title")
	}
	f.Detail = c.PostForm("detail")

	objAttrs, err := models.StoreGCS(f, v)
	if err != nil {
		log.Printf("Failed to upload file. ERROR: %s", err)
		return
	}
	v.Url = fmt.Sprintf("https://storage.googleapis.com/%s/%s", objAttrs.Bucket, objAttrs.Name)

	if err := TxInsertVersionAndInsertFile(v, f, db.Db); err != nil {
		log.Fatalf("Failed to commit. rollback... %s", err.Error())
	}

	c.Redirect(http.StatusFound, "/file")
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

// PUT "/version/"
func UpdateFileAndNewVersionHandler(c *gin.Context) {
	file, fileHeader, _ := c.Request.FormFile("file")

	v,err := models.NewVersion(file, fileHeader)
	if err != nil {
		log.Fatalf("Failed to New File obj. : %s", err)
	}

	f := models.NewFile()
	models.GetFileById(c.PostForm("file_id"), f, db.Db)

	f.Title = fileHeader.Filename
	if c.PostForm("title") != "" {
		f.Title = c.PostForm("title")
	}
	f.Detail = c.PostForm("detail")


	objAttrs, err := models.StoreGCS(f, v)
	if err != nil {
		log.Printf("Failed to upload file. ERROR: %s", err)
		return
	}
	v.Url = fmt.Sprintf("https://storage.googleapis.com/%s/%s", objAttrs.Bucket, objAttrs.Name)

	if err := TxUpdateFileAndInsertVersion(v, f, db.Db); err != nil {
		log.Fatalf("Failed to commit. rollback... %s", err.Error())
	}

	c.Redirect(http.StatusFound, "/file")
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
