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

func (module startupModule) PreInitialize() {
}

func (module startupModule) Initialize() {
}

func (module startupModule) PostInitialize() {
	for i := 1; i <= 10; i++ {
		fSchedule.AddJob(true, "Hello"+strconv.Itoa(i), "测试HelloJob"+strconv.Itoa(i), 1, "0/1 * * * * ?", 1674571566, func(jobContext *fSchedule.JobContext) bool {
			time.Sleep(10 * time.Second)
			return true
		})
	}
}

func (module startupModule) Shutdown() {
}

func TestRun(t *testing.T) {
	fs.Initialize[startupModule]("test fSchedule")
	defer fSchedule.Module{}.Shutdown()

	time.Sleep(300000 * time.Second)
}
