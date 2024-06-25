package student_handler

import (
	"api/model"
	"api/model/student"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func (sh *PoetryPictureHandler) ShellDirList(ctx *gin.Context) {
	makeAllMp4()
}

// step 6  renameMp4 将视频重命名
func renameMp4() {
	var list []student.SPoetryPicture
	model.DefaultStudent().Model(&student.SPoetryPicture{}).Debug().Where("id > 0").
		Find(&list)
	// 目标目录
	dir := "/Users/jiang/demo/result"

	// 遍历目录下的所有文件
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 打印文件路径（只打印文件，不打印目录）
		if !info.IsDir() {
			// 获取文件名（带扩展名）
			fileNameWithExt := filepath.Base(path)

			// 去掉文件扩展名
			fileName := strings.TrimSuffix(fileNameWithExt, filepath.Ext(fileNameWithExt))

			for idx, _ := range list {
				if fileName == list[idx].BookId {

					// 原始文件路径
					oldPath := path

					// 新文件路径
					newPath := "/Users/jiang/demo/result/" + list[idx].Title + ".mp4"

					// 重命名文件
					err := os.Rename(oldPath, newPath)
					if err != nil {
						fmt.Println("文件重命名失败:", err)
					}

					fmt.Println("文件重命名成功")
				}
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("遍历目录失败:", err)
	}

}

// step 5  moveMp4ToResult 将合成好的视频移动到指定的目录
func moveMp4ToResult() {
	root := "/Users/jiang/demo/mp4"

	// 读取根目录下的所有文件和文件夹
	folders, err := ioutil.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}

	// 遍历根目录下的每一个文件夹
	for _, folder := range folders {
		if folder.IsDir() {
			folderPath := filepath.Join(root, folder.Name())

			base := filepath.Base(folderPath)
			srcPath := folderPath + "/" + base + ".mp4"

			// 目标目录
			destDir := "/Users/jiang/demo/result"

			// 获取文件名
			fileName := filepath.Base(srcPath)

			// 目标文件路径
			destPath := filepath.Join(destDir, fileName)

			// 移动文件
			err := os.Rename(srcPath, destPath)
			if err != nil {
				fmt.Println("文件移动失败:", err)
			}

			fmt.Println("文件移动成功")
		}
	}
}

// step 4  makeAllMp4 整合一条命令合并单个的视频
func makeAllMp4() {

	root := "/Users/jiang/demo/shell"

	err := filepath.Walk(root, visitLogFiles)
	if err != nil {
		fmt.Printf("Error walking the path %q: %v\n", root, err)
	}
}

func visitLogFiles(path string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Printf("Error accessing path %q: %v\n", path, err)
		return nil
	}

	if info.IsDir() {
		return nil // 如果是目录，则继续遍历
	}

	if strings.HasSuffix(info.Name(), ".log") {
		// 使用 strings.TrimSuffix 去掉 .log 后缀
		base := strings.TrimSuffix(filepath.Base(path), ".log")
		// 在这里可以对找到的 log 文件进行处理，例如读取内容或者其他操作
		shell := `ffmpeg -f concat -safe 0 -i  /Users/jiang/demo/shell/` + base + `.log -c copy -absf aac_adtstoasc -s 1280x720 -c:v libx264 -pix_fmt yuv420p /Users/jiang/demo/mp4/` + base + `/` + base + `.mp4  2>&1`
		// 打开文件，如果文件不存在则创建文件
		file, err := os.OpenFile("/Users/jiang/demo/shell/mp4.sh", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("无法打开文件: %v\n", err)
		}
		defer file.Close()

		// 写入内容
		if _, err := file.WriteString(shell + "\n"); err != nil {
			fmt.Printf("无法写入内容: %v\n", err)
		}

		fmt.Println("内容已追加写入文件")
	}

	return nil
}

// step 3 makeMp4 合成完整的mp4
func makeMp4() {
	rootDir := "/Users/jiang/demo/file"

	// 读取根目录下的所有文件和文件夹
	folders, err := ioutil.ReadDir(rootDir)
	if err != nil {
		log.Fatal(err)
	}

	// 遍历根目录下的每一个文件夹
	for _, folder := range folders {
		if folder.IsDir() {
			folderPath := filepath.Join(rootDir, folder.Name())
			processFolderMp4(folderPath)
		}
	}
}

func processFolderMp4(folderPath string) {
	// 读取文件夹下的所有文件和文件夹
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		log.Printf("无法读取文件夹: %s\n", folderPath)
		return
	}

	var numbers []int
	re := regexp.MustCompile(`\d+`)

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".mp3" {
			// 提取文件名中的数字
			match := re.FindString(file.Name())
			if match != "" {
				num, err := strconv.Atoi(match)
				if err != nil {
					log.Printf("无法解析文件名中的数字: %s\n", file.Name())
					continue
				}
				numbers = append(numbers, num)
			}
		}
	}

	if len(numbers) == 0 {
		fmt.Printf("文件夹 %s 中没有找到任何包含数字的 .mp3 文件\n", folderPath)
		return
	}

	// 查找最大值和最小值
	max, min := numbers[0], numbers[0]
	for _, num := range numbers {
		if num > max {
			max = num
		}
		if num < min {
			min = num
		}
	}

	//var shell string
	//var ext string
	for i := min; i <= max; i++ {

		//	// 使用 strings.TrimPrefix 删除前缀
		trimmedPath := strings.TrimPrefix(folderPath, "/Users/jiang/demo/file/")
		//	iStr := strconv.Itoa(i)
		//	shell = `ffmpeg  -thread_queue_size 96   -loop 1   -t  5  -y -r 1 -i  /Users/jiang/demo/file/` + trimmedPath + `/` + iStr + `` + ext + `   -i   /Users/jiang/demo/file/` + trimmedPath + `/` + iStr + `.mp3 -x264-params keyint=1:scenecut=0  -vf "scale=2800:-2"   -absf aac_adtstoasc -s 1280x720 -c:v libx264 -pix_fmt yuv420p   /Users/jiang/demo/mp4/` + trimmedPath + `/` + iStr + `.mp4  2>&1`
		//	fmt.Println(shell)
		//
		//	// 打开文件，如果文件不存在则创建文件
		file, err := os.OpenFile("/Users/jiang/demo/shell/"+trimmedPath+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("无法打开文件: %v\n", err)
			return
		}
		defer file.Close()

		//	// 写入内容
		iStr := strconv.Itoa(i)
		content := `file '/Users/jiang/demo/mp4/` + trimmedPath + `/` + iStr + `.mp4'`
		if _, err := file.WriteString(content + "\n"); err != nil {
			fmt.Printf("无法写入内容: %v\n", err)
			return
		}
		//
		//	fmt.Println("内容已追加写入文件")
	}
}

// step 2 makeDirItem 创建对应的目录
func makeDirItem() {
	rootDir := "/Users/jiang/demo/file"

	// 读取根目录下的所有文件和文件夹
	folders, err := ioutil.ReadDir(rootDir)
	if err != nil {
		log.Fatal(err)
	}

	// 遍历根目录下的每一个文件夹
	for _, folder := range folders {
		if folder.IsDir() {
			folderPath := filepath.Join(rootDir, folder.Name())
			base := filepath.Base(folderPath)
			dir := "/Users/jiang/demo/mp4/" + base
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				fmt.Printf("%v 创建目录失败: %v\n", dir, err)
				return
			}

			fmt.Printf("目录 %s 创建成功\n", dir)
		}
	}
}

// step 1: makeAllMp4 制作一个个mp4
func makeAllMp4ItemOneByOne() {
	rootDir := "/Users/jiang/demo/file"

	// 读取根目录下的所有文件和文件夹
	folders, err := ioutil.ReadDir(rootDir)
	if err != nil {
		log.Fatal(err)
	}

	// 遍历根目录下的每一个文件夹
	for _, folder := range folders {
		if folder.IsDir() {
			folderPath := filepath.Join(rootDir, folder.Name())
			processFolder(folderPath)
		}
	}
}

// step 1: processFolder 制作一个个mp4
func processFolder(folderPath string) {
	// 读取文件夹下的所有文件和文件夹
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		log.Printf("无法读取文件夹: %s\n", folderPath)
		return
	}

	var numbers []int
	re := regexp.MustCompile(`\d+`)

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".mp3" {
			// 提取文件名中的数字
			match := re.FindString(file.Name())
			if match != "" {
				num, err := strconv.Atoi(match)
				if err != nil {
					log.Printf("无法解析文件名中的数字: %s\n", file.Name())
					continue
				}
				numbers = append(numbers, num)
			}
		}
	}

	if len(numbers) == 0 {
		fmt.Printf("文件夹 %s 中没有找到任何包含数字的 .mp3 文件\n", folderPath)
		return
	}

	// 查找最大值和最小值
	max, min := numbers[0], numbers[0]
	for _, num := range numbers {
		if num > max {
			max = num
		}
		if num < min {
			min = num
		}
	}

	var shell string
	var ext string
	for i := min; i <= max; i++ {
		if i == min {
			for _, file := range files {
				if !file.IsDir() && strings.Contains(file.Name(), strconv.Itoa(min)) && !strings.HasSuffix(file.Name(), ".mp3") {

					//fmt.Printf("file.Name():%s , 是否有MP3：%#v \n", file.Name(), strings.HasSuffix(file.Name(), ".mp3"))

					//获取后缀
					ext = filepath.Ext(file.Name())
				}
			}
		}

		// 使用 strings.TrimPrefix 删除前缀
		trimmedPath := strings.TrimPrefix(folderPath, "/Users/jiang/demo/file/")
		iStr := strconv.Itoa(i)
		shell = `ffmpeg  -thread_queue_size 96   -loop 1   -t  5  -y -r 1 -i  /Users/jiang/demo/file/` + trimmedPath + `/` + iStr + `` + ext + `   -i   /Users/jiang/demo/file/` + trimmedPath + `/` + iStr + `.mp3 -x264-params keyint=1:scenecut=0  -vf "scale=2800:-2"   -absf aac_adtstoasc -s 1280x720 -c:v libx264 -pix_fmt yuv420p   /Users/jiang/demo/mp4/` + trimmedPath + `/` + iStr + `.mp4  2>&1`
		fmt.Println(shell)

		// 打开文件，如果文件不存在则创建文件
		file, err := os.OpenFile("/Users/jiang/demo/shell/shell.sh", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("无法打开文件: %v\n", err)
			return
		}
		defer file.Close()

		// 写入内容
		if _, err := file.WriteString(shell + "\n"); err != nil {
			fmt.Printf("无法写入内容: %v\n", err)
			return
		}

		fmt.Println("内容已追加写入文件")
	}
}
