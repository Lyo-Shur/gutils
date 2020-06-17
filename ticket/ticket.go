package ticket

import (
	"github.com/Lyo-Shur/gutils/cache"
	"strings"
	"sync"
	"time"
)

// 票据保持者
type Holder interface {
	// 根据用户身份创建票据
	CreateTicket(userId string, duration time.Duration) string
	// 使用票据获得用户身份信息
	GetUserId(ticket string) interface{}
	// 根据用户身份获取票据
	GetTicket(userId string) []string
	// 删除票据
	DeleteTicket(ticket string)
}

// 使用Cache实现的Holder
type CacheHolder struct {
	// UUID算法
	UUID func() string
	// 缓存模块
	Cache cache.Cache
}

var mutex sync.Mutex

// 实现接口 根据用户身份创建票据
func (ch *CacheHolder) CreateTicket(userId string, duration time.Duration) string {
	// 获取票据
	ticket := ch.UUID()

	mutex.Lock()

	// 缓存票据与用户名
	ch.Cache.Set(ticket, userId)
	ch.Cache.SetSurviveTime(ticket, duration)

	// 缓存用户名与票据
	if !ch.Cache.Exist(userId) {
		ch.Cache.Set(userId, ticket)
	} else {
		oldValue := ch.Cache.Get(userId).Data.(string)
		ch.Cache.Set(userId, oldValue+";"+ticket)
	}

	mutex.Unlock()

	// 缓存用户名与票据
	return ticket
}

// 实现接口 使用票据获得用户身份信息
func (ch *CacheHolder) GetUserId(ticket string) interface{} {
	// 使用票据换取对应的用户信息
	return ch.Cache.Get(ticket).Data
}

// 实现接口 根据用户身份获取票据
func (ch *CacheHolder) GetTicket(userId string) []string {
	mutex.Lock()

	tickets := make([]string, 0)
	// 使用用户信息换取对应的票据
	data := ch.Cache.Get(userId).Data.(string)
	cacheTickets := strings.Split(data, ";")
	for i := range cacheTickets {
		ticket := cacheTickets[i]
		if ch.Cache.Exist(ticket) {
			tickets = append(tickets, ticket)
		}
	}
	ch.Cache.Set(userId, strings.Join(tickets, ";"))

	mutex.Unlock()
	return tickets
}

// 实现接口 删除票据
func (ch *CacheHolder) DeleteTicket(ticket string) {
	ch.Cache.Delete(ticket)
}
