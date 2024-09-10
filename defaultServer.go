package fSchedule

import (
	"fmt"
	"github.com/farseer-go/fs/flog"
	"math/rand"
	"time"
)

const tokenName = "FSchedule-ACCESS-TOKEN"

var defaultServer serverVO

// 服务端配置
type serverVO struct {
	Address []string
	Token   string
}

// 随机一个服务端地址
func (receiver *serverVO) getAddress() string {
	count := len(receiver.Address)
	if count == 0 {
		flog.Panic("./farseer.yml配置文件没有找到FSchedule.Server.Address的设置")
	}
	var address string
	if count == 1 {
		address = receiver.Address[0]
	} else {
		//address = receiver.Address[rand.IntN(count-1)]
		address = receiver.Address[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(count-1)]
	}

	return fmt.Sprintf("%s/ws/connect", address)
}
