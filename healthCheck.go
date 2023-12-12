package fSchedule

import (
	"fmt"
	"strings"
)

type healthCheck struct {
}

func (c *healthCheck) Check() (string, error) {
	err := defaultClient.RegistryClient()
	var tips string
	if err == nil {
		tips = fmt.Sprintf("客户端(%d) %s:%d 注册成功！", defaultClient.ClientId, defaultClient.ClientIp, defaultClient.ClientPort)
	} else {
		tips = fmt.Sprintf("客户端(%d) %s:%d ", defaultClient.ClientId, defaultClient.ClientIp, defaultClient.ClientPort)
	}

	return "FSchedule." + strings.Join(defaultServer.Address, ",") + " " + tips, err
}
