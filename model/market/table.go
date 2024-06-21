package market

type ZMAddress struct {
	Id        int          `json:"id"`
	Name      string       `json:"name"`
	Sort      int          `json:"sort"`
	ParentId  int          `json:"parent_id"`
	IsHot     int          `json:"is_hot"`
	ChildList []*ZMAddress `json:"list"`
}

func (ZMAddress) TableName() string {
	return "zm_address"
}

type ZMBanner struct {
	Id     int    `json:"-"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	Status int    `json:"status"`
}

func (ZMBanner) TableName() string {
	return "zm_banner"
}

type ZMOrder struct {
	Id            int     `json:"-"`
	Name          string  `json:"name"`
	OpenId        string  `json:"open_id"`
	UserId        int64   `json:"user_id"`
	OrderId       int64   `json:"order_id"`
	Type          int     `json:"type"`
	CPrice        float64 `json:"c_price"`
	OPrice        float64 `json:"o_price"`
	Number        int     `json:"number"`
	NumberExt     int     `json:"number_ext"`
	Status        int     `json:"status"`
	PayTime       string  `json:"pay_time"`
	PaymentNumber string  `json:"payment_number"`
	RefundTime    int64   `json:"refund_time"`
}

func (ZMOrder) TableName() string {
	return "zm_order"
}

type ZMPay struct {
	Id        int     `json:"id"`
	Name      string  `json:"name"`
	CPrice    float64 `json:"c_price"`
	OPrice    float64 `json:"o_price"`
	Number    int     `json:"number"`
	NumberExt int     `json:"number_ext"`
	Checked   bool    `json:"checked"`
	Type      int     `json:"-"`
}

func (ZMPay) TableName() string {
	return "zm_pay"
}

type ZMTags struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Icon   string `json:"icon"`
	Status int    `json:"-"`
}

func (ZMTags) TableName() string {
	return "zm_tags"
}

type ZMTask struct {
	Id        int    `json:"-"`
	Title     string `json:"title"`
	Desc      string `json:"desc"`
	UserId    int64  `json:"user_id"`
	TagId     int    `json:"tag_id"`
	Status    int    `json:"status"`
	AddTime   int64  `json:"add_time"`
	Address   string `json:"address"`
	Mobile    string `json:"mobile"`
	AddressId int    `json:"address_id"`
	CreatedAt string `json:"created_at" gorm:"-"`
}

func (ZMTask) TableName() string {
	return "zm_task"
}

type ZMUser struct {
	Id          int    `json:"-"`
	UserId      int64  `json:"user_id"`
	OpenId      string `json:"open_id"`
	NickName    string `json:"nick_name"`
	RealName    string `json:"real_name"`
	HeadUrl     string `json:"head_url"`
	Mobile      string `json:"mobile"`
	TagId       int    `json:"tag_id"`
	ParentId    int    `json:"parent_id"`
	IsBest      int    `json:"is_best"`
	BestLimit   int    `json:"best_limit"`
	IsMember    int    `json:"is_member"`
	MemberLimit int    `json:"member_limit"`
	Type        int    `json:"type"`
	LastTime    int64  `json:"last_time"`
}

func (ZMUser) TableName() string {
	return "zm_user"
}

type ZMUserExt struct {
	Id        int    `json:"-"`
	UserId    int64  `json:"user_id"`
	Address   string `json:"address"`
	Desc      string `json:"desc"`
	Demo      string `json:"demo"`
	IsAgree   int    `json:"is_agree"`
	AddressId int    `json:"address_id"`
	ViewCount int    `json:"view_count"`
	LastTime  int64  `json:"last_time"`
}

func (ZMUserExt) TableName() string {
	return "zm_user_ext"
}

type ZMBadWords struct {
	Id   int    `json:"-"`
	Name string `json:"desc"`
}

func (ZMBadWords) TableName() string {
	return "zm_bad_words"
}
