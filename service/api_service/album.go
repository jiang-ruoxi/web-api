package api_service

import (
	"api/common"
	api2 "api/common/web"
	"api/model"
	"api/model/web"
	"github.com/jiang-ruoxi/gopkg/server/api"
	"math"
)

func NewAlbumService() *AlbumService {
	return &AlbumService{}
}

type AlbumService struct {
}

func (srv *AlbumService) AlbumGetList(page int) (response *api2.ChineseAlbumResponse, apiErr api.Error) {
	size := common.DEFAULT_PAGE_SIZE
	offset := size * (page - 1)
	response = &api2.ChineseAlbumResponse{}

	var total int64
	model.DefaultWeb().Model(&web.ChineseBookAlbum{}).Debug().
		Count(&total).
		Order("position desc").
		Limit(size).
		Offset(offset).Find(&response.List)
	response.Total = total
	response.Page = page
	response.TotalPage = math.Ceil(float64(total) / float64(common.DEFAULT_PAGE_SIZE))

	return response, nil
}

func (srv *AlbumService) AlbumGetListInfo(bookId string) (response *api2.ChineseAlbumListInfoResponse, apiErr api.Error) {
	response = &api2.ChineseAlbumListInfoResponse{}
	var total int64
	model.DefaultWeb().Model(&web.ChineseAlbumInfo{}).Debug().
		Where("book_id = ?", bookId).Count(&total).
		Order("position desc").Find(&response.List)
	response.Total = total
	response.Page = 1
	response.TotalPage = math.Ceil(float64(total) / float64(common.DEFAULT_PAGE_SIZE))

	return response, nil
}

func (srv *AlbumService) AlbumGetInfo(id int) (response *api2.ChineseAlbumInfoResponse, apiErr api.Error) {
	response = &api2.ChineseAlbumInfoResponse{}
	model.DefaultWeb().Model(&web.ChineseAlbumInfo{}).Debug().Where("id = ?", id).First(&response.Info)
	return response, nil
}
