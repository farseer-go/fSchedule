package fSchedule

import (
	"encoding/json"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/stopwatch"
)

// 当前正在执行的Job列表
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
			id:           task.Id,
			name:         clientJob.Name,
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
	flog.Infof("开始执行任务：%d %s:Ver%d 计划时间：%s", receiver.jobContext.id, receiver.ClientJob.Name, receiver.ClientJob.Ver)
	defer func() {
		dto := TaskReportDTO{
			Id:           receiver.jobContext.id,
			Name:         receiver.jobContext.name,
			Data:         receiver.jobContext.Data,
			NextTimespan: receiver.jobContext.nextTimespan,
			Progress:     receiver.jobContext.progress,
			Status:       receiver.jobContext.status,
			RunSpeed:     receiver.jobContext.sw.ElapsedMilliseconds(),
		}
		jsonByte, _ := json.Marshal(dto)
		apiResponse, _ := defaultServer.taskReport(jsonByte)
		// 上报成功之后，本地移除
		if apiResponse.StatusCode == 200 {
			jobList.Remove(receiver.jobContext.id)
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

	flog.Infof("任务：%d 运行完成，耗时：%s，结果：%s", receiver.jobContext.id, receiver.jobContext.sw.GetMillisecondsText(), receiver.jobContext.status.String())
}
