package dba

import (
	"bk_reptile/config"
	"bk_reptile/model/websocket"
	"bk_reptile/model/websocket/socketclient"
	"context"
	"fmt"

	"github.com/yangioc/bk_pack/dto"
	"github.com/yangioc/bk_pack/log"
	"github.com/yangioc/bk_pack/proto/dtomsg"
	"github.com/yangioc/bk_pack/util"
)

type IHandle interface {
	Launch() error
	OnClose(token string)

	CreateCoolpcData(uuid string, payload []byte) error
	CreateEfish(uuid string, payload []byte) error
}

type Handle struct {
	websocket      *websocket.Client
	requestTracing dto.FastSyncMap // weboskcet 請求回傳通道管理
}

func New(setting config.Env) *Handle {
	return &Handle{}
}

func (self *Handle) Launch() error {
	if self.websocket == nil {
		self.websocket = websocket.NewClient(self)
	} else if self.websocket.Handler != nil {
		self.websocket.Close(1000, "new socket connect")
	}

	log.Infof("websocket conn to: %s", config.EnvInfo.Service.DBA.Addr)
	defer log.Infof("websocket diconn: %s", config.EnvInfo.Service.DBA.Addr)
	return self.websocket.Launch(config.EnvInfo.Service.DBA.Addr)
}

func (self *Handle) ReceiveMessage(ctx context.Context, socketClient *socketclient.Handler, message []byte) {
	msg, err := util.MsgDecode(message)
	if err != nil {
		panic(err)
	}

	res := &dtomsg.Dto_Msg_Res{}
	if err := util.Unmarshal(msg.Payload, res); err != nil {
		panic(err)
	}

	resChan, ok := self.resChanLoadAndDelete(msg.UUID)
	if !ok {
		panic(fmt.Errorf("requestTracing not find %s", msg.UUID))
	}

	resChan <- res
}

func (self *Handle) OnClose(token string) {}
