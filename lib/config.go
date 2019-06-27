package lib

import (
	"github.com/ctfang/command"
	"github.com/ctfang/network"
)

var BussinessEvent LogicEvent

var Config = config{}

type config struct {
	Gateway   *network.Url// 网关
	Home      *network.Url// 本地地址
	Worker    *network.Url// 业务处理（内部通讯）地址
	Register  *network.Url// 注册中心
	SecretKey string
}

func (self *config) SetInput(input command.Input) {
	if self.Home == nil {
		home := input.GetOption("LanIp")
		if home != "" {
			self.Home = network.NewUrl("pack://" + home)
		} else {
			home := input.GetOption("worker")
			if home!="" {
				self.Home = network.NewUrl("pack://" + home)
			}
		}
	}
	if self.Register == nil {
		register := input.GetOption("register")
		if register!="" {
			self.Register = network.NewUrl("text://" + register)
		}
	}
	if self.Gateway == nil {
		Gateway := input.GetOption("gateway")
		if Gateway!="" {
			self.Gateway = network.NewUrl("ws://" + Gateway)
		}
	}
	if self.Worker == nil {
		Worker := input.GetOption("worker")
		if Worker!="" {
			self.Worker = network.NewUrl("pack://" + Worker)
		}
	}
	if self.SecretKey == "" {
		self.SecretKey = input.GetOption("secret")
	}
}
