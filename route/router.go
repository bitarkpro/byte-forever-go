package route

import (
	handler "FileStore-Server/handler/Gin-handler"
	"FileStore-Server/logger"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	// gin framework, Logger, Recovery
	router := gin.Default()

	//start logger
	logger.OpenLog(router)


	//Wechat login：
	router.POST("/wxuser/code", handler.WxDoGetCodeHandler)
	router.POST("/wxuser/login", handler.WxDoSignInHandler)

	// user info
	router.POST("/user/info", handler.UserInfoHandler)
	// upload file
	router.POST("/file/upload", handler.DoUploadHandler)
	// down file
	router.GET("/file/download", handler.DownloadHandler)
	// query homepage
	router.POST("/file/homepage", handler.GetHomePageHandler)
	//recent upload files
	router.POST("/file/recent", handler.GetRecentHandler)
	//query folder
	router.POST("/file/folder", handler.GetFolderHandler)
	//file details
	router.POST("/file/querydetail", handler.FileDetailQueryHandler)
	// update(coment,tags)
	router.POST("/file/update", handler.FileUpdateHandler)
	// remove file
	router.POST("/file/delete", handler.FileDeleteHandler)
	// star file
	router.POST("/file/starfile", handler.AddStarFileHandler)
	router.POST("/file/cancelstar", handler.CancelStarFileHandler)

	return router
}
