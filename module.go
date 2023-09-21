package fSchedule

import (
	"fmt"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
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
	// 服务端配置
	defaultServer = serverVO{
		Address: configure.GetSlice("FSchedule.Server.Address"),
		Token:   configure.GetString("FSchedule.Server.Token"),
	}

	// 客户端配置
	NewClient()
	timingWheel.Start()

	// 初始化日志队列
	logQueue = make(chan logContent, 2048)

}

func (module Module) PostInitialize() {
	webapi.Area("/api/", func() {
		webapi.RegisterPOST("/check", Check)
		webapi.RegisterPOST("/invoke", Invoke)
		webapi.RegisterPOST("/status", Status)
		webapi.RegisterPOST("/kill", Kill)
	})
	webapi.UseApiResponse()
	webapi.UsePprof()
	go webapi.Run(fmt.Sprintf("%s:%d", defaultClient.ClientIp, defaultClient.ClientPort))

	fs.AddInitCallback("开启上传调度中心日志", func() {
		go enableReportLog()
	})

	// 注册健康检查
	container.RegisterInstance[core.IHealthCheck](&healthCheck{}, "fSchedule")
}

func (module Module) Shutdown() {
	defaultClient.LogoutClient()
}
