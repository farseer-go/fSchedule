# FSchedule 分布式调度中心客户端
> 包：`"github.com/farseer-go/fSchedule"`
> 
> 模块：`fSchedule.Module`

- `Document`
    - [English](https://farseer-go.github.io/doc/#/en-us/)
    - [中文](https://farseer-go.github.io/doc/)
    - [English](https://farseer-go.github.io/doc/#/en-us/)
- Source
    - [github](https://github.com/farseer-go/fs)

![](https://img.shields.io/github/stars/farseer-go?style=social)
![](https://img.shields.io/github/license/farseer-go/fSchedule)
![](https://img.shields.io/github/go-mod/go-version/farseer-go/fSchedule)
![](https://img.shields.io/github/v/release/farseer-go/fSchedule)
[![codecov](https://img.shields.io/codecov/c/github/farseer-go/fSchedule)](https://codecov.io/gh/farseer-go/fSchedule)
![](https://img.shields.io/github/languages/code-size/farseer-go/fSchedule)
[![Build](https://github.com/farseer-go/fSchedule/actions/workflows/build.yml/badge.svg)](https://github.com/farseer-go/fSchedule/actions/workflows/build.yml)
![](https://goreportcard.com/badge/github.com/farseer-go/fSchedule)

## 概述
FSchedule是一款跨语言分布式的调度中心

- `高可用（HA）`：多实例的job客户端。同一个任务、同一个job实例只会被调度一次
- `快速搭建`：服务端可运行于docker或k8s下，1分钟即可把服务部署到您的生产环境中
- `轻量级`：低内存（没有客户端连接的时候130m，有任务的时候250m）、低CPU消耗，依赖少。
- `动态执行`：可定时、间隔时间、Cron、或由业务方job动态设定下次执行时间。
- `快速上手`：Farseer.Net.Job（.NET CORE)、`farseer-go/fSchedule(golang)`，可以快速实现一个job
- `可视化`：使用[FOPS](https://github.com/FarseerNet/fops.go)，可以维护任务组，查看任务进度、耗时、日志。

FSchedule服务端：
- [https://github.com/FSchedule/FSchedule](https://github.com/FSchedule/FSchedule) （go版本，持续更新中）

?> 这个是FSchedule服务端的`客户端组件`，使得接入服务端变得非常简单。


> 包：`"github.com/farseer-go/fSchedule"`
>
> 模块：`fSchedule.Module`

## 定义任务
```go
func job1(jobContext *fSchedule.JobContext) bool {
    return true
}
```

**入参**：`fSchedule.JobContext` 任务上下文，包含任务的相关信息，及控制任务下一次执行时间

**出参**：`true`：本次任务执行成功。`false`：执行失败（失败后，服务端会立即重新调度）。

## 添加一个任务
```go
// isEnable：任务是否开启
// name：任务组名称（英文）
// caption：任务组标题（任务说明）
// ver：任务组版本，初始版本号必须为1
// cron：任务计划时间（间隔时间）
// startAt：任务开始时间，单位：时间戳，秒（在这个时间之后才会开始）
// job：任务执行的函数本体
func AddJob(isEnable bool, name, caption string, ver int, cron string, startAt int64, job JobFunc)
```

示例：
```go
fSchedule.AddJob(true, "Hello"+strconv.Itoa(i), "测试HelloJob"+strconv.Itoa(i), 1, "0/1 * * * * ?", 1674571566, job1)
```
`fSchedule.AddJob`向服务端注册任务。需要放到模块中`PostInitialize`方法执行

任务组版本说明，如果需要向服务端修改任务组属性，则要在原版本号的基础下+1，否则无效。

如果原来服务端版本为3，本次想修改caption，则应将ver改为4。这时服务端才会修改，否则忽略。


## 完整示例
```go

import (
	"github.com/farseer-go/fSchedule"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/modules"
	"strconv"
	"testing"
	"time"
)

// 启动模块
type startupModule struct {}
func (module startupModule) DependsModule() []modules.FarseerModule {
	return []modules.FarseerModule{fSchedule.Module{}}
}
func (module startupModule) PreInitialize() {}
func (module startupModule) Initialize() {}
func (module startupModule) PostInitialize() {
	// 在这里注册任务
    fSchedule.AddJob(true, "Hello1", "测试HelloJob1", 1, "0/1 * * * * ?", 1674571566, job1)
}
func (module startupModule) Shutdown() {}

// 主函数
func main() {
	fs.Initialize[startupModule]("test fSchedule")
	time.Sleep(300000 * time.Second)
}

// 执行任务
func job1(jobContext *fSchedule.JobContext) bool {
	// db.delete()... 比如清理日志数据
    return true
}
```