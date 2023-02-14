package fSchedule

type TaskStatus int

const (
	None         TaskStatus = iota //  未开始
	Scheduling                     //  调度中
	ScheduleFail                   //  调度失败
	Working                        //  执行中
	Fail                           //  成功
	Success                        //  完成
)

func (receiver TaskStatus) String() string {
	switch receiver {
	case Scheduling:
		return "调度中"
	case ScheduleFail:
		return "调度失败"
	case Working:
		return "执行中"
	case Success:
		return "成功"
	case Fail:
		return "失败"
	}
	return "未开始"
}
