package cache

import (
	"api/common"
	"github.com/chenyahui/gin-cache/persist"
	rxRedis "github.com/jiang-ruoxi/gopkg/redis"
	"log"
)

func HttpCache() {
	//判断redis是连接成功
	if rxRedis.ClientDefault("web") == nil {
		log.Printf("%+v", "redis server not connect, http-cache failed.")
		return
	}
	log.Printf("%+v", "http-cache started.")
	common.GVA_HTTP_CACHE = persist.NewRedisStore(rxRedis.ClientDefault("web"))
}
