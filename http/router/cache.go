package router

import (
	"api/cache"
	"api/common"
	"github.com/gin-gonic/gin"
	"time"
)

func routerCache(hour int64) (res gin.HandlerFunc) {
	return cache.CacheByRequestURI(common.GVA_HTTP_CACHE, time.Duration(hour)*time.Hour)
}
