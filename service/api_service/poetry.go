package api_service

import (
	"api/common"
	"api/model"
	"github.com/jiang-ruoxi/gopkg/server/api"
	"math"
	"strings"
)

func NewPoetryService() *PoetryService {
	return &PoetryService{}
}

type PoetryService struct {
}

func (srv *PoetryService) PoetryGetChengList(page, level int) (response *common.ChineseChengYuResponse, apiErr api.Error) {
	size := common.DEFAULT_PAGE_SIZE
	offset := size * (page - 1)
	response = &common.ChineseChengYuResponse{}

	var total int64
	model.DefaultWeb().Model(&model.ChengYU{}).Where("level = ?", level).Debug().
		Count(&total).
		Order("position desc").
		Limit(size).
		Offset(offset).Find(&response.List)
	response.Total = total
	response.Page = page
	response.TotalPage = math.Ceil(float64(total) / float64(common.DEFAULT_PAGE_SIZE))

	return response, nil
}

func (srv *PoetryService) PoetryGetChengInfo(id int) (response *common.ChineseChengYuInfoResponse, apiErr api.Error) {
	response = &common.ChineseChengYuInfoResponse{}
	model.DefaultWeb().Model(&model.ChengYU{}).Where("id = ?", id).Debug().First(&response.Info)
	fields := strings.Fields(response.Info.Story)
	response.Info.StoryList = fields
	return response, nil
}

func (srv *PoetryService) PoetryGetSchoolList(page int) (response *common.SchoolPoetryListResponse, apiErr api.Error) {
	size := common.DEFAULT_PAGE_SIZE
	offset := size * (page - 1)
	response = &common.SchoolPoetryListResponse{}

	var total int64
	model.DefaultWeb().Model(&model.Poetry{}).Debug().
		Count(&total).
		Raw("SELECT id,poetry_id,title,grade_id,grade,grade_level,author,dynasty FROM s_school_poetry limit ? offset ?", size, offset).Scan(&response.List)

	response.Total = total
	response.Page = page
	response.TotalPage = math.Ceil(float64(total) / float64(common.DEFAULT_PAGE_SIZE))

	return response, nil
}
