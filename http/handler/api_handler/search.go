package api_handler

import (
	"api/common"
	"api/service/search_service"
	"api/utils/errs"
	"api/utils/validator"
	"github.com/gin-gonic/gin"
)

func NewSearchHandler() *SearchHandler {
	return &SearchHandler{
		service: search_service.NewSearchService(),
	}
}

type SearchHandler struct {
	service *search_service.SearchService
}

func (sh *SearchHandler) SearchList(ctx *gin.Context) {
	// 获取请求参数
	req := &common.SearchRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(errs.ErrResp(errs.NewError(err.Error())))
		return
	}

	// 验证分页参数
	req.N, req.P = validator.VerifyLimit(req.N, req.P)

	response, err := sh.service.SearchList(ctx, req)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}
