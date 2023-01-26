package fSchedule

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/modules"
	"github.com/farseer-go/fs/snowflake"
	"os"
)

type Module struct {
}

func (module Module) DependsModule() []modules.FarseerModule {
	return nil
}

func (module Module) PreInitialize() {
	// 服务端配置
	defaultServer = serverVO{
		Address: configure.GetSlice("FSchedule.Server.Address"),
	}

	// 客户端配置
	hostname, _ := os.Hostname()
	defaultClient = clientVO{
		ClientId:   snowflake.GenerateId(),
		ClientName: hostname,
		ClientIp:   fs.AppIp,
		ClientPort: 9526,
		ClientJobs: collections.NewList[ClientJob](),
	}

	// 如果手动配置了客户端IP，则覆盖
	clientIp := configure.GetString("FSchedule.ClientIp")
	if clientIp != "" {
		defaultClient.ClientIp = clientIp
	}

	// 如果手动配置了客户端端口，则覆盖
	clientPort := configure.GetInt("FSchedule.ClientPort")
	if clientPort > 0 {
		defaultClient.ClientPort = clientPort
	}
}

func (module Module) Initialize() {
}

func (module Module) PostInitialize() {
	fs.AddInitCallback(func() {
		// 注册客户端
		defaultClient.RegistryClient()
	})
}

func (module Module) Shutdown() {
}
