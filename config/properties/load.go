package properties

import "strings"

// 读取配置
func Load(content string) map[string]string {
	// 按行分割
	sep := "\n"
	if strings.Contains(content, "\r\n") {
		sep = "\r\n"
	}
	lines := strings.Split(content, sep)
	return LoadArray(lines)
}

// 从数组中读取配置
func LoadArray(lines []string) map[string]string {
	attr := make(map[string]string)
	for i := 0; i < len(lines); i++ {
		line := strings.Trim(lines[i], " ")
		// 跳过空行
		if line == "" {
			continue
		}
		// 跳过注释
		if strings.HasPrefix(line, "#") {
			continue
		}
		kvs := strings.Split(line, "=")
		key := strings.Trim(kvs[0], " ")
		value := strings.Join(kvs[1:], "=")
		value = strings.Trim(value, " ")
		attr[key] = value
	}
	return attr
}
