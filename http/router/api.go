package router

import (
	"api/http/handler/api_handler"
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
		prefixRouter.GET("/chinese/getNavList", chineseHandler.ChineseGetNavList)
		prefixRouter.GET("/chinese/getList", chineseHandler.ChineseGetBookList)
		prefixRouter.GET("/chinese/getBookInfo", chineseHandler.ChineseGetBookInfo)
	}

	albumHandler := api_handler.NewAlbumHandler()
	{
		prefixRouter.GET("/chinese/getAlbumList", albumHandler.AlbumGetList)
		prefixRouter.GET("/chinese/getAlbumListInfo", albumHandler.AlbumGetListInfo)
		prefixRouter.GET("/chinese/getAlbumInfo", albumHandler.AlbumGetInfo)
	}

	poetryHandler := api_handler.NewPoetryHandler()
	{
		prefixRouter.GET("/poetry/cheng/getList", poetryHandler.PoetryGetChengList)
		prefixRouter.GET("/poetry/cheng/getPoetryInfo", poetryHandler.PoetryGetChengInfo)
		prefixRouter.GET("/poetry/school/getList", poetryHandler.PoetryGetSchoolList)
		//prefixRouter.GET("/poetry/school/getPoetryInfo", poetryHandler.PoetryGetSchoolInfo)
		//prefixRouter.GET("/poetry/junior/getList", poetryHandler.PoetryGetJuniorList)
		//prefixRouter.GET("/poetry/junior/getPoetryInfo", poetryHandler.PoetryGetJuniorInfo)
	}

}
