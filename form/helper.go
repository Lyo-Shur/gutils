package form

import (
	"errors"
	"github.com/kataras/iris"
	"reflect"
	"strconv"
)

// 表单参数帮助工具
type helper struct {
	ctx iris.Context
}

// 获取表单参数帮助工具
func Helper(ctx iris.Context) *helper {
	h := &helper{}
	h.ctx = ctx
	return h
}

// 将表单中的值绑定到当前参数上
func (h *helper) Binding(dest interface{}) error {
	// 对参数类别进行判断
	it := reflect.TypeOf(dest).Elem()
	iv := reflect.ValueOf(dest).Elem()
	if it.Kind() != reflect.Struct {
		return errors.New("绑定参数必须为结构体")
	}
	// 绑定参数
	var err error
	err = h.bindingParameters(it, iv)
	if err != nil {
		return err
	}
	return nil
}

// 绑定参数到结构体
func (h *helper) bindingParameters(it reflect.Type, iv reflect.Value) error {
	// 首先获取参数
	mss := h.ctx.FormValues()
	// 遍历当前参数
	for i, v := range mss {
		// 当当前值不存在时，跳过本次循环
		if v[0] == "" {
			continue
		}
		// 获取对应的字段名
		name := bigHump(i)
		// 获取结构体字段
		rv := iv.FieldByName(name)

		// 预定义错误
		covError := errors.New("参数" + name + "无法转化为结构体对应类型")
		switch rv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			i64, err := strconv.ParseInt(v[0], 10, rv.Type().Bits())
			if err != nil {
				return covError
			}
			rv.SetInt(i64)
			break
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			u64, err := strconv.ParseUint(v[0], 10, rv.Type().Bits())
			if err != nil {
				return covError
			}
			rv.SetUint(u64)
			break
		case reflect.Float32, reflect.Float64:
			f64, err := strconv.ParseFloat(v[0], rv.Type().Bits())
			if err != nil {
				return covError
			}
			rv.SetFloat(f64)
			break
		case reflect.String:
			rv.SetString(v[0])
			break
		case reflect.Bool:
			if v[0] == "0" {
				rv.SetBool(false)
				break
			}
			if v[0] == "1" {
				rv.SetBool(true)
				break
			}
			return covError
		}
	}
	return nil
}

// 获取文件持有者
func (h *helper) GetFileHolder() *FileHolder {
	fh := FileHolder{}

	// 上传的文件
	form := h.ctx.Request().MultipartForm
	if form != nil {
		fh.m = form.File
	}

	return &fh
}

// 转大驼峰
func bigHump(s string) string {
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
