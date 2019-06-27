package gateway

import (
	"encoding/json"
	"github.com/ctfang/goworker/lib"
	"github.com/ctfang/goworker/lib/worker"
	"github.com/ctfang/network"
	"log"
)

// 处理worker发来的信息命令
type WorkerHandle struct {
}

// 单个用户信息
func (*WorkerHandle) OnSendToOne(c network.Connect, message lib.GatewayMessage) {
	client, err := Router.GetClient(message.ConnectionId)
	if err != nil {
		Router.DeleteClient(message.ConnectionId)
		return
	}
	client.Send(message.Body)
}

// 提出用户
func (*WorkerHandle) OnKick(c network.Connect, message lib.GatewayMessage) {
	// todo
}

// 立即摧毁用户连接
func (*WorkerHandle) OnDestroy(c network.Connect, message lib.GatewayMessage) {
	// todo
}

// 发给gateway的向所有用户发送数据
func (*WorkerHandle) OnSendToAll(c network.Connect, message lib.GatewayMessage) {
	for _, client := range Router.Clients {
		client.Send(message.Body)
	}
}

func (*WorkerHandle) OnWorkerConnect(c network.Connect, message lib.GatewayMessage) {
	WorkerKey := &worker.WorkerKey{}
	err := json.Unmarshal(message.Body, WorkerKey)
	if err != nil || lib.Config.SecretKey != WorkerKey.Secret {
		c.Close()
		return
	}
	Router.AddedWorker(c)
}

func (*WorkerHandle) OnGatewayClientConnect(c network.Connect, message lib.GatewayMessage) {
	WorkerKey := &worker.WorkerKey{}
	err := json.Unmarshal(message.Body, WorkerKey)
	if err != nil || lib.Config.SecretKey != WorkerKey.Secret {
		c.Close()
		return
	}
	Router.AddedWorker(c)
}

// 按照条件查找
func (*WorkerHandle) OnSelect(c network.Connect, message lib.GatewayMessage) {
	// todo
}

// 按照条件查找
func (*WorkerHandle) OnTodo(c network.Connect, message lib.GatewayMessage) {
	// todo
	log.Println("todo ", message)
}

func NewWorkerHandle() *WorkerHandle {
	return &WorkerHandle{}
}
