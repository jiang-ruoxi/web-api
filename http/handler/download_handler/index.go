package download_handler

import (
	"api/model"
	"api/model/web"
	"api/service/download_service"
	"api/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
)

func NewIndexHandler() *IndexHandler {
	return &IndexHandler{
		service: download_service.NewIndexService(),
	}
}

type IndexHandler struct {
	service *download_service.IndexService
}

func (sh *IndexHandler) Index(ctx *gin.Context) {
	var list []web.EEnglishPictureInfo
	var temp web.EEnglishPictureInfo
	model.DefaultWeb().Model(&web.EEnglishPictureInfo{}).Debug().Where("id > 8136").
		Find(&list)
	for idx, _ := range list {
		temp.Pic = "https://oss.58haha.com/english_book/file/" + list[idx].BookId + "/" + strconv.Itoa(list[idx].Position) + ".png"
		temp.Mp3 = "https://oss.58haha.com/english_book/file/" + list[idx].BookId + "/" + strconv.Itoa(list[idx].Position) + ".mp3"
		model.DefaultWeb().Model(&web.EEnglishPictureInfo{}).Where("id = ?", list[idx].Id).Debug().Updates(temp)
	}
}

func (sh *IndexHandler) Index1(ctx *gin.Context) {
	var list []web.EEnglishPicture
	model.DefaultWeb().Model(&web.EEnglishPicture{}).Debug().Where("id > 677").
		Find(&list)

	urls := make(chan web.EEnglishPicture)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go sh.worker(urls, &wg, i)
	}

	for idx, _ := range list {
		//path := "/Users/jiang/Project/english_book/file/" + strconv.Itoa(list[idx].BookId)
		//filepath := "/Users/jiang/Project/english_book/file/" + strconv.Itoa(list[idx].BookId) + "/" + strconv.Itoa(list[idx].Position) + ".png"
		//fmt.Println(list[idx].Pic)
		//fmt.Println(filepath)
		urls <- list[idx]
		//err := sh.downloadFile(list[idx].Pic, path, filepath)
		//fmt.Println(err)
	}

	close(urls)

	wg.Wait()
	fmt.Println("All downloads completed.")

}

func (sh *IndexHandler) worker(urls <-chan web.EEnglishPicture, wg *sync.WaitGroup, id int) {
	fmt.Printf("urls:%#v\n", urls)
	defer wg.Done()
	for url := range urls {
		path := "/Users/jiang/Project/english_book/cover/"
		filepath := "/Users/jiang/Project/english_book/cover/" + "/" + strconv.Itoa(url.Id) + ".png"
		err := sh.downloadFile(url.Icon, path, filepath)
		if err != nil {
			fmt.Printf("Worker %d: Error downloading %s: %v\n", id, url.Icon, err)
		} else {
			fmt.Printf("Worker %d: Successfully downloaded %s\n", id, url.Icon)
		}
	}
}

func (sh *IndexHandler) downloadFile(url, path, filepath string) error {
	utils.ExistDir(path)
	// 创建文件
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 发送 GET 请求
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查 HTTP 响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// 将响应 Body 写入文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
