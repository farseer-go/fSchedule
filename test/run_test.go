package test

import (
	"github.com/farseer-go/fSchedule"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/modules"
	"testing"
	"time"
)

type startupModule struct {
}

func (module startupModule) DependsModule() []modules.FarseerModule {
	return []modules.FarseerModule{fSchedule.Module{}}
}

func (module startupModule) PreInitialize() {
}

func (module startupModule) Initialize() {
}

func (module startupModule) PostInitialize() {
	fSchedule.AddJob(true, "Hello1", "测试HelloJob1", 5, "0/1 * * * * ?", 1674571566, func(jobContext *fSchedule.JobContext) bool {
		return true
	})
	fSchedule.AddJob(true, "Hello2", "测试HelloJob2", 1, "0/1 * * * * ?", 1674571566, func(jobContext *fSchedule.JobContext) bool {
		return true
	})
	fSchedule.AddJob(true, "Hello3", "测试HelloJob3", 1, "0/1 * * * * ?", 1674571566, func(jobContext *fSchedule.JobContext) bool {
		return true
	})
	fSchedule.AddJob(true, "Hello4", "测试HelloJob4", 1, "0/1 * * * * ?", 1674571566, func(jobContext *fSchedule.JobContext) bool {
		return true
	})
	fSchedule.AddJob(true, "Hello5", "测试HelloJob5", 1, "0/1 * * * * ?", 1674571566, func(jobContext *fSchedule.JobContext) bool {
		return true
	})
}

func (module startupModule) Shutdown() {
}

func TestRun(t *testing.T) {
	fs.Initialize[startupModule]("test fSchedule")
	defer fSchedule.Module{}.Shutdown()

	time.Sleep(300000 * time.Second)
}
