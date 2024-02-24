package app

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/yangioc/bk_pack/log"
	"github.com/yangioc/bk_pack/proto/dtomsg"
	"github.com/yangioc/bk_pack/util"
	"google.golang.org/protobuf/proto"
)

// ReceiveMessageMQ 訊號處理
func (self *Handle) ReceiveMessageMQ(natsMsg *nats.Msg) {

	msg, err := util.MsgDecode(natsMsg.Data)
	if err != nil {
		panic(err)
	}

	dtoMsg := dtomsg.Dto_Msg{}
	if err := util.Unmarshal(msg.Payload, &dtoMsg); err != nil {
		panic(err)
	}

	var dbaRes []byte
	switch dtoMsg.Type {
	case "req":
		dbaRes, err = self.handleReqMessage(&dtoMsg)
		if err != nil {
			panic(err)
		}

		// 封裝 Response 格式
		res := dtomsg.Dto_Msg_Res{
			Request: dtoMsg.Request,
			Data:    dbaRes,
			State:   1,
		}

		msg.Payload, err = proto.Marshal(&res)
		if err != nil {
			panic(err)
		}

		var message []byte
		if message, err = util.MsgEncode(msg); err != nil {
			panic(err)
		} else if err = natsMsg.Respond(message); err != nil {
			panic(err)
		}
	case "notic":
		err = self.handleNoticMessage(&dtoMsg)
		if err != nil {
			panic(err)
		}
	}
}

// ReceiveMessageMQChan 訊號處理
func (self *Handle) ReceiveMessageMQChan(ctx context.Context, msgChan chan *nats.Msg) {

	for {
		select {
		case <-ctx.Done():
			log.Info("app ctx done.")
			return

		case natsMsg := <-msgChan:
			msg, err := util.MsgDecode(natsMsg.Data)
			if err != nil {
				panic(err)
			}

			dtoMsg := dtomsg.Dto_Msg{}
			if err := util.Unmarshal(msg.Payload, &dtoMsg); err != nil {
				panic(err)
			}

			var dbaRes []byte
			switch dtoMsg.Type {
			case "req":
				dbaRes, err = self.handleReqMessage(&dtoMsg)
				if err != nil {
					panic(err)
				}

				// 封裝 Response 格式
				res := dtomsg.Dto_Msg_Res{
					Request: dtoMsg.Request,
					Data:    dbaRes,
					State:   1,
				}

				msg.Payload, err = proto.Marshal(&res)
				if err != nil {
					panic(err)
				}

				var message []byte
				if message, err = util.MsgEncode(msg); err != nil {
					panic(err)
				} else if err = natsMsg.Respond(message); err != nil {
					panic(err)
				}
			case "notic":
				err = self.handleNoticMessage(&dtoMsg)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
