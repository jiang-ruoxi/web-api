package api_service

import (
	"api/common"
	api2 "api/common/web"
	"api/model"
	"api/model/web"
	"github.com/jiang-ruoxi/gopkg/server/api"
	"math"
	"strings"
)

func NewPoetryService() *PoetryService {
	return &PoetryService{}
}

type PoetryService struct {
}

func (srv *PoetryService) PoetryGetChengList(page, level int) (response *api2.ChineseChengYuResponse, apiErr api.Error) {
	size := common.DEFAULT_PAGE_SIZE
	offset := size * (page - 1)
	response = &api2.ChineseChengYuResponse{}

	var total int64
	model.DefaultWeb().Model(&web.ChengYU{}).Where("level = ?", level).Debug().
		Count(&total).
		Order("position desc").
		Limit(size).
		Offset(offset).Find(&response.List)
	response.Total = total
	response.Page = page
	response.TotalPage = math.Ceil(float64(total) / float64(common.DEFAULT_PAGE_SIZE))

	return response, nil
}

func (srv *PoetryService) PoetryGetChengInfo(id int) (response *api2.ChineseChengYuInfoResponse, apiErr api.Error) {
	response = &api2.ChineseChengYuInfoResponse{}
	model.DefaultWeb().Model(&web.ChengYU{}).Where("id = ?", id).Debug().First(&response.Info)
	fields := strings.Fields(response.Info.Story)
	response.Info.StoryList = fields
	return response, nil
}

func (srv *PoetryService) PoetryGetSchoolList(page int) (response *api2.SchoolPoetryListResponse, apiErr api.Error) {
	size := common.DEFAULT_PAGE_SIZE
	offset := size * (page - 1)
	response = &api2.SchoolPoetryListResponse{}

	var total int64
	model.DefaultWeb().Model(&web.Poetry{}).Debug().
		Count(&total).
		Raw("SELECT id,poetry_id,title,grade_id,grade,grade_level,author,dynasty FROM s_school_poetry limit ? offset ?", size, offset).Scan(&response.List)

	response.Total = total
	response.Page = page
	response.TotalPage = math.Ceil(float64(total) / float64(common.DEFAULT_PAGE_SIZE))

	return response, nil
}
