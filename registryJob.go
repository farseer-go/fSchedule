package fSchedule

import (
	"github.com/farseer-go/fs/flog"
	"time"
)

var registerNotify = make(chan struct{}, 100)

// RegistryJob 定时10分钟，注册一次任务，防止掉线
func RegistryJob() {
	for {
		select {
		case <-time.NewTicker(10 * time.Minute).C:
		case <-registerNotify:
			// 收到通知的1秒后，执行。（这时还有可能会有多个任务在添加进来）
			time.Sleep(time.Second)
		}

		// 清空消息通知
		registerNotify = make(chan struct{}, 100)

		// 注册
		err := defaultClient.RegistryClient()
		flog.ErrorIfExists(err)
	}
}
