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

type ChineseBookAlbum struct {
	Id       int    `json:"-"`
	BookId   string `json:"book_id"`
	Title    string `json:"title"`
	Icon     string `json:"icon"`
	Position uint8  `json:"position"`
}

func (ChineseBookAlbum) TableName() string {
	return "s_chinese_picture_album"
}

type ChineseAlbumInfo struct {
	Id       int    `json:"id"`
	BookId   string `json:"book_id"`
	Mp3      string `json:"mp3"`
	Pic      string `json:"pic"`
	Position uint8  `json:"position"`
	Title    string `json:"title"`
	Desc     string `json:"desc"`
	Duration string `json:"duration"`
}

func (ChineseAlbumInfo) TableName() string {
	return "s_chinese_picture_album_info"
}

type ChengYU struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Pinyin  string `json:"pinyin"`
	Explain string `json:"explain"`
	Source  string `json:"source"`
	Usage   string `json:"usage"`
	Example string `json:"example"`
	Near    string `json:"near"`
	Antonym string `json:"antonym"`
	Analyse string `json:"analyse"`
	Story   string `json:"story"`
	Level   uint8  `json:"level"`
}

func (ChengYU) TableName() string {
	return "s_chengyu"
}

type Poetry struct {
	Id         int    `json:"-"`
	PoetryId   int    `json:"poetry_id"`
	Title      string `json:"title"`
	GradeId    uint8  `json:"grade_id"`
	Grade      string `json:"grade"`
	GradeLevel uint8  `json:"grade_level"`
	Author     string `json:"author"`
	Dynasty    string `json:"dynasty"`
	Mp3        string `json:"mp3"`
	Content    string `json:"content"`
	Info       string `json:"info"`
}

func (Poetry) TableName() string {
	return "s_school_poetry"
}
