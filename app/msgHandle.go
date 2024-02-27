package app

import (
	"github.com/nats-io/nats.go"
	"github.com/yangioc/bk_pack/log"
	"github.com/yangioc/bk_pack/proto/dtomsg"
	"github.com/yangioc/bk_pack/util"
	"google.golang.org/protobuf/proto"
)

func (self *Handle) messageAll(topic string) {
	if err := self.messageQ.Sub(topic, self.msgallres); err != nil {
		panic(err)
	}
}
func (self *Handle) msgallres(msg *nats.Msg) {
	log.Debug("task:" + msg.Subject)
}

func (self *Handle) messageSub() error {
	topic := "reptile.command.>"

	// self.messageAll(topic)

	readChan := make(chan *nats.Msg, 1024)
	// defer close(readChan)
	if err := self.messageQ.SubChan(topic, readChan); err != nil {
		log.Errorf("messageHandle: %v", err)
		return err
	}
	log.Infof("[MessageQ][Subscribe] %s", topic)
	defer func() {
		if err := self.messageQ.UnSub(topic); err != nil {
			log.Errorf("app messageSub UnSub error: %v", err)
		}

		log.Infof("[MessageQ][UnSubscribe] %s", topic)
		// TODO:可能需要清空chan流程
	}()

	for {
		select {

		case <-self.ctx.Done():
			return ctxDoneError

		case natsMsg := <-readChan:
			msg, err := util.MsgDecode(natsMsg.Data)
			if err != nil {
				panic(err)
			}

			dtoMsg := dtomsg.Dto_Msg{}
			if err := proto.Unmarshal(msg.Payload, &dtoMsg); err != nil {
				panic(err)
			}

			// var dbaRes []byte
			switch dtoMsg.Request {
			case "getefish":
				self.GetEfish()

			case "getoldefish":
				// self.GetOldEfish()

			case "getcoolpc":
				self.GetCoolpc()

			case "getstock":
				self.GetStock()

			}
		}
	}
}
