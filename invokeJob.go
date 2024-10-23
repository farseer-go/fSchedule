package fSchedule

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/farseer-go/fSchedule/executeStatus"
	"github.com/farseer-go/fs/asyncLocal"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/timingWheel"
	"github.com/farseer-go/fs/trace"
)

var workCount int
var lock = &sync.RWMutex{}

// var taskList = make(map[int64]*JobContext)
var taskList = sync.Map{}

// 接受来自服务端的任务
func invokeJob(clientVO ClientVO, task taskDTO) {
	if task.Id == 0 || task.Name == "" {
		return
	}

	ctx, cancel := context.WithCancel(clientVO.client.Ctx)
	clientVO.Data = task.Data
	jobContext := &JobContext{ // 构造上下文
		Id:           task.Id,
		Name:         task.Name,
		Ver:          clientVO.Ver,
		Data:         task.Data,
		Caption:      task.Caption,
		StartAt:      task.StartAt,
		Ctx:          ctx,
		status:       executeStatus.Working,
		cancel:       cancel,
		clientJob:    clientVO,
		traceManager: container.Resolve[trace.IManager](),
	}
	taskList.Store(task.Id, jobContext)

	defer func() {
		// 任务报告完后，移除本次任务
		clientVO.report(jobContext)
		taskList.Delete(task.Id)
		asyncLocal.Release()
	}()

	taskStartAtSince := time.Since(jobContext.StartAt)
	if taskStartAtSince.Microseconds() > 0 {
		flog.Warningf("任务组：%s %d 延迟：%s", jobContext.Name, jobContext.Id, taskStartAtSince.String())
	} else {
		// 为了保证任务不被延迟，服务端会提前下发任务，需要客户端做休眠等待
		<-timingWheel.AddTimePrecision(jobContext.StartAt).C
	}

	// 工作中任务+1
	workCount++
	defer func() {
		// 工作中任务-1
		workCount--
	}()

	// 链路追踪
	entryFSchedule := jobContext.traceManager.EntryFSchedule(jobContext.Name, jobContext.Id, jobContext.Data.ToMap())
	defer entryFSchedule.End(nil)

	// 执行任务并拿到结果
	exception.Try(func() {
		// 通知调度中心，我开始执行了
		clientVO.report(jobContext)
		// 执行任务
		if jobContext.clientJob.jobFunc(jobContext) {
			jobContext.status = executeStatus.Success
		} else {
			jobContext.status = executeStatus.Fail
		}
	}).CatchException(func(exp any) {
		jobContext.status = executeStatus.Fail
		jobContext.Remark("%s", exp)
		jobContext.Error(exp)
		entryFSchedule.Error(fmt.Errorf("%s", exp))
	})
}
