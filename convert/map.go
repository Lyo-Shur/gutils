package convert

import (
	"errors"
	"reflect"
	"strconv"
)

// 将Map中的值，绑定进结构体
func MapBindToStruct(m map[string]string, v interface{}) error {
	// 对参数类别进行判断
	it := reflect.TypeOf(v).Elem()
	iv := reflect.ValueOf(v).Elem()
	if it.Kind() != reflect.Struct {
		return errors.New("绑定参数必须为结构体")
	}
	// 绑定参数
	var err error
	err = bindingParameters(m, iv)
	if err != nil {
		return err
	}
	return nil
}

// 绑定参数到结构体
func bindingParameters(mss map[string]string, iv reflect.Value) error {
	// 遍历当前参数
	for i, v := range mss {
		// 当当前值不存在时，跳过本次循环
		if v == "" {
			continue
		}
		// 获取对应的字段名
		name := ToBigHump(i)
		// 获取结构体字段
		rv := iv.FieldByName(name)

		// 预定义错误
		covError := errors.New("参数" + name + "无法转化为结构体对应类型")
		switch rv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			i64, err := strconv.ParseInt(v, 10, rv.Type().Bits())
			if err != nil {
				return covError
			}
			rv.SetInt(i64)
			break
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			u64, err := strconv.ParseUint(v, 10, rv.Type().Bits())
			if err != nil {
				return covError
			}
			rv.SetUint(u64)
			break
		case reflect.Float32, reflect.Float64:
			f64, err := strconv.ParseFloat(v, rv.Type().Bits())
			if err != nil {
				return covError
			}
			rv.SetFloat(f64)
			break
		case reflect.String:
			rv.SetString(v)
			break
		case reflect.Bool:
			if v == "0" {
				rv.SetBool(false)
				break
			}
			if v == "1" {
				rv.SetBool(true)
				break
			}
			return covError
		}
	}
	return nil
}
