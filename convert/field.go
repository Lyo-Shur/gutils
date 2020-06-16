package convert

// 65是A 90是Z 95是下划线 +32是转为小写

// 转大驼峰
func ToBigHump(s string) string {
	l := len(s)
	var data []byte
	// 遍历字符串进行转化
	for i := 0; i < l; i++ {
		if s[i] != 95 {
			data = append(data, s[i])
			continue
		}
		i++
		if 97 <= s[i] && s[i] <= 122 {
			data = append(data, s[i]-32)
			continue
		}
		data = append(data, s[i])
	}
	// 如果首字符小写则转为大写
	if 97 <= data[0] && data[0] <= 122 {
		data[0] = data[0] - 32
	}
	return string(data)
}

// 转小驼峰
func ToSmallHump(s string) string {
	l := len(s)
	var data []byte
	// 遍历字符串进行转化
	for i := 0; i < l; i++ {
		if s[i] != 95 {
			data = append(data, s[i])
			continue
		}
		i++
		if 97 <= s[i] && s[i] <= 122 {
			data = append(data, s[i]-32)
			continue
		}
		data = append(data, s[i])
	}
	// 如果首字符大写则转为小写
	if 65 <= data[0] && data[0] <= 90 {
		data[0] = data[0] + 32
	}
	return string(data)
}

// 转下划线方法(默认去掉首字符下划线)
func ToUnderline(s string) string {
	l := len(s)
	var data []byte
	// 遍历字符串进行转化
	for i := 0; i < l; i++ {
		if 65 <= s[i] && s[i] <= 90 {
			data = append(data, 95)
			data = append(data, s[i]+32)
		} else {
			data = append(data, s[i])
		}
	}
	// 首字符出现下划线时，截去
	var r string
	if data[0] == 95 {
		r = string(data[1:])
	} else {
		r = string(data)
	}
	return r
}
