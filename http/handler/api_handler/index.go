package api_handler

import (
	"api/service/api_service"
	"api/utils/errs"
	"github.com/gin-gonic/gin"
)

func NewIndexHandler() *IndexHandler {
	return &IndexHandler{
		service: api_service.NewIndexService(),
	}
}

type IndexHandler struct {
	service *api_service.IndexService
}

type IndexRequest struct {
}

func (sh *IndexHandler) Index(ctx *gin.Context) {
	// 获取请求参数
	req := &IndexRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(errs.ErrResp(errs.NewError(err.Error())))
		return
	}

	// 首页
	response, err := sh.service.Index(ctx)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
	return
}
