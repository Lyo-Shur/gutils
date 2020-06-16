package api

import (
	"encoding/json"
	"log"
)

// 状态码模式数据传输对象
type CodeModeDTO struct {
	// 状态码
	Code int64
	// 消息
	Message string
	// 数据体
	Data interface{}
}

// 转化为Json
func (cmd *CodeModeDTO) ToJson() string {
	bs, err := json.Marshal(cmd)
	if err != nil {
		log.Println(err)
		return "JSON转化失败"
	}
	return string(bs)
}
