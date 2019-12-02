package api

import (
	"encoding/json"
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
		return "JSON转化失败"
	}
	return string(bs)
}

// 创建-状态码模式数据传输对象
func JsonCodeModeDTO(code int64, message string, data interface{}) string {
	cmd := CodeModeDTO{}
	cmd.Code = code
	cmd.Message = message
	cmd.Data = data
	return cmd.ToJson()
}
