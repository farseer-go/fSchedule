package fSchedule

import (
	"context"
	"fmt"
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fSchedule/executeStatus"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core/eumLogLevel"
	"github.com/farseer-go/fs/trace"
)

// JobContext 任务在运行时，更新状态
type JobContext struct {
	Id           int64                                  // 主键
	Name         string                                 // 实现Job的特性名称（客户端识别哪个实现类）
	Ver          int                                    // 任务版本
	Caption      string                                 // 任务标题
	Data         collections.Dictionary[string, string] // 数据
	nextTimespan int64                                  // 下次执行时间
	progress     int                                    // 当前进度
	status       executeStatus.Enum                     // 执行状态
	StartAt      time.Time                              // 任务开始时间
	failRemark   string                                 // 失败原因
	cancel       context.CancelFunc                     // 服务端通知Kill时，将调用此函数
	Ctx          context.Context                        // 客户端执行时，需要检查ctx是否被Cancel
	clientJob    ClientVO                               // 客户端
	traceManager trace.IManager                         // 链路追踪
}

// SetNextAt 设置下次运行时间
func (receiver *JobContext) SetNextAt(t time.Time) {
	receiver.nextTimespan = t.UnixMilli()
}

// SetProgress 设置任务进度
func (receiver *JobContext) SetProgress(progress int) {
	receiver.progress = progress
}

// Remark 报告失败原因
func (receiver *JobContext) Remark(format string, a ...any) {
	receiver.failRemark = fmt.Sprintf(format, a...)
	container.Resolve[trace.IManager]().TraceHand(receiver.failRemark).End(nil)
}

// Trace 打印Trace日志
func (receiver *JobContext) Trace(contents ...any) {
	receiver.clientJob.log(receiver, eumLogLevel.Trace, contents...)
}

// Tracef 打印Trace日志
func (receiver *JobContext) Tracef(format string, a ...any) {
	receiver.clientJob.log(receiver, eumLogLevel.Trace, fmt.Sprintf(format, a...))
}

// Debug 打印Debug日志
func (receiver *JobContext) Debug(contents ...any) {
	receiver.clientJob.log(receiver, eumLogLevel.Debug, contents...)
}

// Debugf 打印Debug日志
func (receiver *JobContext) Debugf(format string, a ...any) {
	receiver.clientJob.log(receiver, eumLogLevel.Debug, fmt.Sprintf(format, a...))
}

// Info 打印Info日志
func (receiver *JobContext) Info(contents ...any) {
	receiver.clientJob.log(receiver, eumLogLevel.Information, contents...)
}

// Infof 打印Info日志
func (receiver *JobContext) Infof(format string, a ...any) {
	receiver.clientJob.log(receiver, eumLogLevel.Information, fmt.Sprintf(format, a...))
}

// Warning 打印Warning日志
func (receiver *JobContext) Warning(contents ...any) {
	receiver.clientJob.log(receiver, eumLogLevel.Warning, contents...)
}

// Warningf 打印Warning日志
func (receiver *JobContext) Warningf(format string, a ...any) {
	receiver.clientJob.log(receiver, eumLogLevel.Warning, fmt.Sprintf(format, a...))
}

// Error 打印Error日志
func (receiver *JobContext) Error(contents ...any) {
	receiver.clientJob.log(receiver, eumLogLevel.Error, contents...)
}

// Errorf 打印Error日志
func (receiver *JobContext) Errorf(format string, a ...any) {
	receiver.clientJob.log(receiver, eumLogLevel.Error, fmt.Sprintf(format, a...))
}

// Critical 打印Critical日志
func (receiver *JobContext) Critical(contents ...any) {
	receiver.clientJob.log(receiver, eumLogLevel.Critical, contents...)
}

// Criticalf 打印Critical日志
func (receiver *JobContext) Criticalf(format string, a ...any) {
	receiver.clientJob.log(receiver, eumLogLevel.Critical, fmt.Sprintf(format, a...))
}
