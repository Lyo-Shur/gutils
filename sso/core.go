package sso

// sso票据保持者
type Holder interface {
	// 根据用户身份获得票据
	GetTicket(userInfo interface{}) string
	// 使用票据获得用户身份信息
	GetUserInfo(ticket string) interface{}
}
