package fSchedule

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/exception"
)

// ResourceVO 客户端资源情况
type ResourceVO struct {
	QueueCount    int     // 排队中的任务数量
	WorkCount     int     // 正在处理的任务数量
	CpuUsage      float32 // CPU百分比
	MemoryUsage   float32 // 内存百分比
	AllowSchedule bool    // 是否允许调度
}

// Check 检查客户端存活
func Check(clientId int64) ResourceVO {
	if clientId != defaultClient.ClientId {
		exception.ThrowWebException(403, "客户端ID不一致")
	}
	return getResource()
}

// Invoke 下发任务
func Invoke(task TaskEO) ResourceVO {
	invokeJob(task)
	return getResource()
}

// Status 查询任务状态
func Status(TaskId int64) TaskReportDTO {
	job := getJob(TaskId)
	if job == nil {
		return TaskReportDTO{
			NextTimespan: 0,
			Progress:     0,
			Status:       Fail,
			RunSpeed:     0,
			Data:         collections.Dictionary[string, string]{},
		}
	}
	return job.jobContext.getReport()
}

// Kill 终止任务
func Kill(TaskId int64) {
	// 如果任务无法停止，调用这个异常即可
	exception.ThrowWebException(403, "无法停止任务")
}

// 获取当前客户端的环境信息
func getResource() ResourceVO {
	return ResourceVO{
		QueueCount:    0,
		WorkCount:     0,
		CpuUsage:      0,
		MemoryUsage:   0,
		AllowSchedule: true, // 后面看下这个变量如果控制
	}
}
