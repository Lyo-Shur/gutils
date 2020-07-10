package read

import (
	"io/ioutil"
	"log"
	"net/http"
)

// 根据网络地址读取文件
func HTTP(path string) string {
	// 预创建请求
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		log.Println(err)
	}
	// 设置请求头
	req.Header.Set("Content-Type", "text/html")
	// 创建链接
	client := http.Client{}         // 创建一个httpClient
	response, err := client.Do(req) // 调用rest接口
	if err != nil {
		log.Println(err)
	}
	// 读取响应
	content, err := ioutil.ReadAll(response.Body)
	return string(content)
}
