package router

import (
	"api/common"
	"api/http/handler/api_handler"
	"api/middleware"
	"github.com/gin-gonic/gin"
)

func Api(r *gin.RouterGroup) {

	prefixRouter := r.Group("v2")

	homeHandler := api_handler.NewIndexHandler()
	{
		prefixRouter.GET("/home", homeHandler.Index)
	}

	chineseHandler := api_handler.NewChineseHandler()
	{
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/chinese/getNavList", routerCache(common.RedisURL_CACHE), chineseHandler.ChineseGetNavList)
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/chinese/getList", chineseHandler.ChineseGetBookList)
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/chinese/getBookInfo", chineseHandler.ChineseGetBookInfo)
	}

	albumHandler := api_handler.NewAlbumHandler()
	{
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/chinese/getAlbumList", albumHandler.AlbumGetList)
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/chinese/getAlbumListInfo", albumHandler.AlbumGetListInfo)
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/chinese/getAlbumInfo", albumHandler.AlbumGetInfo)
	}

	poetryHandler := api_handler.NewPoetryHandler()
	{
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/poetry/cheng/getList", poetryHandler.PoetryGetChengList)
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/poetry/cheng/getPoetryInfo", poetryHandler.PoetryGetChengInfo)
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/poetry/school/getList", poetryHandler.PoetryGetSchoolList)
		//prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/poetry/school/getPoetryInfo", poetryHandler.PoetryGetSchoolInfo)
		//prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/poetry/junior/getList", poetryHandler.PoetryGetJuniorList)
		//prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/poetry/junior/getPoetryInfo", poetryHandler.PoetryGetJuniorInfo)
	}

}
