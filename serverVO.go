package fSchedule

import (
	"github.com/farseer-go/collections"
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
	Token   string
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
	var apiResponse core.ApiResponse[any]
	err := http.NewClient(address+"/api/registry").HeadAdd(tokenName, receiver.Token).Body(bodyJson).PostUnmarshal(&apiResponse)
	return apiResponse, err
}

// 服务端下线接口
func (receiver *serverVO) logout(bodyJson []byte) (core.ApiResponse[any], error) {
	address, _ := receiver.getAddress(-1)
	var apiResponse core.ApiResponse[any]
	err := http.NewClient(address+"/api/logout").HeadAdd(tokenName, receiver.Token).Body(bodyJson).PostUnmarshal(&apiResponse)
	return apiResponse, err
}

type TaskReportDTO struct {
	Id           int64                                  // 主键
	Name         string                                 // 实现Job的特性名称（客户端识别哪个实现类）
	Data         collections.Dictionary[string, string] // 数据
	NextTimespan int64                                  // 下次执行时间
	Progress     int                                    // 当前进度
	Status       TaskStatus                             // 执行状态
	RunSpeed     int64                                  // 执行速度
}

// 客户端回调
func (receiver *serverVO) taskReport(bodyJson []byte) (core.ApiResponse[any], error) {
	address, _ := receiver.getAddress(-1)
	var apiResponse core.ApiResponse[any]
	err := http.NewClient(address+"/api/taskReport").HeadAdd(tokenName, receiver.Token).Body(bodyJson).PostUnmarshal(&apiResponse)
	return apiResponse, err
}

// 上传日志
func (receiver *serverVO) logReport(bodyJson []byte) (core.ApiResponse[any], error) {
	address, _ := receiver.getAddress(-1)
	var apiResponse core.ApiResponse[any]
	err := http.NewClient(address+"/api/logReport").HeadAdd(tokenName, receiver.Token).Body(bodyJson).PostUnmarshal(&apiResponse)
	return apiResponse, err
}
