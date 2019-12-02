package validator

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// 校验规则
type Rule interface {
	// 判断此表达式能否解析
	Know(expr string) bool
	// 进行规则判断 表达式 zi
	Check(expr string, fieId reflect.StructField, v reflect.Value) error
}

// 非空规则
type Required struct{}

func (r *Required) Know(expr string) bool {
	return "required" == expr
}
func (r *Required) Check(expr string, fieId reflect.StructField, v reflect.Value) error {
	// 类型判断
	if fieId.Type.Kind() == reflect.String {
		str := v.String()
		if str == "" {
			return errors.New("字段" + fieId.Name + "不能为空")
		}
		return nil
	}
	// 当非string类型字符串使用required时
	return errors.New("字段" + fieId.Name + "类型错误，此处需要string")
}

// 长度规则
type Length struct{}

func (l *Length) Know(expr string) bool {
	return strings.HasPrefix(expr, "length")
}
func (l *Length) Check(expr string, fieId reflect.StructField, v reflect.Value) error {
	// 类型判断
	if fieId.Type.Kind() != reflect.String {
		return errors.New("字段" + fieId.Name + "类型错误，此处需要string")
	}
	str := v.String()
	expr = strings.Replace(expr, "length", "", -1)
	expr = string([]rune(expr)[1 : len(expr)-1])
	// 长度范围
	r := strings.Split(expr, "-")
	start, err := strconv.ParseInt(r[0], 10, 64)
	if err != nil {
		return errors.New("字段" + fieId.Name + "限制规则错误")
	}
	end, err := strconv.ParseInt(r[1], 10, 64)
	if err != nil {
		return errors.New("字段" + fieId.Name + "限制规则错误")
	}
	// 判断长度
	length := int64(len(str))
	if start <= length && length <= end {
		return nil
	}
	return errors.New("字段" + fieId.Name + "长度非法")
}

// 范围规则
type Range struct{}

func (r *Range) Know(expr string) bool {
	return strings.HasPrefix(expr, "range")
}
func (r *Range) Check(expr string, fieId reflect.StructField, v reflect.Value) error {
	var val int64
	// 判断类别
	switch fieId.Type.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val = v.Int()
		expr = strings.Replace(expr, "range", "", -1)
		expr = string([]rune(expr)[1 : len(expr)-1])
		// 长度范围
		r := strings.Split(expr, "-")
		start, err := strconv.ParseInt(r[0], 10, 64)
		if err != nil {
			return errors.New("字段" + fieId.Name + "限制规则错误")
		}
		end, err := strconv.ParseInt(r[1], 10, 64)
		if err != nil {
			return errors.New("字段" + fieId.Name + "限制规则错误")
		}
		// 判断长度
		if start <= val && val <= end {
			return nil
		}
		return errors.New("字段" + fieId.Name + "范围非法")
	}
	return errors.New("字段" + fieId.Name + "类型错误，此处需要int")
}

// 正则规则
type Regexp struct{}

func (r *Regexp) Know(expr string) bool {
	return strings.HasPrefix(expr, "regexp")
}
func (r *Regexp) Check(expr string, fieId reflect.StructField, v reflect.Value) error {
	// 类型判断
	if fieId.Type.Kind() == reflect.String {
		str := v.String()

		// 获取参数
		expr = strings.Replace(expr, "regexp", "", -1)
		expr = string([]rune(expr)[1 : len(expr)-1])

		// 正则匹配
		matched, err := regexp.MatchString(expr, str)
		if err != nil {
			return err
		}
		if !matched {
			return errors.New("字段" + fieId.Name + "格式不匹配")
		}
		return nil
	}
	// 当非string类型字符串使用required时
	return errors.New("字段" + fieId.Name + "类型错误，此处需要string")
}
