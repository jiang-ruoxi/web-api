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
}
