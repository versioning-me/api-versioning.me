package server

import (
	"api-versioning-me/controller"
	"github.com/gin-contrib/cors"
	"api-versioning-me/middlewares"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(SetCors())
	router.MaxMultipartMemory = 8 << 20
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middlewares.SetCors())

	router.POST("/new", controller.NewVersionAndNewFileHandler)
	router.GET("/file", controller.GetFilesHandler)
	router.GET("/version/:file_id", controller.GetVersionsByFileIdHandler)
	router.POST("/version", controller.UpdateFileAndNewVersionHandler)

	return router
}
func SetCors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8081"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
	})
}