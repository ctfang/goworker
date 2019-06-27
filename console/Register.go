package console

import (
	"fmt"
	"github.com/ctfang/command"
	"github.com/ctfang/goworker/lib"
	"github.com/ctfang/goworker/lib/register"
	"github.com/ctfang/network"
	"log"
	"os"
)

type Register struct {
	Name string
}

func (self *Register) Configure() command.Configure {
	self.Name = "register"
	return command.Configure{
		Name:        self.Name,
		Description: "注册中心",
		Input: command.Argument{
			Argument: []command.ArgParam{
				{Name: "runType", Description: "执行操作：start、stop、status"},
			},
			Option: []command.ArgParam{
				{Name: "register", Description: "手动设置地址"},
				{Name: "secret", Default: "", Description: "通讯秘钥"},
			},
		},
	}
}

func (self *Register) Execute(input command.Input) {
	switch input.GetArgument("runType") {
	case "start":
		self.start(input)
	case "stop":
		self.stop(input)
	case "status":
		self.status(input)
	}
}

func (self *Register) start(input command.Input) {
	if input.IsDaemon() {
		logFile, _ := os.OpenFile("register.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		log.SetOutput(logFile)
	}
	command.SavePidToFile(self.Name)
	command.ListenStopSignal(func(sig os.Signal) {
		command.DeleteSavePidToFile(self.Name)
		os.Exit(0)
	})

	lib.Config.SetInput(input)

	server := network.NewServer(lib.Config.Register.Origin)
	server.SetEvent(register.NewRegisterEvent())
	server.ListenAndServe()
}

func (self *Register) status(input command.Input) {
	log.Println("未做")
}

func (self *Register) stop(input command.Input) {
	err := command.StopSignal(self.Name)
	if err != nil {
		fmt.Println("停止失败")
	}
	fmt.Println("停止成功")
}

