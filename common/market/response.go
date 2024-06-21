package market

import "api/model/market"

type ZMAddressHotResponse struct {
	List []market.ZMAddress `json:"list"`
}

type ZMAddressListResponse struct {
	List []market.ZMAddress `json:"list"`
}

type AddressChildListResponse struct {
	List []FormatAddressInfo `json:"list"`
}

type FormatAddressInfo struct {
	CityId   int    `json:"city_id"`
	CityName string `json:"city_name"`
}

type ZMUserResponse struct {
	Info market.ZMUser `json:"info"`
}

type ZMUserExtResponse struct {
	Info market.ZMUserExt `json:"info"`
}

type ZMBannerListResponse struct {
	List []market.ZMBanner `json:"list"`
}

type ZMBannerListNewResponse struct {
	List []string `json:"list"`
}

type ZMTagsListResponse struct {
	List []market.ZMTags `json:"list"`
}

type ZMTagsListSelectResponse struct {
	List []string `json:"list"`
}

type ZMPayListResponse struct {
	List []FormatData `json:"list"`
}
type FormatData struct {
	Value     string  `json:"value"`
	Name      string  `json:"name"`
	Checked   bool    `json:"checked"`
	Number    int     `json:"number"`
	NumberExt int     `json:"number_ext"`
	TotalFee  float64 `json:"total_fee"`
}

type ZMPayResponse struct {
	Info market.ZMPay `json:"info"`
}

type ZMGoodMemberListResponse struct {
	List  []MemberData `json:"list"`
	Count int64        `json:"count"`
}

type ZMGoodUserInfoResponse struct {
	Info MemberData `json:"info"`
}

type MemberData struct {
	Id          int    `json:"id"`
	UserId      int64  `json:"user_id"`
	OpenId      string `json:"open_id"`
	NickName    string `json:"nick_name"`
	RealName    string `json:"real_name"`
	HeadUrl     string `json:"head_url"`
	Mobile      string `json:"mobile"`
	TagId       int    `json:"tag_id"`
	TagName     string `json:"tag_name"`
	Desc        string `json:"desc"`
	IsBest      int    `json:"is_best"`
	IsMember    int    `json:"is_member"`
	MemberLimit int    `json:"member_limit"`
	ViewCount   int    `json:"view_count"`
	Demo        string `json:"demo"`
	Address     string `json:"address"`
}

type ZMTaskListResponse struct {
	List  []FormatTaskData `json:"list"`
	Count int64            `json:"count"`
}

type ZMTaskInfoResponse struct {
	Info FormatTaskData `json:"info"`
}

type FormatTaskData struct {
	Id      int    `json:"id"`
	TagId   int    `json:"tag_id"`
	UserId  int64  `json:"user_id"`
	TagName string `json:"tag_name"`
	Desc    string `json:"desc"`
	Mobile  string `json:"mobile"`
	Date    string `json:"date"`
	Address string `json:"address"`
	Status  int    `json:"status"`
	IsBest  int    `json:"is_best"`
}

type CheckPushTaskResponse struct {
	Info bool `json:"info"`
}

type MemberUpdateDataRequest struct {
	UserId    int    `json:"user_id"`
	AddressId int    `json:"address_id"`
	Address   string `json:"address"`
	NickName  string `json:"nick_name"`
	HeadUrl   string `json:"head_url"`
	Mobile    string `json:"mobile"`
}

type MemberUpdateDataResponse struct {
	Result bool `json:"result"`
}

type TaskUpdateStatusRequest struct {
	TaskId int `json:"task_id"`
	Status int `json:"status"`
}

type TaskUpdateStatusResponse struct {
	Result bool `json:"result"`
}

type MakeTaskDataRequest struct {
	TaskDesc  string `json:"task_desc"`
	TagId     int    `json:"tag_id"`
	UserId    int64  `json:"user_id"`
	AddressId int    `json:"address_id"`
	Title     string `json:"title"`
}

type MakeTaskDataResponse struct {
	Result bool `json:"result"`
}

type MakeTaskOtherDataRequest struct {
	TaskDesc  string `json:"task_desc"`
	TagId     int    `json:"tag_id"`
	AddressId int    `json:"address_id"`
	Mobile    string `json:"mobile"`
}

type MakeTaskOtherDataResponse struct {
	Result bool `json:"result"`
}

type MakeUserDataRequest struct {
	Mobile   string `json:"mobile"`
	OpenId   string `json:"open_id"`
	Type     int    `json:"type"`
	NickName string `json:"nick_name"`
	HeadImg  string `json:"head_img"`
}

type MakeUserDataResponse struct {
	Result bool `json:"result"`
}

type WechatDataResponse struct {
	Data string `json:"data"`
}

type AccessTokenResponse struct {
	Data string `json:"data"`
}

type MakePhotoDataRequest struct {
	Code  string `json:"code"`
	Token string `json:"token"`
}

type MakePhotoDataResponse struct {
	Data string `json:"data"`
}

type WXPayDataRequest struct {
	UserId int `json:"user_id"`
	PayId  int `json:"pay_id"`
}

type WXPayDataResponse struct {
	Data JSPayParam `json:"data"`
}

type JSPayParam struct {
	PrepayID  string `json:"PrepayId"`
	Appid     string `json:"Appid"`
	TimeStamp string `json:"TimeStamp"`
	NonceStr  string `json:"NonceStr"`
	Package   string `json:"Package"`
	SignType  string `json:"SignType"`
	PaySign   string `json:"PaySign"`
}

type OpenGoodPayRequest struct {
	UserID    int      `json:"user_id"`
	OpenID    string   `json:"open_id"`
	UserImage string   `json:"user_image"`
	Address   string   `json:"address"`
	AddressId int      `json:"address_id"`
	NickName  string   `json:"nick_name"`
	UserSelf  string   `json:"user_self"`
	TagID     int      `json:"tag_id"`
	PayId     int      `json:"pay_id"`
	IsAgree   int      `json:"is_agree"`
	UserCase  []string `json:"user_case"`
}

type OpenGoodPayResponse struct {
	Data JSPayParam `json:"data"`
}

type WXCancelPayDataRequest struct {
	UserId int `json:"user_id"`
}

type WXCancelPayDataResponse struct {
	Data bool `json:"data"`
}

type WXRefundsPayDataRequest struct {
	OrderId int `json:"order_id"`
}
