package common

import "api/model"

type ChineseBookNavResponse struct {
	List []model.ChineseBookName `json:"list"`
}

type ChineseBookResponse struct {
	ListModel []model.ChineseBook   `json:"-"`
	List      []ResponseChineseBook `json:"list"`
	Total     int64                 `json:"total"`
	Page      int                   `json:"page"`
	TotalPage float64               `json:"total_page"`
}

type ChineseBookInfoResponse struct {
	Info []model.ChineseBookInfo `json:"info"`
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
