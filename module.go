package fSchedule

import (
	"github.com/farseer-go/fs"
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
	hostname, _ := os.Hostname()
	defaultClient = clientVO{
		ClientId:   snowflake.GenerateId(),
		ClientName: hostname,
		ClientIp:   fs.AppIp,
		ClientPort: 9526,
	}
}

func (module Module) Initialize() {
}

func (module Module) PostInitialize() {
	// 注册客户端

}

func (module Module) Shutdown() {
}
