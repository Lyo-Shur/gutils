package file

import (
	"log"
	"os"
)

// 根据路径读取文件
func Read(path string) string {
	// 尝试打开文件
	fp, err := os.OpenFile(path, os.O_RDONLY, 0755)
	// 延迟关闭
	defer fp.Close()
	if err != nil {
		log.Fatal(err)
	}

	// 结果字符串
	str := ""

	// 创建缓冲区
	data := make([]byte, 100)
	for true {
		// 缓冲读取
		n, err := fp.Read(data)
		// 无更多字节时跳出
		if n == 0 {
			break
		}
		// 错误日志
		if err != nil {
			log.Fatal(err)
		}
		// 结果拼接
		str = str + string(data[:n])
	}
	return str
}
