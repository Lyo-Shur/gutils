package validator

import (
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// 校验规则
type Rule interface {
	// 判断此表达式能否解析
	Know(expr string) bool
	// 进行规则判断 表达式
	Check(expr string, params map[string]string, paramName string) (bool, error)
}

// 非空规则
type Required struct{}

func (r *Required) Know(expr string) bool {
	return "required" == expr
}
func (r *Required) Check(expr string, params map[string]string, paramName string) (bool, error) {
	_, ok := params[paramName]
	if !ok {
		return false, nil
	}
	return true, nil
}

// 非空规则
type Ban struct{}

func (b *Ban) Know(expr string) bool {
	return "ban" == expr
}
func (b *Ban) Check(expr string, params map[string]string, paramName string) (bool, error) {
	_, ok := params[paramName]
	if ok {
		return false, nil
	}
	return true, nil
}

// 长度规则
type Length struct{}

func (l *Length) Know(expr string) bool {
	return strings.HasPrefix(expr, "length")
}
func (l *Length) Check(expr string, params map[string]string, paramName string) (bool, error) {
	_, ok := params[paramName]
	if !ok {
		return true, nil
	}
	expr = strings.Replace(expr, "length", "", -1)
	expr = string([]rune(expr)[1 : len(expr)-1])
	// 长度范围
	r := strings.Split(expr, "-")
	start, err := strconv.ParseInt(r[0], 10, 64)
	if err != nil {
		return false, err
	}
	end, err := strconv.ParseInt(r[1], 10, 64)
	if err != nil {
		return false, err
	}
	// 判断长度
	value := params[paramName]
	length := int64(utf8.RuneCountInString(value))
	if start > length || length > end {
		return false, nil
	}
	return true, nil
}

// 范围规则
type Range struct{}

func (r *Range) Know(expr string) bool {
	return strings.HasPrefix(expr, "range")
}
func (r *Range) Check(expr string, params map[string]string, paramName string) (bool, error) {
	_, ok := params[paramName]
	if !ok {
		return true, nil
	}
	expr = strings.Replace(expr, "range", "", -1)
	expr = string([]rune(expr)[1 : len(expr)-1])
	// 长度范围
	exprs := strings.Split(expr, "-")
	start, err := strconv.ParseInt(exprs[0], 10, 64)
	if err != nil {
		return false, err
	}
	end, err := strconv.ParseInt(exprs[1], 10, 64)
	if err != nil {
		return false, err
	}
	// 对值进行转换
	val, err := strconv.ParseInt(params[paramName], 10, 64)
	if err != nil {
		return false, nil
	}
	if start > val || val > end {
		return false, nil
	}
	return true, nil
}

// 时间日期规则
type DateTime struct{}

func (r *DateTime) Know(expr string) bool {
	return strings.HasPrefix(expr, "datetime")
}
func (r *DateTime) Check(expr string, params map[string]string, paramName string) (bool, error) {
	_, ok := params[paramName]
	if !ok {
		return true, nil
	}
	// 获取参数
	expr = strings.Replace(expr, "datetime", "", -1)
	expr = string([]rune(expr)[1 : len(expr)-1])

	// 尝试转换时间
	loc, err := time.LoadLocation("Local")
	if err != nil {
		return false, err
	}
	_, err = time.ParseInLocation(expr, params[paramName], loc)
	if err != nil {
		return false, err
	}
	return true, nil
}

// 正则规则
type Regexp struct{}

func (r *Regexp) Know(expr string) bool {
	return strings.HasPrefix(expr, "regexp")
}
func (r *Regexp) Check(expr string, params map[string]string, paramName string) (bool, error) {
	_, ok := params[paramName]
	if !ok {
		return true, nil
	}
	// 获取参数
	expr = strings.Replace(expr, "regexp", "", -1)
	expr = string([]rune(expr)[1 : len(expr)-1])

	// 正则匹配
	matched, err := regexp.MatchString(expr, params[paramName])
	if err != nil {
		return false, err
	}
	if !matched {
		return false, nil
	}
	return true, nil
}
