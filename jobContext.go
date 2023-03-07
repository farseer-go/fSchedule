package fSchedule

import (
	"encoding/json"
	"fmt"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/stopwatch"
	"time"
)

// JobContext 任务在运行时，更新状态
type JobContext struct {
	Id           int64                                  // 主键
	Name         string                                 // 实现Job的特性名称（客户端识别哪个实现类）
	Data         collections.Dictionary[string, string] // 数据
	nextTimespan int64                                  // 下次执行时间
	progress     int                                    // 当前进度
	status       TaskStatus                             // 执行状态
	sw           *stopwatch.Watch                       // 运行时间
	StartAt      time.Time                              // 任务开始时间
	LogQueue     chan JobLog                            // 任务日志
}

// SetNextAt 设置下次运行时间
func (receiver *JobContext) SetNextAt(t time.Time) {
	receiver.nextTimespan = t.UnixMicro()
}

// SetProgress 设置任务进度
func (receiver *JobContext) SetProgress(progress int) {
	receiver.progress = progress
}

// SetProgress 设置任务进度
func (receiver *JobContext) report() bool {
	jsonByte, _ := json.Marshal(receiver.getReport())
	apiResponse, _ := defaultServer.taskReport(jsonByte)
	return apiResponse.StatusCode == 200
}

// getReport 获取DTO
func (receiver *JobContext) getReport() TaskReportDTO {
	return TaskReportDTO{
		Id:           receiver.Id,
		Name:         receiver.Name,
		Data:         receiver.Data,
		NextTimespan: receiver.nextTimespan,
		Progress:     receiver.progress,
		Status:       receiver.status,
		RunSpeed:     receiver.sw.ElapsedMilliseconds(),
	}
}

// log 记录日志
func (receiver *JobContext) Log(logLevel Enum, contents ...any) string {
	content := fmt.Sprint(contents...)
	jobLog := JobLog{
		TaskId:   receiver.Id,
		Name:     receiver.Name,
		LogLevel: logLevel,
		Content:  content,
	}
	receiver.LogQueue <- jobLog
	msg := fmt.Sprintf("%s %s %s\r\n", dateTime.Now().ToString("yyyy-MM-dd hh:mm:ss.ffffff"), flog.Colors[logLevel]("["+logLevel.ToString()+"]"), content)
	return msg
}

func (r Enum) ToString() string {
	switch r {
	case Trace:
		return "Trace"
	case Debug:
		return "Debug"
	case Information:
		return "Info"
	case Warning:
		return "Warn"
	case Error:
		return "Error"
	case Critical:
		return "Critical"
	}
	return "Info"
}

// Trace 打印Trace日志
func (receiver *JobContext) Trace(content ...any) {
	receiver.Log(Trace, content)
}

// Tracef 打印Trace日志
func (receiver *JobContext) Tracef(format string, a ...any) {
	content := fmt.Sprintf(format, a...)
	receiver.Log(Trace, content)
}

// Debug 打印Debug日志
func (receiver *JobContext) Debug(contents ...any) {
	receiver.Log(Debug, contents)
}

// Debugf 打印Debug日志
func (receiver *JobContext) Debugf(format string, a ...any) {
	content := fmt.Sprintf(format, a...)
	receiver.Log(Debug, content)
}

// Info 打印Info日志
func (receiver *JobContext) Info(contents ...any) {
	receiver.Log(Information, contents)
}

// Infof 打印Info日志
func (receiver *JobContext) Infof(format string, a ...any) {
	content := fmt.Sprintf(format, a...)
	receiver.Log(Information, content)
}

// Warning 打印Warning日志
func (receiver *JobContext) Warning(content ...any) {
	receiver.Log(Warning, content)
}

// Warningf 打印Warning日志
func (receiver *JobContext) Warningf(format string, a ...any) {
	content := fmt.Sprintf(format, a...)
	receiver.Log(Warning, content)
}

// Error 打印Error日志
func (receiver *JobContext) Error(contents ...any) error {
	return fmt.Errorf(receiver.Log(Error, contents))
}

// Errorf 打印Error日志
func (receiver *JobContext) Errorf(format string, a ...any) error {
	content := fmt.Sprintf(format, a...)
	return fmt.Errorf(receiver.Log(Error, content))
}

// Panic 打印Error日志并panic
func (receiver *JobContext) Panic(contents ...any) {
	if len(contents) > 0 && contents[0] != nil {
		receiver.Log(Error, contents)
		panic(fmt.Sprint(contents...))
	}
}

// Panicf 打印Error日志并panic
func (receiver *JobContext) Panicf(format string, a ...any) {
	content := fmt.Sprintf(format, a...)
	receiver.Log(Error, content)

	panic(content)
}

// Critical 打印Critical日志
func (receiver *JobContext) Critical(contents ...any) {
	receiver.Log(Critical, contents)
}

// Criticalf 打印Critical日志
func (receiver *JobContext) Criticalf(format string, a ...any) {
	content := fmt.Sprintf(format, a...)
	receiver.Log(Critical, content)
}
