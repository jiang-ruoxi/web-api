package api_handler

import (
	"api/common"
	"api/service/api_service"
	"api/utils/errs"
	"github.com/gin-gonic/gin"
	"strconv"
)

func NewAlbumHandler() *AlbumHandler {
	return &AlbumHandler{
		service: api_service.NewAlbumService(),
	}
}

type AlbumHandler struct {
	service *api_service.AlbumService
}

func (rh *AlbumHandler) AlbumGetList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	if page < common.DEFAULT_PAGE {
		page = common.DEFAULT_PAGE
	}
	response, err := rh.service.AlbumGetList(page)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}

func (rh *AlbumHandler) AlbumGetListInfo(ctx *gin.Context) {
	bookId := ctx.Query("book_id")
	response, err := rh.service.AlbumGetListInfo(bookId)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}

func (rh *AlbumHandler) AlbumGetInfo(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Query("id"))
	response, err := rh.service.AlbumGetInfo(id)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}

	ctx.JSON(errs.SucResp(response))
}
