package common

import (
	"github.com/chenyahui/gin-cache/persist"
	"github.com/go-redis/redis/v8"
)

const MINI_WECHAT = "micromessenger"
const DEFAULT_PAGE = 1
const DEFAULT_LEVEL = 1
const DEFAULT_PAGE_SIZE = 20
const RedisURL_CACHE = 7200

const (
	SUCCESS                = 10000
	FAIL                   = 10001
	FORBID                 = 403
	ERR_RES_PARAMS_ILLEGAL = 10002
)

const (
	SUCCESS_MSG                = "success"
	FAIL_MSG                   = "fail"
	FORBID_MSG                 = "非法请求，禁止访问"
	ERR_RES_PARAMS_ILLEGAL_MSG = "参数传递不合法"
)

var (
	GVA_REDIS      *redis.Client
	GVA_HTTP_CACHE *persist.RedisStore
)
