package router

import (
	"api/http/handler/download_handler"
	"api/middleware"
	"github.com/gin-gonic/gin"
)

func Download(r *gin.RouterGroup) {

	prefixRouter := r.Group("v2").Use(middleware.GlobalMiddleware())

	homeHandler := download_handler.NewIndexHandler()
	{
		prefixRouter.GET("/pic", homeHandler.Index)
	}
}
