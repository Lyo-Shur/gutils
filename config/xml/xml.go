package xml

import (
	"encoding/xml"
	"log"
)

// 解析配置文件
func Load(content string, model interface{}) {
	// 解析配置文件
	err := xml.Unmarshal([]byte(content), model)
	if err != nil {
		log.Fatal(err)
	}
}
