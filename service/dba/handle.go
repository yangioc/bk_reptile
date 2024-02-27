package dba

import (
	"context"
	"fmt"
	"time"

	"github.com/yangioc/bk_pack/proto/dtomsg"
	"github.com/yangioc/bk_pack/util"
	"google.golang.org/protobuf/proto"
)

func (self *Handle) CreateCoolpcData(uuid string, payload []byte) error {

	dbaReq, err := proto.Marshal(&dtomsg.Dto_Msg{
		Type:    "notice",
		Request: "create.coolpcdata",
		Data:    payload,
	})
	if err != nil {
		panic(err)
	}

	msg, err := util.MsgEncode(&dtomsg.Dto_Base{
		UUID:           uuid,
		StartTime:      util.ServerTimeNow().UnixMicro(),
		ExpirationTime: util.ServerTimeNow().Add(5 * time.Second).UTC().UnixMicro(),
		Payload:        dbaReq,
	})
	if err != nil {
		panic(err)
	}

	if err := self.websocket.Send(context.TODO(), msg); err != nil {
		panic(err)
	}

	return nil
}

func (self *Handle) resChanNew(uuid string) (chan *dtomsg.Dto_Msg_Res, error) {
	resChan := make(chan *dtomsg.Dto_Msg_Res)
	if _, isLoad := self.requestTracing.LoadOrStore(uuid, resChan); isLoad {
		return nil, fmt.Errorf("[Error][getResChan] %v.", uuid)
	}

	return resChan, nil
}

func (self *Handle) resChanLoadAndDelete(uuid string) (chan *dtomsg.Dto_Msg_Res, bool) {
	resChan, ok := self.requestTracing.LoadAndDelete(uuid)
	if !ok {
		return nil, false
	}
	return resChan.(chan *dtomsg.Dto_Msg_Res), true
}

func (self *Handle) CreateEfish(uuid string, payload []byte) error {
	return self.CommonCreate(uuid, "notice", "create.efish", payload)
}

func (self *Handle) CreateStockAnalysis(uuid string, payload []byte) error {
	return self.CommonCreate(uuid, "notice", "create.stockanalysis", payload)
}

func (self *Handle) CreateStockIndex(uuid string, payload []byte) error {
	return self.CommonCreate(uuid, "notice", "create.stockindex", payload)
}

func (self *Handle) CreateStockMarket(uuid string, payload []byte) error {
	return self.CommonCreate(uuid, "notice", "create.stockmarket", payload)
}

func (self *Handle) CreateStockClosePrice(uuid string, payload []byte) error {
	return self.CommonCreate(uuid, "notice", "create.stockclosePrice", payload)
}

func (self *Handle) CreateStockThreefoundationTotal(uuid string, payload []byte) error {
	return self.CommonCreate(uuid, "notice", "create.stockthreefoundationtotal", payload)
}

func (self *Handle) CreateThreefoundationStockDay(uuid string, payload []byte) error {
	return self.CommonCreate(uuid, "notice", "create.stockthreefoundationstockday", payload)
}
