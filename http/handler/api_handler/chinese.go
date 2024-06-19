package api_handler

import (
	"api/common"
	"api/service/api_service"
	"api/utils/errs"
	"github.com/gin-gonic/gin"
	"strconv"
)

func NewChineseHandler() *ChineseHandler {
	return &ChineseHandler{
		service: api_service.NewChineseService(),
	}
}

type ChineseHandler struct {
	service *api_service.ChineseService
}

func (rh *ChineseHandler) ChineseGetNavList(ctx *gin.Context) {
	response, err := rh.service.ChineseGetNavList()
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}

func (rh *ChineseHandler) ChineseGetBookList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	if page < common.DEFAULT_PAGE {
		page = common.DEFAULT_PAGE
	}
	level, _ := strconv.Atoi(ctx.Query("level"))
	if level < common.DEFAULT_LEVEL {
		level = common.DEFAULT_LEVEL
	}
	response, err := rh.service.ChineseGetBookList(page, level)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}

func (rh *ChineseHandler) ChineseGetBookInfo(ctx *gin.Context) {
	bookId := ctx.Query("book_id")
	response, err := rh.service.ChineseGetBookInfo(bookId)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}
