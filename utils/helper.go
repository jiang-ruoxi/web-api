package utils

import (
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"strings"
)

// TruncateString 截取指定的字符串
func TruncateString(s string, num int) string {
	runes := []rune(s)
	if len(runes) <= num {
		return s
	}
	return string(runes[:num])
}

//ClearMobileText 清除手机号
func ClearMobileText(text string) (cleanedText string) {
	// 定义手机号的正则表达式
	phoneRegex := regexp.MustCompile(`1[3456789]\d{9}`)

	// 查找所有匹配的手机号
	matches := phoneRegex.FindAllString(text, -1)

	if matches != nil {
		//fmt.Println("找到的手机号：", matches)
		// 将手机号去除
		cleanedText = phoneRegex.ReplaceAllString(text, "[手机号敏感数据不予显示]")
		//fmt.Println("去除手机号后的文本：", cleanedText)
	} else {
		cleanedText = text
		//fmt.Println("未找到手机号")
	}

	return
}

//RemoveDuplicates 切片去重
func RemoveDuplicates(slice []int64) []int64 {
	encountered := map[int64]bool{}
	result := []int64{}

	for v := range slice {
		if encountered[slice[v]] != true {
			encountered[slice[v]] = true
			result = append(result, slice[v])
		}
	}

	return result
}

// RegContent 正则匹配敏感词
func RegContent(matchContent string, sensitiveWords []string) string {
	if len(sensitiveWords) < 1 {
		return matchContent
	}
	banWords := make([]string, 0) // 收集匹配到的敏感词

	// 构造正则匹配字符
	regStr := strings.Join(sensitiveWords, "|")
	wordReg := regexp.MustCompile(regStr)
	//println("regStr -> ", regStr)

	textBytes := wordReg.ReplaceAllFunc([]byte(matchContent), func(bytes []byte) []byte {
		banWords = append(banWords, string(bytes))
		textRunes := []rune(string(bytes))
		replaceBytes := make([]byte, 0)
		for i, runeLen := 0, len(textRunes); i < runeLen; i++ {
			replaceBytes = append(replaceBytes, byte('*'))
		}
		return replaceBytes
	})
	//fmt.Println("srcText        -> ", matchContent)
	//fmt.Println("replaceText    -> ", string(textBytes))
	//fmt.Println("sensitiveWords -> ", banWords)
	return string(textBytes)
}

//ExistDir 创建目录
func ExistDir(path string) {
	// 判断路径是否存在
	_, err := os.ReadDir(path)
	if err != nil {
		// 不存在就创建
		err = os.MkdirAll(path, fs.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}
}
