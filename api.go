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

type TaskReportVO struct {
	NextTimespan int64                                  // 下次执行时间
	Progress     int                                    // 当前进度
	Status       TaskStatus                             // 执行状态
	RunSpeed     int64                                  // 执行速度
	Data         collections.Dictionary[string, string] // 数据
}

type TaskStatus int

const (
	None         TaskStatus = iota //  未开始
	Scheduling                     //  调度中
	ScheduleFail                   //  调度失败
	Working                        //  执行中
	Fail                           //  失败
	Success                        //  完成
)

// Check 检查客户端存活
func Check() ResourceVO {
	return getResource()
}

// Invoke 下发任务
func Invoke(task TaskEO) ResourceVO {
	return getResource()
}

// Status 查询任务状态
func Status(TaskId int64) TaskReportVO {
	return TaskReportVO{}
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
