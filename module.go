package fSchedule

import (
	"fmt"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
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
	// 调试状态下，不开启与调度中心的通信
	if configure.GetBool("FSchedule.Debug.Enable") {
		flog.Warning("FSchedule当前为调试状态，将模拟调用任务")
		return
	}
	builder := webapi.NewApplicationBuilder()
	builder.Area("/api/", func() {
		builder.RegisterPOST("/check", Check)
		builder.RegisterPOST("/invoke", Invoke)
		builder.RegisterPOST("/status", Status)
		builder.RegisterPOST("/kill", Kill)
	})
	builder.UseApiResponse()
	go builder.Run(fmt.Sprintf("%s:%d", defaultClient.ClientIp, defaultClient.ClientPort))

	fs.AddInitCallback("开启上传调度中心日志", func() {
		go enableReportLog()
	})

	fs.AddInitCallback("开启定时上传任务，防止掉线", func() {
		go RegistryJob()
	})

	// 注册健康检查
	container.RegisterInstance[core.IHealthCheck](&healthCheck{}, "fSchedule")
}

func (module Module) Shutdown() {
	defaultClient.LogoutClient()
}
