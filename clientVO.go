package fSchedule

import (
	"encoding/json"
	"github.com/farseer-go/fs/flog"
)

var defaultClient clientVO

// 客户端配置
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
func (receiver *clientVO) getHttpHead() map[string]any {
	return map[string]any{
		"ClientIp":   receiver.ClientIp,
		"ClientId":   receiver.ClientId,
		"ClientName": receiver.ClientName,
		"ClientPort": receiver.ClientPort,
	}
}

// RegistryClient 注册客户端
func (receiver *clientVO) RegistryClient() {
	jsonByte, _ := json.Marshal(receiver)
	apiResponse, err := defaultServer.registry(jsonByte)
	flog.Panic(err)
	if apiResponse.StatusCode != 200 {
		flog.Panic("注册失败，服务端状态码为：", apiResponse.StatusCode)
	}
	flog.ComponentInfo("fSchedule", "客户端注册成功！")
}
