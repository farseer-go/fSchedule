package fSchedule

import (
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/utils/http"
	"math/rand"
)

const tokenName = "FSS-ACCESS-TOKEN"

var defaultServer serverVO

type serverIndex = int
type serverAddress = string

// 服务端配置
type serverVO struct {
	Address []string
}

// 随机一个服务端地址
func (receiver *serverVO) getAddress(ignoreIndex int) (serverAddress, serverIndex) {
	count := len(receiver.Address)
	if count == 1 {
		return receiver.Address[0], 0
	}

	index := rand.Intn(count - 1)
	// 如果随机到的索引值与要排除的索引值一样时，则重新随机
	for ignoreIndex == index {
		index = rand.Intn(count - 1)
	}
	return receiver.Address[index], index
}

// 服务端注册接口
func (receiver *serverVO) registry(bodyJson []byte) (core.ApiResponse[any], error) {
	address, _ := receiver.getAddress(-1)
	token := configure.GetString("FSchedule.Server.Token")
	var apiResponse core.ApiResponse[any]
	err := http.NewClient(address+"/api/registry").HeadAdd(tokenName, token).Body(bodyJson).PostUnmarshal(&apiResponse)
	return apiResponse, err
}

// 服务端下线接口
func (receiver *serverVO) logout(bodyJson []byte) (core.ApiResponse[any], error) {
	address, _ := receiver.getAddress(-1)
	token := configure.GetString("FSchedule.Server.Token")
	var apiResponse core.ApiResponse[any]
	err := http.NewClient(address+"/api/logout").HeadAdd(tokenName, token).Body(bodyJson).PostUnmarshal(&apiResponse)
	return apiResponse, err
}
