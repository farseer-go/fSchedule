package fSchedule

import (
	"encoding/json"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/core/eumLogLevel"
	"github.com/farseer-go/fs/timingWheel"
	"time"
)

// 任务日志
var logQueue chan logContent

type logReportDTO struct {
	Logs []logContent
}

type logContent struct {
	TaskId   int64  // 主键
	Name     string // 实现Job的特性名称（客户端识别哪个实现类）
	Ver      int    // 版本
	LogLevel eumLogLevel.Enum
	CreateAt int64
	Content  string
}

// enableReportLog 开启上传日志报告
func enableReportLog() {
	for {
		<-timingWheel.Add(500 * time.Millisecond).C
		lstLogs := collections.NewList[logContent]()
		tw := timingWheel.Add(1 * time.Second)
		isContinue := true // 标记是否一直循环读取，当大于1秒，或者取出10条日志时，上传日志
		for isContinue {
			select {
			case log := <-logQueue: // 从通道中取出日志数据
				lstLogs.Add(log)
				isContinue = lstLogs.Count() < 10
			case <-tw.C: // 1秒时间到了
				isContinue = false
			}
		}

		if lstLogs.Count() > 0 {
			dto := logReportDTO{Logs: lstLogs.ToArray()}
			jsonByte, _ := json.Marshal(dto)
			_, _ = defaultServer.logReport(jsonByte)
		}
	}
}
