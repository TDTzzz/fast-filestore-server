package route

import (
	"fast-filestore-server/handler"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()

	//静态资源处理
	router.Static("/static/", "./static")
	//不需要认证的接口

	router.GET("/file/upload", handler.UploadHandler)
	router.POST("/file/upload", handler.DoUploadHandler)
	router.GET("/file/file/:sha1", handler.GetFileBySha1)
	router.GET("/file/query/:limit", handler.FileQueryHandler)
	router.GET("/file/download/:sha1", handler.FileDownloadHandler)
	router.POST("/file/update", handler.FileUpdateHandler)
	router.POST("/file/delete", handler.FileDeleteHandler)
	//中间件处理（认证相关）

	//文件存储相关(需要认证)

	return router
}
