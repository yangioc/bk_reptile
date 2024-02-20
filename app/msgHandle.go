package app

import (
	"github.com/nats-io/nats.go"
	"github.com/yangioc/bk_pack/log"
	"github.com/yangioc/bk_pack/proto/dtomsg"
	"github.com/yangioc/bk_pack/util"
	"google.golang.org/protobuf/proto"
)

func (self *Handle) messageSub() error {
	readChan := make(chan *nats.Msg, 1024)
	defer close(readChan)
	if err := self.messageQ.SubChan("reptile.command.>", readChan); err != nil {
		log.Errorf("messageHandle: %v", err)
		return err
	}
	defer func() {
		if err := self.messageQ.UnSub("reptile.command.>"); err != nil {
			log.Errorf("app messageSub UnSub error: %v", err)
		}

		log.Infof("UnSub reptile.command.>")
		// TODO:可能需要清空chan流程
	}()

	for {
		select {

		case <-self.ctx.Done():
			return ctxDoneError

		case natsMsg := <-readChan:
			log.Info(natsMsg.Data)

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

			case "getcoolpc":
				self.GetCoolpc()

			}
		}
	}
}
