package server

import (
	"api-versioning-me/controller"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.MaxMultipartMemory = 8 << 20
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.POST("/file", controller.FileUploadHandler)
	router.GET("/files", controller.GetUploadedFilesHandler)
	router.GET("/files/:version_name", controller.GetUploadedFilesByVersionNameHandler)
	router.GET("/versions", controller.GetVersionsHandler)

	return router
}
