package api_service

import (
	"api/common"
	api2 "api/common/web"
	"api/model"
	"api/model/web"
	"github.com/jiang-ruoxi/gopkg/server/api"
	"math"
	"sort"
)

func NewChineseService() *ChineseService {
	return &ChineseService{}
}

type ChineseService struct {
}

func (srv *ChineseService) ChineseGetNavList() (response *api2.ChineseBookNavResponse, apiErr api.Error) {
	response = &api2.ChineseBookNavResponse{}
	model.DefaultWeb().Model(&web.ChineseBookName{}).Where("status = 1").Debug().
		Order("s_sort desc").Order("id asc").
		Find(&response.List)
	return
}

func (srv *ChineseService) ChineseGetBookList(page, level int) (response *api2.ChineseBookResponse, apiErr api.Error) {
	size := common.DEFAULT_PAGE_SIZE
	offset := size * (page - 1)
	response = &api2.ChineseBookResponse{}

	var total int64
	model.DefaultWeb().Model(&web.ChineseBook{}).Where("type = ? and status = 1", level).Debug().
		Count(&total).
		Order("position desc").
		Limit(size).
		Offset(offset).Find(&response.ListModel)
	response.Total = total
	response.Page = page
	response.TotalPage = math.Ceil(float64(total) / float64(common.DEFAULT_PAGE_SIZE))

	var bookInfoCountList []api2.ResponseBookInfoCount
	sql := `SELECT book_id,count(id) as book_count FROM s_chinese_picture_info GROUP BY book_id`
	model.DefaultWeb().Model(&web.ChineseBookInfo{}).Debug().
		Raw(sql).Scan(&bookInfoCountList)

	var temp api2.ResponseChineseBook
	for _, item := range response.ListModel {
		temp.Id = item.Id
		temp.BookId = item.BookId
		temp.Title = item.Title
		temp.Icon = item.Icon
		temp.Level = item.Type
		temp.Position = item.Position
		response.List = append(response.List, temp)
	}

	for index, item := range response.List {
		for _, it := range bookInfoCountList {
			if item.BookId == it.BookId {
				response.List[index].BookCount = it.BookCount
			}
		}
	}

	sort.Slice(response.List, func(i, j int) bool {
		if response.List[i].Position > response.List[j].Position {
			return true
		}
		return response.List[i].Position == response.List[j].Position && response.List[i].Id < response.List[j].Id
	})

	return response, nil
}

func (srv *ChineseService) ChineseGetBookInfo(bookId string) (response *api2.ChineseBookInfoResponse, apiErr api.Error) {
	response = &api2.ChineseBookInfoResponse{}
	model.DefaultWeb().Model(&web.ChineseBookInfo{}).Where("book_id = ?", bookId).Debug().
		Order("position asc").
		Find(&response.Info)
	return
}
