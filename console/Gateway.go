package console

import (
	"fmt"
	"github.com/ctfang/command"
	"github.com/ctfang/goworker/lib"
	"github.com/ctfang/goworker/lib/gateway"
	"github.com/ctfang/network"
	"log"
)

type Gateway struct {
	Name string
}

func (self *Gateway) Configure() command.Configure {
	self.Name = "gateway"
	return command.Configure{
		Name:        self.Name,
		Description: "网关进程gateway",
		Input: command.Argument{
			Argument: []command.ArgParam{
				{Name: "runType", Description: "执行操作：start、stop、status"},
			},
			Option: []command.ArgParam{
				{Name: "gateway", Default: "127.0.0.1:8080", Description: "网关地址websocket"},
				{Name: "register", Default: "127.0.0.1:1238", Description: "注册中心"},
				{Name: "worker", Default: "127.0.0.1:4000", Description: "内部通讯地址"},
				{Name: "secret", Default: "", Description: "通讯秘钥"},
			},
		},
	}
}

func (self *Gateway) Execute(input command.Input) {
	switch input.GetArgument("runType") {
	case "start":
		self.start(input)
	case "stop":
		self.stop(input)
	case "status":
		self.status(input)
	}
}

func (self *Gateway) start(input command.Input) {
	lib.Config.SetInput(input)

	// 启动一个内部通讯tcp server
	worker := network.NewServer(lib.Config.Worker.Origin)
	worker.SetEvent(gateway.NewWorkerEvent())
	go worker.ListenAndServe()

	// 连接到注册中心
	register := network.NewClient(lib.Config.Register.Origin)
	register.SetEvent(gateway.NewRegisterEvent())
	go register.ListenAndServe()

	// 启动对客户端的websocket连接
	// network.WebsocketMessageType = network.BinaryMessage
	server := network.NewServer(lib.Config.Gateway.Origin)
	server.SetEvent(gateway.NewWebSocketEvent())
	server.ListenAndServe()
}

func (self *Gateway) status(input command.Input) {
	log.Println("未做")
}

func (self *Gateway) stop(input command.Input) {
	err := command.StopSignal(self.Name)
	if err != nil {
		fmt.Println("停止失败")
	}
	fmt.Println("停止成功")
}
