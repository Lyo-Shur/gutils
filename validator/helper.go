package validator

import (
	"encoding/json"
	"errors"
)

type Config struct {
	// 表达式
	Expr string
	// 消息
	Message string
}

// 校验帮助工具
type helper struct {
	// 结构体字段
	params map[string]string
	// 校验规则
	configs map[string][]Config
	// 内置规则
	rule []Rule
}

// 获取校验帮助工具
func Helper(params map[string]string, configJson string) *helper {
	configs := make(map[string][]Config)

	// 解析配置
	err := json.Unmarshal([]byte(configJson), &configs)
	if err != nil {
		panic(err)
	}

	h := helper{}
	h.params = params
	h.configs = configs
	// 添加默认规则
	h.rule = make([]Rule, 6)
	h.rule[0] = &Required{}
	h.rule[1] = &Ban{}
	h.rule[2] = &Length{}
	h.rule[3] = &Range{}
	h.rule[4] = &DateTime{}
	h.rule[5] = &Regexp{}
	return &h
}

// 添加自定义规则
func (h *helper) AddRule(r Rule) {
	h.rule = append(h.rule, r)
}

// 执行校验
func (h *helper) Check() (bool, string, error) {
	for key, configs := range h.configs {
		for i := 0; i < len(configs); i++ {
			checked := false
			for j := 0; j < len(h.rule); j++ {
				r := h.rule[j]
				// 如果当前规则匹配格式
				if r.Know(configs[i].Expr) {
					b, err := r.Check(configs[i].Expr, h.params, key)
					if err != nil {
						return false, "", err
					}
					checked = true
					if !b {
						if configs[i].Message != "" {
							return false, configs[i].Message, nil
						}
						return false, "Parameter failed validation [rule=" + configs[i].Expr + ", name=" + key + "]", nil
					}
					break
				}
			}
			// 如果规则匹配完仍然没人认识这个表达式
			if !checked {
				msg := "Unknown expression, please check validation rule [rule=" + configs[i].Expr + ", name=" + key + "]"
				return false, "", errors.New(msg)
			}
		}
	}
	return true, "", nil
}
