package student_handler

import (
	"api/model"
	"api/model/student"
	"api/service/student_service"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"path/filepath"
	"strings"
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
	var list []student.SEnglishPictureInfo
	var temp student.SEnglishPictureInfo
	model.DefaultStudent().Model(&student.SEnglishPictureInfo{}).Debug().Where("id > 8136").
		Find(&list)
	for idx, _ := range list {

		split := strings.Split(list[idx].BookId, "-")
		fmt.Println(split[0])
		////替换字符串
		//original := "https://static.58haha.com/"
		//replacement := "https://oss.58haha.com/"
		//pic := strings.Replace(list[idx].Icon, original, replacement, -1)
		temp.BookId = split[0]
		model.DefaultStudent().Model(&student.SEnglishPictureInfo{}).Where("id = ?", list[idx].Id).Debug().Updates(temp)
	}
}

func (sh *PoetryPictureHandler) moveFile(src, destDir string) error {
	// 创建目标目录（如果不存在）
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", destDir, err)
	}

	// 获取源文件的文件名
	srcFileName := filepath.Base(src)
	dest := filepath.Join(destDir, srcFileName)

	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %v", src, err)
	}
	defer srcFile.Close()

	// 创建目标文件
	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %v", dest, err)
	}
	defer destFile.Close()

	// 复制文件内容
	if _, err := io.Copy(destFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file from %s to %s: %v", src, dest, err)
	}

	//// 删除源文件
	//if err := os.Remove(src); err != nil {
	//	return fmt.Errorf("failed to remove source file %s: %v", src, err)
	//}

	return nil
}
