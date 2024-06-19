package cache

import (
	"github.com/chenyahui/gin-cache/persist"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	//这里是发起http请求时候的缓存的时长，1秒--作用是获取到http的body内容中间过渡，设置为正常使用的缓存中
	GIN_CACHE_TEMP = 1
	//这个参数是发起http get请求的中间操作，就是根据是否存在这个参数，判断是否不进行缓存操作
	GIN_CACHE_EXT = "no_cache_ext"
)

// redisCacheForcedRefresh redis中的路由缓存key强制更新
func redisCacheForcedRefresh(cacheKey, ginBaseUrl string, cfg *Config, cacheDuration time.Duration, cacheStore persist.CacheStore, respCache *ResponseCache) {
	// 组装完整的url路径
	url := makeHttpUrlPath(cacheKey, ginBaseUrl)
	// 发起临时的http请求，获取新的数据信息
	httpGetContent, httpCode := sendHttpGet(url)

	// only cache 2xx response and httpGetContent length gt 0
	if httpCode == http.StatusOK && len(httpGetContent) > 0 {
		var respCacheNew = respCache
		respCacheNew.Data = httpGetContent
		if err := cacheStore.Set(cacheKey, respCacheNew, cacheDuration); err != nil {
			cfg.logger.Errorf("set cache key error: %s, cache key: %s", err, cacheKey)
		}
	}
}

// makeHttpUrlPath 组装完整的url路径
func makeHttpUrlPath(cacheKey, ginBaseUrl string) (url string) {
	cacheUrlPath := cacheKey
	if isHave := strings.HasPrefix(cacheKey, "/"); !isHave {
		cacheUrlPath = "/" + cacheKey
	}

	url = ginBaseUrl + cacheUrlPath + "&=" + GIN_CACHE_EXT + "=" + strconv.FormatInt(time.Now().Unix(), 10)
	return
}

func sendHttpGet(url string) (bodyContent []byte, httpStatusCode int) {
	// 创建一个新的HTTP请求
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	// 设置自定义的Header
	request.Header.Set("User-Agent", "micromessenger")
	// 创建一个HTTP客户端
	client := &http.Client{}
	// 发起请求
	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close() // 确保在函数退出前关闭响应体
	httpStatusCode = response.StatusCode
	// 读取响应体内容
	bodyContent, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	return
}
