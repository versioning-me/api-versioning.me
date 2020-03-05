package server

import (
	"api-versioning-me/controller"
	"api-versioning-me/middlewares"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.MaxMultipartMemory = 8 << 20
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middlewares.SetCors())

	router.POST("/new", controller.NewVersionAndNewFileHandler)
	router.GET("/files", controller.GetFilesHandler)
	router.GET("/versions/:file_id", controller.GetVersionsByFileIdHandler)
	router.POST("/version/:file_id", controller.UpdateFileAndNewVersionHandler)

	return router
}
