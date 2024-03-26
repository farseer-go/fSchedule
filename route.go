// 该文件由fsctl route命令自动生成，请不要手动修改此文件
package fSchedule

import (
	"github.com/farseer-go/webapi"
	"github.com/farseer-go/webapi/context"
)

var route = []webapi.Route{
	{"POST", "/api/check", Check, "", []context.IFilter{}, []string{"clientId"}},
	{"POST", "/api/invoke", Invoke, "", []context.IFilter{}, []string{"task"}},
	{"POST", "/api/status", Status, "", []context.IFilter{}, []string{"TaskId"}},
	{"POST", "/api/kill", Kill, "", []context.IFilter{}, []string{"taskId"}},
}
