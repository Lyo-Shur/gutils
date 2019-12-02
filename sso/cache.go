package sso

import (
	"github.com/Lyo-Shur/gutils/cache"
	"github.com/satori/go.uuid"
	"log"
	"time"
)

// 使用Cache实现的Holder
type CacheHolder struct{}

// 实现接口 根据用户身份获得票据
func (ch *CacheHolder) GetTicket(userInfo interface{}) string {
	// 生成票据
	uid, err := uuid.NewV4()
	// 此处就不向上抛这个异常了
	// 生成UUID不考虑出错的情况
	if err != nil {
		log.Println(err)
	}

	// 获取票据
	ticket := uid.String()

	// 缓存票据与用户名(5分钟失效)
	c := cache.GetHolder()
	c.Set(ticket, userInfo)
	c.SetSurviveTime(ticket, time.Second*5)

	return ticket
}

// 实现接口 使用票据获得用户身份信息
func (ch *CacheHolder) GetUserInfo(ticket string) interface{} {
	// 使用票据换取对应的用户信息
	holder := cache.GetHolder()
	data := holder.Get(ticket).Data
	// 使用后删除票据
	holder.Delete(ticket)
	return data
}
