package fSchedule

import (
	"fmt"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/stopwatch"
	"github.com/farseer-go/fs/timingWheel"
	"github.com/farseer-go/fs/trace"
	"sync"
	"time"
)

// 当前正在执行的Job列表
// key:taskId
var jobList = collections.NewDictionary[int64, *Job]()
var lock = &sync.RWMutex{}

type Job struct {
	ClientJob  ClientJob
	jobContext *JobContext
	// 链路追踪
	traceManager trace.IManager
}

// 接受来自服务端的任务
func invokeJob(task TaskEO) {
	clientJob := defaultClient.ClientJobs.Where(func(item ClientJob) bool {
		return item.Name == task.Name
	}).First()

	if task.Id == 0 || task.TaskGroupId == 0 || task.Name == "" || clientJob.Name == "" {
		return
	}
	job := &Job{
		ClientJob: clientJob,
		jobContext: &JobContext{ // 构造上下文
			Id:           task.Id,
			TaskGroupId:  task.TaskGroupId,
			Ver:          clientJob.Ver,
			Name:         clientJob.Name,
			Data:         task.Data,
			StartAt:      task.StartAt,
			nextTimespan: 0,
			progress:     0,
			status:       Working,
			sw:           stopwatch.New(),
		},
		traceManager: container.Resolve[trace.IManager](),
	}
	go job.Run()

	lock.Lock()
	defer lock.Unlock()
	jobList.Add(task.Id, job)
}

func (receiver *Job) Run() {
	defer func() {
		if receiver.jobContext.report() {
			lock.Lock()
			jobList.Remove(receiver.jobContext.Id)
			lock.Unlock()
		}
	}()

	taskStartAtSince := time.Since(receiver.jobContext.StartAt)
	if taskStartAtSince.Microseconds() > 0 {
		flog.Warningf("任务组：%s（%d） %d 延迟：%s", receiver.jobContext.Name, receiver.jobContext.TaskGroupId, receiver.jobContext.Id, taskStartAtSince.String())
	}

	// 为了保证任务不被延迟，服务端会提前下发任务，需要客户端做休眠等待
	<-timingWheel.AddTimePrecision(receiver.jobContext.StartAt).C

	// 链路追踪
	entryFSchedule := receiver.traceManager.EntryFSchedule(receiver.jobContext.Name, receiver.jobContext.TaskGroupId, receiver.jobContext.Id)
	// 执行任务并拿到结果
	exception.Try(func() {
		receiver.jobContext.sw.Start()
		// 执行任务
		if receiver.ClientJob.jobFunc(receiver.jobContext) {
			receiver.jobContext.status = Success
		} else {
			receiver.jobContext.status = Fail
		}
	}).CatchException(func(exp any) {
		receiver.jobContext.status = Fail
		entryFSchedule.Error(fmt.Errorf("%s", exp))
	})
	entryFSchedule.End()
	flog.ComponentInfof("fSchedule", "任务：%s（%d） %d，耗时：%s，结果：%s", receiver.jobContext.Name, receiver.jobContext.TaskGroupId, receiver.jobContext.Id, receiver.jobContext.sw.GetMillisecondsText(), receiver.jobContext.status.String())
}

func getJob(taskId int64) *Job {
	lock.RLock()
	defer lock.RUnlock()
	return jobList.GetValue(taskId)
}
