// @area /api/
package fSchedule

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fSchedule/executeStatus"
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
// @post check
func Check(clientId int64) ResourceVO {
	v := getResource()
	flog.Debugf("收到调度中心的存活检查： %+v", v)
	if clientId != defaultClient.ClientId {
		// 说明当前应用ID变了，IP:Port没变。一般出现在本机调试的时候。进程重启后会出现这个情况。
		exception.ThrowWebExceptionf(403, "客户端ID不一致，当前clientId：%d，接收到的是：%d", defaultClient.ClientId, clientId)
	}
	return v
}

// Invoke 下发任务
// @post invoke
func Invoke(task TaskEO) ResourceVO {
	invokeJob(task)
	return getResource()
}

// Status 查询任务状态
// @post status
func Status(TaskId int64) TaskReportDTO {
	job := getJob(TaskId)
	if job == nil {
		return TaskReportDTO{
			Data:       collections.Dictionary[string, string]{},
			Status:     executeStatus.Fail,
			ResourceVO: getResource(),
		}
	}
	dto := job.jobContext.getReport()

	// 任务报告完了，删除当前任务
	if dto.Status == executeStatus.Fail || dto.Status == executeStatus.Success {
		jobList.Remove(TaskId)
	}
	return dto
}

// Kill 终止任务
// @post kill
func Kill(taskId int64) {
	job := getJob(taskId)
	if job == nil {
		// 如果任务无法停止，调用这个异常即可
		exception.ThrowWebException(403, "无法停止任务")
	}
	// 通知取消任务
	job.jobContext.cancel()
}

// 获取当前客户端的环境信息
func getResource() ResourceVO {
	resource := system.GetResource()
	return ResourceVO{
		QueueCount:    jobList.Count() - defaultClient.WorkCount,
		WorkCount:     defaultClient.WorkCount,
		CpuUsage:      resource.CpuUsagePercent,
		MemoryUsage:   resource.MemoryUsagePercent,
		AllowSchedule: true, // 后面看下这个变量如何控制
	}
}
