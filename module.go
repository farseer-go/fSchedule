package fSchedule

import (
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/modules"
	"github.com/farseer-go/fs/timingWheel"
	"github.com/farseer-go/webapi"
)

type Module struct {
}

func (module Module) DependsModule() []modules.FarseerModule {
	return []modules.FarseerModule{webapi.Module{}}
}

func (module Module) PreInitialize() {
	// 配置时间轮
	timingWheel.Start()

	// 调试状态下，不开启与调度中心的通信
	if configure.GetBool("FSchedule.Debug.Enable") {
		flog.Warning("FSchedule当前为调试状态，将模拟调用任务")
		return
	}

	// 服务端配置
	defaultServer = serverVO{
		Address: configure.GetSlice("FSchedule.Server.Address"),
		Token:   configure.GetString("FSchedule.Server.Token"),
	}

	if len(defaultServer.Address) < 1 {
		panic("调度中心的地址[FSchedule.Server.Address]未配置")
	}

}

func (module Module) PostInitialize() {

}

func (module Module) Shutdown() {
}
