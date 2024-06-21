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
		r.MaxMultipartMemory = 10 << 20 // 10MB

		prefixRouter := r.Group("/")

		// 绘本项目
		Api(prefixRouter)

		// Market项目
		Market(prefixRouter)
	}
}
