package fSchedule

import (
	"fmt"
	"strings"
)

type healthCheck struct {
}

func (c *healthCheck) Check() (string, error) {
	if defaultClient.ClientJobs.Count() == 0 {
		return "FSchedule." + strings.Join(defaultServer.Address, ",") + " 启动时没有任务，跳过检查", nil
	}
	err := defaultClient.RegistryClient()
	var tips string
	if err == nil {
		tips = fmt.Sprintf("客户端(%d) %s:%d 注册成功！", defaultClient.ClientId, defaultClient.ClientIp, defaultClient.ClientPort)
	} else {
		tips = fmt.Sprintf("客户端(%d) %s:%d 注册失败", defaultClient.ClientId, defaultClient.ClientIp, defaultClient.ClientPort)
	}
	err = nil

	return "FSchedule." + strings.Join(defaultServer.Address, ",") + " " + tips, err
}
