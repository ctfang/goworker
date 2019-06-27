package gateway

import (
	"github.com/ctfang/goworker/lib"
	"github.com/ctfang/network"
	"log"
)

/**
客户端信息
*/
type GatewayHeader struct {
	// 内部通讯地址 , 对应本机地址
	LocalIp      uint32
	LocalPort    uint16
	ClientIp     uint32
	ClientPort   uint16
	GatewayPort  uint16
	ConnectionId uint32
	flag         uint8
}

type WebSocketEvent struct {
	// 内部通讯地址
	WorkerServerIp   string
	WorkerServerPort uint16
}

func (ws *WebSocketEvent) OnError(listen network.ListenTcp, err error) {

}

func (ws *WebSocketEvent) OnStart(listen network.ListenTcp) {
	ws.WorkerServerIp = lib.Config.Worker.Ip
	ws.WorkerServerPort = lib.Config.Worker.Port

	log.Println("ws server listening at: ", listen.Url().Origin)
}

func (ws *WebSocketEvent) GetClientId(client network.Connect) string {
	return network.Bin2hex(network.Ip2long(ws.WorkerServerIp), ws.WorkerServerPort, client.Id())
}

func (ws *WebSocketEvent) OnConnect(client network.Connect) {
	client.SetUid(ws.GetClientId(client))
	// 添加连接池
	Router.AddedClient(client)

	ws.SendToWorker(client, lib.CMD_ON_CONNECT, []byte(""))
}

func (ws *WebSocketEvent) OnMessage(c network.Connect, message []byte) {
	ws.SendToWorker(c, lib.CMD_ON_MESSAGE, message)
}

func (ws *WebSocketEvent) OnClose(c network.Connect) {
	ws.SendToWorker(c, lib.CMD_ON_CLOSE, []byte(""))
	Router.DeleteClient(c.Id())
}

func (ws *WebSocketEvent) SendToWorker(client network.Connect, cmd uint8, body []byte) {
	msg := lib.GatewayMessage{
		PackageLen:   28 + uint32(len(body)),
		Cmd:          cmd,
		LocalIp:      network.Ip2long(ws.WorkerServerIp),
		LocalPort:    ws.WorkerServerPort,
		ClientIp:     client.GetIp(),
		ClientPort:   client.GetPort(),
		ConnectionId: client.Id(),
		Flag:         1,
		GatewayPort:  lib.Config.Gateway.Port,
		ExtLen:       0,
		ExtData:      "",
		Body:         body,
	}

	worker, err := Router.GetWorker(client)
	if err != nil {
		// worker 找不到 获取连接
		log.Println("主动断开客户端连接 err:", err)
		client.Close()
		Router.DeleteClient(client.Id())
		return
	}
	worker.SendByte(lib.GatewayMessageToByte(msg))
}

func NewWebSocketEvent() network.Event {
	return &WebSocketEvent{}
}
