package fSchedule

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/stopwatch"
)

// 当前正在执行的Job列表
// key:taskId
var jobList = collections.NewDictionary[int64, *Job]()

type Job struct {
	ClientJob  ClientJob
	jobContext *JobContext
}

// 接受来自服务端的任务
func invokeJob(task TaskEO) {
	clientJob := defaultClient.ClientJobs.Where(func(item ClientJob) bool {
		return item.Name == task.Name
	}).First()

	job := &Job{
		ClientJob: clientJob,
		jobContext: &JobContext{ // 构造上下文
			Id:           task.Id,
			Name:         clientJob.Name,
			Data:         task.Data,
			nextTimespan: 0,
			progress:     0,
			status:       Working,
			sw:           stopwatch.New(),
		},
	}
	jobList.Add(task.Id, job)
	go job.Run()
}

func (receiver *Job) Run() {
	defer func() {
		if receiver.jobContext.report() {
			jobList.Remove(receiver.jobContext.Id)
		}
	}()

	// 执行任务并拿到结果
	exception.Try(func() {
		if receiver.ClientJob.jobFunc(receiver.jobContext) {
			receiver.jobContext.status = Success
		} else {
			receiver.jobContext.status = Fail
		}
	}).CatchException(func(exp any) {
		receiver.jobContext.status = Fail
	})

	flog.Infof("任务：%s %d，耗时：%s，结果：%s", receiver.jobContext.Name, receiver.jobContext.Id, receiver.jobContext.sw.GetMillisecondsText(), receiver.jobContext.status.String())
}
