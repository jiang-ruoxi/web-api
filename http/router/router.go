package router

import (
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/jiang-ruoxi/gopkg/log"
)

func All() func(r *gin.Engine) {
	return func(r *gin.Engine) {

		// panic日志
		r.Use(ginzap.RecoveryWithZap(log.Sugar().Desugar(), true))

		prefixRouter := r.Group("/")

		// 网站前台
		Api(prefixRouter)
	}
}
