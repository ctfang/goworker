package apps

import (
	"github.com/ctfang/goworker/lib"
	"log"
)

type MainEvent struct {
}

func (*MainEvent) OnWebSocketConnect(clientId string, header []byte) {

}

func (*MainEvent) OnStart() {
	log.Println("OnStart ")
}

func (*MainEvent) OnConnect(clientId string) {
	log.Println("OnConnect ", clientId)
}

func (*MainEvent) OnMessage(clientId string, body []byte) {
	log.Println("OnMessage ", clientId, string(body))
	lib.SendToAll([]byte(clientId+"中文"))
}

func (*MainEvent) OnClose(clientId string) {
	log.Println("OnClose ", clientId)
}
