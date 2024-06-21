package market_handler

import (
	"api/common"
	api2 "api/common/market"
	"api/service/market_service"
	"api/utils"
	"api/utils/errs"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

func NewMarketHandler() *MarketHandler {
	return &MarketHandler{
		service: market_service.NewMarketService(),
	}
}

type MarketHandler struct {
	service *market_service.MarketService
}

func (mh *MarketHandler) ApiGetAddressHot(ctx *gin.Context) {
	response, err := mh.service.ApiGetAddressHot()
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}

func (mh *MarketHandler) ApiGetAddressList(ctx *gin.Context) {
	response, err := mh.service.ApiGetAddressList()
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}

func (mh *MarketHandler) ApiGetAddressChildList(ctx *gin.Context) {
	response, err := mh.service.ApiGetAddressChildList()
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}

func (mh *MarketHandler) ApiGetCheckLogin(ctx *gin.Context) {
	response, err := mh.service.ApiGetCheckLogin(ctx.Query("open_id"))
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}

func (mh *MarketHandler) ApiGetUserExtInfo(ctx *gin.Context) {
	userId, _ := strconv.Atoi(ctx.Query("user_id"))
	response, err := mh.service.ApiGetUserExtInfo(userId)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}

func (mh *MarketHandler) ApiUploadFileData(ctx *gin.Context) {
	//获取其他的参数
	userId := ctx.PostForm("user_id")
	// 获取上传的文件
	file, err := ctx.FormFile("file")
	if err != nil {
		return
	}
	// 获取文件名
	fileName := filepath.Base(file.Filename)
	// 获取文件扩展名
	extension := filepath.Ext(fileName)

	// 生成文件名（使用时间戳）
	fName := time.Now().Format("20060102150405") + extension

	path := "/data/static/market/" + userId
	utils.ExistDir(path)

	// 保存文件到服务器
	ctx.SaveUploadedFile(file, filepath.Join(path, fName))

	dst := "https://oss.58haha.com/market/" + userId + "/" + fName

	ctx.JSON(200, gin.H{
		"dst": dst,
	})
}

func (mh *MarketHandler) ApiGetBannerList(ctx *gin.Context) {
	tType, _ := strconv.Atoi(ctx.Query("type"))
	response, err := mh.service.ApiGetBannerList(tType)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}

func (mh *MarketHandler) ApiGetBannerListNew(ctx *gin.Context) {
	tType, _ := strconv.Atoi(ctx.Query("type"))
	response, err := mh.service.ApiGetBannerListNew(tType)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}

func (mh *MarketHandler) ApiGetTagList(ctx *gin.Context) {
	response, err := mh.service.ApiGetTagList()
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}

func (mh *MarketHandler) ApiGetTagSelect(ctx *gin.Context) {
	response, err := mh.service.ApiGetTagSelect()
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}

func (mh *MarketHandler) ApiGetPayList(ctx *gin.Context) {
	response, err := mh.service.ApiGetPayList()
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}

func (mh *MarketHandler) ApiGetGoodPay(ctx *gin.Context) {
	response, err := mh.service.ApiGetGoodPay()
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}

func (mh *MarketHandler) ApiGetGoodMemberList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	tType, _ := strconv.Atoi(ctx.Query("type"))
	response, err := mh.service.ApiGetGoodMemberList(page, tType)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}

func (mh *MarketHandler) ApiGetMemberInfo(ctx *gin.Context) {
	userId, _ := strconv.Atoi(ctx.Query("user_id"))
	response, err := mh.service.ApiGetMemberInfo(userId)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}

func (mh *MarketHandler) ApiUpdateMemberData(ctx *gin.Context) {
	var json api2.MemberUpdateDataRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		return
	}
	response, err := mh.service.ApiUpdateMemberData(json)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}

func (mh *MarketHandler) ApiGetTaskList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	tType, _ := strconv.Atoi(ctx.Query("type"))
	addressId, _ := strconv.Atoi(ctx.Query("addressId"))
	response, err := mh.service.ApiGetTaskList(page, tType, addressId)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}

func (mh *MarketHandler) ApiGetMyTaskList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	userId, _ := strconv.Atoi(ctx.Query("user_id"))
	response, err := mh.service.ApiGetMyTaskList(page, userId)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}

func (mh *MarketHandler) ApiUpdateTaskStatus(ctx *gin.Context) {
	var json api2.TaskUpdateStatusRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := mh.service.ApiUpdateTaskStatus(json)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}

func (mh *MarketHandler) ApiGetTaskInfo(ctx *gin.Context) {
	taskId, _ := strconv.Atoi(ctx.Query("task_id"))
	response, err := mh.service.ApiGetTaskInfo(taskId)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}

func (mh *MarketHandler) ApiDoMakeTaskData(ctx *gin.Context) {
	var json api2.MakeTaskDataRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := mh.service.ApiDoMakeTaskData(json)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}

func (mh *MarketHandler) ApiDoMakeOtherTaskData(ctx *gin.Context) {
	var json api2.MakeTaskOtherDataRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := mh.service.ApiDoMakeOtherTaskData(json)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}

func (mh *MarketHandler) ApiCheckPushTask(ctx *gin.Context) {
	userId, _ := strconv.Atoi(ctx.Query("user_id"))
	response, err := mh.service.ApiCheckPushTask(userId)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}

func (mh *MarketHandler) ApiDoMakeUserData(ctx *gin.Context) {
	var json api2.MakeUserDataRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := mh.service.ApiDoMakeUserData(json)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}

func (mh *MarketHandler) ApiGetWechatData(ctx *gin.Context) {
	code := ctx.Query("code")
	response, err := mh.service.ApiGetWechatData(code)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}

func (mh *MarketHandler) ApiGetWxAccessToken(ctx *gin.Context) {
	response, err := mh.service.ApiGetWxAccessToken()
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}

func (mh *MarketHandler) ApiGetWxUserPhoneNumber(ctx *gin.Context) {
	var json api2.MakePhotoDataRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := mh.service.ApiGetWxUserPhoneNumber(json)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}

func (mh *MarketHandler) ApiGetWxPay(ctx *gin.Context) {
	var json api2.WXPayDataRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := mh.service.ApiGetWxPay(json)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}

func (mh *MarketHandler) ApiGetWxOpenPay(ctx *gin.Context) {
	var json api2.OpenGoodPayRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := mh.service.ApiGetWxOpenPay(json)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}

func (mh *MarketHandler) ApiGetWxPayCallback(ctx *gin.Context) {
	notifyReq, err := mh.service.ApiDealWxPayCallback(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": "FAIL", "message": err.Error()})
		return
	}
	//解析Plaintext参数，这里面可以拿到订单的基本信息
	var data = []byte(notifyReq.Resource.Plaintext)
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		fmt.Println("JSON转换失败：", err)
		return
	}
	// 处理通知内容
	fmt.Println("解密结果", notifyReq)
	fmt.Println("解密结果Plaintext", result)

	//将解密结果进行处理
	mh.service.ApiDealUserPaySuccess(notifyReq, result)
	ctx.JSON(http.StatusOK, gin.H{"code": "SUCCESS"})
}
func (mh *MarketHandler) ApiGetWxPayCancel(ctx *gin.Context) {
	var json api2.WXCancelPayDataRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := mh.service.ApiGetWxPayCancel(json)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}

func (mh *MarketHandler) ApiGetWxPayRefunds(ctx *gin.Context) {
	var json api2.WXRefundsPayDataRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	message, code := mh.service.ApiGetWxPayRefunds(json)
	common.ReturnResponse(common.SUCCESS, map[string]interface{}{
		"message": message,
		"code":    code,
	}, common.SUCCESS_MSG, ctx)
}
