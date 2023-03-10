package fSchedule

import (
	"encoding/json"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/flog"
	"os"
)

var defaultClient *clientVO

// 客户端配置
type clientVO struct {
	ClientId   int64                       // 客户端ID
	ClientName string                      // 客户端名称
	ClientIp   string                      // 客户端IP
	ClientPort int                         // 客户端端口
	ClientJobs collections.List[ClientJob] // 客户端动态注册任务
}

func NewClient() {
	hostname, _ := os.Hostname()
	defaultClient = &clientVO{
		ClientId:   fs.AppId,
		ClientName: hostname,
		ClientIp:   fs.AppIp,
		ClientPort: 8888, // 先填写默认值
		ClientJobs: collections.NewList[ClientJob](),
	}

	// 如果手动配置了客户端IP，则覆盖
	clientIp := configure.GetString("FSchedule.ClientIp")
	if clientIp != "" {
		defaultClient.ClientIp = clientIp
	}

	// 如果手动配置了客户端端口，则覆盖
	clientPort := configure.GetInt("FSchedule.ClientPort")
	if clientPort > 0 {
		defaultClient.ClientPort = clientPort
	}
}

type JobFunc func(jobContext *JobContext) bool

type ClientJob struct {
	Name     string // 任务名称
	Caption  string // 任务标题
	Ver      int    // 任务版本
	Cron     string // 任务执行表达式
	StartAt  int64  // 任务开始时间（时间戳秒）
	IsEnable bool   // 任务是否启用
	jobFunc  JobFunc
}

func GetClient() *clientVO {
	return defaultClient
}

// AddJob 客户端支持的任务
func AddJob(isEnable bool, name, caption string, ver int, cron string, startAt int64, job JobFunc) {
	clientJob := ClientJob{
		Name:     name,
		IsEnable: isEnable,
		Caption:  caption,
		Ver:      ver,
		Cron:     cron,
		StartAt:  startAt,
		jobFunc:  job,
	}
	defaultClient.ClientJobs.Add(clientJob)
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
	apiResponse, _ := defaultServer.registry(jsonByte)
	if apiResponse.StatusCode != 200 {
		flog.Panic("注册失败，服务端状态码为：", apiResponse.StatusCode)
	}
	flog.ComponentInfo("fSchedule", "客户端注册成功！")
}

// LogoutClient 客户端下线
func (receiver *clientVO) LogoutClient() {
	jsonByte, _ := json.Marshal(map[string]any{"clientId": receiver.ClientId})
	apiResponse, err := defaultServer.logout(jsonByte)
	flog.Panic(err)
	if apiResponse.StatusCode != 200 {
		flog.Panic("下线失败，服务端状态码为：", apiResponse.StatusCode)
	}
	flog.ComponentInfo("fSchedule", "客户端下线成功！")
}
