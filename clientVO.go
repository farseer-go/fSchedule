package fSchedule

import (
	"encoding/json"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/stopwatch"
	"strings"
	"time"
)

var defaultClient *clientVO
var isRegistryJobCount int // 向调度中心注册的JOB数量

// 客户端配置
type clientVO struct {
	ClientId   int64                       // 客户端ID
	ClientName string                      // 客户端名称
	ClientIp   string                      // 客户端IP
	ClientPort int                         // 客户端端口
	ClientJobs collections.List[ClientJob] // 客户端动态注册任务
}

func NewClient() {
	defaultClient = &clientVO{
		ClientId:   fs.AppId,
		ClientName: fs.AppName,
		ClientIp:   "",
		ClientPort: 8888, // 先填写默认值
		ClientJobs: collections.NewList[ClientJob](),
	}

	// 优先使用本地IP
	if strings.HasPrefix(fs.AppIp, "192.168.") || strings.HasPrefix(fs.AppIp, "172.20.") || strings.HasPrefix(fs.AppIp, "10.") {
		defaultClient.ClientIp = fs.AppIp
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

// JobFunc 客户端要执行的JOB
type JobFunc func(jobContext *JobContext) bool

type ClientJob struct {
	Name     string                                 // 任务名称
	Ver      int                                    // 任务版本
	Caption  string                                 // 任务标题
	Cron     string                                 // 任务执行表达式
	StartAt  int64                                  // 任务开始时间（时间戳秒）
	IsEnable bool                                   // 任务是否启用
	Data     collections.Dictionary[string, string] // 第一次注册时使用
	jobFunc  JobFunc
}

func GetClient() *clientVO {
	return defaultClient
}

type Option struct {
	StartAt int64                                  // 任务开始时间（时间戳秒）
	Data    collections.Dictionary[string, string] // 第一次注册时使用
}
type options func(opt *Option)

// AddJob 客户端支持的任务
func AddJob(isEnable bool, name, caption string, ver int, cron string, job JobFunc, ops ...options) {
	// 说明没有启用调度中心（没有依赖模块）
	if defaultClient == nil {
		return
	}
	// 设置额度参数
	opt := &Option{Data: collections.NewDictionary[string, string]()}
	for _, op := range ops {
		op(opt)
	}
	defaultClient.ClientJobs.Add(ClientJob{
		Name:     name,
		IsEnable: isEnable,
		Caption:  caption,
		Ver:      ver,
		Cron:     cron,
		jobFunc:  job,
		StartAt:  opt.StartAt,
		Data:     opt.Data,
	})

	// 如果是调试状态，则模拟调度
	if configure.GetBool("FSchedule.Debug.Enable") {
		jobContext := &JobContext{
			Id:           888,
			TaskGroupId:  888,
			Ver:          888,
			Name:         name,
			Data:         collections.NewDictionary[string, string](),
			nextTimespan: 0,
			progress:     0,
			status:       Working,
			sw:           stopwatch.New(),
			StartAt:      time.Now(),
		}
		for k, v := range configure.GetSubNodes("FSchedule.Debug." + name) {
			jobContext.Data.Add(k, parse.ToString(v))
		}
		job(jobContext)
	}

	// 说明已经向调度中心注册过，之后又添加了新的任务
	if isRegistryJobCount > 0 {
		err := defaultClient.RegistryClient()
		flog.ErrorIfExists(err)
	}
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
func (receiver *clientVO) RegistryClient() error {
	jsonByte, _ := json.Marshal(receiver)
	apiResponse, _ := defaultServer.registry(jsonByte)
	if apiResponse.StatusCode != 200 {
		return flog.Errorf("注册失败：%d %s", apiResponse.StatusCode, apiResponse.StatusMessage)
	}
	receiver.ClientIp = apiResponse.Data.ClientIp
	receiver.ClientPort = apiResponse.Data.ClientPort

	// 向调度中心注册的JOB数量
	isRegistryJobCount = receiver.ClientJobs.Count()
	return nil
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
