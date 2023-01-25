package fSchedule

var defaultClient clientVO

type clientVO struct {
	ClientId   int64  // 客户端ID
	ClientName string // 客户端名称
	ClientIp   string // 客户端IP
	ClientPort int64  // 客户端端口
}

func GetClient() clientVO {
	return defaultClient
}

// 转换成http head
func (receiver clientVO) getHttpHead() map[string]any {
	return map[string]any{
		"ClientIp":   defaultClient.ClientIp,
		"ClientId":   defaultClient.ClientId,
		"ClientName": defaultClient.ClientName,
		"ClientPort": defaultClient.ClientPort,
	}
}

// RegistryClient 注册客户端
func (receiver clientVO) RegistryClient() {
	
}
