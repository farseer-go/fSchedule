package test

import (
	"github.com/farseer-go/fSchedule"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/modules"
	"strconv"
	"testing"
	"time"
)

type startupModule struct {
}

func (module startupModule) DependsModule() []modules.FarseerModule {
	return []modules.FarseerModule{fSchedule.Module{}}
}

func (module startupModule) PostInitialize() {
	for i := 1; i <= 1; i++ {
		fSchedule.AddJob(true, "Hello"+strconv.Itoa(i), "测试HelloJob"+strconv.Itoa(i), 1, "0/1 * * * * ?", 1674571566, func(jobContext *fSchedule.JobContext) bool {
			jobContext.Debug("测试日志")
			jobContext.Tracef("测试日志")
			jobContext.Info("测试日志")
			return true
		})
	}
}

func TestRun(t *testing.T) {
	fs.Initialize[startupModule]("test fSchedule")
	defer fSchedule.Module{}.Shutdown()

	time.Sleep(300000 * time.Second)
}
