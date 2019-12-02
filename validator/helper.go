package validator

import (
	"errors"
	"log"
	"reflect"
)

// 校验帮助工具
type helper struct {
	// 结构体字段
	i interface{}
	// 校验规则
	config map[string]string
	// 内置规则
	rule []Rule
}

// 获取校验帮助工具
func Helper(i interface{}, config map[string]string) *helper {
	h := helper{}
	h.i = i
	h.config = config
	// 添加默认规则
	h.rule = make([]Rule, 4)
	h.rule[0] = &Required{}
	h.rule[1] = &Length{}
	h.rule[2] = &Range{}
	h.rule[3] = &Regexp{}
	return &h
}

// 添加自定义规则
func (h *helper) AddRule(r Rule) {
	h.rule = append(h.rule, r)
}

// 执行校验
func (h *helper) Check() error {
	it := reflect.TypeOf(h.i).Elem()
	iv := reflect.ValueOf(h.i).Elem()
	for k, expr := range h.config {
		itf, b := it.FieldByName(k)
		if !b {
			log.Fatal("校验控件找不到字段" + k)
		}
		ivf := iv.FieldByName(k)

		// 当前字段已 经过检查？
		checked := false
		for i := 0; i < len(h.rule); i++ {
			r := h.rule[i]
			// 如果当前规则匹配格式
			if r.Know(expr) {
				err := r.Check(expr, itf, ivf)
				if err != nil {
					return err
				}
				checked = true
				continue
			}
		}
		// 如果规则匹配完仍然没人认识这个表达式
		if !checked {
			return errors.New("表达式" + expr + "无法被识别")
		}
	}
	return nil
}
