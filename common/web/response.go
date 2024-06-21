package web

import (
	"api/model/web"
)

type ChineseBookNavResponse struct {
	List []web.ChineseBookName `json:"list"`
}

type ChineseBookResponse struct {
	ListModel []web.ChineseBook     `json:"-"`
	List      []ResponseChineseBook `json:"list"`
	Total     int64                 `json:"total"`
	Page      int                   `json:"page"`
	TotalPage float64               `json:"total_page"`
}

type ChineseBookInfoResponse struct {
	Info []web.ChineseBookInfo `json:"info"`
}

type ResponseChineseBook struct {
	Id        int    `json:"-"`
	BookId    string `json:"book_id"`
	Title     string `json:"title"`
	Icon      string `json:"icon"`
	Level     uint8  `json:"-"`
	Position  uint8  `json:"-"`
	BookCount string `json:"book_count"`
}

type ResponseBookInfoCount struct {
	BookId    string `json:"book_id"`
	BookCount string `json:"book_count"`
}

type ChineseAlbumResponse struct {
	List      []web.ChineseBookAlbum `json:"list"`
	Total     int64                  `json:"total"`
	Page      int                    `json:"page"`
	TotalPage float64                `json:"total_page"`
}

type ChineseAlbumListInfoResponse struct {
	List      []web.ChineseAlbumInfo `json:"list"`
	Total     int64                  `json:"total"`
	Page      int                    `json:"page"`
	TotalPage float64                `json:"total_page"`
}

type ChineseAlbumInfoResponse struct {
	Info web.ChineseAlbumInfo `json:"info"`
}

type ChineseChengYuResponse struct {
	List      []web.ChengYU `json:"list"`
	Total     int64         `json:"total"`
	Page      int           `json:"page"`
	TotalPage float64       `json:"total_page"`
}

type ChineseChengYuInfoResponse struct {
	Info CYdATA `json:"info"`
}

type CYdATA struct {
	Id        int      `json:"id"`
	Title     string   `json:"title"`
	Pinyin    string   `json:"pinyin"`
	Explain   string   `json:"explain"`
	Source    string   `json:"source"`
	Usage     string   `json:"usage"`
	Example   string   `json:"example"`
	Near      string   `json:"near"`
	Antonym   string   `json:"antonym"`
	Analyse   string   `json:"analyse"`
	Story     string   `json:"story"`
	Level     uint8    `json:"level"`
	StoryList []string `json:"story_list"`
}

type SchoolPoetryListResponse struct {
	List      []ResponseSchoolPoetry `json:"list"`
	Total     int64                  `json:"total"`
	Page      int                    `json:"page"`
	TotalPage float64                `json:"total_page"`
}

type ResponseSchoolPoetry struct {
	Id         int    `json:"id"`
	PoetryId   int    `json:"poetry_id"`
	Title      string `json:"title"`
	GradeId    uint8  `json:"grade_id" `
	Grade      string `json:"grade" `
	GradeLevel uint8  `json:"grade_level" `
	Author     string `json:"author" `
	Dynasty    string `json:"dynasty"`
}
