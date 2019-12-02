package config

import (
	"encoding/json"
	"github.com/Lyo-Shur/gutils/file"
	"log"
	"strings"
)

// Json配置文件所在的路径
// 此变量建议谨慎设置，尽量避免使用绝对路径
// Json配置文件默认放置在config文件夹下
// 配置文件默认名称为config.json
var JsonConfigPath = "./config/config.json"

// 配置结构体
type jsonConfig struct {
	f interface{}
}

// 初始化配置
func JsonConfig() jsonConfig {
	// 预创建配置结构体
	c := jsonConfig{}
	// 读取配置文件并缓存到结构体
	c.f = fromJson(file.Read(JsonConfigPath))
	return c
}

// 读取内容
func (c *jsonConfig) Get(key string) string {
	// 切分key
	ss := strings.Split(key, ".")
	// 缓存当前的键
	var k string
	// 缓存当前interface
	v := c.f
	// 循环切分结果
	for i := 0; i < len(ss); i++ {
		// 当前待获取的键
		k = ss[i]
		v = v.(map[string]interface{})[k]
	}
	return v.(string)
}

// json转换为接口
func fromJson(sjson string) interface{} {
	// 预返回接口
	var f interface{}
	// 解析json内容到接口
	err := json.Unmarshal([]byte(sjson), &f)
	// 异常输出
	if err != nil {
		log.Fatal(err)
	}
	return f
}
