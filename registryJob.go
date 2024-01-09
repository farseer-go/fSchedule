package fSchedule

import (
	"github.com/farseer-go/fs/flog"
	"time"
)

// RegistryJob 定时10分钟，注册一次任务，防止掉线
func RegistryJob() {
	for range time.NewTicker(10 * time.Minute).C {
		err := defaultClient.RegistryClient()
		flog.ErrorIfExists(err)
	}
}
