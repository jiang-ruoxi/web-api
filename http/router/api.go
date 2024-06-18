package router

import (
	"api/http/handler/api_handler"
	"github.com/gin-gonic/gin"
)

func Api(r *gin.RouterGroup) {

	prefixRouter := r.Group("/v1")

	homeHandler := api_handler.NewIndexHandler()
	{
		prefixRouter.GET("/home", homeHandler.Index)
	}
}
