package fSchedule

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fSchedule/executeStatus"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/core/eumLogLevel"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/utils/ws"
	"sync"
	"time"
)

var mapClient = sync.Map{}

func connectFScheduleServer(job ClientVO) {
	for {
		address := defaultServer.getAddress()
		var err error
		job.client, err = ws.Connect(address, 8192)
		job.client.AutoExit = false
		if err != nil {
			flog.Warningf("[%s]调度中心连接失败：%s", job.Name, err.Error())
			time.Sleep(3 * time.Second)
			continue
		}
		mapClient.Store(job.Name, job)
		// 连接成功后，需要先注册
		err = job.client.Send(sendDTO{Type: -1, Registry: registryDTO{ClientId: core.AppId, ClientName: core.AppName, Job: job}})
		if err != nil {
			flog.Warningf("[%s]调度中心注册失败：%s", job.Name, err.Error())
			time.Sleep(3 * time.Second)
			continue
		}

		for {
			// 接收调度请求
			var dto receiverDTO
			err = job.client.Receiver(&dto)
			if err != nil {
				if job.client.IsClose() {
					mapClient.Delete(job.Name)
					flog.Warningf("[%s]调度中心服务端：%s 已断开连接，将在3秒后重连", job.Name, address)
					break
				}
				flog.Warningf("[%s]接收调度中心数据时失败：%s", job.Name, err.Error())
				continue
			}

			switch dto.Type {
			// 新任务
			case 0:
				go invokeJob(job, dto.Task)
			// 停止任务
			case 1:
				flog.Infof("[%s]%d收到Kill请求", job.Name, dto.Task.Id)
				if jContext, exists := taskList.Load(dto.Task.Id); exists {
					jobContext := jContext.(*JobContext)
					jobContext.Remark("FOPS主动停止任务")
					jobContext.status = executeStatus.Fail
					jobContext.clientJob.report(jobContext)
					jobContext.cancel()
					flog.Infof("[%s]%d主动停止任务", job.Name, dto.Task.Id)
				}
			}
		}

		// 断开后重连
		<-job.client.Ctx.Done()
		time.Sleep(3 * time.Second)
	}
}

type registryDTO struct {
	ClientId   int64  // 客户端ID
	ClientName string // 客户端名称
	Job        ClientVO
}

// 从服务端发送
type sendDTO struct {
	Type       int // 发送消息的类型
	Registry   registryDTO
	Log        logDTO
	TaskReport taskReportDTO
}

// 从服务端接收
type receiverDTO struct {
	Type int // 接收消息的类型
	Task taskDTO
}

// taskDTO 任务记录
type taskDTO struct {
	Id      int64                                  // 主键
	Caption string                                 // 任务组标题
	Name    string                                 // 实现Job的特性名称（客户端识别哪个实现类）
	StartAt time.Time                              // 开始时间
	Data    collections.Dictionary[string, string] // 本次执行任务时的Data数据
}

type taskReportDTO struct {
	Id           int64                                  // 主键
	Name         string                                 // 实现Job的特性名称（客户端识别哪个实现类）
	Data         collections.Dictionary[string, string] // 数据
	NextTimespan int64                                  // 下次执行时间
	Progress     int                                    // 当前进度
	Status       executeStatus.Enum                     // 执行状态
	FailRemark   string                                 // 失败原因
	resourceVO
}

// resourceVO 客户端资源情况
type resourceVO struct {
	QueueCount int // 排队中的任务数量
	WorkCount  int // 正在处理的任务数量
}

type logDTO struct {
	TaskId   int64                                  // 主键
	Ver      int                                    // 版本
	Name     string                                 // 实现Job的特性名称（客户端识别哪个实现类）
	Caption  string                                 // 任务标题
	Data     collections.Dictionary[string, string] // 本次执行任务时的Data数据
	LogLevel eumLogLevel.Enum
	CreateAt int64
	Content  string
}
