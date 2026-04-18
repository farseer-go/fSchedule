package fSchedule

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/core/eumLogLevel"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/utils/ws"
)

type ClientVO struct {
	Name     string                                 // 任务名称
	Ver      int                                    // 任务版本
	Caption  string                                 // 任务标题
	Cron     string                                 // 任务执行表达式
	StartAt  int64                                  // 任务开始时间（时间戳秒）
	IsEnable bool                                   // 任务是否启用
	Data     collections.Dictionary[string, string] // 第一次注册时使用

	jobFunc JobFunc    // 任务执行函数
	client  *ws.Client // ws客户端
}

func (receiver *ClientVO) registry() error {
	return receiver.client.Send(sendDTO{Type: -1, Registry: registryDTO{ClientName: core.AppName, ClientIp: core.AppIp, Job: *receiver}})
}

// SetProgress 报告任务结果
func (receiver *ClientVO) report(jobContext *JobContext) {
	dto := sendDTO{
		Type: 0,
		TaskReport: taskReportDTO{
			Id:           jobContext.Id,
			Name:         jobContext.Name,
			Data:         jobContext.Data,
			NextTimespan: jobContext.nextTimespan,
			Progress:     jobContext.progress,
			Status:       jobContext.status,
			FailRemark:   jobContext.failRemark,
			resourceVO:   getResource(),
		},
	}
	err := receiver.client.Send(dto)
	if err != nil {
		// 当前持有的连接可能是旧连接（调度中心重启后重连，但 invokeJob goroutine 仍持有旧 ClientVO 副本）
		// 从 mapClient 取最新的 ClientVO，用新连接重试一次
		if latest, ok := mapClient.Load(receiver.Name); ok {
			latestVO := latest.(ClientVO)
			if latestVO.client != receiver.client {
				if retryErr := latestVO.client.Send(dto); retryErr == nil {
					return
				}
			}
		}
		flog.Warningf("向调度中心报告任务结果时失败：%s", err.Error())
	}
}

// log 记录日志
func (receiver *ClientVO) log(jobContext *JobContext, logLevel eumLogLevel.Enum, contents ...any) {
	content := fmt.Sprint(contents...)
	container.Resolve[trace.IManager]().TraceHand(content).End(nil)
	dto := sendDTO{
		Type: 1,
		Log: logDTO{
			TaskId:   jobContext.Id,
			Name:     jobContext.Name,
			Ver:      jobContext.Ver,
			Caption:  jobContext.Caption,
			Data:     jobContext.Data,
			LogLevel: logLevel,
			CreateAt: time.Now().UnixMilli(),
			Content:  content,
		},
	}
	err := receiver.client.Send(dto)
	if err != nil {
		// 与 report 相同：重连后旧 goroutine 持有旧连接，尝试从 mapClient 取新连接重发
		if latest, ok := mapClient.Load(receiver.Name); ok {
			latestVO := latest.(ClientVO)
			if latestVO.client != receiver.client {
				if retryErr := latestVO.client.Send(dto); retryErr == nil {
					return
				}
			}
		}
		flog.Warningf("向调度中心上报日志时失败：%s", err.Error())
	}
}

// 获取当前客户端的环境信息
func getResource() resourceVO {
	// 计算长度
	taskListLength := 0
	taskList.Range(func(k, v interface{}) bool {
		taskListLength++
		return true
	})
	wc := int(atomic.LoadInt64(&workCount))
	return resourceVO{
		QueueCount: taskListLength - wc,
		WorkCount:  wc,
	}
}
