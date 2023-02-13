package fSchedule

import (
	"encoding/json"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/flog"
)

var defaultClient clientVO

// 客户端配置
type clientVO struct {
	ClientId   int64                       // 客户端ID
	ClientName string                      // 客户端名称
	ClientIp   string                      // 客户端IP
	ClientPort int                         // 客户端端口
	ClientJobs collections.List[ClientJob] // 客户端动态注册任务
}

type ClientJob struct {
	Name    string // 任务名称
	Caption string // 任务标题
	Ver     int    // 任务版本
	Cron    string // 任务执行表达式
	StartAt int64  // 任务开始时间（时间戳秒）
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

// LogoutClient 客户端下线
func (receiver *clientVO) LogoutClient() {
	jsonByte, _ := json.Marshal(receiver)
	apiResponse, err := defaultServer.logout(jsonByte)
	flog.Panic(err)
	if apiResponse.StatusCode != 200 {
		flog.Panic("下线失败，服务端状态码为：", apiResponse.StatusCode)
	}
	flog.ComponentInfo("fSchedule", "客户端下线成功！")
}
