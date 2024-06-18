package api_service

import (
	"api/model"
	"api/utils/errs"
	"context"
	"github.com/jiang-ruoxi/gopkg/server/api"
)

func NewIndexService() *IndexService {
	return &IndexService{}
}

type IndexService struct {
}

type IndexResponse struct {
	LexiconList []*model.Lexicon `json:"list"`
}

func (srv *IndexService) Index(ctx context.Context) (response *IndexResponse, apiErr api.Error) {
	response = &IndexResponse{}
	model.Default().Where("status = ?", 1)
	if err := model.Default().Where("status = ?", 1).Order("id DESC").Limit(10).Find(&response.LexiconList).Error; err != nil {
		return response, errs.NewError(err.Error())
	}
	return response, nil
}
