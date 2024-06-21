package router

import (
	"api/http/handler/market_handler"
	"api/middleware"
	"github.com/gin-gonic/gin"
)

func Market(r *gin.RouterGroup) {

	prefixRouter := r.Group("api/v1").Use(middleware.GlobalMiddleware())
	marketHandler := market_handler.NewMarketHandler()
	{
		//获取地址列表
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/getAddressHot", marketHandler.ApiGetAddressHot)
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/getAddressList", marketHandler.ApiGetAddressList)
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/getAddressChild", marketHandler.ApiGetAddressChildList)
		//根据openId获取用户是否登录
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/getCheckLogin", marketHandler.ApiGetCheckLogin)
		//根据userId获取用户附属信息
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/getUserExt", marketHandler.ApiGetUserExtInfo)
		//上传文件
		prefixRouter.Use(middleware.CheckWechatMiddleware()).POST("/uploadFile", marketHandler.ApiUploadFileData)
		//获取banner列表
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/getBannerList", marketHandler.ApiGetBannerList)
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/getBannerListNew", marketHandler.ApiGetBannerListNew)
		//获取工种列表
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/getTagList", marketHandler.ApiGetTagList)
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/getTagSelect", marketHandler.ApiGetTagSelect)
		//获取会员价格列表
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/getPayList", marketHandler.ApiGetPayList)
		//获取优选工匠的价格
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/getGoodPay", marketHandler.ApiGetGoodPay)
		//获取优选工匠列表
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/getMemberList", marketHandler.ApiGetGoodMemberList)
		//获取会员详情
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/getMemberInfo", marketHandler.ApiGetMemberInfo)
		//更新用户资料信息
		prefixRouter.Use(middleware.CheckWechatMiddleware()).POST("/updateMemberData", marketHandler.ApiUpdateMemberData)
		//获取任务列表
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/getTaskList", marketHandler.ApiGetTaskList)
		//获取已发布的任务列表
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/getMyTaskList", marketHandler.ApiGetMyTaskList)
		//更新任务状态
		prefixRouter.Use(middleware.CheckWechatMiddleware()).POST("/updateTaskStatus", marketHandler.ApiUpdateTaskStatus)
		//获取任务详情
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/getTaskInfo", marketHandler.ApiGetTaskInfo)
		//发布任务
		prefixRouter.Use(middleware.CheckWechatMiddleware()).POST("/doMakeTaskData", marketHandler.ApiDoMakeTaskData)
		//代发任务
		prefixRouter.Use(middleware.CheckWechatMiddleware()).POST("/other/doMakeTaskData", marketHandler.ApiDoMakeOtherTaskData)
		//校验是否可发布
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/checkPushTask", marketHandler.ApiCheckPushTask)
		//创建用户
		prefixRouter.Use(middleware.CheckWechatMiddleware()).POST("/doMakeUserData", marketHandler.ApiDoMakeUserData)
	}

	prefixRouter = r.Group("api/wechat").Use(middleware.GlobalMiddleware())
	marketHandler = market_handler.NewMarketHandler()
	{
		//获取openid
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/getWxData", marketHandler.ApiGetWechatData)
		//获取access_token
		prefixRouter.Use(middleware.CheckWechatMiddleware()).GET("/getWxAccessToken", marketHandler.ApiGetWxAccessToken)
		//获取用户手机号
		prefixRouter.Use(middleware.CheckWechatMiddleware()).POST("/getWxUserPhoneNumber", marketHandler.ApiGetWxUserPhoneNumber)
		//微信支付
		prefixRouter.Use(middleware.CheckWechatMiddleware()).POST("/pay", marketHandler.ApiGetWxPay)
		//微信开通优选工匠
		prefixRouter.Use(middleware.CheckWechatMiddleware()).POST("/open/pay", marketHandler.ApiGetWxOpenPay)
		//微信支付回调
		prefixRouter.Use(middleware.CheckWechatMiddleware()).POST("/pay/notice", marketHandler.ApiGetWxPayCallback)
		//微信支付更新为取消
		prefixRouter.Use(middleware.CheckWechatMiddleware()).POST("/pay/cancel", marketHandler.ApiGetWxPayCancel)
		//微信支付退款
		prefixRouter.Use(middleware.CheckWechatMiddleware()).POST("/pay/refunds", marketHandler.ApiGetWxPayRefunds)
	}
}
