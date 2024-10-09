package fSchedule

import (
	"fmt"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/core/eumLogLevel"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/utils/ws"
	"time"
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
	return receiver.client.Send(sendDTO{Type: -1, Registry: registryDTO{ClientName: core.AppName, Job: *receiver}})
}

// SetProgress 报告任务结果
func (receiver *ClientVO) report(jobContext *JobContext) {
	err := receiver.client.Send(sendDTO{
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
	})
	if err != nil {
		flog.Warningf("向调度中心报告任务结果时失败：%s", err.Error())
	}
}

// log 记录日志
func (receiver *ClientVO) log(jobContext *JobContext, logLevel eumLogLevel.Enum, contents ...any) {
	err := receiver.client.Send(sendDTO{
		Type: 1,
		Log: logDTO{
			TaskId:   jobContext.Id,
			Name:     jobContext.Name,
			Ver:      jobContext.Ver,
			Caption:  jobContext.Caption,
			Data:     jobContext.Data,
			LogLevel: logLevel,
			CreateAt: time.Now().UnixMilli(),
			Content:  fmt.Sprint(contents...),
		},
	})

	if err != nil {
		flog.Warningf("向调度中心报告任务结果时失败：%s", err.Error())
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
	return resourceVO{
		QueueCount: taskListLength - workCount,
		WorkCount:  workCount,
	}
}
