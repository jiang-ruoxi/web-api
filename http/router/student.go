package router

import (
	"api/http/handler/student_handler"
	"api/middleware"
	"github.com/gin-gonic/gin"
)

func Student(r *gin.RouterGroup) {

	prefixRouter := r.Group("v2").Use(middleware.GlobalMiddleware())

	poetryHandler := student_handler.NewPoetryPictureInfoHandler()
	{
		prefixRouter.GET("/chinese/getNavList", poetryHandler.PoetryPictureInfoList)
	}
}
