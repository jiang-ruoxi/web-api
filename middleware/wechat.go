package middleware

import (
	"api/common"
	"github.com/gin-gonic/gin"
	"strings"
)

func GlobalMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("gin-base-url", common.GIN_BASE_URL)
		ctx.Set("gin-less-value", common.GIN_LESS_VALUE)
		ctx.Next()
	}
}

// CheckWechatMiddleware 验证是否为微信小程序访问
func CheckWechatMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !checkRequestUserAgent(ctx) {
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func checkRequestUserAgent(c *gin.Context) bool {
	uaText := c.Request.Header.Get("User-Agent")
	isFlag := strings.Contains(strings.ToLower(uaText), common.MINI_WECHAT)
	if !isFlag {
		ReturnResponse(common.FORBID, map[string]interface{}{}, common.FORBID_MSG, c)
		return false
	}
	return true
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func ReturnResponse(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(200, Response{
		code,
		msg,
		data,
	})
}
