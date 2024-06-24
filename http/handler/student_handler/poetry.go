package student_handler

import (
	"api/model"
	"api/model/web"
	"api/service/student_service"
	"github.com/gin-gonic/gin"
	"strconv"
)

func NewPoetryPictureInfoHandler() *PoetryPictureHandler {
	return &PoetryPictureHandler{
		service: student_service.NewPoetryPictureService(),
	}
}

type PoetryPictureHandler struct {
	service *student_service.PoetryPictureService
}

func (sh *PoetryPictureHandler) PoetryPictureInfoList(ctx *gin.Context) {
	var list []web.EEnglishPictureInfo
	var temp web.EEnglishPictureInfo
	model.DefaultStudent().Model(&web.EEnglishPictureInfo{}).Debug().Where("id > 8136").
		Find(&list)
	for idx, _ := range list {
		temp.Pic = "https://oss.58haha.com/english_book/file/" + list[idx].BookId + "/" + strconv.Itoa(list[idx].Position) + ".png"
		temp.Mp3 = "https://oss.58haha.com/english_book/file/" + list[idx].BookId + "/" + strconv.Itoa(list[idx].Position) + ".mp3"
		model.DefaultWeb().Model(&web.EEnglishPictureInfo{}).Where("id = ?", list[idx].Id).Debug().Updates(temp)
	}
}
