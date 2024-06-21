package market_service

import (
	"api/common"
	api2 "api/common/market"
	"api/model"
	"api/model/market"
	"api/utils"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	rxRedis "github.com/jiang-ruoxi/gopkg/redis"
	"github.com/jiang-ruoxi/gopkg/server/api"
	"github.com/spf13/viper"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
	apiUtils "github.com/wechatpay-apiv3/wechatpay-go/utils"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func NewMarketService() *MarketService {
	return &MarketService{}
}

type MarketService struct {
}

func (srv *MarketService) ApiGetAddressHot() (response *api2.ZMAddressHotResponse, apiErr api.Error) {
	response = &api2.ZMAddressHotResponse{}
	model.DefaultMarket().Model(&market.ZMAddress{}).Debug().
		Where("is_deleted = 0 and is_hot = 1 and parent_id > 0").
		Order("sort desc").Order("id asc").Limit(5).
		Find(&response.List)
	return
}

func (srv *MarketService) ApiGetAddressList() (response *api2.ZMAddressListResponse, apiErr api.Error) {
	var addressParentList []market.ZMAddress
	model.DefaultMarket().Model(&market.ZMAddress{}).Debug().
		Where("is_deleted = 0 and parent_id = 0").
		Order("sort desc").Order("id asc").Find(&addressParentList)
	for idx, _ := range addressParentList {
		var addressChildList []*market.ZMAddress
		model.DefaultMarket().Model(&market.ZMAddress{}).Debug().Where("is_deleted = 0 and parent_id = ?", addressParentList[idx].Id).Order("sort desc").Order("id asc").Find(&addressChildList)
		addressParentList[idx].ChildList = addressChildList
	}
	response.List = addressParentList
	return
}

func (srv *MarketService) ApiGetAddressChildList() (response *api2.AddressChildListResponse, apiErr api.Error) {
	response = &api2.AddressChildListResponse{}
	var addressList []market.ZMAddress
	model.DefaultMarket().Model(&market.ZMAddress{}).Where("is_deleted = 0 and parent_id > 0").Debug().
		Order("sort desc").Order("id asc").
		Find(&addressList)
	var temp api2.FormatAddressInfo
	for idx, _ := range addressList {
		temp.CityId = addressList[idx].Id
		temp.CityName = addressList[idx].Name
		response.List = append(response.List, temp)
	}
	return
}

func (srv *MarketService) ApiGetCheckLogin(openId string) (response *api2.ZMUserResponse, apiErr api.Error) {
	response = &api2.ZMUserResponse{}
	model.DefaultMarket().Model(&market.ZMUser{}).Where("status=1 AND open_id = ?", openId).Debug().
		First(&response.Info)
	response.Info.NickName = utils.TruncateString(response.Info.NickName, 5)
	if response.Info.IsBest == 1 {
		//当前时间
		currentYMD, _ := strconv.Atoi(utils.GetCurrentDateYMD())
		if currentYMD > response.Info.BestLimit {
			response.Info.IsBest = 0
		}
	}
	return
}

func (srv *MarketService) ApiGetUserExtInfo(userId int) (response *api2.ZMUserExtResponse, apiErr api.Error) {
	response = &api2.ZMUserExtResponse{}
	model.DefaultMarket().Model(&market.ZMUserExt{}).Where("user_id = ?", userId).Debug().
		First(&response.Info)
	return
}

func (srv *MarketService) ApiGetBannerList(tType int) (response *api2.ZMBannerListResponse, apiErr api.Error) {
	response = &api2.ZMBannerListResponse{}
	model.DefaultMarket().Model(&market.ZMBanner{}).Where("is_deleted = 0 and status=1 and type = ?", tType).Debug().
		Order("id asc").Limit(4).
		Find(&response.List)
	return
}

func (srv *MarketService) ApiGetBannerListNew(tType int) (response *api2.ZMBannerListNewResponse, apiErr api.Error) {
	response = &api2.ZMBannerListNewResponse{}
	var banners []market.ZMBanner
	model.DefaultMarket().Model(&market.ZMBanner{}).Where("is_deleted = 0 and status=1 and type = ?", tType).Debug().
		Order("id asc").Limit(4).
		Find(&response.List)
	for idx, _ := range banners {
		response.List = append(response.List, banners[idx].Image)
	}
	return
}

func (srv *MarketService) ApiGetTagList() (response *api2.ZMTagsListResponse, apiErr api.Error) {
	response = &api2.ZMTagsListResponse{}
	model.DefaultMarket().Model(&market.ZMTags{}).Where("is_deleted = 0 and status=1").Debug().
		Order("id asc").Limit(20).
		Find(&response.List)

	return
}

func (srv *MarketService) ApiGetTagSelect() (response *api2.ZMTagsListSelectResponse, apiErr api.Error) {
	response = &api2.ZMTagsListSelectResponse{}
	var tagList []market.ZMTags
	model.DefaultMarket().Model(&market.ZMTags{}).Where("is_deleted = 0 and status=1").Debug().
		Order("sort desc").Order("id asc").Limit(20).
		Find(&tagList)
	response.List = append(response.List, "请选择分类")
	for idx, _ := range tagList {
		response.List = append(response.List, tagList[idx].Name)
	}
	return
}

func (srv *MarketService) ApiGetPayList() (response *api2.ZMPayListResponse, apiErr api.Error) {
	response = &api2.ZMPayListResponse{}
	var payList []market.ZMPay
	model.DefaultMarket().Model(&market.ZMPay{}).Debug().
		Where("type=1 and status=1 and is_deleted = 0").
		Order("sort desc").Limit(6).
		Find(&payList)
	var temp api2.FormatData
	for idx, _ := range payList {
		var payCPrice float64
		payCPrice = payList[idx].CPrice / 100
		temp.Checked = payList[idx].Checked
		temp.TotalFee = payCPrice
		temp.Number = payList[idx].Number
		temp.NumberExt = payList[idx].NumberExt
		temp.Value = strconv.Itoa(payList[idx].Id)
		strValue := fmt.Sprintf("%.2f", payCPrice)
		temp.Name = payList[idx].Name + "（￥" + strValue + "）"
		response.List = append(response.List, temp)
	}
	return
}

func (srv *MarketService) ApiGetGoodPay() (response *api2.ZMPayResponse, apiErr api.Error) {
	response = &api2.ZMPayResponse{}
	var payList []market.ZMPay
	model.DefaultMarket().Model(&market.ZMPay{}).Where("type=2 and status=1 and is_deleted = 0").Debug().
		Order("sort desc").Limit(1).
		Find(&payList)
	if len(payList) > 0 {
		response.Info = payList[0]
		var payCPrice float64
		payCPrice = response.Info.CPrice / 100
		response.Info.CPrice = payCPrice
	}
	return
}

func (srv *MarketService) ApiGetGoodMemberList(page int, tType int) (response *api2.ZMGoodMemberListResponse, apiErr api.Error) {
	response = &api2.ZMGoodMemberListResponse{}
	size := common.DEFAULT_PAGE_SIZE
	offset := size * (page - 1)
	var tagDataList []market.ZMTags
	model.DefaultMarket().Model(&market.ZMTags{}).Find(&tagDataList)

	var count int64
	var memberList []market.ZMUser
	odb := model.DefaultMarket().Model(&market.ZMAddress{}).Debug().
		Where("is_deleted = 0 and status= 1 and is_best = 1")
	if tType > 0 && tType < 15 {
		odb = odb.Where(" tag_id = ?", tType)
	}
	odb.Count(&count)
	odb = odb.Order("id desc").Limit(size).Offset(offset)
	odb.Find(&memberList)

	response.Count = count

	//组合userId的集合
	var userIds []int64
	for idx, _ := range memberList {
		userIds = append(userIds, memberList[idx].UserId)
	}
	var memberExtList []market.ZMUserExt

	odbExt := model.DefaultMarket().Model(&market.ZMUserExt{}).Debug()
	odb = odbExt.Where("user_id in(?)", userIds).Find(&memberExtList)

	var temp api2.MemberData
	for idx, _ := range memberList {
		temp.Desc = ""
		temp.ViewCount = 1
		temp.Id = memberList[idx].Id
		temp.UserId = memberList[idx].UserId
		temp.OpenId = memberList[idx].OpenId
		temp.NickName = utils.TruncateString(memberList[idx].NickName, 3)
		temp.RealName = memberList[idx].RealName
		temp.HeadUrl = memberList[idx].HeadUrl
		temp.Mobile = memberList[idx].Mobile
		temp.TagId = memberList[idx].TagId
		for dIndex, _ := range tagDataList {
			if memberList[idx].TagId == tagDataList[dIndex].Id {
				temp.TagName = tagDataList[dIndex].Name
			}
		}

		for dIndex, _ := range memberExtList {
			if memberList[idx].UserId == memberExtList[dIndex].UserId {
				temp.Desc = utils.TruncateString(utils.ClearMobileText(utils.RegContent(memberExtList[dIndex].Desc, srv.getBadWordsList())), 50) + "......"
				temp.ViewCount = memberExtList[dIndex].ViewCount
			}
		}
		temp.IsBest = memberList[idx].IsBest
		temp.IsMember = memberList[idx].IsMember
		temp.MemberLimit = memberList[idx].MemberLimit
		response.List = append(response.List, temp)
	}

	return
}

func (srv *MarketService) getBadWordsList() (badWordsSlice []string) {
	var badWords []market.ZMBadWords
	model.DefaultMarket().Model(&market.ZMBadWords{}).Debug().Find(&badWords)
	for idx, _ := range badWords {
		badWordsSlice = append(badWordsSlice, badWords[idx].Name)
	}
	return
}

func (srv *MarketService) ApiGetMemberInfo(userId int) (response *api2.ZMGoodUserInfoResponse, apiErr api.Error) {
	response = &api2.ZMGoodUserInfoResponse{}
	var user market.ZMUser
	odb := model.DefaultMarket().Model(&market.ZMUser{}).Debug()
	odb.Where("user_id=?", userId).First(&user)

	var userExt market.ZMUserExt
	odbExt := model.DefaultMarket().Model(&market.ZMUserExt{}).Debug()
	odbExt.Where("user_id=?", userId).First(&userExt)

	model.DefaultMarket().Model(&market.ZMUserExt{}).Debug().Where("user_id=?", userId).Update("view_count", userExt.ViewCount+1)

	var tagInfo market.ZMTags
	model.DefaultMarket().Model(&market.ZMTags{}).Debug().Where("id=?", response.Info.TagId).First(&tagInfo)

	response.Info.Id = user.Id
	response.Info.UserId = user.UserId
	response.Info.OpenId = user.OpenId
	response.Info.NickName = utils.TruncateString(user.NickName, 5)
	response.Info.RealName = utils.TruncateString(user.RealName, 5)
	response.Info.HeadUrl = user.HeadUrl
	response.Info.Mobile = user.Mobile
	response.Info.TagId = user.TagId
	response.Info.TagName = tagInfo.Name
	response.Info.Desc = utils.ClearMobileText(utils.RegContent(userExt.Desc, srv.getBadWordsList()))
	response.Info.Demo = userExt.Demo
	response.Info.ViewCount = userExt.ViewCount

	//查询城市详情
	var addressInfo market.ZMAddress
	model.DefaultMarket().Model(&market.ZMAddress{}).Debug().Where("id = ?", userExt.AddressId).First(&addressInfo)

	response.Info.Address = addressInfo.Name
	return
}

func (srv *MarketService) ApiUpdateMemberData(memberData api2.MemberUpdateDataRequest) (response *api2.MemberUpdateDataResponse, apiErr api.Error) {
	if memberData.UserId <= 0 {
		return
	}
	var member market.ZMUser
	if memberData.NickName != "" {
		member.NickName = memberData.NickName
	}
	if memberData.Mobile != "" {
		member.Mobile = memberData.Mobile
	}
	if memberData.HeadUrl != "" {
		member.HeadUrl = memberData.HeadUrl
	}
	model.DefaultMarket().Model(&market.ZMUser{}).Debug().Where("user_id=?", memberData.UserId).Updates(&member)

	var memberExt market.ZMUserExt
	memberExt.AddressId = memberData.AddressId
	memberExt.Address = memberData.Address
	model.DefaultMarket().Model(&market.ZMUserExt{}).Debug().Where("user_id=?", memberData.UserId).Updates(&memberExt)
	response.Result = true

	return
}

func (srv *MarketService) ApiGetTaskList(page, tType, addressId int) (response *api2.ZMTaskListResponse, apiErr api.Error) {
	response = &api2.ZMTaskListResponse{}

	var tagDataList []market.ZMTags
	model.DefaultMarket().Model(&market.ZMTags{}).Find(&tagDataList)

	size := common.DEFAULT_PAGE_SIZE
	offset := size * (page - 1)
	var count int64
	var taskList []market.ZMTask
	odb := model.DefaultMarket().Model(&market.ZMTask{}).Debug()
	odb = odb.Where("is_deleted = 0 and status > 0")
	if addressId > 0 {
		odb = odb.Where(" address_id = ?", addressId)
	}
	if tType > 1 {
		odb = odb.Where(" tag_id = ?", tType)
	}
	odb.Count(&count)
	odb = odb.Order("id desc").Limit(size).Offset(offset)
	odb.Find(&taskList)

	//组合userId的集合
	var userIds []int64
	for idx, _ := range taskList {
		userIds = append(userIds, taskList[idx].UserId)
	}
	userIdsResult := utils.RemoveDuplicates(userIds)
	var memberList []market.ZMUser
	odbUser := model.DefaultMarket().Model(&market.ZMUser{}).Debug()
	odb = odbUser.Where("user_id in(?)", userIdsResult).Find(&memberList)

	//获取所有的地址
	//查询父级ID
	var addressList []market.ZMAddress
	model.DefaultMarket().Model(&market.ZMAddress{}).Debug().Where("is_deleted = 0 and parent_id > 0").Find(&addressList)

	badWordList := srv.getBadWordsList()
	var temp api2.FormatTaskData
	for idx, _ := range taskList {
		temp.Id = taskList[idx].Id
		temp.TagId = taskList[idx].TagId
		temp.UserId = taskList[idx].UserId
		for dIndex, _ := range tagDataList {
			if taskList[idx].TagId == tagDataList[dIndex].Id {
				temp.TagName = tagDataList[dIndex].Name
			}
		}

		//过滤敏感数据并清除手机号
		tempDesc := utils.ClearMobileText(utils.RegContent(taskList[idx].Desc, badWordList))
		temp.Desc = utils.TruncateString(tempDesc, 60) + "......"

		temp.Mobile = ""
		temp.IsBest = 0
		for dIndex, _ := range memberList {
			if memberList[dIndex].UserId == taskList[idx].UserId {
				temp.Mobile = memberList[dIndex].Mobile
				temp.IsBest = memberList[dIndex].IsBest
			}
		}

		//判断是否后台填入
		if taskList[idx].Mobile != "" {
			temp.Mobile = taskList[idx].Mobile
		}

		temp.Date = utils.GetUnixTimeToDateTime1(taskList[idx].AddTime)
		temp.Address = ""
		for aIdx, _ := range addressList {
			if addressList[aIdx].Id == taskList[idx].AddressId {
				temp.Address = addressList[aIdx].Name
			}
		}
		temp.Status = taskList[idx].Status

		response.List = append(response.List, temp)
	}
	response.Count = count

	return
}

func (srv *MarketService) ApiGetMyTaskList(page, userId int) (response *api2.ZMTaskListResponse, apiErr api.Error) {
	response = &api2.ZMTaskListResponse{}

	var tagDataList []market.ZMTags
	model.DefaultMarket().Model(&market.ZMTags{}).Find(&tagDataList)

	size := common.DEFAULT_PAGE_SIZE
	offset := size * (page - 1)
	var count int64
	var taskList []market.ZMTask
	odb := model.DefaultMarket().Model(&market.ZMTask{}).Debug()
	if userId > 0 {
		odb = odb.Where("is_deleted = 0 and user_id = ?", userId)
	}
	odb.Count(&count)
	odb = odb.Order("id desc").Limit(size).Offset(offset)
	odb.Find(&taskList)

	//组合userId的集合
	var userIds []int64
	for idx, _ := range taskList {
		userIds = append(userIds, taskList[idx].UserId)
	}
	var memberList []market.ZMUser
	odbUser := model.DefaultMarket().Model(&market.ZMUser{}).Debug()
	odb = odbUser.Where("user_id in(?)", userIds).Find(&memberList)

	//查询父级ID
	var addressList []market.ZMAddress
	model.DefaultMarket().Model(&market.ZMAddress{}).Debug().Where("is_deleted = 0 and parent_id > 0").Find(&addressList)

	var temp api2.FormatTaskData
	for idx, _ := range taskList {
		temp.Id = taskList[idx].Id
		temp.TagId = taskList[idx].TagId
		for dIndex, _ := range tagDataList {
			if taskList[idx].TagId == tagDataList[dIndex].Id {
				temp.TagName = tagDataList[dIndex].Name
			}
		}

		temp.Desc = utils.TruncateString(taskList[idx].Desc, 50) + "......"
		for dIndex, _ := range memberList {
			if taskList[idx].UserId == memberList[dIndex].UserId {
				temp.Mobile = memberList[dIndex].Mobile
			}
		}

		temp.Date = utils.GetUnixTimeToDateTime1(taskList[idx].AddTime)

		temp.Address = ""
		for aIdx, _ := range addressList {
			if addressList[aIdx].Id == taskList[idx].AddressId {
				temp.Address = addressList[aIdx].Name
			}
		}
		temp.Status = taskList[idx].Status
		response.List = append(response.List, temp)
	}
	response.Count = count
	return
}

func (srv *MarketService) ApiUpdateTaskStatus(taskData api2.TaskUpdateStatusRequest) (response *api2.TaskUpdateStatusResponse, apiErr api.Error) {
	response = &api2.TaskUpdateStatusResponse{}
	if taskData.TaskId < 0 {
		return
	}
	var task market.ZMTask
	task.Status = taskData.Status
	model.DefaultMarket().Model(&market.ZMTask{}).Debug().Where("id=?", taskData.TaskId).Updates(&task)
	response.Result = true

	return
}

func (srv *MarketService) ApiGetTaskInfo(taskId int) (response *api2.ZMTaskInfoResponse, apiErr api.Error) {
	response = &api2.ZMTaskInfoResponse{}
	var task market.ZMTask
	odb := model.DefaultMarket().Model(&market.ZMTask{}).Debug()
	odb.Where("id=?", taskId).First(&task)

	var user market.ZMUser
	odbUser := model.DefaultMarket().Model(&market.ZMUser{}).Debug()
	odbUser.Where("user_id=?", task.UserId).First(&user)

	var tagInfo market.ZMTags
	model.DefaultMarket().Model(&market.ZMTags{}).Debug().Where("id=?", task.TagId).First(&tagInfo)

	response.Info.Id = task.Id
	response.Info.UserId = task.UserId
	response.Info.TagId = task.TagId
	response.Info.TagName = tagInfo.Name
	response.Info.Desc = utils.ClearMobileText(utils.RegContent(task.Desc, srv.getBadWordsList()))
	response.Info.Mobile = user.Mobile
	response.Info.Date = utils.GetUnixTimeToDateTime(task.AddTime)

	var addressInfo market.ZMAddress
	model.DefaultMarket().Model(&market.ZMAddress{}).Debug().Where("id = ?", task.AddressId).First(&addressInfo)

	response.Info.Address = addressInfo.Name
	response.Info.Status = task.Status
	return
}

func (srv *MarketService) ApiDoMakeTaskData(taskData api2.MakeTaskDataRequest) (response *api2.MakeTaskDataResponse, apiErr api.Error) {
	if taskData.Title == "" || taskData.TaskDesc == "" || taskData.AddressId == 0 || taskData.TagId == 0 || taskData.UserId == 0 {
		return
	}
	var result bool
	response = &api2.MakeTaskDataResponse{}
	//判断用户发的次数，5分钟设置最大3次
	cacheCount := fmt.Sprintf("request_count:%d", taskData.UserId)
	count, _ := rxRedis.ClientDefault("web").Get(context.Background(), cacheCount).Result()
	countInt, _ := strconv.Atoi(count)
	if countInt >= 3 {
		return
	}

	var task market.ZMTask
	odb := model.DefaultMarket().Model(&market.ZMTask{}).Debug()
	task.TagId = taskData.TagId
	task.UserId = taskData.UserId
	task.Title = taskData.Title
	task.Desc = taskData.TaskDesc
	task.AddressId = taskData.AddressId
	task.Status = 1
	task.AddTime = utils.GetCurrentUnixTimestamp()
	affected := odb.Create(&task).RowsAffected
	if affected > 0 {
		result = true
		//键值增加1
		rxRedis.ClientDefault("web").Incr(context.Background(), cacheCount).Result()
		//设置5分钟的有效期
		rxRedis.ClientDefault("web").Expire(context.Background(), cacheCount, time.Duration(300)*time.Second)
		//单独的防止连续点击
		rxRedis.ClientDefault("web").SetNX(context.Background(), fmt.Sprintf("userPushTask_%d", task.UserId), 1, time.Duration(3)*time.Second)
	}
	response.Result = result
	return
}

func (srv *MarketService) ApiDoMakeOtherTaskData(taskData api2.MakeTaskOtherDataRequest) (response *api2.MakeTaskOtherDataResponse, apiErr api.Error) {
	if taskData.Mobile == "" || taskData.TaskDesc == "" || taskData.AddressId == 0 || taskData.TagId == 0 {
		return
	}
	response = &api2.MakeTaskOtherDataResponse{}
	s, _ := rxRedis.ClientDefault("web").Get(context.Background(), "userPushTask_lock").Result()
	atoInt, _ := strconv.Atoi(s)
	if atoInt > 1 {
		return
	}
	var task market.ZMTask
	odb := model.DefaultMarket().Model(&market.ZMTask{}).Debug()
	task.TagId = taskData.TagId
	task.UserId = 1712141769
	task.Title = "代发任务"
	task.Desc = taskData.TaskDesc
	task.AddressId = taskData.AddressId
	task.Mobile = taskData.Mobile
	task.Status = 1
	task.AddTime = utils.GetCurrentUnixTimestamp()
	affected := odb.Create(&task).RowsAffected
	var result bool
	if affected > 0 {
		result = true
		//单独的防止连续点击
		rxRedis.ClientDefault("web").SetNX(context.Background(), "userPushTask_lock", 1, time.Duration(5)*time.Second)
	}

	response.Result = result
	return
}

func (srv *MarketService) ApiCheckPushTask(userId int) (response *api2.CheckPushTaskResponse, apiErr api.Error) {
	response = &api2.CheckPushTaskResponse{}
	s, _ := rxRedis.ClientDefault("web").Get(context.Background(), fmt.Sprintf("userPushTask_%d", userId)).Result()
	var result bool
	if atoi, _ := strconv.Atoi(s); atoi < 1 {
		result = true
	}
	response.Info = result
	return
}

func (srv *MarketService) ApiDoMakeUserData(userData api2.MakeUserDataRequest) (response *api2.MakeUserDataResponse, apiErr api.Error) {
	if userData.Type < 0 || userData.Mobile == "" || userData.OpenId == "" {
		return
	}
	response = &api2.MakeUserDataResponse{}
	//判断用户的openID是否存在
	var userTemp market.ZMUser
	model.DefaultMarket().Model(&market.ZMUser{}).Debug().Where("open_id=?", userData.OpenId).First(&userTemp)
	if userTemp.UserId > 0 {
		return
	}

	var user market.ZMUser
	odb := model.DefaultMarket().Model(&market.ZMUser{}).Debug()
	var userIdTemp = utils.GetCurrentUnixTimestamp()
	user.UserId = userIdTemp
	user.OpenId = userData.OpenId
	if userData.Type == 2 {
		user.OpenId = "app-" + strconv.FormatInt(userIdTemp, 10)
	}
	user.Mobile = userData.Mobile
	user.Type = userData.Type
	if userData.NickName != "" && len(userData.NickName) > 1 {
		user.NickName = userData.NickName
	}
	if userData.HeadImg != "" && len(userData.HeadImg) > 10 {
		user.HeadUrl = userData.HeadImg
	}
	affected := odb.Create(&user).RowsAffected
	var result bool
	if affected > 0 {
		result = true
	}
	response.Result = result

	return
}

func (srv *MarketService) ApiGetWechatData(code string) (response *api2.WechatDataResponse, apiErr api.Error) {
	response = &api2.WechatDataResponse{}
	appId := viper.GetString("wechat.appid")
	secret := viper.GetString("wechat.secret")
	urlFormat := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", appId, secret, code)
	var wxInfo string

	resp, err := http.Get(urlFormat)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	wxInfo = string(body)

	response.Data = wxInfo
	return
}

func (srv *MarketService) ApiGetWxAccessToken() (response *api2.AccessTokenResponse, apiErr api.Error) {
	response = &api2.AccessTokenResponse{}
	appId := viper.GetString("wechat.appid")
	secret := viper.GetString("wechat.secret")
	urlFormat := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", appId, secret)
	var wxInfo string

	resp, err := http.Get(urlFormat)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	wxInfo = string(body)

	response.Data = wxInfo
	return
}

func (srv *MarketService) ApiGetWxUserPhoneNumber(photoData api2.MakePhotoDataRequest) (response *api2.MakePhotoDataResponse, apiErr api.Error) {
	if photoData.Token == "" || photoData.Code == "" {
		return
	}
	requestBody, err := json.Marshal(map[string]string{
		"code": photoData.Code,
	})
	if err != nil {
		return
	}

	// 发起 POST 请求
	resp, err := http.Post(fmt.Sprintf("https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=%s", photoData.Token), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// 读取响应数据
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var wxInfo string
	wxInfo = string(body)

	response.Data = wxInfo
	return
}

func (srv *MarketService) ApiGetWxPay(payData api2.WXPayDataRequest) (response *api2.WXPayDataResponse, apiErr api.Error) {
	response = &api2.WXPayDataResponse{}
	if payData.UserId <= 0 || payData.PayId <= 0 {
		return
	}
	//创建订单
	orderInfo := srv.createOrderData(payData, 1)
	//创建JSPay订单参数
	JSPayParam := srv.creatJsApi(orderInfo)
	response.Data = JSPayParam
	return
}

func (srv *MarketService) createOrderData(payData api2.WXPayDataRequest, orderType int) (orderInfo market.ZMOrder) {
	//根据用户id查询用户信息
	var userInfo market.ZMUser
	model.DefaultMarket().Model(&market.ZMUser{}).Debug().Where("user_id=?", payData.UserId).First(&userInfo)

	//根据payId查询相关价格
	var payItem market.ZMPay
	model.DefaultMarket().Model(&market.ZMPay{}).Debug().Where("id=?", payData.PayId).First(&payItem)

	tempOrderId := utils.GetCurrentMilliseconds()
	var order market.ZMOrder
	odb := model.DefaultMarket().Model(&market.ZMOrder{}).Debug()
	order.Name = payItem.Name
	order.OpenId = userInfo.OpenId
	order.UserId = userInfo.UserId
	order.OrderId = tempOrderId
	order.Type = orderType //1普通会员,2优选工匠, 3积分兑换
	order.CPrice = payItem.CPrice
	order.OPrice = payItem.OPrice
	order.Number = payItem.Number
	order.NumberExt = payItem.NumberExt
	order.Status = 0 //支付状态,1支付完成,0待支付
	affected := odb.Create(&order).RowsAffected
	if affected > 0 {
		orderInfo = order
	}
	return
}

func (srv *MarketService) creatJsApi(orderInfo market.ZMOrder) (JSPayParam api2.JSPayParam) {

	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := apiUtils.LoadPrivateKeyWithPath("/data/web/market-api/run/wx_market_cert/apiclient_key.pem")
	if err != nil {
		log.Fatal("load merchant private key error")
	}
	appId := viper.GetString("wechat.appid")
	mchId := viper.GetString("wechat.mch-id")
	mchCert := viper.GetString("wechat.mch-id")
	mchIv3 := viper.GetString("wechat.mch-iv3")

	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(mchId, mchCert, mchPrivateKey, mchIv3),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		log.Fatalf("new wechat pay client err:%s", err)
	}

	description := "千皓优选（" + orderInfo.Name + "）"
	var cPrice int64 = int64(orderInfo.CPrice)
	svc := jsapi.JsapiApiService{Client: client}
	// 得到prepay_id，以及调起支付所需的参数和签名
	resp, _, err := svc.PrepayWithRequestPayment(ctx,
		jsapi.PrepayRequest{
			Appid:       core.String(appId),
			Mchid:       core.String(mchId),
			Description: core.String(description),
			OutTradeNo:  core.String(strconv.FormatInt(orderInfo.OrderId, 10)),
			Attach:      core.String("千皓优选用工好选择"),
			NotifyUrl:   core.String("https://market.58haha.com/api/wechat/pay/notice"),
			Amount: &jsapi.Amount{
				Total: core.Int64(cPrice),
			},
			Payer: &jsapi.Payer{
				Openid: core.String(orderInfo.OpenId),
			},
		},
	)
	if err == nil {
		log.Println(resp)
		JSPayParam.PrepayID = *resp.PrepayId
		JSPayParam.Appid = *resp.Appid
		JSPayParam.TimeStamp = *resp.TimeStamp
		JSPayParam.NonceStr = *resp.NonceStr
		JSPayParam.Package = *resp.Package
		JSPayParam.SignType = *resp.SignType
		JSPayParam.PaySign = *resp.PaySign
	} else {
		log.Println(err)
	}

	return
}

func (srv *MarketService) ApiGetWxOpenPay(openData api2.OpenGoodPayRequest) (response *api2.OpenGoodPayResponse, apiErr api.Error) {
	response = &api2.OpenGoodPayResponse{}

	//验证数据
	if openData.UserID < 0 || openData.PayId <= 0 || openData.OpenID == "" || openData.UserImage == "" || openData.AddressId <= 0 || openData.NickName == "" || openData.UserSelf == "" || openData.TagID <= 0 || openData.IsAgree <= 0 {
		return
	}
	//开通优选工匠的前置业务处理
	res := srv.openPayPreCreatData(openData)
	if res == true {
		//创建订单
		var payData api2.WXPayDataRequest
		payData.UserId = openData.UserID
		payData.PayId = openData.PayId
		orderInfo := srv.ApiCreateOrderData(payData, 2)
		//创建JSPay订单参数
		JSPayParam := srv.creatJsApi(orderInfo)
		response.Data = JSPayParam
	}

	return
}

func (srv *MarketService) openPayPreCreatData(openData api2.OpenGoodPayRequest) (result bool) {
	//根据用户id查询用户信息
	var userInfo market.ZMUser
	model.DefaultMarket().Model(&market.ZMUser{}).Debug().Where("user_id=?", openData.UserID).First(&userInfo)
	if userInfo.Id <= 0 {
		return
	}
	var lastTime int64 = utils.GetCurrentUnixTimestamp()
	var userTemp market.ZMUser
	userTemp.NickName = openData.NickName
	userTemp.RealName = openData.NickName
	userTemp.HeadUrl = openData.UserImage
	userTemp.TagId = openData.TagID
	userTemp.LastTime = lastTime
	obj := model.DefaultMarket().Model(&market.ZMUser{}).Debug().Where("user_id=?", openData.UserID)
	affected := obj.Updates(&userTemp).RowsAffected
	//如果更新成功，才可以处理用户附属信息
	if affected > 0 {
		//查询是否存在用户的附属信息
		var userExtInfo market.ZMUserExt
		model.DefaultMarket().Model(&market.ZMUserExt{}).Debug().Where("user_id=?", openData.UserID).First(&userExtInfo)
		var userExtCreateData market.ZMUserExt
		userExtCreateData.UserId = int64(openData.UserID)
		userExtCreateData.Address = openData.Address
		userExtCreateData.AddressId = openData.AddressId
		tempDesc := utils.ClearMobileText(openData.UserSelf)
		userExtCreateData.Desc = tempDesc
		if len(openData.UserCase) > 0 {
			caseInfo, _ := json.Marshal(openData.UserCase)
			userExtCreateData.Demo = string(caseInfo)
		}
		userExtCreateData.IsAgree = openData.IsAgree
		userExtCreateData.LastTime = lastTime
		if userExtInfo.Id > 0 {
			//更新操作
			affected2 := model.DefaultMarket().Model(&market.ZMUserExt{}).Debug().Where("user_id=?", openData.UserID).Updates(&userExtCreateData).RowsAffected
			if affected2 > 0 {
				result = true
			}
		} else {
			//添加操作
			affected1 := model.DefaultMarket().Model(&market.ZMUserExt{}).Debug().Create(&userExtCreateData).RowsAffected
			if affected1 > 0 {
				result = true
			}
		}
	}

	return
}

func (srv *MarketService) ApiCreateOrderData(payData api2.WXPayDataRequest, orderType int) (orderInfo market.ZMOrder) {
	//根据用户id查询用户信息
	var userInfo market.ZMUser
	model.DefaultMarket().Model(&market.ZMUser{}).Debug().Where("user_id=?", payData.UserId).First(&userInfo)

	//根据payId查询相关价格
	var payItem market.ZMPay
	model.DefaultMarket().Model(&market.ZMPay{}).Debug().Where("id=?", payData.PayId).First(&payItem)

	tempOrderId := utils.GetCurrentMilliseconds()
	var order market.ZMOrder
	odb := model.DefaultMarket().Model(&market.ZMOrder{}).Debug()
	order.Name = payItem.Name
	order.OpenId = userInfo.OpenId
	order.UserId = userInfo.UserId
	order.OrderId = tempOrderId
	order.Type = orderType //1普通会员,2优选工匠, 3积分兑换
	order.CPrice = payItem.CPrice
	order.OPrice = payItem.OPrice
	order.Number = payItem.Number
	order.NumberExt = payItem.NumberExt
	order.Status = 0 //支付状态,1支付完成,0待支付
	affected := odb.Create(&order).RowsAffected
	if affected > 0 {
		orderInfo = order
	}
	return
}

// ApiDealWxPayCallback 微信支付成功毁掉处理
func (srv *MarketService) ApiDealWxPayCallback(c *gin.Context) (notifyReq *notify.Request, err error) {
	ctx := c             //这个参数是context.Background()
	request := c.Request //这个值是*http.Request
	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := apiUtils.LoadPrivateKeyWithPath("/data/web/market-api/run/wx_market_cert/apiclient_key.pem")
	if err != nil {
		log.Println("load merchant private key error")
		return nil, err
	}

	mchId := viper.GetString("wechat.mch-id")
	mchCert := viper.GetString("wechat.mch-id")
	mchIv3 := viper.GetString("wechat.mch-iv3")
	// 1. 使用 `RegisterDownloaderWithPrivateKey` 注册下载器
	err = downloader.MgrInstance().RegisterDownloaderWithPrivateKey(ctx, mchPrivateKey, mchCert, mchId, mchIv3)
	if err != nil {
		return nil, err
	}
	// 2. 获取商户号对应的微信支付平台证书访问器
	certificateVisitor := downloader.MgrInstance().GetCertificateVisitor(mchId)
	// 3. 使用证书访问器初始化 `notify.Handler`
	handler := notify.NewNotifyHandler(mchIv3, verifiers.NewSHA256WithRSAVerifier(certificateVisitor))

	transaction := new(payments.Transaction)
	notifyReq, err = handler.ParseNotifyRequest(ctx, request, transaction)
	// 如果验签未通过，或者解密失败
	if err != nil {
		fmt.Println(err)
		//return
	}
	// 处理通知内容
	fmt.Println(notifyReq.Summary)
	fmt.Println(transaction.TransactionId)
	// 如果验签未通过，或者解密失败
	if err != nil {
		return nil, err
	}

	return notifyReq, nil
}

// ApiDealUserPaySuccess 将解密结果进行处理
func (srv *MarketService) ApiDealUserPaySuccess(notifyReq *notify.Request, result map[string]interface{}) {
	//判断是否支付成功
	if notifyReq.EventType == "TRANSACTION.SUCCESS" && notifyReq.ResourceType == "encrypt-resource" {
		//支付成功处理订单状态
		if orderId, ok := result["out_trade_no"]; ok {
			fmt.Println("订单id 存在，值为：", orderId)
			//处理订单信息--更新订单的支付时间，支付状态
			var orderTemp market.ZMOrder
			ob := model.DefaultMarket().Model(&market.ZMOrder{}).Debug().Where("order_id=?", orderId)
			ob.Find(&orderTemp)
			//判断是否为已支付状态,只有是未支付成功的状态才可操作
			if orderTemp.Id > 0 && orderTemp.Status == 0 {
				var order market.ZMOrder
				order.Status = 1 //支付完成
				order.PaymentNumber = result["transaction_id"].(string)
				order.PayTime = utils.GetCurrentDateTime()
				obj := model.DefaultMarket().Model(&market.ZMOrder{}).Debug().Where("order_id=?", orderId)
				affected := obj.Updates(&order).RowsAffected
				//如果更新成功，才可以处理用户信息
				if affected > 0 {
					//根据订单类型进行业务处理-类型,1普通会员,2优选工匠,3积分兑换'
					//处理普通会员业务逻辑
					if orderTemp.Type == 1 {
						//处理用户信息--增加会员标识，标识有效期
						var userTemp market.ZMUser
						var user market.ZMUser
						obu := model.DefaultMarket().Model(&market.ZMUser{}).Debug().Where("user_id=? and open_id = ?", orderTemp.UserId, orderTemp.OpenId)
						obu.Find(&userTemp)
						//当前时间
						currentYMD, _ := strconv.Atoi(utils.GetCurrentDateYMD())
						//判断用户是否存在有效期
						var total int = orderTemp.Number + orderTemp.NumberExt
						if userTemp.MemberLimit > 0 && userTemp.MemberLimit >= currentYMD {
							user.MemberLimit = utils.CalculateAfterDate(userTemp.MemberLimit, total) //会员截止日期
						} else {
							//已失效的会员有效期进行重置
							user.MemberLimit = utils.CalculateAfterDate(currentYMD, total) //会员截止日期
						}
						user.IsMember = 1 //标记为会员
						obu.Updates(&user)
					}
					//处理优选工匠业务逻辑
					if orderTemp.Type == 2 {
						//处理用户信息--增加会员标识，标识有效期
						var userTemp market.ZMUser
						obu := model.DefaultMarket().Model(&market.ZMUser{}).Debug().Where("user_id=? and open_id = ?", orderTemp.UserId, orderTemp.OpenId)
						userTemp.IsBest = 1
						//当前时间
						currentYMD, _ := strconv.Atoi(utils.GetCurrentDateYMD())
						userTemp.BestLimit = utils.CalculateAfterDate(currentYMD, 365)
						userTemp.LastTime = utils.GetCurrentUnixTimestamp()
						obu.Updates(&userTemp)
					}
				}
			}
		}
	}
}

func (srv *MarketService) ApiGetWxPayCallback() (response *api2.ZMAddressListResponse, apiErr api.Error) {
	response = &api2.ZMAddressListResponse{}
	model.DefaultMarket().Model(&market.ZMAddress{}).Where("is_deleted = 0 and is_hot = 1 and parent_id > 0").Debug().
		Order("sort desc").Order("id asc").Limit(5).
		Find(&response.List)
	return
}

func (srv *MarketService) ApiGetWxPayCancel(cancelData api2.WXCancelPayDataRequest) (response *api2.WXCancelPayDataResponse, apiErr api.Error) {
	response = &api2.WXCancelPayDataResponse{}
	if cancelData.UserId < 0 {
		return
	}
	//获取用户最后一条没有支付的订单
	var orderTemp market.ZMOrder
	obj := model.DefaultMarket().Model(&market.ZMOrder{}).Debug().Where("user_id=? and status = 0", cancelData.UserId)
	obj.Order("id desc").First(&orderTemp)
	var result bool
	if orderTemp.Id > 0 {
		//调用微信关单
		srv.JsApiCloseOrder(orderTemp)

		var order market.ZMOrder
		order.Status = -1 //取消支付
		affected := model.DefaultMarket().Model(&market.ZMOrder{}).Debug().Where("id = ?", orderTemp.Id).Updates(&order).RowsAffected
		if affected > 0 {
			result = true
		}
	}
	response.Data = result
	return
}

func (srv *MarketService) JsApiCloseOrder(orderInfo market.ZMOrder) (JSPayParam api2.JSPayParam) {
	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := apiUtils.LoadPrivateKeyWithPath("/data/web/market-api/run/wx_market_cert/apiclient_key.pem")
	if err != nil {
		log.Fatal("load merchant private key error")
	}

	mchId := viper.GetString("wechat.mch-id")
	mchCert := viper.GetString("wechat.mch-id")
	mchIv3 := viper.GetString("wechat.mch-iv3")

	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(mchId, mchCert, mchPrivateKey, mchIv3),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		log.Fatalf("new wechat pay client err:%s", err)
	}
	svc := jsapi.JsapiApiService{Client: client}
	result, err := svc.CloseOrder(ctx,
		jsapi.CloseOrderRequest{
			OutTradeNo: core.String(strconv.FormatInt(orderInfo.OrderId, 10)),
			Mchid:      core.String(mchId),
		},
	)

	if err != nil {
		// 处理错误
		log.Printf("关闭订单call CloseOrder err:%s", err)
	} else {
		// 处理返回结果
		log.Printf("关闭订单status=%d", result.Response.StatusCode)
	}

	return
}

func (srv *MarketService) ApiGetWxPayRefunds(params api2.WXRefundsPayDataRequest) (message string, code int) {
	fmt.Printf("退费操作的参数:%#v \n", params)
	//查询当前订单是否已支付成功的状态
	var orderTemp market.ZMOrder
	obj := model.DefaultMarket().Model(&market.ZMOrder{}).Debug().Where("order_id=? and status = 1 and is_deleted = 0", params.OrderId)
	obj.Order("id desc").First(&orderTemp)
	//说明存在可以退费的数据
	if orderTemp.Id > 0 {
		//进行退费操作
		err, status := srv.refundsApiServiceCreate(&orderTemp)
		if err != nil || status != http.StatusOK {
			message = err.Error()
			return message, status
		}

		//退款成功，判断订单类型，处理会员情况
		srv.dealRefundsAfterData(&orderTemp)

		message = "success"
		code = http.StatusOK
	}
	return
}

func (srv *MarketService) refundsApiServiceCreate(orderInfo *market.ZMOrder) (err error, status int) {
	fmt.Printf("退费操作的参数:%#v \n", orderInfo)
	//return
	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := apiUtils.LoadPrivateKeyWithPath("/data/web/market-api/run/wx_market_cert/apiclient_key.pem")
	if err != nil {
		log.Fatal("load merchant private key error")
	}

	mchId := viper.GetString("wechat.mch-id")
	mchCert := viper.GetString("wechat.mch-id")
	mchIv3 := viper.GetString("wechat.mch-iv3")
	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(mchId, mchCert, mchPrivateKey, mchIv3),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		log.Printf("new wechat pay client err:%s", err)
	}

	svc := refunddomestic.RefundsApiService{Client: client}
	resp, result, err := svc.Create(ctx,
		refunddomestic.CreateRequest{
			OutRefundNo: core.String(strconv.FormatInt(orderInfo.OrderId, 10)),
			OutTradeNo:  core.String(strconv.FormatInt(orderInfo.OrderId, 10)),
			Reason:      core.String("客服处理退款_" + orderInfo.Name),
			NotifyUrl:   core.String("https://weixin.qq.com"),
			Amount: &refunddomestic.AmountReq{
				Currency: core.String("CNY"),
				Refund:   core.Int64(int64(orderInfo.CPrice)), //必填，退款金额，单位为分
				Total:    core.Int64(int64(orderInfo.CPrice)), //必填，原支付交易订单金额
			},
		},
	)
	if err != nil {
		// 处理错误
		log.Printf("call Create err:%s", err)
	} else {
		// 处理返回结果
		log.Printf("status=%d resp=%s", result.Response.StatusCode, resp)
	}

	return err, result.Response.StatusCode
}

func (srv *MarketService) dealRefundsAfterData(orderInfo *market.ZMOrder) {
	//将订单状态改成已退费
	var order market.ZMOrder
	order.Status = -3                                  //已退费
	order.RefundTime = utils.GetCurrentUnixTimestamp() //退费时间
	model.DefaultMarket().Model(&market.ZMOrder{}).Debug().Where("id = ?", orderInfo.Id).Updates(&order)

	//判断订单类型
	orderType := orderInfo.Type //1普通会员,2优选工匠
	//查询用户信息
	var userTemp market.ZMUser
	obj := model.DefaultMarket().Model(&market.ZMUser{}).Debug().Where("user_id=? and status = 1 and is_deleted = 0", orderInfo.UserId)
	obj.First(&userTemp)

	fmt.Printf("用户信息：%#v \n", userTemp)
	//用户模型
	var user market.ZMUser
	//如果是普通会员
	if orderType == 1 {
		//判断会员是否过期
		today := utils.GetCurrentDateYMD()
		todayInt, _ := strconv.Atoi(today)
		fmt.Printf("%#v \n", todayInt)
		if todayInt > userTemp.MemberLimit {
			user.IsMember = 0    //已退费,就不是会员了
			user.MemberLimit = 0 //已退费,就不是会员了
			model.DefaultMarket().Model(&market.ZMUser{}).Debug().Where("user_id = ?", userTemp.UserId).Updates(&user)
			return
		}
		if todayInt < userTemp.MemberLimit {
			//判断赠送的总天数
			totalDay := orderInfo.Number + orderInfo.NumberExt
			//当前的有效期减去时间,就是剩余的会员期限
			surplus := utils.CalculateBeforeDate(userTemp.MemberLimit, totalDay)
			surplusInt, _ := strconv.Atoi(surplus)
			fmt.Printf("surplusInt:%#v \n", surplusInt)
			fmt.Printf("todayInt:%#v \n", todayInt)
			//如果剩余的有效期，小于今天，则证明，会员到期
			if surplusInt <= todayInt {
				user.IsMember = 0    //已退费,就不是会员了
				user.MemberLimit = 0 //已退费,就不是会员了
				model.DefaultMarket().Model(&market.ZMUser{}).Debug().Where("user_id = ?", userTemp.UserId).Updates(&user)
				return
			} else {
				//如果剩余有效期大于今天，则将有效期的截止日期跟新为当前剩余的有效期
				user.MemberLimit = surplusInt //更新有效期
				model.DefaultMarket().Model(&market.ZMUser{}).Debug().Where("user_id = ?", userTemp.UserId).Updates(&user)
				return
			}
		}
	}

	//如果是优选工匠
	if orderType == 2 {
		user.IsBest = 0    //已退费，改为非优选工匠
		user.BestLimit = 0 //已退费，改为非优选工匠
		model.DefaultMarket().Model(&market.ZMUser{}).Debug().Where("user_id = ?", orderInfo.UserId).Updates(&order)
	}
}
