package cache

import (
	"github.com/Lyo-Shur/gutils/task"
	"sync"
	"time"
)

func init() {
	// 启动每分钟一次的定时清理任务
	task.Run(1*time.Minute, func() {
		m := GetHolder().M
		m.Range(func(k interface{}, v interface{}) bool {
			key := k.(string)
			value := v.(data)
			// 死亡时间大于当前时间
			if (!value.DeathTime.IsZero()) && value.DeathTime.Before(time.Now()) {
				m.Delete(key)
			}
			return true
		})
	})
}

// 缓存数据包装类
type data struct {
	// 真实的数据
	Data interface{}
	// 销毁时间
	DeathTime time.Time
}

// 数据持有结构体
var h *holder

// 数据持有结构体
type holder struct {
	M sync.Map
}

var once sync.Once

// 获取数据持有结构体
func GetHolder() *holder {
	once.Do(func() {
		h = &holder{}
	})
	return h
}

// 获取缓存数据
func (h *holder) Get(key string) data {
	v, ok := h.M.Load(key)
	if ok {
		d := v.(data)
		return d
	}
	return data{}
}

// 设置缓存数据
// 不设置存活时间的数据为永生状态
func (h *holder) Set(key string, value interface{}) {
	// 创建数据包装体
	d := data{}
	d.Data = value
	h.M.Store(key, d)
}

// 删除缓存的数据
func (h *holder) Delete(key string) {
	h.M.Delete(key)
}

// 设置缓存数据带生命周期
func (h *holder) SetSurviveTime(key string, duration time.Duration) {
	v, ok := h.M.Load(key)
	if !ok {
		return
	}
	d := v.(data)
	d.DeathTime = time.Now().Add(duration)
	h.M.Store(key, d)
}

// 设置为永生
func (h *holder) SetEternalLife(key string) {
	v, ok := h.M.Load(key)
	if !ok {
		return
	}
	d := v.(data)
	d.DeathTime = time.Time{}
	h.M.Store(key, d)
}
