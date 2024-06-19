package model

type ChineseBookName struct {
	Id         int    `json:"-"`
	CategoryId int    `json:"category_id"`
	Name       string `json:"name"`
	SSort      int    `json:"s_sort"`
}

func (ChineseBookName) TableName() string {
	return "s_chinese_name"
}

type ChineseBook struct {
	Id       int    `json:"-"`
	BookId   string `json:"book_id"`
	Title    string `json:"title"`
	Icon     string `json:"icon"`
	Type     uint8  `json:"type"`
	Position uint8  `json:"position"`
}

func (ChineseBook) TableName() string {
	return "s_chinese_picture"
}

type ChineseBookInfo struct {
	Id       int    `json:"id"`
	BookId   string `json:"book_id"`
	Mp3      string `json:"mp3"`
	Pic      string `json:"pic"`
	Position uint8  `json:"position"`
}

func (ChineseBookInfo) TableName() string {
	return "s_chinese_picture_info"
}
