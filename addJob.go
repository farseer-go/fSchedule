package fSchedule

import (
	"fmt"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fSchedule/executeStatus"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/parse"
	"github.com/robfig/cron/v3"
	"regexp"
	"time"
)

// JobFunc 客户端要执行的JOB
type JobFunc func(jobContext *JobContext) bool

var StandardParser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

type Option struct {
	StartAt int64                                  // 任务开始时间（时间戳秒）
	Data    collections.Dictionary[string, string] // 第一次注册时使用
}
type options func(opt *Option)

// AddJob 客户端支持的任务
func AddJob(isEnable bool, name, caption string, ver int, cronString string, job JobFunc, ops ...options) {
	matched, err := regexp.MatchString("[a-zA-Z0-9\\-_]+", name)
	if err != nil {
		panic(fmt.Sprintf("任务组:%s %s，name格式错误:%s", name, caption, err.Error()))
	}
	if !matched {
		panic(fmt.Sprintf("任务组:%s %s，name格式错误，只允许【字母、数字、_、-】", name, caption))
	}
	_, err = StandardParser.Parse(cronString)
	if err != nil {
		panic(fmt.Sprintf("任务组:%s %s，Cron格式[%s]错误:%s", name, caption, cronString, err.Error()))
	}

	// 说明没有启用调度中心（没有依赖模块）
	if len(defaultServer.Address) < 1 {
		return
	}
	// 设置额度参数
	opt := &Option{Data: collections.NewDictionary[string, string]()}
	for _, op := range ops {
		op(opt)
	}

	// 如果是调试状态，则模拟调度
	if configure.GetBool("FSchedule.Debug.Enable") {
		jobContext := &JobContext{
			Id:           888,
			Name:         name,
			Data:         collections.NewDictionary[string, string](),
			Caption:      caption,
			nextTimespan: 0,
			progress:     0,
			status:       executeStatus.Working,
			StartAt:      time.Now(),
		}
		for k, v := range configure.GetSubNodes("FSchedule.Debug." + name) {
			jobContext.Data.Add(k, parse.ToString(v))
		}
		job(jobContext)
		return
	}

	fs.AddInitCallback("向调度中心注册任务组："+name, func() {
		// 向调度中心注册
		go connectFScheduleServer(ClientVO{
			Name:     name,
			IsEnable: isEnable,
			Caption:  caption,
			Ver:      ver,
			Cron:     cronString,
			jobFunc:  job,
			StartAt:  opt.StartAt,
			Data:     opt.Data,
		})
	})
}
