package bean

import (
	"reflect"
	"sync"
)

// bean的工厂
var f *Factory

// 工厂结构体
// 此结构体用来维持Service层和DAO层引用
type Factory struct {
	m map[string]interface{}
}

// 使用sync.Once保证factory单例
var once sync.Once

// 获取bean工厂
func GetFactory() *Factory {
	once.Do(func() {
		f = &Factory{}
		f.m = make(map[string]interface{})
	})
	return f
}

// 注册bean
// 内部使用map进行维持数据结构
// 所以当注册的键重复时可能导致bean被替换，而发生错误
func (f *Factory) Register(s string, i interface{}) {
	injection(i, f)
	f.m[s] = i
}

// 获取bean
// 根据注册时的键获取bean
func (f *Factory) Get(s string) interface{} {
	return f.m[s]
}

// 执行自动注入
// 将i中所需的接口通过反射注入到i
func injection(i interface{}, f *Factory) {
	// 首先反射获取此接口的Type Value
	it := reflect.TypeOf(i)
	iv := reflect.ValueOf(i)

	// 遍历字段
	// 因为注册的都是结构体的指针 所以此处需要先.Elem
	for j := 0; j < it.Elem().NumField(); j++ {
		// 当前字段的类型
		item := it.Elem().Field(j).Type

		// 对已注册的bean进行遍历，查找是否有当前字段需要注入的部分
		for _, v := range f.m {
			// 获取当前bean的Type Value
			vt := reflect.TypeOf(v)
			vv := reflect.ValueOf(v)

			// 如果当前bean能够转化为item
			if vt.ConvertibleTo(item) {
				// 将当前bean转化并注入到i中
				cv := vv.Convert(item)
				iv.Elem().Field(j).Set(cv)
			}
		}
	}
}
