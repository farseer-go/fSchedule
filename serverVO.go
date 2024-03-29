package fSchedule

import (
	"fmt"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fSchedule/executeStatus"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/utils/http"
	"math/rand"
)

const tokenName = "FSchedule-ACCESS-TOKEN"

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

type RegistryResponse struct {
	ClientIp   string // 客户端IP
	ClientPort int    // 客户端端口
}

// 服务端注册接口
func (receiver *serverVO) registry(bodyJson []byte) (core.ApiResponse[RegistryResponse], error) {
	traceContext := container.Resolve[trace.IManager]().EntryTask("向调度中心注册客户端")
	defer traceContext.End()

	address, _ := receiver.getAddress(-1)
	var apiResponse core.ApiResponse[RegistryResponse]
	_, err := http.NewClient(address+"/api/registry").HeadAdd(tokenName, receiver.Token).Body(bodyJson).Timeout(5000).PostUnmarshal(&apiResponse)
	if err != nil {
		err = fmt.Errorf("注册调度中心失败：%s", err.Error())
	} else if apiResponse.StatusCode != 200 {
		err = fmt.Errorf("注册调度中心失败：%d %s", apiResponse.StatusCode, apiResponse.StatusMessage)
	}

	traceContext.Error(err)
	return apiResponse, err
}

// 服务端下线接口
func (receiver *serverVO) logout(bodyJson []byte) (core.ApiResponse[any], error) {
	traceContext := container.Resolve[trace.IManager]().EntryTask("向调度中心注销客户端")
	defer traceContext.End()

	address, _ := receiver.getAddress(-1)
	var apiResponse core.ApiResponse[any]
	_, err := http.NewClient(address+"/api/logout").HeadAdd(tokenName, receiver.Token).Body(bodyJson).PostUnmarshal(&apiResponse)
	traceContext.Error(err)

	return apiResponse, err
}

type TaskReportDTO struct {
	Id           int64                                  // 主键
	Ver          int                                    // 任务版本
	Name         string                                 // 实现Job的特性名称（客户端识别哪个实现类）
	Data         collections.Dictionary[string, string] // 数据
	NextTimespan int64                                  // 下次执行时间
	Progress     int                                    // 当前进度
	Status       executeStatus.Enum                     // 执行状态
	FailRemark   string                                 // 失败原因
	ResourceVO
}

// 客户端回调
func (receiver *serverVO) taskReport(bodyJson []byte) (core.ApiResponse[any], error) {
	address, _ := receiver.getAddress(-1)
	var apiResponse core.ApiResponse[any]
	_, err := http.NewClient(address+"/api/taskReport").HeadAdd(tokenName, receiver.Token).Body(bodyJson).PostUnmarshal(&apiResponse)
	return apiResponse, err
}

// 上传日志
func (receiver *serverVO) logReport(bodyJson []byte) (core.ApiResponse[any], error) {
	address, _ := receiver.getAddress(-1)
	var apiResponse core.ApiResponse[any]
	_, err := http.NewClient(address+"/api/logReport").HeadAdd(tokenName, receiver.Token).Body(bodyJson).PostUnmarshal(&apiResponse)
	return apiResponse, err
}
