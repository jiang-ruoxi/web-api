package api_handler

import (
	"api/common"
	"api/service/api_service"
	"api/utils/errs"
	"github.com/gin-gonic/gin"
	"strconv"
)

func NewPoetryHandler() *PoetryHandler {
	return &PoetryHandler{
		service: api_service.NewPoetryService(),
	}
}

type PoetryHandler struct {
	service *api_service.PoetryService
}

func (rh *PoetryHandler) PoetryGetChengList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	if page < common.DEFAULT_PAGE {
		page = common.DEFAULT_PAGE
	}
	level, _ := strconv.Atoi(ctx.Query("level"))
	if level < common.DEFAULT_LEVEL {
		level = common.DEFAULT_LEVEL
	}
	response, err := rh.service.PoetryGetChengList(page, level)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}

func (rh *PoetryHandler) PoetryGetChengInfo(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Query("id"))
	response, err := rh.service.PoetryGetChengInfo(id)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}

func (rh *PoetryHandler) PoetryGetSchoolList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	if page < common.DEFAULT_PAGE {
		page = common.DEFAULT_PAGE
	}
	response, err := rh.service.PoetryGetSchoolList(page)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}
