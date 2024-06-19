package router

import (
	"api/common"
	"api/http/handler/api_handler"
	"api/middleware"
	"github.com/gin-gonic/gin"
)

func Api(r *gin.RouterGroup) {

	prefixRouter := r.Group("v2").Use(middleware.GlobalMiddleware())

	homeHandler := api_handler.NewIndexHandler()
	{
		prefixRouter.GET("/home", homeHandler.Index)
	}

	chineseHandler := api_handler.NewChineseHandler()
	{
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/chinese/getNavList", routerCache(common.RedisURL_CACHE), chineseHandler.ChineseGetNavList)
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/chinese/getList", routerCache(common.RedisURL_CACHE), chineseHandler.ChineseGetBookList)
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/chinese/getBookInfo", routerCache(common.RedisURL_CACHE), chineseHandler.ChineseGetBookInfo)
	}

	albumHandler := api_handler.NewAlbumHandler()
	{
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/chinese/getAlbumList", routerCache(common.RedisURL_CACHE), albumHandler.AlbumGetList)
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/chinese/getAlbumListInfo", routerCache(common.RedisURL_CACHE), albumHandler.AlbumGetListInfo)
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/chinese/getAlbumInfo", routerCache(common.RedisURL_CACHE), albumHandler.AlbumGetInfo)
	}

	poetryHandler := api_handler.NewPoetryHandler()
	{
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/poetry/cheng/getList", routerCache(common.RedisURL_CACHE), poetryHandler.PoetryGetChengList)
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/poetry/cheng/getPoetryInfo", routerCache(common.RedisURL_CACHE), poetryHandler.PoetryGetChengInfo)
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/poetry/school/getList", routerCache(common.RedisURL_CACHE), poetryHandler.PoetryGetSchoolList)
		//prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/poetry/school/getPoetryInfo", poetryHandler.PoetryGetSchoolInfo)
		//prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/poetry/junior/getList", poetryHandler.PoetryGetJuniorList)
		//prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/poetry/junior/getPoetryInfo", poetryHandler.PoetryGetJuniorInfo)
	}
}
