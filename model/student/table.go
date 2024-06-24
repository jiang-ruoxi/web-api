package student

type SPoetryPicture struct {
	Id     int    `json:"id"`
	BookId int    `json:"book_id"`
	Title  string `json:"title"`
	Icon   string `json:"icon"`
	TypeId int    `json:"type_id"`
	Author string `json:"author"`
}

func (SPoetryPicture) TableName() string {
	return "s_poetry_picture"
}

type SPoetryPictureInfo struct {
	Id       int    `json:"id"`
	BookId   int    `json:"book_id"`
	Cn       string `json:"cn"`
	Pic      string `json:"pic"`
	Mp3      string `json:"mp3"`
	Position int    `json:"position"`
}

func (SPoetryPictureInfo) TableName() string {
	return "s_poetry_picture_info"
}

type SEnglishPicture struct {
	Id       int    `json:"id"`
	BookId   string `json:"book_id"`
	Title    string `json:"title"`
	Icon     string `json:"icon"`
	Level    int    `json:"level"`
	Position int    `json:"position"`
}

func (SEnglishPicture) TableName() string {
	return "s_english_picture"
}

type SEnglishPictureInfo struct {
	Id       int    `json:"id"`
	BookId   string `json:"book_id"`
	Pic      string `json:"pic"`
	Mp3      string `json:"mp3"`
	En       string `json:"en"`
	Zh       string `json:"zh"`
	Position int    `json:"position"`
}

func (SEnglishPictureInfo) TableName() string {
	return "s_english_picture_info"
}
