package fSchedule

import (
	"encoding/json"
	"fmt"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/stopwatch"
	"time"
)

// 日志结构
type JobLog struct {
	TaskId   int64
	Name     string
	LogLevel Enum
	Content  string // 日志内容
}

type logDto struct {
	TaskId int64
	Name   string
	Log    collections.List[logBody] // 客户端动态注册任务
}

type logBody struct {
	LogLevel Enum
	Content  string
	CreateAt time.Time
}

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

// logReport 上传日志报告
func (receiver *JobContext) logReport() {
	select {
	case log := <-receiver.LogQueue:
		logBody := logBody{
			LogLevel: log.LogLevel,
			Content:  log.Content,
			CreateAt: time.Now(),
		}
		logMsg := logDto{
			TaskId: log.TaskId,
			Name:   log.Name,
			Log:    collections.NewList(logBody),
		}
		jsonByte, _ := json.Marshal(logMsg)
		defaultServer.logReport(jsonByte)
	}
}

// SetNextAt 设置下次运行时间
func (receiver *JobContext) SetNextAt(t time.Time) {
	receiver.nextTimespan = t.UnixMicro()
}

// SetProgress 设置任务进度
func (receiver *JobContext) SetProgress(progress int) {
	receiver.progress = progress
}

// SetProgress 报告任务结果
func (receiver *JobContext) report() bool {
	jsonByte, _ := json.Marshal(receiver.getReport())
	apiResponse, _ := defaultServer.taskReport(jsonByte)

	// 上传日志
	receiver.logReport()
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
func (receiver *JobContext) log(logLevel Enum, contents ...any) {
	jobLog := JobLog{
		TaskId:   receiver.Id,
		Name:     receiver.Name,
		LogLevel: logLevel,
		Content:  fmt.Sprint(contents...),
	}
	receiver.LogQueue <- jobLog
}

// Trace 打印Trace日志
func (receiver *JobContext) Trace(content ...any) {
	receiver.log(Trace, content)
	flog.Trace(content...)
}

// Tracef 打印Trace日志
func (receiver *JobContext) Tracef(format string, a ...any) {
	receiver.log(Trace, fmt.Sprintf(format, a...))
	flog.Tracef(format, a...)
}

// Debug 打印Debug日志
func (receiver *JobContext) Debug(contents ...any) {
	receiver.log(Debug, contents)
	flog.Debug(contents...)
}

// Debugf 打印Debug日志
func (receiver *JobContext) Debugf(format string, a ...any) {
	receiver.log(Debug, fmt.Sprintf(format, a...))
	flog.Debugf(format, a...)
}

// Info 打印Info日志
func (receiver *JobContext) Info(contents ...any) {
	receiver.log(Information, contents)
	flog.Info(contents...)
}

// Infof 打印Info日志
func (receiver *JobContext) Infof(format string, a ...any) {
	receiver.log(Information, fmt.Sprintf(format, a...))
	flog.Infof(format, a...)
}

// Warning 打印Warning日志
func (receiver *JobContext) Warning(contents ...any) {
	receiver.log(Warning, contents)
	flog.Warning(contents...)
}

// Warningf 打印Warning日志
func (receiver *JobContext) Warningf(format string, a ...any) {
	receiver.log(Warning, fmt.Sprintf(format, a...))
	flog.Warningf(format, a...)
}

// Error 打印Error日志
func (receiver *JobContext) Error(contents ...any) error {
	receiver.log(Error, contents)
	return flog.Error(contents...)
}

// Errorf 打印Error日志
func (receiver *JobContext) Errorf(format string, a ...any) error {
	receiver.log(Error, fmt.Sprintf(format, a...))
	return flog.Errorf(format, a...)
}

// Critical 打印Critical日志
func (receiver *JobContext) Critical(contents ...any) {
	receiver.log(Critical, contents)
	flog.Critical(contents...)
}

// Criticalf 打印Critical日志
func (receiver *JobContext) Criticalf(format string, a ...any) {
	receiver.log(Critical, fmt.Sprintf(format, a...))
	flog.Criticalf(format, a...)
}
