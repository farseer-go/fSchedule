package fSchedule

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/utils/system"
)

// ResourceVO 客户端资源情况
type ResourceVO struct {
	QueueCount    int     // 排队中的任务数量
	WorkCount     int     // 正在处理的任务数量
	CpuUsage      float64 // CPU百分比
	MemoryUsage   float64 // 内存百分比
	AllowSchedule bool    // 是否允许调度
}

// Check 检查客户端存活
func Check(clientId int64) ResourceVO {
	v := getResource()
	flog.Debugf("收到调度中心的存活检查： %+v", v)
	if clientId != defaultClient.ClientId {
		exception.ThrowWebExceptionf(403, "客户端ID不一致，当前clientId：%d，接收到的是：%d", defaultClient.ClientId, clientId)
	}
	return v
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
			Data:   collections.Dictionary[string, string]{},
			Status: Fail,
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
	resource := system.GetResource()
	return ResourceVO{
		QueueCount:    defaultClient.QueueCount,
		WorkCount:     defaultClient.WorkCount,
		CpuUsage:      resource.CpuUsagePercent,
		MemoryUsage:   resource.MemoryUsagePercent,
		AllowSchedule: true, // 后面看下这个变量如何控制
	}
}
