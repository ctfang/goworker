package console

import (
	"github.com/ctfang/command"
	"github.com/ctfang/goworker/apps"
	"github.com/ctfang/goworker/lib"
	"github.com/ctfang/goworker/lib/worker"
	"github.com/ctfang/network"
	"time"
)

type Worker struct {
}

func (Worker) Configure() command.Configure {
	return command.Configure{
		Name:        "worker",
		Description: "业务worker进程",
		Input: command.Argument{
			Argument: []command.ArgParam{
				{Name: "runType", Description: "执行操作：start、stop、status"},
			},
			Option: []command.ArgParam{
				{Name: "register", Default: "127.0.0.1:1238", Description: "注册中心"},
				{Name: "secret", Default: "", Description: "通讯秘钥"},
			},
		},
	}
}

func (Worker) Execute(input command.Input) {
	lib.Config.SetInput(input)
	lib.BussinessEvent = &apps.MainEvent{}

	// 连接到注册中心
	register := network.NewClient(lib.Config.Register.Origin)
	register.SetEvent(worker.NewRegisterEvent())
	register.ListenAndServe()

	// 断线重连
	for {
		ticker := time.NewTicker(time.Second * 2)
		select {
		case <-ticker.C:
			register.ListenAndServe()
		}
	}

}