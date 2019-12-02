package form

// 文件保存配置
type SaveConfig struct {
	// 服务器访问路径
	ServerPath string
	// 磁盘保存路径
	DiskPath string
	// 计算文件保存路径
	// 参数 请求中的文件名
	// 返回值 次级路径 文件名
	GetSavePath func(FileName string) (string, string, error)
}

// 复制-文件保存配置-方法
func (s *SaveConfig) Clone() SaveConfig {
	sc := SaveConfig{}
	sc.ServerPath = s.ServerPath
	sc.DiskPath = s.DiskPath
	sc.GetSavePath = s.GetSavePath
	return sc
}
