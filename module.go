package fSchedule

import (
	"fmt"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/modules"
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
}

func (module Module) Initialize() {

}

func (module Module) PostInitialize() {
	webapi.Area("/api/", func() {
		webapi.RegisterPOST("/check", Check)
		webapi.RegisterPOST("/invoke", Invoke)
		webapi.RegisterPOST("/status", Status)
		webapi.RegisterPOST("/kill", Kill)
	})
	webapi.UseApiResponse()
	webUrl := fmt.Sprintf("%s:%d", defaultClient.ClientIp, defaultClient.ClientPort)
	go webapi.Run(webUrl)

	fs.AddInitCallback(func() {
		// 注册客户端
		defaultClient.RegistryClient()
	})
}

func (module Module) Shutdown() {
	defaultClient.LogoutClient()
}
